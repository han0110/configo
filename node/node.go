package node

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

const (
	// defaultTagName defines default tag key for name.
	defaultTagName = "yaml"
	// defaultTagDescription defines default tag key for description.
	defaultTagDescription = "description"
)

// New encodes element into node.
func New(element interface{}, option EncoderOption) (*Node, error) {
	rValue := reflect.ValueOf(element)
	rType := reflect.TypeOf(element)

	if ok := IsSupportedType(rType); !ok {
		return nil, fmt.Errorf("type %s is not supported", rType)
	}
	node := &Node{Value: rValue}

	if option.TagName == "" {
		option.TagName = defaultTagName
	}
	if option.TagDescription == "" {
		option.TagDescription = defaultTagDescription
	}

	return node, (&encoder{option}).setNode(node, node.Value)
}

// Node defines struct for a field in struct.
type Node struct {
	Key         string
	Name        string
	Description string
	FiledName   string
	Value       reflect.Value
	Children    []*Node
}

// WalkCallback defines function called when walk.
type WalkCallback func(*Node) error

// Walk walks through node and all children.
func (node *Node) Walk(callback WalkCallback) error {
	if len(node.Children) == 0 || node.isDynamic() {
		if err := callback(node); err != nil {
			return err
		}
	} else {
		for _, child := range node.Children {
			if err := child.Walk(callback); err != nil {
				return err
			}
		}
	}

	return nil
}

// FillNode fills data into node.
func (node *Node) FillNode(flattenMap FlattenMap) error {
	return (&nodeFiller{FlattenMap: flattenMap}).fill(node)
}

// Flat gets slice of all nodes sorted by name.
func (node *Node) Flat() []*Node {
	var nodes []*Node
	_ = node.Walk(func(n *Node) error {
		nodes = append(nodes, n)
		return nil
	})
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Key < nodes[j].Key })
	return nodes
}

// SerializeValue serializes node's value.
func (node *Node) SerializeValue() string {
	switch node.Value.Kind() {
	case reflect.Slice:
		var str []string
		for i, l := 0, node.Value.Len(); i < l; i++ {
			str = append(str, fmt.Sprint(node.Value.Index(i)))
		}
		return fmt.Sprintf("[%s]", strings.Join(str, ","))
	case reflect.Map:
		var str []string
		for iter := node.Value.MapRange(); iter.Next(); {
			str = append(str, fmt.Sprintf("%s=%s", iter.Key().String(), fmt.Sprint(iter.Value())))
		}
		return fmt.Sprintf("{%s}", strings.Join(str, ","))
	default:
		return fmt.Sprint(node.Value.Interface())
	}
}

func (node *Node) clone() *Node {
	clone := &Node{
		Key:         node.Key,
		Name:        node.Name,
		Description: node.Description,
		FiledName:   node.FiledName,
		Value:       reflect.New(node.Value.Type()).Elem(),
	}
	if node.Value.Kind() == reflect.Ptr {
		clone.Value = reflect.New(node.Value.Type().Elem())
	}
	if node.isDynamic() {
		clone.Children = []*Node{node.Children[0].clone()}
	} else {
		for _, child := range node.Children {
			childClone := child.clone()
			childClone.Value = reflect.Indirect(clone.Value).FieldByName(childClone.FiledName)
			clone.Children = append(clone.Children, childClone)
		}
	}
	return clone
}

func (node *Node) reKey(oldkey, newkey string) {
	node.Key = strings.ReplaceAll(node.Key, oldkey, newkey)
	for _, child := range node.Children {
		child.reKey(oldkey, newkey)
	}
}

func (node *Node) isDynamic() bool {
	kind := reflect.Indirect(node.Value).Kind()
	return kind == reflect.Map || kind == reflect.Slice
}
