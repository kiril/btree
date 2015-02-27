package btree

import "fmt"

func (tree *Btree) left() ([]byte, error) {
    var value []byte
    tnode, err := tree.getTreeNode(tree.GetRoot())
    if err != nil {
        return value, err
    }
    return tnode.left(tree)
}

func (node *TreeNode) left(tree *Btree) ([]byte, error) {
    var value []byte
    if node.GetNodeType() == isNode {
        nextNode, error := tree.getTreeNode(node.Childrens[0])
        if error != nil {
            return value, error
        }

        return nextNode.left(tree)

    } else {
        if len(node.Keys) == 0 {
            return nil, fmt.Errorf("empty")
        } else {
            return node.Values[0], nil
        }
    }
}
