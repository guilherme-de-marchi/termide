package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type FileNode struct {
	Info     fs.FileInfo
	Path     string
	AbsPath  string
	Tree     *FileTree
	children map[string]*FileNode
}

func NewFileNodeFrom(absPath string) (*FileNode, error) {
	if !filepath.IsAbs(absPath) {
		return nil, errors.New("argument absPath is not an absolute path")
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}
	return &FileNode{
		Info:    info,
		Path:    info.Name(),
		AbsPath: absPath,
	}, nil
}

func (fn *FileNode) AddChild(child *FileNode) error {
	if _, ok := fn.children[child.Path]; ok {
		return errors.New(fmt.Sprintf("file node %s already registered", fn.Path))
	}
	fn.children[child.Path] = child
	fn.Tree.RegisterNode(child)
	return nil
}

func (fn *FileNode) HasChildren() bool {
	return len(fn.children) > 0
}

func (fn *FileNode) LoadEntries() error {
	childrenCopy := fn.children
	fn.children = make(map[string]*FileNode)
	entries, err := os.ReadDir(fn.AbsPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}
		path := filepath.Join(fn.Path, info.Name())
		absPath := filepath.Join(fn.AbsPath, info.Name())
		if child, ok := childrenCopy[path]; ok {
			// preserves the child.children map
			err = fn.AddChild(child)
			if err != nil {
				return err
			}
			continue
		}
		err = fn.AddChild(&FileNode{
			Info:    info,
			Path:    path,
			AbsPath: absPath,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (fn *FileNode) UpdateChildren() error {
	for _, child := range fn.children {
		delete(fn.Tree.FileNodeMap, child.Path)
	}
	return fn.LoadEntries()
}

func (fn *FileNode) IterateOverChildren(f func(*FileNode)) {
	for _, child := range fn.children {
		f(child)
	}
}

type FileTree struct {
	Root        *FileNode
	FileNodeMap map[string]*FileNode
}

func NewFileTree(root *FileNode) *FileTree {
	tree := &FileTree{
		Root:        root,
		FileNodeMap: make(map[string]*FileNode),
	}
	tree.RegisterNode(root)
	return tree
}

func (ft *FileTree) RegisterNode(fn *FileNode) {
	fn.Tree = ft
	ft.FileNodeMap[fn.Path] = fn
}
