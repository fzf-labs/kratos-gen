package biz

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
	"github.com/fzf-labs/kratos-gen/biz/tpl"
	"github.com/fzf-labs/kratos-gen/proto"
	"github.com/fzf-labs/kratos-gen/utils"
	"github.com/samber/lo"
)

type Biz struct {
	inPutPbPath   string // pb输入路径
	outPutBizPath string // biz输出路径
}

func NewBiz(inPutPbPath, outPutBizPath string) *Biz {
	return &Biz{
		inPutPbPath:   inPutPbPath,
		outPutBizPath: outPutBizPath,
	}
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

func (l *Biz) Run() {
	// 获取pb文件
	pbFiles, err := proto.GetPathPbFiles(l.inPutPbPath)
	if err != nil {
		log.Println("pb getPathPbFiles err ", err)
		return
	}
	l.Biz(pbFiles)
}

func (l *Biz) Biz(pbFiles []string) { //nolint:funlen
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
