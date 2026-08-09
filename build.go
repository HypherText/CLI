package main

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

func GetAssetsPath() (string, error) {
	pwd, _ := os.Getwd()
	composerFile := path.Join(pwd, "composer.json")

	if _, err := os.Stat(composerFile); os.IsNotExist(err) {
		return "", errors.New("No composer.json file in current working directory.")
	}

	bytes, err := os.ReadFile(composerFile)
	if err != nil {
		return "", errors.Join(errors.New("Failed to read composer.json file."), err)
	}

	pkg := ComposerPackage{}
	err = json.Unmarshal(bytes, &pkg)
	if err != nil {
		return "", errors.Join(errors.New("Failed to read composer.json file."), err)
	}

	namespaces, ok := pkg.Autoload["psr-4"].(map[string]any)
	if !ok {
		return "", errors.New("Found no PSR-4 namespaces in composer.json.")
	}

	var appPath string
	for ns, path := range namespaces {
		if ns == "App\\" {
			appPath, ok = path.(string)
			break
		}
	}

	if appPath == "" {
		return "", errors.New("composer.json must have a namespace called \"App\\\".")
	}

	if !ok {
		return "", errors.New("\"App\\\" namespace does not resolve to a path")
	}

	return path.Join(appPath, "assets"), nil
}
