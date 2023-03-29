//go:build integration
// +build integration

package integration_test

import (
	"net/http"
	"os"
	"os/exec"
	"testing"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

func TestRetrievePublicIP(t *testing.T) {
	g := gomega.NewWithT(t)

	discordNotifyIPChangeCLI, err := gexec.Build("github.com/dustinspecker/discord-notify-ip-change/cmd/discord-notify-ip-change")
	g.Expect(err).To(gomega.BeNil(), "failed to build discord-notify-ip-change")

	ghttptest := ghttp.NewGHTTPWithGomega(g)

	ipServer := ghttp.NewServer()
	defer ipServer.Close()

	ipServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttptest.RespondWith(http.StatusOK, `{"ip": "192.168.0.1"}"`),
		),
	)

	discordServer := ghttp.NewServer()
	defer discordServer.Close()

	discordServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttptest.VerifyContentType("application/json"),
			ghttptest.VerifyRequest(http.MethodPost, "/"),
			ghttptest.VerifyBody([]byte(`{"content": "192.168.0.1"}`)),
		),
	)

	command := exec.Command(discordNotifyIPChangeCLI, "-ip-url", ipServer.URL(), "-discord-webhook-url", discordServer.URL(), "-timeout", "5s")
	session, err := gexec.Start(command, os.Stdout, os.Stderr)
	g.Expect(err).To(gomega.BeNil(), "error while running command")

	g.Eventually(session).Should(gexec.Exit(0), "command should exit with 0 return code")

	g.Expect(ipServer.ReceivedRequests()).Should(gomega.HaveLen(1), "expected only 1 request made to server")
	g.Expect(discordServer.ReceivedRequests()).Should(gomega.HaveLen(1), "expected only 1 message to be sent to discord")
}

func TestErrorIsReturnedWhenUnableToGetIP(t *testing.T) {
	g := gomega.NewWithT(t)

	discordNotifyIPChangeCLI, err := gexec.Build("github.com/dustinspecker/discord-notify-ip-change/cmd/discord-notify-ip-change")
	g.Expect(err).To(gomega.BeNil(), "failed to build discord-notify-ip-change")

	command := exec.Command(discordNotifyIPChangeCLI, "-ip-url", "", "-timeout", "5s")
	session, err := gexec.Start(command, os.Stdout, os.Stderr)
	g.Expect(err).To(gomega.BeNil(), "error while running command")

	g.Eventually(session).Should(gexec.Exit(1), "command should exit with non-zero return code")

	g.Expect(session.Err).To(gbytes.Say(`.*error getting public IP: error getting URL.*`), "should print error")
}

func TestErrorIsReturnedWhenUnableToSendMessage(t *testing.T) {
	g := gomega.NewWithT(t)

	discordNotifyIPChangeCLI, err := gexec.Build("github.com/dustinspecker/discord-notify-ip-change/cmd/discord-notify-ip-change")
	g.Expect(err).To(gomega.BeNil(), "failed to build discord-notify-ip-change")

	ghttptest := ghttp.NewGHTTPWithGomega(g)

	ipServer := ghttp.NewServer()
	defer ipServer.Close()

	ipServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttptest.RespondWith(http.StatusOK, `{"ip": "192.168.0.1"}"`),
		),
	)

	command := exec.Command(discordNotifyIPChangeCLI, "-ip-url", ipServer.URL(), "-discord-webhook-url", "", "-timeout", "5s")
	session, err := gexec.Start(command, os.Stdout, os.Stderr)
	g.Expect(err).To(gomega.BeNil(), "error while running command")

	g.Eventually(session).Should(gexec.Exit(1), "command should exit with 1 return code")
	g.Expect(session.Err).To(gbytes.Say(`.*error sending message to discord: error sending message: Post "": unsupported protocol scheme "".*`), "should print error")

	g.Expect(ipServer.ReceivedRequests()).Should(gomega.HaveLen(1), "expected only 1 request made to server")
}

func TestErrorIsReturnedWhenUnableToRenderMessage(t *testing.T) {
	g := gomega.NewWithT(t)

	discordNotifyIPChangeCLI, err := gexec.Build("github.com/dustinspecker/discord-notify-ip-change/cmd/discord-notify-ip-change")
	g.Expect(err).To(gomega.BeNil(), "failed to build discord-notify-ip-change")

	ghttptest := ghttp.NewGHTTPWithGomega(g)

	ipServer := ghttp.NewServer()
	defer ipServer.Close()

	ipServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttptest.RespondWith(http.StatusOK, `{"ip": "192.168.0.1"}"`),
		),
	)

	command := exec.Command(discordNotifyIPChangeCLI, "-ip-url", ipServer.URL(), "-discord-webhook-url", "", "-timeout", "5s", "-format", "{{ .}")
	session, err := gexec.Start(command, os.Stdout, os.Stderr)
	g.Expect(err).To(gomega.BeNil(), "error while running command")

	g.Eventually(session).Should(gexec.Exit(1), "command should exit with 1 return code")
	g.Expect(session.Err).To(gbytes.Say(`.*error rendering message: error parsing template: template: message:1: bad character.*`), "should print error")

	g.Expect(ipServer.ReceivedRequests()).Should(gomega.HaveLen(1), "expected only 1 request made to server")
}

func TestErrorIsReturnedWhenUnableToParseTimeout(t *testing.T) {
	g := gomega.NewWithT(t)

	discordNotifyIPChangeCLI, err := gexec.Build("github.com/dustinspecker/discord-notify-ip-change/cmd/discord-notify-ip-change")
	g.Expect(err).To(gomega.BeNil(), "failed to build discord-notify-ip-change")

	command := exec.Command(discordNotifyIPChangeCLI, "-ip-url", "", "-timeout", "10f")
	session, err := gexec.Start(command, os.Stdout, os.Stderr)
	g.Expect(err).To(gomega.BeNil(), "error while running command")

	g.Eventually(session).Should(gexec.Exit(1), "command should exit with non-zero return code")

	g.Expect(session.Err).To(gbytes.Say(`.*unable to parse timeout: time: unknown unit "f" in duration "10f".*`), "should print error")
}
