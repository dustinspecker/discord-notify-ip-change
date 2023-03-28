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

	server := ghttp.NewServer()
	defer server.Close()

	server.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.RespondWith(http.StatusOK, `{"ip": "192.168.0.1"}"`),
		),
	)

	command := exec.Command(discordNotifyIPChangeCLI, "-ip-url", server.URL())
	session, err := gexec.Start(command, os.Stdout, os.Stderr)
	g.Expect(err).To(gomega.BeNil(), "error while running command")

	g.Eventually(session).Should(gexec.Exit(0), "command should exit with 0 return code")

	g.Expect(session.Out).To(gbytes.Say("192.168.0.1"), "should print ip address")

	g.Expect(server.ReceivedRequests()).Should(gomega.HaveLen(1), "expected only 1 request made to server")
}

func TestErrorIsReturnedWhenUnableToGetIP(t *testing.T) {
	g := gomega.NewWithT(t)

	discordNotifyIPChangeCLI, err := gexec.Build("github.com/dustinspecker/discord-notify-ip-change/cmd/discord-notify-ip-change")
	g.Expect(err).To(gomega.BeNil(), "failed to build discord-notify-ip-change")

	command := exec.Command(discordNotifyIPChangeCLI, "-ip-url", "")
	session, err := gexec.Start(command, os.Stdout, os.Stderr)
	g.Expect(err).To(gomega.BeNil(), "error while running command")

	g.Eventually(session).Should(gexec.Exit(1), "command should exit with non-zero return code")

	g.Expect(session.Err).To(gbytes.Say(`.*error getting public IP: error getting URL.*`), "should print error")
}
