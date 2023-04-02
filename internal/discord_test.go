package internal_test

import (
	"bytes"
	"testing"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/dustinspecker/discord-notify-ip-change/internal"
)

func TestSend(t *testing.T) {
	g := gomega.NewWithT(t)

	ghttptest := ghttp.NewGHTTPWithGomega(g)

	server := ghttp.NewServer()
	defer server.Close()
	server.AppendHandlers(
		ghttp.CombineHandlers(
			ghttptest.VerifyContentType("application/json"),
			ghttptest.VerifyBody([]byte("hey")),
		),
	)

	err := internal.SendMessage(server.URL(), bytes.NewReader([]byte("hey")))
	g.Expect(err).To(gomega.BeNil(), "error sending message")

	g.Expect(server.ReceivedRequests()).To(gomega.HaveLen(1), "expected message to only be sent once")
}

func TestSendReturnsErrorForInvalidURL(t *testing.T) {
	g := gomega.NewWithT(t)

	err := internal.SendMessage("", nil)
	g.Expect(err).ToNot(gomega.BeNil(), "expected an error for invalid URL")

	g.Expect(err.Error()).To(gomega.Equal(`error sending message: Post "": unsupported protocol scheme ""`))
}
