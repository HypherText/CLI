package main

import "encoding/json"

// Broadly inferred from https://getcomposer.org/schema.json
type ComposerPackage struct {
	Name                string              `json:"name,omitempty"`
	Description         string              `json:"description,omitempty"`
	License             json.RawMessage     `json:"license,omitempty"` // string or []string
	Type                string              `json:"type,omitempty"`
	Abandoned           json.RawMessage     `json:"abandoned,omitempty"` // bool or string
	Version             string              `json:"version,omitempty"`
	DefaultBranch       bool                `json:"default-branch,omitempty"`
	NonFeatureBranches  []string            `json:"non-feature-branches,omitempty"`
	Keywords            []string            `json:"keywords,omitempty"`
	Readme              string              `json:"readme,omitempty"`
	Time                string              `json:"time,omitempty"`
	Authors             []Author            `json:"authors,omitempty"`
	Homepage            string              `json:"homepage,omitempty"`
	Support             map[string]string   `json:"support,omitempty"`
	Funding             []Funding           `json:"funding,omitempty"`
	Source              map[string]string   `json:"source,omitempty"`
	Dist                map[string]string   `json:"dist,omitempty"`
	Require             map[string]string   `json:"require,omitempty"`
	RequireDev          map[string]string   `json:"require-dev,omitempty"`
	Replace             map[string]string   `json:"replace,omitempty"`
	Conflict            map[string]string   `json:"conflict,omitempty"`
	Provide             map[string]string   `json:"provide,omitempty"`
	Suggest             map[string]string   `json:"suggest,omitempty"`
	Repositories        json.RawMessage     `json:"repositories,omitempty"`
	MinimumStability    string              `json:"minimum-stability,omitempty"`
	PreferStable        bool                `json:"prefer-stable,omitempty"`
	Autoload            map[string]any      `json:"autoload,omitempty"`
	AutoloadDev         map[string]any      `json:"autoload-dev,omitempty"`
	TargetDir           string              `json:"target-dir,omitempty"`
	IncludePath         []string            `json:"include-path,omitempty"`
	Bin                 json.RawMessage     `json:"bin,omitempty"` // string or []string
	Archive             map[string]any      `json:"archive,omitempty"`
	PhpExt              map[string]any      `json:"php-ext,omitempty"`
	Config              map[string]any      `json:"config,omitempty"`
	Extra               json.RawMessage     `json:"extra,omitempty"`
	Scripts             map[string][]string `json:"scripts,omitempty"`
	ScriptsDescriptions map[string]string   `json:"scripts-descriptions,omitempty"`
	ScriptsAliases      map[string][]string `json:"scripts-aliases,omitempty"`
}

type Author struct {
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Homepage string `json:"homepage,omitempty"`
	Role     string `json:"role,omitempty"`
}

type Funding struct {
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
}
