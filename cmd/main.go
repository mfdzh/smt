package main

import (
	// "bytes"
	psmt "github.com/celestiaorg/smt"
	"crypto/sha256"
	"fmt"
)

type Array struct {
	data [1000]uint64
	index int
}

func (a *Array) Push(tmp uint64) {
	a.data[a.index] = tmp
	a.index = (a.index+1)%1000
}

func (a *Array) Sum() uint64 {
	var sum uint64 = 0
	for i := 0; i < 1000; i++ {
		sum += a.data[i]
	}
	return sum
}

func addOne(tmp *[4]byte) {
	for i, _ := range tmp {
		tmp[i] += 1
		if tmp[i] != 0 {
			break
		}
	}
}

func main() {
	// smn, err := psmt.NewBadgerStore("/mnt/2004/kv-tmp/nodes", nil)
	// if err != nil {
	// 	fmt.Printf("%+v\n", err)
	// 	return
	// }
	// smv, err := psmt.NewBadgerStore("/mnt/2004/kv-tmp/values", nil)
	// if err != nil {
	// 	fmt.Printf("%+v\n", err)
	// 	return
	// }
	smn, smv := psmt.NewSimpleMap(), psmt.NewSimpleMap()
	smt := psmt.NewSparseMerkleTree(smn, smv, sha256.New())
	var tmp [4]byte
	var newRoot []byte
	var oldRoot []byte
	var maxDepth uint64
	var depth uint64
	var depthArr *Array = &Array{
		data: [1000]uint64{},
		index: 0,
	}
	var rdepthArr *Array = &Array{
		data: [1000]uint64{},
		index: 0,
	}
	var totalNode uint64
	var nodeSize uint64
	var tmpSize uint64 = 100
	var err error
	newRoot, totalNode, err = smt.Update(tmp[:], []byte(fmt.Sprintf("%d", 0)))
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	nodeSize = totalNode
	for i := uint64(2); true; i++ {
		addOne(&tmp)

		oldRoot = newRoot
		newRoot, depth, err = smt.Update(tmp[:], []byte(fmt.Sprintf("%d", i)))
		if err != nil {
			fmt.Printf("%+v\n", err)
			return
		}
		if depth > maxDepth {
			maxDepth = depth
		}

		depthArr.Push(depth)
		totalNode = totalNode + depth
		nodeSize = nodeSize + depth

		depth, err = smt.RemovePathForRoot(tmp[:], oldRoot)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return
		}
		rdepthArr.Push(depth)
		nodeSize = nodeSize - depth

		if i == 10*tmpSize || i % 200000000 == 0 {
			fmt.Println(i)
			fmt.Println("Node Size: ", smt.NodeSize())
			fmt.Println("maxDepth: ", maxDepth)
			fmt.Println("Node Size(all removed): ", nodeSize)
			fmt.Println("Node Size(remain 1000 key): ", nodeSize + rdepthArr.Sum())
			fmt.Println("Node Size(remain all key): ", totalNode)
			fmt.Printf("avage Write:%.4f\n", float64(depthArr.Sum())/float64(1000))
			fmt.Println("Value Size: ", smt.ValueSize())
			tmpSize = i
		}
	}
}