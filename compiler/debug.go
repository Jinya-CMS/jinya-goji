package compiler

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
)

var gojiCompilerSocket = "/var/run/goji-compiler.sock"

func init() {
	if os.Getenv("ENV") == "dev" {
		gojiCompilerSocket = "goji-compiler.sock"
	}
}

func getCompilerClient() (*http.Client, error) {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", gojiCompilerSocket)
			},
		},
	}

	return &httpc, nil
}

func debugCompileTheme(themeId int) error {
	client, err := getCompilerClient()
	if err != nil {
		return err
	}

	res, err := client.Post("http://localhost/compile/"+strconv.Itoa(themeId), "", nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to compile theme")
	}

	return nil
}
