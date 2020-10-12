package configo

import (
	"os"
	"strings"

	"github.com/han0110/configo/node"
	"github.com/han0110/configo/util"
	"github.com/pkg/errors"
)

// EnvLoader loads from environment variables.
type EnvLoader struct {
	Prefix         string
	DisallowUnused bool
}

var _ Loader = (*EnvLoader)(nil)

// Load implements ConfigLoader.
func (loader *EnvLoader) Load(n *node.Node, args []string) error {
	if loader.Prefix != "" {
		loader.Prefix += util.CharUnderscore
	}

	flattenMap := util.NewFlattenMap()
	for _, env := range os.Environ() {
		key := strings.SplitN(env, "=", 2)[0]
		if !strings.HasPrefix(env, loader.Prefix) {
			continue
		}
		flattenMap.Set(key[len(loader.Prefix):], os.Getenv(key))
	}

	// Fill data into node
	if err := n.FillNode(flattenMap); err != nil {
		return err
	}

	// Check whether there are unused keys
	if loader.DisallowUnused {
		if unusedKeys := flattenMap.UnusedKeys(nil); len(unusedKeys) > 0 {
			if loader.Prefix != "" {
				for i := range unusedKeys {
					unusedKeys[i] = loader.Prefix + unusedKeys[i]
				}
			}
			return errors.Errorf("unused env %s", strings.Join(unusedKeys, ", "))
		}
	}

	return nil
}
