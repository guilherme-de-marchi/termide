package main

import (
	"log"

	"github.com/Guilherme-De-Marchi/termide/internal"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GetFileTreeView(path string) (*tview.TreeView, error) {
	root, err := internal.NewFileNodeFrom(path)
	if err != nil {
		return nil, err
	}
	tree := internal.NewFileTree(root)

	rootView := tview.NewTreeNode(root.Info.Name()).
		SetReference(root.Path).
		SetExpanded(false).
		SetColor(tcell.ColorGreen)

	treeView := tview.NewTreeView().
		SetRoot(rootView).
		SetCurrentNode(rootView)

	// App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 	switch event.Key() {
	// 	case tcell.KeyF5:
	// 		tree.UpdateOpenedNodes()
	// 	}
	// 	return event
	// })

	treeView.SetSelectedFunc(func(nodeView *tview.TreeNode) {
		ref, ok := nodeView.GetReference().(string)
		if !ok {
			log.Println("could not type assert TreeNode reference to string")
			return
		}
		fn, ok := tree.FileNodeMap[ref]
		if !ok {
			log.Printf("file node %s not registered on the file tree\n", ref)
			return
		}

		if fn.Info.IsDir() {
			nodeView.SetExpanded(!nodeView.IsExpanded())

			if nodeView.IsExpanded() {
				err = fn.UpdateChildren()
				if err != nil {
					log.Println("could not update children:", err)
					return
				}
				nodeView.ClearChildren()
				fn.IterateOverChildren(func(child *internal.FileNode) {
					childView := tview.NewTreeNode(child.Info.Name()).
						SetReference(child.Path)
					if child.Info.IsDir() {
						childView.SetExpanded(false).SetColor(tcell.ColorGreen)
					}
					nodeView.AddChild(childView)
				})
			}
		}
	})

	return treeView, nil
}
