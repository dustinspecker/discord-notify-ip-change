package internal_test

import (
	"io"
	"testing"

	"github.com/onsi/gomega"

	"github.com/dustinspecker/discord-notify-ip-change/internal"
)

func TestRender(t *testing.T) {
	testCases := []struct {
		name     string
		template string
		data     any
		expected []byte
	}{
		{
			name:     "simple string",
			template: "hello",
			expected: []byte("hello"),
		},
		{
			name:     "simple template",
			template: "{{ .IP }}",
			data:     struct{ IP string }{IP: "192.168.0.1"},
			expected: []byte("192.168.0.1"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			g := gomega.NewWithT(t)

			renderedMessage, err := internal.RenderMessage(testCase.template, testCase.data)
			g.Expect(err).To(gomega.BeNil(), "error rendering message")

			renderedMessageStr, err := io.ReadAll(renderedMessage)
			g.Expect(err).To(gomega.BeNil(), "error converting renderedMessage to string")
			g.Expect(renderedMessageStr).To(gomega.Equal(testCase.expected))
		})
	}
}

func TestRenderReturnsErrorForInvalidTemplate(t *testing.T) {
	g := gomega.NewWithT(t)

	_, err := internal.RenderMessage("{{ . }", nil)

	g.Expect(err).ToNot(gomega.BeNil(), "expected an error for a bad template")

	g.Expect(err.Error()).To(gomega.Equal(`error parsing template: template: message:1: unexpected "}" in operand`))
}

func TestRenderReturnsErrorWhenUnableToExecute(t *testing.T) {
	g := gomega.NewWithT(t)

	_, err := internal.RenderMessage("{{ .IP }}", struct{}{})

	g.Expect(err).ToNot(gomega.BeNil(), "expected an error when unable to execute template")

	g.Expect(err.Error()).To(gomega.Equal(`error executing template: template: message:1:3: executing "message" at <.IP>: can't evaluate field IP in type struct {}`))
}
