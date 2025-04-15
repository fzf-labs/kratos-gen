package service

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
	"github.com/fzf-labs/kratos-gen/proto"
	"github.com/fzf-labs/kratos-gen/service/tpl"
	"github.com/fzf-labs/kratos-gen/utils"
	"github.com/samber/lo"
)

type Service struct {
	inPutPbPath       string // pb输入路径
	outPutServicePath string // service输出路径
}

func NewService(inPutPbPath, outPutServicePath string) *Service {
	return &Service{
		inPutPbPath:       inPutPbPath,
		outPutServicePath: outPutServicePath,
	}
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

func (l *Service) Run() {
	pbFiles, err := proto.GetPathPbFiles(l.inPutPbPath)
	if err != nil {
		log.Println("pb getPathPbFiles err ", err)
		return
	}
	l.Service(pbFiles)
}

func (l *Service) Service(pbFiles []string) { //nolint:funlen,gocyclo
	wires := make([]string, 0)
	var wg sync.WaitGroup
	wg.Add(len(pbFiles))
	for _, file := range pbFiles {
		go func(file string) {
			defer wg.Done()
			dir := strings.TrimPrefix(filepath.Dir(file), "api/")
			dirSnakeName := utils.GetDirSnakeName(dir)
			dirUpperName := utils.GetDirUpperName(dir)
			pbInfos := proto.Translate(file)
			for _, s := range pbInfos {
				if len(s.Methods) == 0 {
					continue
				}
				serviceName := dirUpperName + s.UpperName
				upperServiceName := utils.GetUpperName(serviceName)
				lowerServiceName := utils.GetLowerName(serviceName)
				firstChar := strings.ToLower(serviceName[0:1])
				wires = append(wires, fmt.Sprintf("New%sService", upperServiceName))
				toService := filepath.Join(l.outPutServicePath, strings.ToLower(dirSnakeName+"_"+s.UpperName)+utils.GoExt)
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
					methodToProto := make(map[string]*proto.Method)
					for _, v := range s.Methods {
						methodToProto[v.Name] = v
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
