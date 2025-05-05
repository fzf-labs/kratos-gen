package tpl

import _ "embed"

//go:embed service_method.tpl
var ServiceMethod string

//go:embed service_method_create.tpl
var ServiceMethodCreate string

//go:embed service_method_update.tpl
var ServiceMethodUpdate string

//go:embed service_method_update_status.tpl
var ServiceMethodUpdateStatus string

//go:embed service_method_delete.tpl
var ServiceMethodDelete string

//go:embed service_method_info.tpl
var ServiceMethodInfo string

//go:embed service_method_list.tpl
var ServiceMethodList string

//go:embed service_new.tpl
var ServiceNew string

//go:embed service_wire.tpl
var ServiceWire string
