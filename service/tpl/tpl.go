package tpl

import _ "embed"

//go:embed service_method.tpl
var ServiceMethod string

//go:embed service_new.tpl
var ServiceNew string

//go:embed service_wire.tpl
var ServiceWire string
