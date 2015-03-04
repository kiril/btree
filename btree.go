package btree

import (
	"github.com/golang/protobuf/proto"
	"time"
)

// Btree metadata
type Btree struct {
	BtreeMetadata
	gcIndex     int64
	dupnodelist map[int64]int
	opChan      chan *treeOperation
	exitChan    chan int
}

const (
	// TreeSize is  tree size
	TreeSize = 1 << 10
	// LeafSize is leaf size
	LeafSize = 1 << 5
	// NodeSize is node size
	NodeSize = 1 << 6
)

// isNode, isLeaf is treenode tag
const (
	isNode = iota
	isLeaf
)

// NewBtree create a btree
func NewBtree() *Btree {
	return NewBtreeSize(LeafSize, NodeSize)
}

// NewBtreeSize create new btree with custom leafsize/nodesize
func NewBtreeSize(leafsize int64, nodesize int64) *Btree {
	tree := &Btree{
		dupnodelist: make(map[int64]int),
		opChan:      make(chan *treeOperation),
		BtreeMetadata: BtreeMetadata{
			Root:        proto.Int64(0),
			Size:        proto.Int64(TreeSize),
			LeafMax:     proto.Int64(leafsize),
			NodeMax:     proto.Int64(nodesize),
			IndexCursor: proto.Int64(0),
			Nodes:       make([][]byte, TreeSize),
		},
	}
	go tree.run()
	return tree
}

func (t *Btree) run() {
	tick := time.Tick(time.Second * 2)
	for {
		select {
		case <-t.exitChan:
			break
		case op := <-t.opChan:
			switch op.GetAction() {
			case "insert":
				op.errChan <- t.insert(op.TreeLog)
			case "delete":
				op.errChan <- t.dodelete(op.Key)
			case "update":
				op.errChan <- t.update(op.TreeLog)
			case "search":
				rst, err := t.search(op.Key)
				op.valueChan <- rst
				op.errChan <- err
            case "left":
                rst, err := t.left()
                op.valueChan <- rst
                op.errChan <- err
            case "count":
                count, err := t.count()
                op.valueChan <- []byte{byte(count)}
                op.errChan <- err
			}
			t.Index = proto.Int64(t.GetIndexCursor())
		case <-tick:
			t.gc()
		}
	}
	t.Marshal("treedump.tmp")
}

func (t *Btree) Sync(file string) {
	//t.Marshal(file)
}

// Insert can insert record into a btree
func (t *Btree) Insert(key, value []byte) error {
	q := &treeOperation{
		valueChan: make(chan []byte),
		errChan:   make(chan error),
	}
	q.Action = proto.String("insert")
	q.Key = key
	q.Value = value
	t.opChan <- q
	return <-q.errChan
}

// Delete can delete record
func (t *Btree) Delete(key []byte) error {
	q := &treeOperation{
		valueChan: make(chan []byte),
		errChan:   make(chan error),
	}
	q.Action = proto.String("delete")
	q.Key = key
	t.opChan <- q
	return <-q.errChan
}

// Search return value
func (t *Btree) Search(key []byte) ([]byte, error) {
	q := &treeOperation{
		valueChan: make(chan []byte),
		errChan:   make(chan error),
	}
	q.Action = proto.String("search")
	q.Key = key
	t.opChan <- q
	return <-q.valueChan, <-q.errChan
}

// Update is used to update key/value
func (t *Btree) Update(key, value []byte) error {
	q := &treeOperation{
		valueChan: make(chan []byte),
		errChan:   make(chan error),
	}
	q.Action = proto.String("update")
	q.Key = key
	q.Value = value
	t.opChan <- q
	return <-q.errChan
}

// get the left-most value, just for shits and giggles
func (t *Btree) Left() ([]byte, error) {
    q := &treeOperation {
        valueChan: make(chan []byte),
        errChan: make(chan error),
    }
    q.Action = proto.String("left")
    t.opChan <- q
    return <-q.valueChan, <-q.errChan
}

func (t *Btree) Count() (int, error) {
    q := &treeOperation {
        valueChan: make(chan []byte),
        errChan: make(chan error),
    }
    q.Action = proto.String("count")
    t.opChan <- q
    return int((<-q.valueChan)[0]), <-q.errChan
}
