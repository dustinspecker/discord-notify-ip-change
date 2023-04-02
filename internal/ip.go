package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ipResponse struct {
	IP string
}

func GetIP(url string, timeout time.Duration) (string, error) {
	client := http.Client{
		Timeout: timeout,
	}
	response, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("error getting URL %q: %w", url, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting URL %q: encountered bad status: 404", url)
	}

	var ipResponse ipResponse
	err = json.NewDecoder(response.Body).Decode(&ipResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}

	return ipResponse.IP, nil
}
