// Package configo provides an auto configuration loading tool without
// bothering to define a lot of flags, which is also highly extensible.
//
// For the most simple usage, we define our config structure, and then
// let configo help to load configuration from env, file, flag by build-in
// loaders as follows:
//
//	// ./examples/simple/main.go
//
//	type Config struct {
//		Foo string
//	}
//
//	func main() {
//		var config Config
//		_ = configo.Default().Load(&config, os.Args[1:])
//	}
//
// Then we can run ./examples/simple/main.go like:
//
//	// with env
//	FOO=bar go run ./examples/simple/main.go
//	// with flag
//	go run ./examples/simple/main.go --foo bar
//	// with file config.yaml
//	go run ./examples/simple/main.go -f ./examples/simple/config.yaml
//
// You will see more usage in ./examples
package configo
