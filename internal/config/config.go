package config

import (
	"deploy-lander/internal/structs"
	"errors"
	"os"
	"strings"
)

const (
	CommonPrefix = "deploylander."
)

func LoadFromEnv() structs.ProjectMap {
	projects := make(map[string]structs.Project)

	for _, kv := range os.Environ() {
		if !strings.HasPrefix(kv, CommonPrefix) {
			continue
		}

		alias, param, value, err := parsePair(kv)
		if err != nil {
			continue
		}

		var ok bool
		var project structs.Project
		if project, ok = projects[alias]; !ok {
			project = structs.Project{Alias: alias}
		}
		assignParamValue(&project, param, value)

		projects[alias] = project
	}

	return projects
}

func assignParamValue(project *structs.Project, param, value string) {
	switch param {
	case "path":
		project.Path = value
	case "tags":
		project.TagsRe = value
	case "branches":
		project.BranchesRe = value
	case "deploy_command":
		project.DeployCommand = value
	}
}

func parsePair(pair string) (string, string, string, error) {
	kvArr := strings.SplitN(pair, "=", 2)

	if len(kvArr) != 2 {
		return "", "", "", errors.New("len != 2")
	}

	key := kvArr[0]
	if len(key) == 0 {
		return "", "", "", errors.New("len=0")
	}

	keyParts := strings.Split(key, ".")
	if len(keyParts) <= 3 {
		return "", "", "", errors.New("keyparts len <= 3")
	}

	return keyParts[1], keyParts[2], kvArr[1], nil
}