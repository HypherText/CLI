package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

type Builder struct {
	workdir string
}

type Asset struct {
	file     string
	importer string
}

func (b *Builder) CheckValidWd() error {
	composerFile := filepath.Join(b.workdir, "composer.json")
	_, err := os.Stat(composerFile)
	if os.IsNotExist(err) {
		return fmt.Errorf("Working directory is not a valid HypherText project: %w.", err)
	}
	return nil
}

func (b *Builder) GetAppPath() (string, error) {
	if err := b.CheckValidWd(); err != nil {
		return "", err
	}

	bytes, err := os.ReadFile(filepath.Join(b.workdir, "composer.json"))
	if err != nil {
		return "", fmt.Errorf("Failed to read composer.json file: %w", err)
	}

	pkg := ComposerPackage{}
	err = json.Unmarshal(bytes, &pkg)
	if err != nil {
		return "", fmt.Errorf("Failed to read composer.json file: %w", err)
	}

	namespaces, ok := pkg.Autoload["psr-4"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("Found no PSR-4 namespaces in composer.json.")
	}

	var appPath string
	for ns, path := range namespaces {
		if ns == "App\\" {
			appPath, ok = path.(string)
			break
		}
	}

	if appPath == "" {
		return "", fmt.Errorf("composer.json must have a namespace called \"App\\\".")
	}

	if !ok {
		return "", fmt.Errorf("\"App\\\" namespace does not resolve to a path")
	}

	return path.Join(b.workdir, appPath), nil
}

func (b *Builder) GetCachePath() (string, error) {
	if err := b.CheckValidWd(); err != nil {
		return "", err
	}

	return filepath.Join(b.workdir, ".hypher"), nil
}

func (b *Builder) LocateAssets(root string) ([]Asset, error) {
	assets := []Asset{}
	importRegex := regexp.MustCompile(`(?s)<\?(?:|=|php) .*?import\( *"(.+?)" *\)`)
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path.Ext(p) != ".php" {
			return nil
		}

		bytes, err := os.ReadFile(p)
		if err != nil {
			return err
		}

		matches := importRegex.FindAllSubmatch(bytes, -1)
		for _, match := range matches {
			assets = append(assets, Asset{
				file:     filepath.Join(p, "..", string(match[1])),
				importer: p,
			})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error locating files: %w", err)
	}

	return assets, nil
}

func (b *Builder) BuildAssets() {
	appPath, err := b.GetAppPath()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	assets, err := b.LocateAssets(appPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cachePath, err := b.GetCachePath()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	distPath := path.Join(cachePath, "dist")
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		os.Mkdir(distPath, 0755)
	}

	entrypoints := []string{}
	for _, asset := range assets {
		entrypoints = append(entrypoints, asset.file)
	}

	result := api.Build(api.BuildOptions{
		EntryPoints:       entrypoints,
		Outdir:            distPath,
		Bundle:            true,
		MinifySyntax:      true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		Write:             true,
		Sourcemap:         api.SourceMapExternal,
	})

	for _, file := range result.OutputFiles {
		fileName := filepath.Base(file.Path)
		fullExt := fileName[strings.Index(fileName, "."):]
		if fullExt == ".module.css" {
			if _, err = os.Stat(file.Path + ".map"); os.IsNotExist(err) {
				fmt.Printf("Warning: Found CSS module with no SourceMap, skipping.")
				continue
			}
			GenerateStubs(file.Path)
		}
	}
}

func GenerateStubs(outputPath string) error {
	mapPath := outputPath + ".map"
	data, err := os.ReadFile(mapPath)
	if err != nil {
		return fmt.Errorf("Could not read %s: %w\n", mapPath, err)
	}

	sourceMap, err := ParseSourceMap(data)
	if err != nil {
		return fmt.Errorf("Could not parse SourceMap %s: %w\n", mapPath, err)
	}

	data, err = os.ReadFile(outputPath)
	if err != nil {
		return fmt.Errorf("Could not read %s: %w\n", mapPath, err)
	}

	generatedNameRegex := regexp.MustCompile(`-?[_a-zA-Z]+[_a-zA-Z0-9-]*`)
	nameMap := map[string]string{}
	for _, vlq := range sourceMap.Mappings {
		if vlq.NameIndex != -1 {
			name := sourceMap.Names[vlq.NameIndex]
			generatedName := string(generatedNameRegex.Find(data[vlq.Column:]))
			nameMap[name] = generatedName
		}
	}

	// TODO: Make php type stubs out of this name map

	return nil
}
