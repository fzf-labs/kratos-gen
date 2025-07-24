package utils

import (
	"bytes"
	"html/template"
	"strings"
)

// toUpperCamelCase converts a string to UpperCamelCase
func toUpperCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[0:1]) + strings.ToLower(p[1:])
		}
	}
	return strings.Join(parts, "")
}

// toLowerCamelCase converts a string to lowerCamelCase
func toLowerCamelCase(s string) string {
	upperCamel := toUpperCamelCase(s)
	if len(upperCamel) > 0 {
		return strings.ToLower(upperCamel[0:1]) + upperCamel[1:]
	}
	return ""
}

// toSnakeCase converts a string to snake_case
func toSnakeCase(s string) string {
	// 先处理驼峰命名法，在大写字母前添加下划线
	var result string
	for i, c := range s {
		if i > 0 && c >= 'A' && c <= 'Z' {
			result += "_" + string(c)
		} else {
			result += string(c)
		}
	}
	return strings.ToLower(result)
}

// 首字母大写
func UCFirst(s string) string {
	if len(s) > 0 {
		return strings.ToUpper(s[0:1]) + s[1:]
	}
	return ""
}

// 首字母小写
func LCFirst(s string) string {
	if len(s) > 0 {
		return strings.ToLower(s[0:1]) + s[1:]
	}
	return ""
}

func PBToDB(s string) string {
	if len(s) > 0 {
		s = strings.ReplaceAll(s, "Id", "ID") // 将Id替换为ID
		s = strings.ToUpper(s[0:1]) + s[1:]   // 将首字母大写
		return s
	}
	return ""
}

// TemplateExecute 执行模板
func TemplateExecute(t string, data any) ([]byte, error) {
	buf := new(bytes.Buffer)
	funcMap := template.FuncMap{
		"ToCamel":      toUpperCamelCase,
		"ToLowerCamel": toLowerCamelCase,
		"ToSnake":      toSnakeCase,
		"UCFirst":      UCFirst,
		"LCFirst":      LCFirst,
		"HasPrefix":    strings.HasPrefix,
		"HasSuffix":    strings.HasSuffix,
		"Contains":     strings.Contains,
		"Replace":      strings.Replace,
		"Join":         strings.Join,
		"Split":        strings.Split,
		"ToLower":      strings.ToLower,
		"ToUpper":      strings.ToUpper,
		"TrimSpace":    strings.TrimSpace,
		"PBToDB":       PBToDB,
	}
	tmpl, err := template.New("").Funcs(funcMap).Parse(t)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
