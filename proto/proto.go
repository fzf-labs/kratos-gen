package proto

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/emicklei/proto"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const empty = "google.protobuf.Empty"

const (
	unaryType          uint8 = 1
	twoWayStreamsType  uint8 = 2
	requestStreamsType uint8 = 3
	returnsStreamsType uint8 = 4
)

type Proto struct {
	GoPackage   string
	UpperName   string
	LowerName   string
	FirstChar   string
	GoogleEmpty bool
	UseIO       bool
	UseContext  bool
	Methods     []*Method
}

type Method struct {
	Name    string
	Request string
	Reply   string
	Type    uint8
	Comment string
}

func Translate(path string) []*Proto {
	res := make([]*Proto, 0)
	pkg := ""
	reader, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer reader.Close()
	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		log.Println(err)
		return nil
	}
	proto.Walk(definition,
		proto.WithOption(func(o *proto.Option) {
			if o.Name == "go_package" {
				pkg = strings.Split(o.Constant.Source, ";")[0]
			}
		}),
		proto.WithService(func(s *proto.Service) {
			cs := &Proto{
				GoPackage: pkg,
				UpperName: GetUpperName(s.Name),
				LowerName: GetLowerName(s.Name),
				FirstChar: strings.ToLower(s.Name[0:1]),
			}
			for _, e := range s.Elements {
				r, ok := e.(*proto.RPC)
				if !ok {
					continue
				}
				if r.Comment == nil {
					log.Fatalf("rpc %s comment is nil", r.Name)
					return
				}
				method := &Method{
					Name:    GetUpperName(r.Name),
					Request: GetParametersName(r.RequestType),
					Reply:   GetParametersName(r.ReturnsType),
					Type:    GetMethodType(r.StreamsRequest, r.StreamsReturns),
					Comment: strings.TrimSpace(r.Comment.Message()),
				}
				if (method.Type == unaryType && (method.Request == empty || method.Reply == empty)) ||
					(method.Type == returnsStreamsType && method.Request == empty) {
					cs.GoogleEmpty = true
				}
				if method.Type == twoWayStreamsType || method.Type == requestStreamsType {
					cs.UseIO = true
				}
				if method.Type == unaryType {
					cs.UseContext = true
				}
				cs.Methods = append(cs.Methods, method)
			}
			res = append(res, cs)
		}),
	)
	return res
}

// GetMethodType 获取方法类型
func GetMethodType(streamsRequest, streamsReturns bool) uint8 {
	if streamsRequest && streamsReturns {
		return twoWayStreamsType
	}
	if streamsRequest {
		return requestStreamsType
	}
	if streamsReturns {
		return returnsStreamsType
	}
	return unaryType
}

// GetParametersName 获取参数名称
func GetParametersName(name string) string {
	return strings.ReplaceAll(name, ".", "_")
}

// GetUpperName 获取大驼峰名称
func GetUpperName(name string) string {
	s := strings.Split(name, ".")[0]
	s = strings.ReplaceAll(s, "_", " ")
	s = cases.Title(language.Und, cases.NoLower).String(s)
	return strings.ReplaceAll(s, " ", "")
}

// GetLowerName 获取小驼峰名称
func GetLowerName(name string) string {
	s := strings.Split(name, ".")[0]
	s = strings.ReplaceAll(s, "_", " ")
	s = LCFirst(cases.Title(language.Und, cases.NoLower).String(s))
	return strings.ReplaceAll(s, " ", "")
}

// LCFirst 首字母小写
func LCFirst(s string) string {
	if s == "" {
		return s
	}
	rs := []rune(s)
	f := rs[0]
	if 'A' <= f && f <= 'Z' {
		return string(unicode.ToLower(f)) + string(rs[1:])
	}
	return s
}

// GetPathPbFiles 获取目录下指定后缀文件
func GetPathPbFiles(p string) ([]string, error) {
	var files []string
	err := filepath.Walk(p, func(path string, _ os.FileInfo, _ error) error {
		if strings.HasSuffix(path, ".proto") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return files, nil
}
