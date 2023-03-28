//go:build integration
// +build integration

package integration_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

func TestOutput(t *testing.T) {
	g := gomega.NewWithT(t)

	discordNotifyIPChangeCLI, err := gexec.Build("github.com/dustinspecker/discord-notify-ip-change/cmd/discord-notify-ip-change")
	g.Expect(err).To(gomega.BeNil(), "failed to build discord-notify-ip-change")

	command := exec.Command(discordNotifyIPChangeCLI)
	session, err := gexec.Start(command, os.Stdout, os.Stderr)
	g.Expect(err).To(gomega.BeNil(), "error while running command")

	g.Eventually(session).Should(gexec.Exit(0), "command should exit with 0 return code")

	g.Expect(session.Out).To(gbytes.Say("discord-notify-ip-change"), "should print command name")
}
