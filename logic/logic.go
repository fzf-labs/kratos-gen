package logic

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/fzf-labs/kratos-gen/logic/proto"
	"github.com/fzf-labs/kratos-gen/logic/tpl"
	"github.com/fzf-labs/kratos-gen/utils"
	"github.com/samber/lo"
)

type Logic struct {
	inPutPbPath       string // pb输入路径
	outPutBizPath     string // biz输出路径
	outPutServicePath string // service输出路径
}

func NewLogic(inPutPbPath, outPutBizPath, outPutServicePath string) *Logic {
	return &Logic{inPutPbPath: inPutPbPath, outPutBizPath: outPutBizPath, outPutServicePath: outPutServicePath}
}

func (l *Logic) Run() {
	if l.inPutPbPath == "" {
		log.Println("inPutPbPath is empty")
		return
	}
	// 获取pb文件
	pbFiles, err := proto.GetPathPbFiles(l.inPutPbPath)
	if err != nil {
		log.Println("pb getPathPbFiles err ", err)
		return
	}
	l.Biz(pbFiles)
	l.Service(pbFiles)
}

type BizNewTplParam struct {
	GoPackage   string
	UpperName   string
	LowerName   string
	FirstChar   string
	GoogleEmpty bool
	UseIO       bool
	UseContext  bool
}

type BizMethodTplParam struct {
	GoPackage   string
	UpperName   string
	LowerName   string
	FirstChar   string
	GoogleEmpty bool
	UseIO       bool
	UseContext  bool
	Name        string
	Request     string
	Reply       string
	Type        uint8
	Comment     string
}

type BizWireTplParam struct {
	Wires []string
}

const (
	GoExt = ".go"
)

func (l *Logic) Biz(pbFiles []string) { //nolint:funlen
	wires := make([]string, 0)
	var wg sync.WaitGroup
	wg.Add(len(pbFiles))
	for _, file := range pbFiles {
		go func(file string) {
			defer wg.Done()
			dir := strings.TrimPrefix(filepath.Dir(file), "api/")
			dirSnakeName := getDirSnakeName(dir)
			dirUpperName := getDirUpperName(dir)
			pbInfos := proto.Translate(file)
			for _, s := range pbInfos {
				if len(s.Methods) == 0 {
					continue
				}
				bizName := dirUpperName + s.UpperName
				upperBizName := proto.GetUpperName(bizName)
				lowerBizName := proto.GetUpperName(bizName)
				firstChar := strings.ToLower(bizName[0:1])
				wires = append(wires, fmt.Sprintf("New%sUseCase", bizName))
				toService := filepath.Join(l.outPutBizPath, strings.ToLower(dirSnakeName+"_"+s.UpperName)+GoExt)
				if _, err2 := os.Stat(toService); !os.IsNotExist(err2) {
					log.Printf("biz new already exists: %s\n", toService)
				} else {
					b, err3 := utils.TemplateExecute(tpl.BizNew, &BizNewTplParam{
						GoPackage:   s.GoPackage,
						UpperName:   upperBizName,
						LowerName:   lowerBizName,
						FirstChar:   firstChar,
						GoogleEmpty: s.GoogleEmpty,
						UseIO:       s.UseIO,
						UseContext:  s.UseContext,
					})
					if err3 != nil {
						return
					}
					if err3 = utils.Output(toService, b); err3 != nil {
						log.Fatal(err3)
						return
					}
					log.Printf("biz new generated successfully: %s\n", toService)
				}
				for _, method := range s.Methods {
					toMethod := filepath.Join(l.outPutBizPath, strings.ToLower(dirSnakeName+"_"+s.UpperName+"_"+method.Name)+GoExt)
					if _, err3 := os.Stat(toMethod); !os.IsNotExist(err3) {
						log.Printf("biz method already exists: %s\n", toMethod)
						continue
					}
					b, err3 := utils.TemplateExecute(tpl.BizMethod, BizMethodTplParam{
						GoPackage:   s.GoPackage,
						UpperName:   upperBizName,
						LowerName:   lowerBizName,
						FirstChar:   firstChar,
						GoogleEmpty: s.GoogleEmpty,
						UseIO:       s.UseIO,
						UseContext:  s.UseContext,
						Name:        method.Name,
						Request:     method.Request,
						Reply:       method.Reply,
						Type:        method.Type,
						Comment:     method.Comment,
					})
					if err3 != nil {
						log.Fatal(err3)
						return
					}
					if err3 = utils.Output(toMethod, b); err3 != nil {
						log.Fatal(err3)
						return
					}
					log.Printf("biz method generated successfully: %s\n", toMethod)
				}
			}
		}(file)
	}
	wg.Wait()
	toBizWire := filepath.Join(l.outPutBizPath, "biz.go")
	if _, err := os.Stat(toBizWire); !os.IsNotExist(err) {
		// 语法树解析
		fileSet := token.NewFileSet()
		// 这里取绝对路径，方便打印出来的语法树可以转跳到编辑器
		path, _ := filepath.Abs(toBizWire)
		f, err := decorator.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return
		}
		dst.Inspect(f, func(node dst.Node) bool {
			fn, ok := node.(*dst.CallExpr)
			if ok && fn != nil {
				for _, arg := range fn.Args {
					if ident, ok2 := arg.(*dst.Ident); ok2 {
						wires = append(wires, ident.Name)
					}
				}
				wires = lo.Union(wires)
				sort.Strings(wires)
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
			return true
		})
		// 创建一个缓冲区来存储格式化后的代码
		buf := &bytes.Buffer{}
		err = decorator.Fprint(buf, f)
		if err != nil {
			return
		}
		err = utils.Output(toBizWire, buf.Bytes())
		if err != nil {
			log.Printf("pb.Output err4 %v", err)
			return
		}
	} else {
		execute, err := utils.TemplateExecute(tpl.BizWire, &BizWireTplParam{
			Wires: wires,
		})
		if err != nil {
			log.Fatal(err)
			return
		}
		if err = utils.Output(toBizWire, execute); err != nil {
			log.Fatal(err)
			return
		}
	}
	log.Printf("biz wire generated successfully: %s\n", toBizWire)
}

type ServiceNewTplParam struct {
	GoPackage        string
	UpperName        string
	LowerName        string
	UpperServiceName string
	FirstChar        string
	GoogleEmpty      bool
	UseIO            bool
	UseContext       bool
}

type ServiceMethodTplParam struct {
	GoPackage        string
	UpperName        string
	LowerName        string
	UpperServiceName string
	FirstChar        string
	GoogleEmpty      bool
	UseIO            bool
	UseContext       bool
	Name             string
	Request          string
	Reply            string
	Type             uint8
	Comment          string
}

type ServiceWireTplParam struct {
	Wires []string
}

func (l *Logic) Service(pbFiles []string) { //nolint:funlen,gocyclo
	wires := make([]string, 0)
	var wg sync.WaitGroup
	wg.Add(len(pbFiles))
	for _, file := range pbFiles {
		go func(file string) {
			defer wg.Done()
			dir := strings.TrimPrefix(filepath.Dir(file), "api/")
			dirSnakeName := getDirSnakeName(dir)
			dirUpperName := getDirUpperName(dir)
			pbInfos := proto.Translate(file)
			for _, s := range pbInfos {
				if len(s.Methods) == 0 {
					continue
				}
				serviceName := dirUpperName + s.UpperName
				upperServiceName := proto.GetUpperName(serviceName)
				lowerServiceName := proto.GetLowerName(serviceName)
				firstChar := strings.ToLower(serviceName[0:1])
				wires = append(wires, fmt.Sprintf("New%sService", upperServiceName))
				toService := filepath.Join(l.outPutServicePath, strings.ToLower(dirSnakeName+"_"+s.UpperName)+GoExt)
				if _, err3 := os.Stat(toService); os.IsNotExist(err3) {
					// 创建
					service, err := utils.TemplateExecute(tpl.ServiceNew, &ServiceNewTplParam{
						GoPackage:        s.GoPackage,
						UpperName:        upperServiceName,
						LowerName:        lowerServiceName,
						UpperServiceName: s.UpperName,
						FirstChar:        firstChar,
						GoogleEmpty:      s.GoogleEmpty,
						UseIO:            s.UseIO,
						UseContext:       s.UseContext,
					})
					if err != nil {
						log.Printf("pb.TemplateExecute err %v", err)
						return
					}
					for _, v := range s.Methods {
						method, err2 := utils.TemplateExecute(tpl.ServiceMethod, &ServiceMethodTplParam{
							GoPackage:        s.GoPackage,
							UpperName:        upperServiceName,
							LowerName:        lowerServiceName,
							UpperServiceName: s.UpperName,
							FirstChar:        firstChar,
							GoogleEmpty:      s.GoogleEmpty,
							UseIO:            s.UseIO,
							UseContext:       s.UseContext,
							Name:             v.Name,
							Request:          v.Request,
							Reply:            v.Reply,
							Type:             v.Type,
							Comment:          v.Comment,
						})
						if err2 != nil {
							return
						}
						// 换行
						service = append(service, []byte("\n")...)
						service = append(service, method...)
					}
					err = utils.Output(toService, service)
					if err != nil {
						log.Printf("pb.Output err %v", err)
						return
					}
					log.Printf("service generated successfully: %s\n", toService)
				} else {
					// 更新
					if len(s.Methods) == 0 {
						return
					}
					methodToProto := make(map[string]*proto.Proto)
					for _, v := range s.Methods {
						methodToProto[v.Name] = s
					}
					fileSet := token.NewFileSet()
					// 这里取绝对路径，方便打印出来的语法树可以转跳到编辑器
					path, _ := filepath.Abs(toService)
					f, err := decorator.ParseFile(fileSet, path, nil, parser.ParseComments)
					if err != nil {
						return
					}
					dst.Inspect(f, func(node dst.Node) bool {
						fn, ok := node.(*dst.FuncDecl)
						if ok {
							delete(methodToProto, fn.Name.Name)
						}
						return true
					})
					if len(methodToProto) == 0 {
						return
					}
					for _, v := range s.Methods {
						if _, ok := methodToProto[v.Name]; !ok {
							continue
						}
						method, err2 := utils.TemplateExecute(tpl.ServiceNew+tpl.ServiceMethod, &ServiceMethodTplParam{
							GoPackage:        s.GoPackage,
							UpperName:        upperServiceName,
							LowerName:        lowerServiceName,
							UpperServiceName: s.UpperName,
							FirstChar:        firstChar,
							GoogleEmpty:      s.GoogleEmpty,
							UseIO:            s.UseIO,
							UseContext:       s.UseContext,
							Name:             v.Name,
							Request:          v.Request,
							Reply:            v.Reply,
							Type:             v.Type,
							Comment:          v.Comment,
						})
						if err2 != nil {
							return
						}
						decls, err2 := getAstFuncDecl(method)
						if err2 != nil {
							log.Printf("getAstFuncDecl err2 %v", err2)
							return
						}
						for _, funcDecl := range decls {
							f.Decls = append(f.Decls, funcDecl)
						}
					}
					// 创建一个缓冲区来存储格式化后的代码
					buf := &bytes.Buffer{}
					err = decorator.Fprint(buf, f)
					if err != nil {
						return
					}
					err = utils.Output(toService, buf.Bytes())
					if err != nil {
						log.Printf("pb.Output err4 %v", err)
						return
					}
					log.Printf("service updated successfully: %s\n", toService)
				}
			}
		}(file)
	}
	wg.Wait()
	toServiceWire := filepath.Join(l.outPutServicePath, "service.go")
	if _, err := os.Stat(toServiceWire); !os.IsNotExist(err) {
		// 语法树解析
		fileSet := token.NewFileSet()
		// 这里取绝对路径，方便打印出来的语法树可以转跳到编辑器
		path, _ := filepath.Abs(toServiceWire)
		f, err := decorator.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return
		}
		dst.Inspect(f, func(node dst.Node) bool {
			fn, ok := node.(*dst.CallExpr)
			if ok && fn != nil {
				// 判断是否已经存在
				if fn.Fun != nil && fn.Fun.(*dst.SelectorExpr).Sel != nil && fn.Fun.(*dst.SelectorExpr).Sel.Name == "NewSet" {
					for _, arg := range fn.Args {
						if ident, ok2 := arg.(*dst.Ident); ok2 {
							wires = append(wires, ident.Name)
						}
					}
					wires = lo.Union(wires)
					sort.Strings(wires)
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
			return true
		})
		// 创建一个缓冲区来存储格式化后的代码
		buf := &bytes.Buffer{}
		err = decorator.Fprint(buf, f)
		if err != nil {
			return
		}
		err = utils.Output(toServiceWire, buf.Bytes())
		if err != nil {
			log.Printf("pb.Output err4 %v", err)
			return
		}
	} else {
		execute, err := utils.TemplateExecute(tpl.ServiceWire, &ServiceWireTplParam{
			Wires: lo.Union(wires),
		})
		if err != nil {
			log.Printf("pb.TemplateExecute err %v", err)
			return
		}
		err = utils.Output(toServiceWire, execute)
		if err != nil {
			log.Printf("pb.Output err %v", err)
			return
		}
	}
	log.Printf("service wire generated successfully: %s\n", toServiceWire)
}

// getAstFuncDecl 获取函数
func getAstFuncDecl(code []byte) ([]*dst.FuncDecl, error) {
	var methods []*dst.FuncDecl
	f, err := decorator.Parse(code)
	if err != nil {
		return nil, err
	}
	dst.Inspect(f, func(node dst.Node) bool {
		fn, ok := node.(*dst.FuncDecl)
		if ok && !strings.Contains(fn.Name.String(), "New") {
			methods = append(methods, fn)
		}
		return true
	})
	return methods, nil
}

// getDirUpperName 获取目录大写名
func getDirUpperName(dir string) string {
	str := ""
	splitList := strings.Split(filepath.ToSlash(dir), "/")
	for _, v := range splitList {
		str += proto.GetUpperName(v)
	}
	return str
}

func getDirSnakeName(dir string) string {
	list := make([]string, 0)
	splitList := strings.Split(filepath.ToSlash(dir), "/")
	for _, v := range splitList {
		list = append(list, strings.ToLower(proto.GetUpperName(v)))
	}
	return strings.Join(list, "_")
}
