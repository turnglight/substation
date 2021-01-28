package main

import (
	"bufio"
	list2 "container/list"
	"go.uber.org/zap"
	"net"
	"substation.com/database"
	"substation.com/logger"
	"substation.com/protocol"
	"time"
)

const (
	Region = "hubei"
	Code   = "SheathCurrent"
	Productor   = "伏佳安达"
)

// 自定义日志归档器
var logx *zap.Logger = logger.NewInstance()

func main() {
	defer logx.Sync() // flushes buffer, if any

	InitSocket("tcp", "10.190.5.78:7001")
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
			logx.Info("received data",
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
	pushToMysql(list)
}

func pushToMysql(list *list2.List) {
	db, err := database.NewConnection("mysql", "root:roottest@tcp(192.168.0.2:3111)/monitor_data_center")
	if err != nil {
		logx.Error("connect to mysql unsuccessfully", zap.String("remoteAddress", "root:roottest@tcp(192.168.0.2:3111)/monitor_data_center"))
		return
	}
	db.SetMaxOpenConns(20)
	defer db.Close()
///////////
////lkjlkj
	i := 0
	for e := list.Front(); e != nil; e = e.Next() {
		sheath := e.Value
		cmdType := sheath.(protocol.Sheath).CmdType
		monitorId := sheath.(protocol.Sheath).MonitorId
		if i == 0 {
			querySql := "select id from equipment_info where monior_id =? and code =? and region=?"
			rows, err := db.Query(querySql, monitorId, Code, Region)
			if err != nil {
				logx.Error("")
			}
			// 如果还没有初始化此设备，则进行初始化
			if !rows.Next() {
				name :=  "设备"+string(monitorId)+"("+Region+Code+")"
				now := time.Now()
				now.Format("2006-01-02 15:04:05")
				insertSql := "insert into equipment_info(monitor_id, region, code, name, productor, state, create_time) value(?, ?, ?, ?, ?, ?, ?)"
				db.Exec(insertSql, monitorId, Region, Code, name, Productor, 1, now.Format("2006-01-02 15:04:05"))
			}
		}
		// 先判断是否存在表
	}
}
