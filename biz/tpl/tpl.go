package tpl

import _ "embed"

//go:embed biz_method.tpl
var BizMethod string

//go:embed biz_new.tpl
var BizNew string

//go:embed biz_wire.tpl
var BizWire string
