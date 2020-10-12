package configo

import (
	"github.com/han0110/configo/node"
	"go.uber.org/fx"
)

// Default provides a common usage of Configo with env, file, flag loaders.
func Default() *Configo {
	return &Configo{
		Loader: Loaders{
			&EnvLoader{},
			&FileLoader{DisallowUnused: true},
			&FlagLoader{DisallowUnused: true},
		},
	}
}

// Configo is the main structure of configo which wraps all utilities for quick usage.
type Configo struct {
	Loader         Loader
	TagName        string
	TagDescription string
}

// Load loads configurations into config, optionally used with arguments to do so.
func (configo *Configo) Load(config interface{}, args []string) error {
	// Encode config to node for loaders to load data into it
	n, err := node.New(config, node.EncoderOption{
		TagName:        configo.TagName,
		TagDescription: configo.TagDescription,
	})
	if err != nil {
		return err
	}

	// Load data into conifg
	if err := configo.Loader.Load(n, args); err != nil {
		return err
	}

	return nil
}

// AsFxOption provides as a fx.Option to help load and supply config.
func (configo *Configo) AsFxOption(config interface{}, args []string) fx.Option {
	if err := configo.Load(config, args); err != nil {
		return fx.Error(err)
	}
	return fx.Supply(fx.Annotated{Target: config})
}
