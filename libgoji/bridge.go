package libgoji

import (
	_ "embed"
	"log/slog"
	"os"

	"github.com/ebitengine/purego"
)

var (
	nativeFiles []string
)

//go:embed target/release/libgoji.so
var libgoji []byte

func SetupLibGoji() error {
	slog.Info("Loading libgoji.so from embed")

	tmpFile, err := os.CreateTemp("", "libgoji")
	if err != nil {
		slog.Error("Failed to create temp file", "error", err)
		return err
	}

	defer tmpFile.Close()

	slog.Info("Write libgoji.so to temp directory", "path", tmpFile)
	_, err = tmpFile.Write(libgoji)
	if err != nil {
		slog.Error("Failed to write libgoji.so", "error", err)
		return err
	}

	err = tmpFile.Close()
	if err != nil {
		slog.Error("Failed to close libgoji.so", "error", err)
		return err
	}

	slog.Info("Perform Dlopen for libgoji.so", "path", tmpFile.Name())
	loadedLib, err := purego.Dlopen(tmpFile.Name(), purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		slog.Error("Failed to open the templating library", "error", err)
		return err
	}

	nativeFiles = append(nativeFiles, tmpFile.Name())

	bindHelper(loadedLib)

	bindCompiler(loadedLib)
	bindEngine(loadedLib)
	bindMailing(loadedLib)

	return nil
}

func Cleanup() {
	for _, file := range nativeFiles {
		err := os.Remove(file)
		if err != nil {
			slog.Error("Failed to remove libcompiler.so", "error", err)
		}
	}
}
