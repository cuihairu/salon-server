package utils

import (
	"encoding/json"
	"testing"
)

type Images []string

func TestJsonField_MarshalJSON(t *testing.T) {
	images := Images{"https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7f947a enthusiasm.png", "https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7f947a enthusiasm.png", "b", "c"}
	imgs := NewJsonField[Images](&images)
	jsonData, err := imgs.MarshalJSON()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("images marshal: %v", string(jsonData))
	images2 := Images{}
	imgs2 := NewJsonField[Images](&images2)
	jsonData2, err := imgs2.MarshalJSON()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("images2 marshal: %v", string(jsonData2))
	json.Unmarshal(jsonData, &images2)
	t.Logf("images unmarshal: %v", images2)
}
