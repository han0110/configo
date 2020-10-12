# Configo

Configo provides an auto configuration loading tool without bothering to define a lot of flags, which is also highly extensible.

## Usage

For the most simple usage, we define our config structure, and then let configo help to load configuration from env, file, flag by build-in loaders as follows:

```golang
// ./examples/simple/main.go

type Config struct {
        Foo string `yaml:"foo"`
}

func main() {
        var config Config
        _ = configo.Default().Load(&config, os.Args[1:])
}
```

Then we can run it like:

```bash
# with env
FOO=bar go run ./examples/simple/main.go

# with flag
go run ./examples/simple/main.go --foo bar

# with file config.yaml
go run ./examples/simple/main.go -f ./examples/simple/config.yaml
```

## Custom Configuration Loader

TODO

## Acknowledgement

The whole idea came from the awesome project [traefik](https://github.com/traefik/traefik)'s package [`config`](https://github.com/traefik/traefik/tree/master/pkg/config), which has much more features to do configuration loading than current work. Configo is just another lightweight choice to do so.
