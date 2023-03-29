package message

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
)

func Render(templateStr string, data any) (io.Reader, error) {
	var buffer bytes.Buffer

	tmpl, err := template.New("message").Parse(templateStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %v", err)
	}

	if err := tmpl.Execute(&buffer, data); err != nil {
		return nil, fmt.Errorf("error executing template: %v", err)
	}

	return &buffer, nil
}
