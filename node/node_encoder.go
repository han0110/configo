package node

import (
	"fmt"
	"reflect"

	"github.com/han0110/configo/util"
	"github.com/pkg/errors"
)

// EncoderOption defines option for encoder.
type EncoderOption struct {
	TagName        string
	TagDescription string
}

type encoder struct {
	EncoderOption
}

func (coder *encoder) setNode(node *Node, rValue reflect.Value) error {
	switch rValue.Kind() {
	case reflect.Ptr:
		if rValue.IsNil() {
			rValue.Set(reflect.New(rValue.Type().Elem()))
		}
		return coder.setNode(node, rValue.Elem())
	case reflect.Struct:
		return coder.setStruct(node, rValue)
	case reflect.Map, reflect.Slice:
		return coder.setDynamic(node, rValue)
	case reflect.String,
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64:
		return nil
	default:
		panic(errors.Errorf("unexpected call on setNode by node of kind %s", rValue.Kind()))
	}
}

func (coder *encoder) setStruct(node *Node, rValue reflect.Value) error {
	rType := rValue.Type()

	for i := 0; i < rValue.NumField(); i++ {
		childField := rType.Field(i)
		childValue := rValue.Field(i)

		if !IsExported(&childField) {
			continue
		}

		if IsIgnored(&childField, coder.TagName) {
			continue
		}

		if ok := IsSupportedType(childField.Type); !ok {
			return errors.Errorf("type %s is not supported", childField.Type)
		}

		// Includes anonymous fileds into node's children
		if childField.Anonymous {
			if err := coder.setNode(node, childValue); err != nil {
				return err
			}
			continue
		}

		childName := childField.Tag.Get(coder.TagName)
		if childName == "" {
			childName = childField.Name
		}

		childKey := childName
		if node.Key != "" {
			childKey = node.Key + util.CharDot + childKey
		}

		child := &Node{
			Key:         util.ToDotCase(childKey),
			Name:        childName,
			Description: childField.Tag.Get(coder.TagDescription),
			FiledName:   childField.Name,
			Value:       childValue,
		}

		if err := coder.setNode(child, child.Value); err != nil {
			return err
		}

		node.Children = append(node.Children, child)
	}

	return nil
}

func (coder *encoder) setDynamic(node *Node, rValue reflect.Value) error {
	var childKey, childDescription string

	if node.Key != "" {
		childKey = node.Key + util.CharDot + childKey
	}

	switch rValue.Kind() {
	case reflect.Map:
		childDescription = "indexed value by key in map"
	case reflect.Slice:
		childDescription = "nth item in list"
	default:
		panic(errors.Errorf("unexpected call on setDynamic by node of kind %s", rValue.Kind()))
	}

	child := &Node{
		Key:         childKey,
		Description: childDescription,
		Value:       reflect.New(rValue.Type().Elem()).Elem(),
	}
	if err := coder.setNode(child, child.Value); err != nil {
		return err
	}

	// Add a template child for clone when filling data
	node.Children = []*Node{child}

	return nil
}

// IsExported checks whether field is exported..
func IsExported(rField *reflect.StructField) bool {
	return rField.PkgPath == ""
}

// IsIgnored checks whether field is ignored..
func IsIgnored(rField *reflect.StructField, tagName string) bool {
	return rField.Tag.Get(tagName) == "-"
}

// IsSupportedType checks whether type is supported.
func IsSupportedType(rType reflect.Type) bool {
	switch rType.Kind() {
	case reflect.String,
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Struct,
		reflect.Ptr:
		return true
	case reflect.Slice:
		return IsSupportedType(rType.Elem())
	case reflect.Map:
		if rType.Key().Kind() == reflect.String {
			return true
		}
	}

	fmt.Println(rType.Kind())

	return false
}
