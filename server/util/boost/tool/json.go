package tool

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
	jsondecoder "github.com/json-iterator/go/extra"
)

var Json *jsonTool

func init() {
	jsondecoder.RegisterFuzzyDecoders()
	Json = &jsonTool{
		pbUnmarshaler: &jsonpb.Unmarshaler{},
		pbMarshaller: &jsonpb.Marshaler{
			EnumsAsInts:  true,
			EmitDefaults: true,
			Indent:       "",
			OrigName:     true,
			AnyResolver:  nil,
		},
		normalJson: jsoniter.ConfigCompatibleWithStandardLibrary,
		numberJson: jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			UseNumber:              true,
		}.Froze(),
	}
}

type jsonTool struct {
	pbUnmarshaler *jsonpb.Unmarshaler
	pbMarshaller  *jsonpb.Marshaler
	normalJson    jsoniter.API
	numberJson    jsoniter.API
}

// Unmarshal ...
func (j *jsonTool) Unmarshal(data []byte, v interface{}) error {
	pbVal, ok := v.(proto.Message)
	if !ok {
		return j.normalJson.Unmarshal(data, v)
	}
	return j.pbUnmarshaler.Unmarshal(bytes.NewBuffer(data), pbVal)
}

// UnmarshalFromString ...
func (j *jsonTool) UnmarshalFromString(str string, v interface{}) error {
	pbVal, ok := v.(proto.Message)
	if !ok {
		return j.normalJson.UnmarshalFromString(str, v)
	}
	return j.pbUnmarshaler.Unmarshal(strings.NewReader(str), pbVal)
}

// Marshal ...
func (j *jsonTool) Marshal(v interface{}) ([]byte, error) {
	// check nil
	if v == nil {
		return []byte("null"), nil
	}

	pbVal, ok := v.(proto.Message)
	if !ok {
		return j.normalJson.Marshal(v)
	}

	var buf bytes.Buffer
	err := j.pbMarshaller.Marshal(&buf, pbVal)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalToString ...
func (j *jsonTool) MarshalToString(v interface{}) (string, error) {
	if v == nil {
		return "null", nil
	}
	pbVal, ok := v.(proto.Message)
	if !ok {
		return j.normalJson.MarshalToString(v)
	}
	return j.pbMarshaller.MarshalToString(pbVal)
}

// UseNumber ...
func (j *jsonTool) UseNumber() jsoniter.API {
	return j.numberJson
}

// ToJson 将结构体转换为 json 字符串，如果转换失败则返回空字符串
func ToJson(v any) string {
	str, _ := Json.MarshalToString(v)
	return str
}

func ToIndentJson(v any) string {
	js := ToJson(v)
	buf := bytes.NewBuffer(make([]byte, 0, len(js)))
	err := json.Indent(buf, []byte(js), "", "\t")
	if err != nil {
		return js
	}
	return buf.String()
}

// LoadJson 将 json 字符串转换为结构体，如果 withDefault 为 true, 则在转换失败时返回结构体的默认值
func LoadJson[T any](json string, withDefault bool) *T {
	v := new(T)
	if err := Json.UnmarshalFromString(json, v); err != nil && !withDefault {
		return nil
	}
	return v
}
