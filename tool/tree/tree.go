package tree

import (
	"errors"
	"regexp"
	"strings"
	"sync"
)

func NewTree() *Tree {
	return &Tree{
		Root: newNode(),
	}
}

type Tree struct {
	Root *Node
}

func (tree *Tree) parse(key string) (keys []string) {
	r := regexp.MustCompile(`(\w+)|'(.*?)'|"(.*?)"|(\$\{.*?\})`)
	temp := r.FindAllString(key, -1)
	for _, key := range temp {
		if strings.HasPrefix(key, `'`) && strings.HasSuffix(key, `'`) {
			key = key[1 : len(key)-1]
		}
		if strings.HasPrefix(key, `"`) && strings.HasSuffix(key, `"`) {
			key = key[1 : len(key)-1]
		}
		keys = append(keys, key)
	}
	return
}

func (tree *Tree) Insert(key string, value []byte) error {
	return tree.Root.Insert(tree.parse(key), value)
}

func (tree *Tree) Update(key string, value []byte) error {
	return tree.Root.Update(tree.parse(key), value)
}

func (tree *Tree) Delete(key string) error {
	return tree.Root.Delete(tree.parse(key))
}

func (tree *Tree) Search(key string) ([]byte, error) {
	return tree.Root.Search(tree.parse(key))
}

func newNode() *Node {
	return &Node{
		Children: make(map[string]*Node),
	}
}

type Node struct {
	sync.Mutex
	Index         []byte
	Children      map[string]*Node
	ParamterChild *Node
}

func (node *Node) isParmeter(key string) bool {
	ok, _ := regexp.MatchString(`(\$\{.*?\})`, key)
	return ok
}

func (node *Node) Insert(keys []string, value []byte) error {
	node.Lock()
	defer node.Unlock()

	if len(keys) <= 0 {
		if node.Index != nil {
			return errors.New("key already exists")
		}
		node.Index = value
		// stop recursive
		return nil
	}

	key := keys[0]
	if key == "" {
		return errors.New("empty key name")
	}

	var children *Node

	if node.isParmeter(key) {
		if node.ParamterChild == nil {
			node.ParamterChild = newNode()
		}
		children = node.ParamterChild
		return children.Insert(keys[1:], value)
	}

	if node.Children[key] != nil {
		children = node.Children[key]
	} else {
		children = newNode()
		node.Children[key] = children
	}

	return children.Insert(keys[1:], value)
}

func (node *Node) Update(keys []string, value []byte) error {
	node.Lock()
	defer node.Unlock()

	if len(keys) <= 0 {
		node.Index = value
		return nil
	}

	key := keys[0]
	if key == "" {
		return errors.New("empty key name")
	}

	var children *Node

	if node.isParmeter(key) {
		if node.ParamterChild == nil {
			node.ParamterChild = newNode()
		}
		children = node.ParamterChild
		return children.Update(keys[1:], value)
	}

	if node.Children[key] != nil {
		children = node.Children[key]
	} else {
		children = newNode()
		node.Children[key] = children
	}

	return children.Update(keys[1:], value)
}

func (node *Node) Delete(keys []string) error {
	node.Lock()
	defer node.Unlock()

	if len(keys) <= 0 {
		node.Index = nil
		return nil
	}
	key := keys[0]

	var children *Node

	isParmeter := node.isParmeter(key)
	if isParmeter {
		children = node.ParamterChild
	} else {
		children = node.Children[key]
	}

	var err error
	if children != nil {
		err = children.Delete(keys[1:])
	} else {
		err = errors.New("wrong key")
	}

	if err != nil {
		return err
	}

	// if children child size is zero , remove this children
	if len(children.Children) < 1 {
		if isParmeter {
			node.ParamterChild = nil
		} else {
			delete(node.Children, key)
		}
	}
	return nil
}

func (node *Node) Search(keys []string) ([]byte, error) {
	node.Lock()
	defer node.Unlock()

	if len(keys) <= 0 {
		if node.Index != nil {
			return node.Index, nil
		}
		return nil, errors.New("key is not exists")
	}

	key := keys[0]
	children := node.Children[key]
	if children == nil {
		if node.ParamterChild == nil {
			return nil, errors.New("wrong key")
		}
		children = node.ParamterChild
	}

	return children.Search(keys[1:])
}
