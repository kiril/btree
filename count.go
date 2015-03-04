package btree

func (tree *Btree) count() (int, error) {
    var value int
    tnode, err := tree.getTreeNode(tree.GetRoot())
    if err != nil {
        return value, err
    }
    return tnode.count(tree)
}

func (node *TreeNode) count(tree *Btree) (int, error) {
    count := 0

    if node.GetNodeType() == isNode {
        for i := 0; i < len(node.Childrens); i++ {
            subTree, error := tree.getTreeNode(node.Childrens[i])
            if error != nil {
                return count, error
            }

            subCount, countError := subTree.count(tree)
            if countError != nil {
                return count, countError
            }

            count += subCount
        }

        return count, nil

    } else {
        return len(node.Keys), nil
    }
}
