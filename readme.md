

+ go get -u github.com/gin-gonic/gin
+ go get -u go.uber.org/zap
+ go get gopkg.in/natefinch/lumberjack.v2
+ go get -u github.com/Knetic/govaluate
+ go get -u github.com/go-sql-driver/mysql


~~~sql
CREATE TABLE `equipment_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
	`region` varchar(20) NOT NULL,
	`code` varchar(20) NOT NULL,
  `monitor_id` int(11) NOT NULL,
  `name` varchar(20) NOT NULL,
  `productor` varchar(50) NOT NULL,
	`state` int(2) NOT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;

~~~