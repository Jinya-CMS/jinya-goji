package native

import (
	"log/slog"
	"os"

	"github.com/ebitengine/purego"
)

var (
	nativeFiles []string
)

func SetupNativeBridge(lib []byte, libName string) (*uintptr, error) {
	slog.Info("Loading " + libName + ".so from embed")

	tmpFile, err := os.CreateTemp("", libName)
	if err != nil {
		slog.Error("Failed to create temp file", "error", err)
		return nil, err
	}

	defer tmpFile.Close()

	slog.Info("Write "+libName+".so to temp directory", "path", tmpFile)
	_, err = tmpFile.Write(lib)
	if err != nil {
		slog.Error("Failed to write "+libName+".so", "error", err)
		return nil, err
	}

	err = tmpFile.Close()
	if err != nil {
		slog.Error("Failed to close "+libName+".so", "error", err)
		return nil, err
	}

	slog.Info("Perform Dlopen for "+libName+".so", "path", tmpFile.Name())
	loadedLib, err := purego.Dlopen(tmpFile.Name(), purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		slog.Error("Failed to open the templating library", "error", err)
		return nil, err
	}

	nativeFiles = append(nativeFiles, tmpFile.Name())

	return &loadedLib, nil
}

func Cleanup() {
	for _, file := range nativeFiles {
		err := os.Remove(file)
		if err != nil {
			slog.Error("Failed to remove libcompiler.so", "error", err)
		}
	}
}
