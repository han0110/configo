package configo

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/han0110/configo/util"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// ParseFile parses files into a map.
func ParseFile(filepaths []string) (*util.FlattenMap, error) {
	if len(filepaths) == 0 {
		return util.NewFlattenMap(), nil
	}

	flattenMap := util.NewFlattenMap()
	for _, path := range filepaths {
		var flattener flattener

		if path == "" {
			return nil, errors.New("expected config file path, but got empty string")
		}

		// Check whether file extension is supported
		switch strings.ToLower(filepath.Ext(path)) {
		case ".yaml", ".yml":
			flattener = &yamlFlattener{}
		// TODO: Support json with ordered map library
		// case ".json":
		case "":
			return nil, errors.Errorf(
				"unrecognized config file extension, file: %s",
				path,
			)
		default:
			return nil, errors.Errorf(
				"unsupported config file extension: %s",
				filepath.Ext(path),
			)
		}

		// Read file
		data, err := ioutil.ReadFile(filepath.Clean(path))
		if err != nil {
			return nil, err
		}

		if err := flattener.Flatten(flattenMap, data); err != nil {
			return nil, err
		}
	}

	return flattenMap, nil
}

type flattener interface {
	Flatten(*util.FlattenMap, []byte) error
}

type yamlFlattener struct {
}

func (flattener *yamlFlattener) Flatten(flattenMap *util.FlattenMap, data []byte) error {
	dec := yaml.NewDecoder(bytes.NewReader([]byte(data)))
	for {
		err := dec.Decode(&yamlCursor{setter: flattenMap})
		if err == nil {
			continue
		}
		if err == io.EOF {
			break
		}
		return err
	}
	return nil
}

type yamlCursor struct {
	setter interface {
		Set(key, value string)
	}
	key string
}

func (cur *yamlCursor) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.SequenceNode:
		for i := range node.Content {
			itemNode := node.Content[i]
			key := strconv.Itoa(i)
			if cur.key != "" {
				key = cur.key + util.CharDot + key
			}
			if err := (&yamlCursor{setter: cur.setter, key: key}).UnmarshalYAML(itemNode); err != nil {
				return err
			}
		}
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			keyNode, valNode := node.Content[i], node.Content[i+1]
			key := keyNode.Value
			if cur.key != "" {
				key = cur.key + util.CharDot + key
			}
			if err := (&yamlCursor{setter: cur.setter, key: key}).UnmarshalYAML(valNode); err != nil {
				return err
			}
		}
	case yaml.ScalarNode:
		cur.setter.Set(cur.key, node.Value)
	}
	return nil
}
