package frontend

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
)

func getClient() (*http.Client, error) {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/goji-frontend.sock")
			},
		},
	}

	return &httpc, nil
}

func CompileTheme(themeId int) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	res, err := client.Post("/compile/"+strconv.Itoa(themeId), "", nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to compile theme")
	}

	return nil
}
