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

// Field 表示消息中的字段
type Field struct {
	Name          string // 字段名称
	Type          string // 字段类型
	Comment       string // 字段注释 (来自字段上方的注释)
	InlineComment string // 字段行内注释 (来自字段后面的注释)
	Repeated      bool   // 是否是重复字段
}

// Message 表示 proto 文件中的消息定义
type Message struct {
	Name   string   // 消息名称
	Fields []*Field // 消息字段
}

type Proto struct {
	GoPackage   string
	UpperName   string
	LowerName   string
	FirstChar   string
	GoogleEmpty bool
	UseIO       bool
	UseContext  bool
	Methods     []*Method
	Messages    []*Message // 添加消息定义列表
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

	// 先解析所有消息定义
	messages := make(map[string]*Message)
	proto.Walk(definition,
		proto.WithMessage(func(m *proto.Message) {
			message := &Message{
				Name:   m.Name,
				Fields: make([]*Field, 0),
			}

			for _, e := range m.Elements {
				if f, ok := e.(*proto.NormalField); ok {
					field := &Field{
						Name:     f.Name,
						Type:     f.Type,
						Repeated: f.Repeated,
					}
					// 获取字段上方的注释
					if f.Comment != nil {
						field.Comment = strings.TrimSpace(f.Comment.Message())
					}
					// 获取字段后面的行内注释
					if f.InlineComment != nil {
						field.InlineComment = strings.TrimSpace(f.InlineComment.Message())
					}
					message.Fields = append(message.Fields, field)
				}
			}

			messages[m.Name] = message
		}),
	)

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
				Messages:  make([]*Message, 0),
			}

			// 添加消息定义到 Proto 结构体
			for _, msg := range messages {
				cs.Messages = append(cs.Messages, msg)
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

				requestType := GetParametersName(r.RequestType)
				replyType := GetParametersName(r.ReturnsType)
				method := &Method{
					Name:    GetUpperName(r.Name),
					Request: requestType,
					Reply:   replyType,
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
