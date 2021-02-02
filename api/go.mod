module substation.com/api

go 1.15

replace substation.com/models => ../models
replace substation.com/setting => ../setting

require (
	github.com/gin-gonic/gin v1.6.3
	substation.com/models v0.0.0-00010101000000-000000000000
)
