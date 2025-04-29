package data

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	"gorm.io/gorm"

	"github.com/fzf-labs/kratos-gen/data/tpl"
	"github.com/fzf-labs/kratos-gen/utils"
)

type Data struct {
	db             string
	dsn            string // 数据库连接
	targetTables   string // 指定表
	outPutDataPath string // data输出路径
}

func NewData(db, dsn, targetTables, outPutDataPath string) *Data {
	return &Data{
		db:             db,
		dsn:            dsn,
		targetTables:   targetTables,
		outPutDataPath: outPutDataPath,
	}
}

func (d *Data) Run() {
	orm, err := utils.NewDB(d.db, d.dsn)
	if err != nil {
		log.Fatal(err)
	}
	dbName := orm.Migrator().CurrentDatabase()
	tables, err := orm.Migrator().GetTables()
	if err != nil {
		log.Printf("get tables err: %v", err)
	}
	if d.targetTables != "" {
		targetTables := strings.Split(d.targetTables, ",")
		for _, table := range targetTables {
			if !slices.Contains(tables, table) {
				log.Printf("table %s not found", table)
			}
		}
	}
	if _, err := os.Stat(d.outPutDataPath); os.IsNotExist(err) {
		log.Printf("Target directory: %s does not exsit\n", d.outPutDataPath)
		return
	}
	for _, table := range tables {
		tmp := map[string]string{
			"dbName":    dbName,
			"upperName": upperName(orm, table),
			"lowerName": lowerName(orm, table),
		}
		toData := filepath.Join(d.outPutDataPath, strings.ToLower(tmp["upperName"])+".go")
		if _, err := os.Stat(toData); !os.IsNotExist(err) {
			log.Printf("data new already exists: %s\n", toData)
		} else {
			b, err2 := utils.TemplateExecute(tpl.DataNew, tmp)
			if err2 != nil {
				return
			}
			if err3 := utils.Output(toData, b); err3 != nil {
				log.Fatal(err3)
			}
			log.Printf("data new generated successfully: %s\n", toData)
		}
	}
}

// upperName 大写
func upperName(db *gorm.DB, s string) string {
	return db.NamingStrategy.SchemaName(s)
}

// lowerName 小写
func lowerName(db *gorm.DB, s string) string {
	str := upperName(db, s)
	if str == "" {
		return str
	}
	words := []string{
		"API", "ASCII", "CPU", "CSS", "DNS",
		"EOF", "GUID", "HTML", "HTTP", "HTTPS",
		"ID", "IP", "JSON", "LHS", "QPS",
		"RAM", "RHS", "RPC", "SLA", "SMTP",
		"SSH", "TLS", "TTL", "UID", "UI",
		"UUID", "URI", "URL", "UTF8", "VM",
		"XML", "XSRF", "XSS",
	}
	// 如果第一个单词命中  则不处理
	for _, v := range words {
		if strings.HasPrefix(str, v) {
			return str
		}
	}
	rs := []rune(str)
	f := rs[0]
	if 'A' <= f && f <= 'Z' {
		str = string(unicode.ToLower(f)) + string(rs[1:])
	}
	return str
}
