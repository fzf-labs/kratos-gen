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
	GoPackage        string // 包名
	UpperName        string // 服务名
	LowerName        string // 服务名小写
	UpperServiceName string // 服务名大写
	LowerServiceName string // 服务名小写
	FirstChar        string // 服务名首字母
	GoogleEmpty      bool   // 是否使用google.protobuf.Empty
	UseIO            bool   // 是否使用io.Reader
	UseContext       bool   // 是否使用context.Context
}

type ServiceMethodTplParam struct {
	GoPackage        string           // 包名
	UpperName        string           // 服务名
	LowerName        string           // 服务名小写
	UpperServiceName string           // 服务名大写
	LowerServiceName string           // 服务名小写
	FirstChar        string           // 服务名首字母
	GoogleEmpty      bool             // 是否使用google.protobuf.Empty
	UseIO            bool             // 是否使用io.Reader
	UseContext       bool             // 是否使用context.Context
	Name             string           // 方法名
	Request          string           // 请求参数
	Reply            string           // 响应参数
	Type             uint8            // 方法类型
	Comment          string           // 方法注释
	Messages         []*proto.Message // proto文件中的消息定义
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

				if _, err3 := os.Stat(toService); !os.IsNotExist(err3) {
					log.Printf("service already exists: %s\n", toService)
				} else {
					// 创建
					service, err := utils.TemplateExecute(tpl.ServiceNew, &ServiceNewTplParam{
						GoPackage:        s.GoPackage,
						UpperName:        s.UpperName,
						LowerName:        s.LowerName,
						UpperServiceName: upperServiceName,
						LowerServiceName: lowerServiceName,
						FirstChar:        firstChar,
						GoogleEmpty:      s.GoogleEmpty,
						UseIO:            s.UseIO,
						UseContext:       s.UseContext,
					})
					if err != nil {
						log.Printf("pb.TemplateExecute err %v", err)
						return
					}
					err = utils.Output(toService, service)
					if err != nil {
						log.Printf("pb.Output err %v", err)
						return
					}
					log.Printf("service generated successfully: %s\n", toService)
				}
				for _, v := range s.Methods {
					toMethod := filepath.Join(l.outPutServicePath, strings.ToLower(dirSnakeName+"_"+s.UpperName+"_"+v.Name)+utils.GoExt)
					if _, err3 := os.Stat(toMethod); !os.IsNotExist(err3) {
						log.Printf("service method already exists: %s\n", toMethod)
						continue
					}
					serviceMethodTplParam := &ServiceMethodTplParam{
						GoPackage:        s.GoPackage,
						UpperName:        s.UpperName,
						LowerName:        s.LowerName,
						UpperServiceName: upperServiceName,
						LowerServiceName: lowerServiceName,
						FirstChar:        firstChar,
						GoogleEmpty:      s.GoogleEmpty,
						UseIO:            s.UseIO,
						UseContext:       s.UseContext,
						Name:             v.Name,
						Request:          v.Request,
						Reply:            v.Reply,
						Type:             v.Type,
						Comment:          v.Comment,
						Messages:         s.Messages,
					}
					var serviceMethod string
					switch v.Name {
					case "Create" + s.UpperName:
						serviceMethod = tpl.ServiceMethodCreate
					case "Update" + s.UpperName:
						serviceMethod = tpl.ServiceMethodUpdate
					case "Update" + s.UpperName + "Status":
						serviceMethod = tpl.ServiceMethodUpdateStatus
					case "Delete" + s.UpperName:
						serviceMethod = tpl.ServiceMethodDelete
					case "Get" + s.UpperName + "Info":
						serviceMethod = tpl.ServiceMethodInfo
					case "Get" + s.UpperName + "List":
						serviceMethod = tpl.ServiceMethodList
					default:
						serviceMethod = tpl.ServiceMethod
					}
					method, err := utils.TemplateExecute(serviceMethod, serviceMethodTplParam)
					if err != nil {
						log.Fatal(err)
						return
					}
					if err := utils.Output(toMethod, method); err != nil {
						log.Fatal(err)
						return
					}
					log.Printf("service method generated successfully: %s\n", toMethod)
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
