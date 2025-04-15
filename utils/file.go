package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/tools/imports"
)

const (
	GoExt = ".go"
)

// ReadDirFilesWithSuffix 遍历指定目录下的指定后缀文件
func ReadDirFilesWithSuffix(root, suffix string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, suffix) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// ReadFileToString 读取文件到string
func ReadFileToString(dir string) (string, error) {
	file, err := os.OpenFile(dir, os.O_RDWR, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	filesize := fileInfo.Size()
	buffer := make([]byte, filesize)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}

// WriteContentCover 数据写入，不存在则创建
func WriteContentCover(filePath, content string) error {
	fileDir := filepath.Dir(filePath)
	if err := os.MkdirAll(fileDir, 0775); err != nil {
		return err
	}
	dstFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0665)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = dstFile.WriteString(content)
	if err != nil {
		return err
	}
	return err
}

// Output 输出文件
func Output(fileName string, content []byte) error {
	result, err := imports.Process(fileName, content, nil)
	if err != nil {
		lines := strings.Split(string(content), "\n")
		errLine, _ := strconv.Atoi(strings.Split(err.Error(), ":")[1])
		startLine, endLine := errLine-5, errLine+5
		fmt.Println("Format fail:", errLine, err)
		if startLine < 0 {
			startLine = 0
		}
		if endLine > len(lines)-1 {
			endLine = len(lines) - 1
		}
		for i := startLine; i <= endLine; i++ {
			fmt.Println(i, lines[i])
		}
		return fmt.Errorf("cannot format file: %w", err)
	}
	fileDir := filepath.Dir(fileName)
	err = os.MkdirAll(fileDir, 0755)
	if err != nil {
		return err
	}
	dstFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = dstFile.Write(result)
	if err != nil {
		return err
	}
	return err
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
	s = LcFirst(cases.Title(language.Und, cases.NoLower).String(s))
	return strings.ReplaceAll(s, " ", "")
}

// LcFirst 首字母小写
func LcFirst(s string) string {
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

// GetDirUpperName 获取目录大写名
func GetDirUpperName(dir string) string {
	str := ""
	splitList := strings.Split(filepath.ToSlash(dir), "/")
	for _, v := range splitList {
		str += GetUpperName(v)
	}
	return str
}

// GetDirSnakeName 获取目录小写名
func GetDirSnakeName(dir string) string {
	list := make([]string, 0)
	splitList := strings.Split(filepath.ToSlash(dir), "/")
	for _, v := range splitList {
		list = append(list, strings.ToLower(GetUpperName(v)))
	}
	return strings.Join(list, "_")
}
