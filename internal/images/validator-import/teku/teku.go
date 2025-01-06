package teku

import (
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//go:embed context
var tekuContext embed.FS

func InitContext() (string, error) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sedge-validator-import-context-*")
	if err != nil {
		return "", fmt.Errorf("error creating teku validator import dir context: %s", err.Error())
	}

	contextDir, err := tekuContext.ReadDir("context")
	if err != nil {
		return "", err
	}
	for _, item := range contextDir {
		itemData, err := tekuContext.ReadFile(strings.Join([]string{"context", item.Name()}, "/"))
		if err != nil {
			return "", err
		}
		if err := ioutil.WriteFile(filepath.Join(tempDir, item.Name()), itemData, 0o755); err != nil {
			return "", err
		}
	}
	return tempDir, nil
}
