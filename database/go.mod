module substation.com/database

go 1.15

replace substation.com/logger => ../logger

require (
	github.com/go-sql-driver/mysql v1.5.0
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	substation.com/logger v0.0.0-00010101000000-000000000000
)
