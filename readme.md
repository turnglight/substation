

+ go get -u github.com/gin-gonic/gin
+ go get -u go.uber.org/zap
+ go get gopkg.in/natefinch/lumberjack.v2
+ go get -u github.com/Knetic/govaluate
+ go get -u github.com/go-sql-driver/mysql


~~~sql
CREATE TABLE `equipment_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `region` varchar(20) NOT NULL,
  `code` varchar(20) DEFAULT NULL,
  `monitor_id` int(11) NOT NULL,
  `name` varchar(40) NOT NULL,
  `productor` varchar(50) DEFAULT NULL,
  `state` int(2) NOT NULL,
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

CREATE TABLE `sheath_equepment_4` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `monitor_id` int(6) NOT NULL,
  `cmd_type` int(4) NOT NULL,
  `seq_num` int(4) NOT NULL,
  `receive_time` datetime NOT NULL,
  `device_id` int(4) NOT NULL,
  `formula` varchar(40) NOT NULL,
  `data` varchar(10) NOT NULL,
  `final_value` varchar(10) DEFAULT NULL,
  `state` int(4) DEFAULT NULL,
  `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=149 DEFAULT CHARSET=utf8;

~~~