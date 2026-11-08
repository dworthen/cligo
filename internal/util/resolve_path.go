package util

import (
	"os"
	"path"
	"slices"
	"strings"
)

func ResolvePath(elem ...string) (string, error) {

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cwd = strings.ReplaceAll(cwd, "\\", "/")
	home = strings.ReplaceAll(home, "\\", "/")
	paths := []string{cwd}

	for _, p := range elem {
		p = strings.ReplaceAll(p, "\\", "/")
		if strings.HasPrefix(p, "~/") {
			p = strings.Replace(p, "~", home, 1)
		}

		paths = append(paths, p)
	}

	startIndex := 0
	for ind, v := range slices.Backward(paths) {
		if path.IsAbs(v) {
			startIndex = ind
			break
		}
	}
	return path.Join(paths[startIndex:]...), nil
}
