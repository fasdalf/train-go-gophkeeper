package convertors

import (
	"github.com/fasdalf/train-go-gophkeeper/internal/client/model/entity"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/view/bubbleforms"
)

// PasswordToFields for entity.Password
func PasswordToFields(sData entity.SecretData) (fields []bubbleforms.InputField) {
	fields = []bubbleforms.InputField{
		{
			Type:  bubbleforms.InputTextPlain,
			Title: bubbleforms.UITextInputsModelName,
			Value: sData.Name,
		},
		{
			Type:  bubbleforms.InputTextPlain,
			Title: bubbleforms.UITextInputsModelPassword,
			Value: sData.Value.(entity.Password).Pass,
		},
	}
	return
}

// FieldsToPassword for entity.Password
func FieldsToPassword(fields []bubbleforms.InputField, md map[string]string) (sData entity.SecretData) {
	p := entity.Password{fields[1].Value}
	sData = entity.SecretData{
		Name:     fields[0].Value,
		Value:    p,
		Metadata: md,
	}
	return
}
