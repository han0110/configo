package node

import (
	"reflect"
	"strconv"

	"github.com/han0110/configo/util"
	"github.com/pkg/errors"
)

// FlattenMap provides data to fill a node.
type FlattenMap interface {
	Value(key string) string
	ChildrenByPrefix(prefix string) []string
}

type nodeFiller struct {
	FlattenMap
}

func (filler *nodeFiller) fill(node *Node) error {
	callback := func(n *Node) error {
		return filler.fillNode(n, n.Value)
	}

	if err := node.Walk(callback); err != nil {
		return err
	}

	return nil
}

func (filler *nodeFiller) fillNode(node *Node, rValue reflect.Value) error {
	switch rValue.Kind() {
	case reflect.Map:
		return filler.fillMap(node, rValue)
	case reflect.Slice:
		return filler.fillSlice(node, rValue)
	case reflect.Ptr:
		return filler.fillNode(node, rValue.Elem())
	default:
		if value := filler.FlattenMap.Value(node.Key); value != "" {
			return filler.fillSingle(rValue, value)
		}
	}
	return nil
}

func (filler *nodeFiller) fillMap(node *Node, rValue reflect.Value) error {
	if rValue.IsZero() {
		rValue.Set(reflect.MakeMapWithSize(rValue.Type(), 0))
	}
	prefix := node.Key + util.CharDot
	for _, key := range filler.FlattenMap.ChildrenByPrefix(prefix) {
		child := node.Children[0].clone()
		child.reKey(prefix, prefix+key)
		if err := filler.fill(child); err != nil {
			return err
		}
		rValue.SetMapIndex(reflect.ValueOf(key), child.Value)
	}
	return nil
}

func (filler *nodeFiller) fillSlice(node *Node, rValue reflect.Value) error {
	length := rValue.Len()
	children := make(map[int64]reflect.Value)
	prefix := node.Key + util.CharDot
	for _, key := range filler.FlattenMap.ChildrenByPrefix(prefix) {
		childIndex, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			return err
		}
		if int(childIndex)+1 > length {
			length = int(childIndex) + 1
		}
		child := node.Children[0].clone()
		child.reKey(prefix, prefix+key)
		if err := filler.fill(child); err != nil {
			return err
		}
		children[childIndex] = child.Value
	}
	newSlice := reflect.MakeSlice(rValue.Type(), length, length)
	reflect.Copy(newSlice, rValue)
	for index, value := range children {
		newSlice.Index(int(index)).Set(value)
	}
	rValue.Set(newSlice)
	return nil
}

func (filler *nodeFiller) fillSingle(rValue reflect.Value, value string) error {
	switch rValue.Kind() {
	case reflect.String:
		rValue.SetString(value)
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		rValue.SetBool(val)
	case reflect.Int8:
		return filler.setInt(rValue, value, 8)
	case reflect.Int16:
		return filler.setInt(rValue, value, 16)
	case reflect.Int32:
		return filler.setInt(rValue, value, 32)
	case reflect.Int64, reflect.Int:
		return filler.setInt(rValue, value, 64)
	case reflect.Uint8:
		return filler.setUint(rValue, value, 8)
	case reflect.Uint16:
		return filler.setUint(rValue, value, 16)
	case reflect.Uint32:
		return filler.setUint(rValue, value, 32)
	case reflect.Uint64, reflect.Uint:
		return filler.setUint(rValue, value, 64)
	case reflect.Float32:
		return filler.setFloat(rValue, value, 32)
	case reflect.Float64:
		return filler.setFloat(rValue, value, 64)
	default:
		panic(errors.Errorf("unexpected call on fillSingle by node of kind %s", rValue.Kind()))
	}
	return nil
}

func (filler *nodeFiller) setInt(rValue reflect.Value, value string, bitSize int) error {
	val, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		return err
	}
	rValue.SetInt(val)
	return nil
}

func (filler *nodeFiller) setUint(rValue reflect.Value, value string, bitSize int) error {
	val, err := strconv.ParseUint(value, 10, bitSize)
	if err != nil {
		return err
	}
	rValue.SetUint(val)
	return nil
}

func (filler *nodeFiller) setFloat(rValue reflect.Value, value string, bitSize int) error {
	val, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		return err
	}
	rValue.SetFloat(val)
	return nil
}
