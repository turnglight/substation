package models

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Equipment struct{
	Id int32
	Region string
	Code string
	Name string
	MonitorId int32
	Productor string
	LineName string
	State int8
	CreateTime string
	UpdateTime string
}

// 数据包结构
// 其中包头34字节，后面的是数据区
type Sheath struct {
	Type string
	// 包头，2字节
	HeaderTag int16
	// 设备ID，4字节
	MonitorId int32
	// 设备token, 16字节
	Token int64
	// 命令字, 4字节
	CmdType int32
	// 序列号，4字节
	SeqNumber int32
	// 数据区长度 = 数据对应的传感器ID+数据长度+数据 = 4 + 4 + x
	DataTotalLength int32
	// 传感器ID
	DeviceId int32
	// 数据长度
	DataLength int32
	// 数据区
	Data string
	// 公式
	Formula string
	// 最终的数据值
	FinalValue float64
	// 数据接收时间
	ReceiveTime string
	State int8
	Tag string
}

func (Equipment) TableName() string {
	return "monitor_equipment_info"
}

func GetDeviceList() *[]Equipment{
	var list []Equipment
	result := db.Find(&list)
	rowsAffected := result.RowsAffected
	if rowsAffected > 0 {
		return &list
	}
	return nil
}

func GetDataForLine(startTime string, monitorId, cmdType int32) *[]Sheath{
	now := time.Now()
	duration, _ := time.ParseDuration("-24h")
	lastTime := now.Add(duration)
	stime := lastTime.Format("2006-01-02 15:04:05")
	var list []Sheath
	tableName := "monitor_sheath_equipment_" + strconv.Itoa(int(monitorId))
	var result *gorm.DB
	if len(startTime)>0 {
		stime = startTime
	}
	result = db.Table(tableName).Where(" monitor_id = ? and cmd_type = ? and create_time > ?",  monitorId, cmdType, stime).Find(&list)
	fmt.Println(monitorId, cmdType)
	fmt.Println(result.RowsAffected)
	rowsAffected := result.RowsAffected
	if rowsAffected > 0 {
		return &list
	}
	return nil
}




