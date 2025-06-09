package proto

import (
	"encoding/json"
	"os"
	"testing"
)

func TestTranslate(t *testing.T) {
	protos := Translate("../example/api/v1/admin_demo.proto")
	marshal, _ := json.Marshal(protos)
	// 写入到文件中
	os.WriteFile("../example/api/v1/admin_demo.json", marshal, 0644)
}
