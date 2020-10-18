package main

import (
	"fmt"
	"os"

	"github.com/han0110/configo"
)

// Config defines app's configuration
type Config struct {
	Foo string `yaml:"foo"`
}

func main() {
	var config Config
	_ = configo.Default().Load(&config, os.Args[1:])
	fmt.Printf("foo: %s\n", config.Foo)
}
