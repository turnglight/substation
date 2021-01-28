package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Knetic/govaluate"
	"go.uber.org/zap"
	"net"
	"strconv"
	"substation.com/logger"
)

// 8通道信号采集板通讯协议

type Parser interface {
	Handle(buffer []byte)(*[]byte, error)
	Compute()
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
	Data []byte
	// 公式
	Formula string
	// 最终的数据值
	FinalValue float64
}

// WhatError
//	8001: data's format illegal
//	8002: header's length too short
type ProtocolError struct{
	ErrorCode int16
	WhatError string
}


func (e *ProtocolError) Error() string{
	return fmt.Sprintf("error code: %v", e.WhatError)
}

func (sheath *Sheath) Handle(buffer []byte, conn net.Conn) (*[]byte, error){
	logx := logger.NewInstance()
	defer logx.Sync()
	headLength := 34
	var rtnBytes []byte
	if len(buffer) > headLength {
		reader := bytes.NewReader(buffer)
		sheath.Type = "SheathProtocol"
		binary.Read(reader, binary.BigEndian, &sheath.HeaderTag)
		binary.Read(reader, binary.BigEndian, &sheath.MonitorId)
		skip := make([]byte, 8)
		binary.Read(reader, binary.BigEndian, &skip)
		// 在16个字节的token中，前面8个字节跳过，取后面的8个字节为ID
		binary.Read(reader, binary.BigEndian, &sheath.Token)
		binary.Read(reader, binary.BigEndian, &sheath.CmdType)

		binary.Read(reader, binary.BigEndian, &sheath.SeqNumber)
		binary.Read(reader, binary.BigEndian, &sheath.DataTotalLength)

		// 解析完header后， 判断数据区的长度
		if sheath.DataTotalLength < 8 {
			logx.Error("data's format illegal",
				zap.Int("data length", len(buffer)),
				zap.ByteString("data bytes", buffer),
				zap.String("remoteAddr", conn.RemoteAddr().String()),
				zap.String("localAddr", conn.LocalAddr().String()))
			err := ProtocolError{ErrorCode: 8001, WhatError: "data's format illegal"}
			return nil, &err
		}

		binary.Read(reader, binary.BigEndian, &sheath.DeviceId)
		binary.Read(reader, binary.BigEndian, &sheath.DataLength)

		sheath.Data = make([]byte, sheath.DataLength)

		// 数据区总长=4 + 4 + x, 第一个4是传感器的4个字节，第二个4是存储数据长度的4个字段，x是第二个4中存储的数值大小
		// 所以最终的数据实际长度应该等于第二个4中存储的数据大小，而且也等于数据区的长度减8
		binary.Read(reader, binary.LittleEndian, &sheath.Data)
		// 根据cmdType确认数据的计算表达式
		if sheath.CmdType == 5 {
			// 环流电流计算公式
			sheath.Formula = "value*20"
		} else if sheath.CmdType == 10 {
			// 温度计算公式
			sheath.Formula = "-15+(100.00/16)*(value-4)"
		}
		// 如果buffer的长度与数据解析的长度相等，则代表解析结束，否则返回剩余的字节，循环继续解析
		if len(buffer) == (int)(34+sheath.DataTotalLength) {
			return nil, nil
		}
		if len(buffer) < (int)(34+sheath.DataTotalLength) {
			logx.Error("data's format illegal",
				zap.Int("data length", len(buffer)),
				zap.String("data bytes", string(buffer[:])),
				zap.String("remoteAddr", conn.RemoteAddr().String()),
				zap.String("localAddr", conn.LocalAddr().String()))
			err := ProtocolError{ErrorCode: 8001, WhatError: "data's format illegal"}
			return nil, &err
		}
		rtnBytes = buffer[34+sheath.DataTotalLength:]
		return &rtnBytes, nil
	}else{
		logx.Error("data's format illegal",
			zap.String("error", "header's length too short"),
			zap.ByteString("data bytes", buffer),
			zap.String("remoteAddr", conn.RemoteAddr().String()),
			zap.String("localAddr", conn.LocalAddr().String()))
		err := ProtocolError{ErrorCode: 8002, WhatError: "header's length too short"}
		return nil, &err
	}
}

func (sheath *Sheath) Compute() {
	// 初始化表达式
	var expression *govaluate.EvaluableExpression
	expression, _ = govaluate.NewEvaluableExpression(sheath.Formula)
	// byte -> string -> float
	sData := string(sheath.Data[:])
	sheath.FinalValue, _ = strconv.ParseFloat(sData, 64)
	// 传入参数，得到表达式的计算结果
	parameters := make(map[string]interface{}, 8)
	parameters["value"] = sheath.FinalValue
	result, _ := expression.Evaluate(parameters)
	sheath.FinalValue = result.(float64)
}


