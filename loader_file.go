package configo

import (
	"os"
	"strings"

	"github.com/han0110/configo/node"
	"github.com/pkg/errors"
)

const (
	// defaultConfigFileFlag defines config file flag key.
	defaultConfigFileFlag = "f"
	// defaultConfigFileEnv defines config file env key.
	defaultConfigFileEnv = "CONFIG_FILE"
)

// FileLoader loads config from file.
type FileLoader struct {
	DisallowUnused bool
	ConfigFileFlag string
	ConfigFileEnv  string
}

var _ Loader = (*FileLoader)(nil)

// Load implements ConfigLoader.
func (loader *FileLoader) Load(n *node.Node, args []string) error {
	if loader.ConfigFileFlag == "" {
		loader.ConfigFileFlag = defaultConfigFileFlag
	}
	if loader.ConfigFileEnv == "" {
		loader.ConfigFileEnv = defaultConfigFileEnv
	}

	// Find config filepaths from flags and environments.
	filepaths := loader.findConfigFilePaths(args)

	// Parse files into map[string]string.
	flattenMap, err := ParseFile(filepaths)
	if err != nil {
		return err
	}

	// Fill data into node
	if err := n.FillNode(flattenMap); err != nil {
		return err
	}

	// Check whether there are unused keys
	if loader.DisallowUnused {
		if unusedKeys := flattenMap.UnusedKeys(nil); len(unusedKeys) > 0 {
			return errors.Errorf("unused keys %s", strings.Join(unusedKeys, ", "))
		}
	}

	return nil
}

func (loader *FileLoader) findConfigFilePaths(args []string) []string {
	var filepaths []string

	// Find config file from flags.
	if loader.ConfigFileFlag != "-" {
		if flattenMap, err := ParseFlag(args); err == nil {
			if value, ok := flattenMap.Value(loader.ConfigFileFlag); ok && value != "" {
				filepaths = append(filepaths, strings.Split(value, ",")...)
			}
		}
	}

	// Find config file from environments.
	if loader.ConfigFileEnv != "-" {
		if value, ok := os.LookupEnv(loader.ConfigFileEnv); ok {
			filepaths = append(filepaths, strings.Split(value, ",")...)
		}
	}

	return filepaths
}
