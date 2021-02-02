module substation.com/rest

go 1.15

replace substation.com/setting => ../setting
replace substation.com/routers => ../routers
replace substation.com/api => ../api
replace substation.com/models => ../models

require (
	github.com/gin-gonic/gin v1.6.3
	substation.com/models v0.0.0-00010101000000-000000000000
	substation.com/routers v0.0.0-00010101000000-000000000000
	substation.com/setting v0.0.0-00010101000000-000000000000
)
