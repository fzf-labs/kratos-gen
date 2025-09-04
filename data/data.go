package data

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/fzf-labs/kratos-gen/data/tpl"
	"github.com/fzf-labs/kratos-gen/utils"
	"gorm.io/gorm"
)

type Data struct {
	db             string
	dsn            string // 数据库连接
	targetTables   string // 指定表
	outPutDataPath string // data输出路径
	partitionTable bool   // 是否处理分区表
}

func NewData(db, dsn, targetTables, outPutDataPath string, partitionTable bool) *Data {
	return &Data{
		db:             db,
		dsn:            dsn,
		targetTables:   targetTables,
		outPutDataPath: outPutDataPath,
		partitionTable: partitionTable,
	}
}

func (d *Data) Run() {
	orm, err := utils.NewDB(d.db, d.dsn)
	if err != nil {
		log.Printf("NewDB err: %v", err)
		return
	}
	dbName := orm.Migrator().CurrentDatabase()
	tables, err := orm.Migrator().GetTables()
	if err != nil {
		log.Printf("get tables err: %v", err)
		return
	}
	if d.targetTables != "" {
		targetTables := strings.Split(d.targetTables, ",")
		for _, table := range targetTables {
			if !slices.Contains(tables, table) {
				log.Printf("table %s not found", table)
				return
			}
		}
		tables = targetTables
	}
	if d.partitionTable {
		// 查询分区表父级到子表的映射
		partitionTableToChildTables, err := utils.GetPartitionTableToChildTables(orm)
		if err != nil {
			return
		}
		partitionChildTables := make([]string, 0)
		for _, v := range partitionTableToChildTables {
			partitionChildTables = append(partitionChildTables, v...)
		}
		// 去掉tables中的分区子表
		tables = utils.SliRemove(tables, partitionChildTables)
	}
	if _, err := os.Stat(d.outPutDataPath); os.IsNotExist(err) {
		log.Printf("Target directory: %s does not exsit\n", d.outPutDataPath)
		return
	}
	dataWires := make([]string, 0)
	for _, table := range tables {
		tmp := map[string]string{
			"dbName":    dbName,
			"upperName": upperName(orm, table),
			"lowerName": lowerName(orm, table),
		}
		dataWires = append(dataWires, fmt.Sprintf("New%sRepo", tmp["upperName"]), fmt.Sprintf("%s_repo.New%sRepo", tmp["dbName"], tmp["upperName"]))
		toData := filepath.Join(d.outPutDataPath, strings.ToLower(tmp["upperName"])+".go")
		if _, err := os.Stat(toData); !os.IsNotExist(err) {
			log.Printf("data new already exists: %s\n", toData)
		} else {
			b, err2 := utils.TemplateExecute(tpl.DataRepoNew, tmp)
			if err2 != nil {
				return
			}
			if err3 := utils.Output(toData, b); err3 != nil {
				log.Fatal(err3)
			}
			log.Printf("data new generated successfully: %s\n", toData)
		}
	}
	toDataWire := filepath.Join(d.outPutDataPath, "data.go")
	if _, err := os.Stat(toDataWire); os.IsNotExist(err) {
		execute, err := utils.TemplateExecute(tpl.DataNew, map[string]string{})
		if err != nil {
			log.Printf("pb.TemplateExecute err %v", err)
			return
		}
		err = utils.Output(toDataWire, execute)
		if err != nil {
			log.Printf("pb.Output err %v", err)
			return
		}
	}
	if len(dataWires) > 0 {
		// 语法树解析
		fileSet := token.NewFileSet()
		// 这里取绝对路径，方便打印出来的语法树可以转跳到编辑器
		path, _ := filepath.Abs(toDataWire)
		f, err := decorator.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return
		}
		dst.Inspect(f, func(node dst.Node) bool {
			fn, ok := node.(*dst.CallExpr)
			if ok && fn != nil {
				// 判断Fun 是 *dst.SelectorExpr
				selector, ok := fn.Fun.(*dst.SelectorExpr)
				if ok && selector != nil {
					// 判断是否已经存在
					if selector.Sel != nil && selector.Sel.Name == "NewSet" {
						wires := make([]string, 0)
						for _, arg := range fn.Args {
							// 判断arg 是否是 *dst.Ident
							if ident, ok2 := arg.(*dst.Ident); ok2 {
								if ident.Name != "" {
									wires = append(wires, ident.Name)
								}
							}
							// 判断arg 是否是 *dst.SelectorExpr
							if selector, ok2 := arg.(*dst.SelectorExpr); ok2 {
								if selector.X.(*dst.Ident).Name != "" {
									wires = append(wires, selector.X.(*dst.Ident).Name+"."+selector.Sel.Name)
								}
							}
						}
						for _, wire := range dataWires {
							if !slices.Contains(wires, wire) {
								wires = append(wires, wire)
							}
						}
						args := make([]dst.Expr, 0)
						for _, wire := range wires {
							args = append(args, &dst.Ident{
								Name: wire,
								Obj:  nil,
								Path: "",
								Decs: dst.IdentDecorations{
									NodeDecs: dst.NodeDecs{
										Before: 1,
										After:  1,
									},
								},
							})
						}
						fn.Args = args
						return true
					}
				}
			}
			return true
		})
		// 创建一个缓冲区来存储格式化后的代码
		buf := &bytes.Buffer{}
		err = decorator.Fprint(buf, f)
		if err != nil {
			return
		}
		err = utils.Output(toDataWire, buf.Bytes())
		if err != nil {
			log.Printf("pb.Output err4 %v", err)
			return
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
