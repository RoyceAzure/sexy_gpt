package gpt_error

import (
	"encoding/json"
	"strings"
)

type ErrorField struct {
	Field           string `json:"field"`
	ErrFieldMessage string `json:"err_message"`
}
type Message struct {
	ErrMessage []*ErrorField `json:"message"`
}

func NewErrField(field string, errMsg ...error) *ErrorField {
	var errs []string
	for _, err := range errMsg {
		errs = append(errs, err.Error())
	}
	return &ErrorField{Field: field, ErrFieldMessage: strings.Join(errs, ",")}
}

func AddErrFieldString(msg *Message, toAdd ...*ErrorField) {
	msg.ErrMessage = append(msg.ErrMessage, toAdd...)
}

func (m *Message) ToJson() (string, error) {
	jsonData, err := json.Marshal(m.ErrMessage)
	if err != nil {
		return "", ErrInvalidArgument
	}
	return string(jsonData), nil
}
