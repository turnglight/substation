module substation.com/models

go 1.15

replace substation.com/setting => ../setting

require (
	gorm.io/driver/mysql v1.0.4
	gorm.io/gorm v1.20.12
	substation.com/setting v0.0.0-00010101000000-000000000000
)
