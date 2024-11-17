package tpl

import _ "embed"

//go:embed biz_method.tpl
var BizMethod string

//go:embed biz_new.tpl
var BizNew string

//go:embed biz_wire.tpl
var BizWire string

//go:embed service_method.tpl
var ServiceMethod string

//go:embed service_new.tpl
var ServiceNew string

//go:embed service_wire.tpl
var ServiceWire string
