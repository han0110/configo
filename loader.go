package configo

import (
	"github.com/han0110/configo/node"
)

// Loader is a configuration store.
type Loader interface {
	// Load fills node, optionally used with arguments to do so.
	Load(n *node.Node, args []string) error
}

// Loaders wraps multiple Loaders as Loader interface
type Loaders []Loader

var _ Loader = (Loaders)(nil)

// Load implements Loader
func (loaders Loaders) Load(n *node.Node, args []string) error {
	for _, loader := range loaders {
		if err := loader.Load(n, args); err != nil {
			return err
		}
	}
	return nil
}
