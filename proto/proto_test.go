package proto

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestTranslate(t *testing.T) {
	protos := Translate("../example/admin_demo.proto")
	j, _ := json.Marshal(protos)
	fmt.Println(string(j))
}
