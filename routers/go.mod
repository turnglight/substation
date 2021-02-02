module substation.com/routers

go 1.15

replace substation.com/setting => ../setting
replace substation.com/api => ../api
replace substation.com/models => ../models
replace substation.com/logger => ../logger
replace substation.com/protocol => ../protocol

require (
	github.com/gin-gonic/gin v1.6.3
	substation.com/api v0.0.0-00010101000000-000000000000
	substation.com/setting v0.0.0-00010101000000-000000000000
)
