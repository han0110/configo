package configo

import (
	"strings"

	"github.com/han0110/configo/node"
	"github.com/pkg/errors"
)

var (
	// defaultEscapeUnused defines default escaped unused keys.
	defaultEscapeUnused = []string{defaultConfigFileFlag} // for FileLoader.
)

// FlagLoader loads config from flags.
type FlagLoader struct {
	DisallowUnused bool
	EscapeUnused   []string
}

var _ Loader = (*FlagLoader)(nil)

// Load implements ConfigLoader.
func (loader *FlagLoader) Load(n *node.Node, args []string) error {
	if len(loader.EscapeUnused) == 0 {
		loader.EscapeUnused = defaultEscapeUnused
	}

	// Parse data
	flattenMap, err := ParseFlag(args)
	if err != nil {
		return err
	}

	// Fill data into node
	if err := n.FillNode(flattenMap); err != nil {
		return err
	}

	// Check whether there are unused keys
	if loader.DisallowUnused {
		if unusedKeys := flattenMap.UnusedKeys(loader.EscapeUnused); len(unusedKeys) > 0 {
			for i := range unusedKeys {
				unusedKeys[i] = "--" + unusedKeys[i]
			}
			return errors.Errorf("unused flag %s", strings.Join(unusedKeys, ", "))
		}
	}

	return nil
}
