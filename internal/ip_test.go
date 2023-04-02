package internal_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/onsi/gomega"

	"github.com/dustinspecker/discord-notify-ip-change/internal"
)

func TestGetReturnsErrorWhenURLInvalid(t *testing.T) {
	g := gomega.NewWithT(t)

	_, err := internal.GetIP("", time.Second)
	g.Expect(err).ToNot(gomega.BeNil(), "error should be returned when invalid URL")

	g.Expect(err.Error()).To(gomega.Equal(`error getting URL "": Get "": unsupported protocol scheme ""`))
}

func TestGetReturnsErrorWhenURLNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	g := gomega.NewWithT(t)

	_, err := internal.GetIP(server.URL, time.Second)
	g.Expect(err).ToNot(gomega.BeNil(), "error should be returned when URL not found")

	g.Expect(err.Error()).To(gomega.Equal(fmt.Sprintf(`error getting URL %q: encountered bad status: 404`, server.URL)))
}

func TestGetReturnsErrorWhenInvalidResponseFormat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ip"`)
	}))
	defer server.Close()

	g := gomega.NewWithT(t)

	_, err := internal.GetIP(server.URL, time.Second)
	g.Expect(err).ToNot(gomega.BeNil(), "error should be returned when URL has invalid response format")

	g.Expect(err.Error()).To(gomega.Equal("error unmarshalling response: unexpected EOF"))
}

func TestGetReturnsErrorWhenContextTimeoutReached(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ip"`)
	}))
	defer server.Close()

	g := gomega.NewWithT(t)

	_, err := internal.GetIP(server.URL, time.Nanosecond)
	g.Expect(err).ToNot(gomega.BeNil(), "error should be returned when timeout is reached")

	g.Expect(err.Error()).To(gomega.Equal(fmt.Sprintf(`error getting URL %q: Get %q: context deadline exceeded (Client.Timeout exceeded while awaiting headers)`, server.URL, server.URL)))
}

func TestGetReturnsIP(t *testing.T) {
	expectedIP := "192.168.0.1"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fmt.Sprintf(`{"ip": "%s"}`, expectedIP))
	}))
	defer server.Close()

	g := gomega.NewWithT(t)

	output, err := internal.GetIP(server.URL, time.Second)
	g.Expect(err).To(gomega.BeNil(), "no error should be returned when URL is retrieved")

	g.Expect(output).To(gomega.Equal(expectedIP), "Get should return output from server")
}
