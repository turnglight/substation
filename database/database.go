package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)
//
//const (
//	USER_NAME = "root"
//	PASS_WORD = "roottest"
//	//HOST      = "10.190.5.78"
//	HOST      = "192.168.0.2"
//	PORT      = "3111"
//	DATABASE  = "monitor_data_center"
//	CHARSET   = "utf8"
//)

var logx *zap.Logger

func NewConnection(driver, dataSource string) (*sql.DB, error) {
	return sql.Open(driver, dataSource)
}

//func main(){
//	db, err := sql.Open("mysql", "root:roottest@tcp(192.168.0.2:3111)/monitor_data_center")
//	defer db.Close()
//	if err != nil {
//		fmt.Printf("err: %v", err.Error())
//	}
//	rows, _ := db.Query("select id, monitor_id from monitor_info ")
//	for rows.Next() {
//		var id int32
//		var monitorId int32
//		rows.Scan(&id, &monitorId)
//		fmt.Printf("%v  %v", id, monitorId)
//	}
//}







