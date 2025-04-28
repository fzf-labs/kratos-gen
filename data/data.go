package data

import (
	"bytes"
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
	db                string
	dsn               string // 数据库连接
	targetTables      string // 指定表
	outPutDataPath    string // data输出路径
	outPutServicePath string // service输出路径
}

func NewData(db, dsn, targetTables, outPutDataPath, outPutServicePath string) *Data {
	return &Data{
		db:                db,
		dsn:               dsn,
		targetTables:      targetTables,
		outPutDataPath:    outPutDataPath,
		outPutServicePath: outPutServicePath,
	}
}

func (d *Data) Run() {
	orm, err := utils.NewDB(d.db, d.dsn)
	if err != nil {
		log.Fatal(err)
	}
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
	interfaces := make(map[string]struct{})
	for _, table := range tables {
		tmp := map[string]string{
			"upperName": upperName(orm, table),
			"lowerName": lowerName(orm, table),
		}
		interfaces[tmp["upperName"]+"Repo"] = struct{}{}
		toService := filepath.Join(d.outPutDataPath, strings.ToLower(tmp["upperName"])+".go")
		if _, err := os.Stat(toService); !os.IsNotExist(err) {
			log.Printf("data new already exists: %s\n", toService)
		} else {
			b, err2 := utils.TemplateExecute(tpl.DataNew, tmp)
			if err2 != nil {
				return
			}
			if err3 := utils.Output(toService, b); err3 != nil {
				log.Fatal(err3)
			}
			log.Printf("data new generated successfully: %s\n", toService)
		}
	}
	if len(interfaces) > 0 {
		toService := filepath.Join(d.outPutServicePath, "service.go")
		if _, err := os.Stat(toService); !os.IsNotExist(err) {
			// 语法树解析
			fileSet := token.NewFileSet()
			// 这里取绝对路径，方便打印出来的语法树可以转跳到编辑器
			path, _ := filepath.Abs(toService)
			f, err2 := decorator.ParseFile(fileSet, path, nil, parser.ParseComments)
			if err2 != nil {
				return
			}
			dst.Inspect(f, func(node dst.Node) bool {
				fn, ok := node.(*dst.GenDecl)
				if ok && fn != nil {
					if fn.Tok == token.TYPE {
						for _, spec := range fn.Specs {
							typeSpec := spec.(*dst.TypeSpec)
							delete(interfaces, typeSpec.Name.Name)
						}
					}
				}
				return true
			})
			if len(interfaces) > 0 {
				for k := range interfaces {
					f.Decls = append(f.Decls, &dst.GenDecl{
						Tok:    token.TYPE,
						Lparen: false,
						Specs: []dst.Spec{
							&dst.TypeSpec{
								Name: &dst.Ident{
									Name: k,
									Obj:  dst.NewObj(dst.Typ, k),
									Path: "",
									Decs: dst.IdentDecorations{
										NodeDecs: dst.NodeDecs{
											Before: dst.EmptyLine,
											After:  dst.EmptyLine,
										},
									},
								},
								TypeParams: nil,
								Assign:     false,
								Type: &dst.InterfaceType{
									Methods: &dst.FieldList{
										Opening: true,
										Closing: true,
										Decs: dst.FieldListDecorations{
											NodeDecs: dst.NodeDecs{},
											Opening: dst.Decorations{
												"\n",
											},
										},
									},
									Incomplete: false,
									Decs:       dst.InterfaceTypeDecorations{},
								},
								Decs: dst.TypeSpecDecorations{
									NodeDecs: dst.NodeDecs{
										Before: dst.EmptyLine,
										After:  dst.EmptyLine,
									},
								},
							},
						},
						Rparen: false,
						Decs: dst.GenDeclDecorations{
							NodeDecs: dst.NodeDecs{
								Before: dst.EmptyLine,
								After:  dst.EmptyLine,
							},
						},
					})
				}
				// 创建一个缓冲区来存储格式化后的代码
				buf := &bytes.Buffer{}
				err = decorator.Fprint(buf, f)
				if err != nil {
					return
				}
				err = utils.Output(toService, buf.Bytes())
				if err != nil {
					log.Printf("pb.Output err %v", err)
					return
				}
			}
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
