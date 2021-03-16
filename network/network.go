package main

import (
	"bufio"
	list2 "container/list"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net"
	"strconv"
	"substation.com/database"
	"substation.com/logger"
	"substation.com/protocol"
	"time"
)

const (
	Region     = "tangshan"
	Short     = "TS"
	Code       = "SheathCurrent"
	Productor  = "伏佳安达"
	DbDriver   = "mysql"
	DataSource = "root:whfjad888@tcp(localhost:3306)/monitor_tangshan"
)

// 文件日志处理器
var logx *zap.Logger = logger.NewInstanceForHook()

// stdout日志处理器
var logf *zap.Logger = logger.NewInstanceForStdout()

func main() {
	defer logx.Sync() // flushes buffer, if any

	//InitSocket("tcp", "10.190.5.78:7001")
	InitSocket("tcp", "10.190.21.34:7020")
}

func InitSocket(netType, netAddress string) {
	listen, error := net.Listen(netType, netAddress)
	defer listen.Close()

	if error != nil {
		logx.Error("failed to open connection",
			zap.String("netType", netType),
			zap.String("netAddress", netAddress),
		)
	}
	logx.Info("listen connection successful",
		zap.String("netType", netType),
		zap.String("netAddress", netAddress))

	for {
		// 等待客户端建立连接
		conn, err := listen.Accept()
		if err != nil {
			logx.Info("accept failed", zap.String("error", err.Error()))
			continue
		}

		// 启动一个单独的 goroutine 去处理连接
		go process(conn)
	}
}

func process(conn net.Conn) {
	logx.Info("successful connected", zap.String("remoteAddr", conn.RemoteAddr().String()), zap.String("localAddr", conn.LocalAddr().String()))
	defer conn.Close()

	list := new(list2.List)
	// 针对当前连接做发送和接受操作
	for {
		// 获取当前连接的reader
		reader := bufio.NewReader(conn)
		// 当前最大接收包的大小设置为2048个字节
		var buf [2048]byte
		n, err := reader.Read(buf[:])
		// 连接结束，EOF -> break
		if err != nil {
			break
		}

		data := buf[:n]

		// 数据包解析
		for {
			sheath := new(protocol.Sheath)
			data, protocolError := sheath.Handle(data, conn)

			if protocolError != nil {
				_, err = conn.Write([]byte(protocolError.Error()))
				return
			}

			// 解析后计算最终value
			sheath.Compute()
			logf.Info("received data",
				zap.Int16("headerTag", sheath.HeaderTag),
				zap.Int32("monitorId", sheath.MonitorId),
				zap.Int64("token", sheath.Token),
				zap.Int32("cmdType", sheath.CmdType),
				zap.Int32("seqNumber", sheath.SeqNumber),
				zap.Int32("dataTotalLength", sheath.DataTotalLength),
				zap.Int32("deviceId", sheath.DeviceId),
				zap.Int32("dataLength", sheath.DataLength),
				zap.String("data", string(sheath.Data[:])),
				zap.Float64("finalValue", sheath.FinalValue))
			list.PushBack(sheath)

			// 如果剩余需要被解析的数据，那么break
			if data == nil {
				break
			}
		}
		// 将接受到的数据返回给客户端
		_, err = conn.Write([]byte("successful"))
		if err != nil {
			logx.Info("write from conn failed",
				zap.String("error", err.Error()),
				zap.String("remoteAddr", conn.RemoteAddr().String()),
				zap.String("localAddr", conn.LocalAddr().String()))
			break
		}
	}
	go pushToMysql(list)
}

func pushToMysql(list *list2.List) {
	db, err := database.NewConnection(DbDriver, DataSource)
	defer db.Close()
	if err != nil {
		logx.Error("connect to mysql unsuccessfully", zap.String("remoteAddress", DataSource))
		return
	}
	db.SetMaxOpenConns(50)
	db.SetConnMaxIdleTime(50)
	pingErr := db.Ping()
	if pingErr != nil {
		logx.Error("sorry, can't connect to mysql", zap.String("remoteAddress", DataSource))
	}

	tx, _ := db.Begin()
	for i, e := 0, list.Front(); e != nil; i, e = i + 1, e.Next() {
		sheath := e.Value
		cmdType := sheath.(*protocol.Sheath).CmdType
		monitorId := sheath.(*protocol.Sheath).MonitorId
		seqNumber := sheath.(*protocol.Sheath).SeqNumber
		receiveTime := sheath.(*protocol.Sheath).ReceiveTime
		deviceId := sheath.(*protocol.Sheath).DeviceId
		data := sheath.(*protocol.Sheath).Data
		formula := sheath.(*protocol.Sheath).Formula
		finalValue := sheath.(*protocol.Sheath).FinalValue
		state := sheath.(*protocol.Sheath).State
		tag := sheath.(*protocol.Sheath).Tag
		tableName := "monitor_sheath_equipment_" + strconv.Itoa(int(monitorId))
		if i == 0 {
			createTable(tableName)
			querySql := "select id from monitor_equipment_info where monitor_id =? and code =? and region=?"
			rows, err := db.Query(querySql, monitorId, Code, Region)
			if err != nil {
				logx.Error("execute query sql failure", zap.String("error", err.Error()))
				panic(err.Error())
			}
			defer rows.Close()
			// 如果还没有初始化此设备，则进行初始化
			if !rows.Next() {
				name := "设备" + strconv.Itoa(int(monitorId)) + "(" + Short + "-" + Code + ")"
				now := time.Now()
				insertSql := "insert into monitor_equipment_info(monitor_id, region, code, name, productor, state, create_time) value(?, ?, ?, ?, ?, ?, ?)"
				rs, err := tx.Exec(insertSql, monitorId, Region, Code, name, Productor, 1, now.Format("2006-01-02 15:04:05"))
				if err != nil {
					logx.Fatal("failed to execute insertSql", zap.String("error", err.Error()))
				}
				_, err = rs.RowsAffected()
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
		insertSql := "insert into " + tableName +
			"(monitor_id, cmd_type, seq_num, receive_time, device_id, tag, data, formula, final_value, state, create_time) " +
			"value(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		rs, err := tx.Exec(insertSql, monitorId, cmdType, seqNumber, receiveTime, deviceId, tag, data, formula, finalValue, state, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			logx.Fatal("failed to execute insertSql", zap.String("error", err.Error()))
		}
		rowAffected, err := rs.RowsAffected()
		if err != nil {
			log.Fatalln(err)
		}
		logx.Info("successful", zap.Int64("effected rows", rowAffected))
	}
	tx.Commit()
}

func createTable(tableName string) (bool, *error) {
	db, err := database.NewConnection(DbDriver, DataSource)
	defer db.Close()
	if err != nil {
		logx.Error("connect to mysql unsuccessfully", zap.String("remoteAddress", DataSource))
		return false, &err
	}
	createTableSql := "create table if not exists " + tableName + "(" +
		"id int(11) not null auto_increment," +
		"monitor_id int(6) not null," +
		"cmd_type int(4) not null," +
		"seq_num int(4) not null," +
		"receive_time datetime not null," +
		"device_id int(4) not null," +
		"tag varchar(20) not null," +
		"formula varchar(40) not null," +
		"data varchar(10) not null," +
		"final_value varchar(10)," +
		"state int(4)," +
		"create_time datetime," +
		"primary key(id)" +
		")"
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("%T\n", tx)
		logx.Error("failed to get tx", zap.String("error", err.Error()))
		return false, nil
	}
	_, err = tx.Exec(createTableSql)
	if err != nil {
		return false, &err
	}
	tx.Commit()
	return true, nil
}
