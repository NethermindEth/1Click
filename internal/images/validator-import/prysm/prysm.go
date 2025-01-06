package prysm

import (
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//go:embed context
var prysmContext embed.FS

func InitContext() (string, error) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sedge-validator-import-context-*")
	if err != nil {
		return "", fmt.Errorf("error creating prysm validator import dir context: %s", err.Error())
	}

	contextDir, err := prysmContext.ReadDir("context")
	if err != nil {
		return "", err
	}
	for _, item := range contextDir {
		itemData, err := prysmContext.ReadFile(strings.Join([]string{"context", item.Name()}, "/"))
		if err != nil {
			return "", err
		}
		if err := ioutil.WriteFile(filepath.Join(tempDir, item.Name()), itemData, 0o755); err != nil {
			return "", err
		}
	}
	return tempDir, nil
}
