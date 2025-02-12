package huffmancoding

import (
	"container/heap"
	"fmt"
)

type byteFreqMap map[byte]int

func buildByteFreqMap(bs []byte) byteFreqMap {
	bfm := make(byteFreqMap)

	for _, b := range bs {
		bfm[b]++
	}

	return bfm
}

type Node struct {
	Byte  byte
	Freq  int
	Left  *Node
	Right *Node
}

type priorityQueue []*Node

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].Freq < pq[j].Freq ||
		(pq[i].Freq == pq[j].Freq && pq[i].Byte < pq[j].Byte)
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Node))
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)

	item := old[n-1]
	*pq = old[:n-1]

	return item
}

func buildPriorityQueue(bfm byteFreqMap) priorityQueue {
	pq := make(priorityQueue, len(bfm))
	i := 0

	for b, f := range bfm {
		pq[i] = &Node{
			Byte: b,
			Freq: f,
		}
		i++
	}

	return pq
}

func buildTreeNode(pq priorityQueue) *Node {
	heap.Init(&pq)

	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*Node)
		right := heap.Pop(&pq).(*Node)

		merged := &Node{
			Byte:  '$',
			Freq:  left.Freq + right.Freq,
			Left:  left,
			Right: right,
		}

		heap.Push(&pq, merged)
	}

	return heap.Pop(&pq).(*Node)
}

func getBitFromTreeNode(node *Node, c byte, r *uint64, l *uint) bool {
	if node.Left == nil && node.Right == nil {
		if node.Byte == c {
			return true
		}
	}

	if node.Left != nil {
		rc := *r
		rc <<= 1
		lc := *l
		lc++

		if getBitFromTreeNode(node.Left, c, &rc, &lc) {
			*r = rc
			*l = lc
			return true
		}
	}

	if node.Right != nil {
		rc := *r
		rc <<= 1
		rc |= 0x1
		lc := *l
		lc++

		if getBitFromTreeNode(node.Right, c, &rc, &lc) {
			*r = rc
			*l = lc
			return true
		}
	}
	return false
}

func Compress(bs []byte) ([]byte, uint, *Node, error) {
	var compressed = make([]byte, 0)
	bfm := buildByteFreqMap(bs)
	pq := buildPriorityQueue(bfm)
	tn := buildTreeNode(pq)

	bitCount := 0
	var totalBitCount uint = 0
	var currentByte byte = 0

	for k, b := range bs {
		var r uint64 = 0
		var l uint = 0

		if !getBitFromTreeNode(tn, b, &r, &l) {
			return []byte{}, 0, tn, fmt.Errorf("cannot get compressed bit of %b", b)
		}

		if l == 0 {
			return []byte{}, 0, tn, fmt.Errorf("cannot get compressed bit of %b", b)
		}

		r <<= 64 - l

		for i := 0; i < int(l); i++ {
			bitCount++
			totalBitCount++

			currentByte <<= 1
			currentByte |= uint8((r & 0x8000000000000000) >> 63)
			r <<= 1

			if bitCount == 8 {
				compressed = append(compressed, currentByte)
				currentByte = 0
				bitCount = 0
			}
		}

		if k == len(bs)-1 && bitCount > 0 {
			currentByte <<= (8 - bitCount)
			compressed = append(compressed, currentByte)
		}
	}

	return compressed, totalBitCount, tn, nil
}

func findOriginalByte(b uint64, node *Node, l uint, n uint) (byte, bool) {
	if node.Left == nil && node.Right == nil {
		return node.Byte, true
	}

	if n < l {
		r := b & 0x8000000000000000
		r >>= 63
		n++
		b <<= 1

		if r == 1 {
			if node.Right != nil {
				return findOriginalByte(b, node.Right, l, n)
			}
		} else {
			if node.Left != nil {
				return findOriginalByte(b, node.Left, l, n)
			}
		}
	}

	return 0, false
}

func shiftLeftBytes(bb []byte, n uint) []byte {
	if n == 0 || len(bb) == 0 {
		return bb
	}

	byteShift := int(n / 8)
	bitShift := n % 8

	if byteShift >= len(bb) {
		return make([]byte, len(bb))
	}

	copy(bb, append(bb[byteShift:], make([]byte, byteShift)...))

	if bitShift > 0 {
		var prev byte
		for i := len(bb) - 1; i >= 0; i-- {
			temp := bb[i] >> (8 - bitShift)
			bb[i] <<= bitShift
			bb[i] |= prev
			prev = temp
		}
	}

	return bb
}

func getNBitFromRight(bs []byte, n uint) (uint64, error) {
	var result uint64
	var tl uint = 0

	if n > uint(len(bs)*8) {
		return 0, fmt.Errorf("n cannot bigger than bytes length. got: %d", n)
	}

	if n > 64 {
		return 0, fmt.Errorf("n cannot bigger than 64. got: %d", n)
	}

	var loopC uint = 0
	for _, b := range bs {
		if tl >= n {
			break
		}

		var d uint = 8
		if tl+8 > n {
			b &= 0xFF << (8 - (n - tl))
			d = 8 - ((tl + 8) - n)
		}

		result <<= 8
		result |= uint64(b)
		tl += d
		loopC++
	}

	result <<= 64 - (8 * loopC)
	return result, nil
}

func Decompress(bs []byte, length uint, tNode *Node) ([]byte, error) {
	var original = make([]byte, 0)

	var i uint
	for i = 0; i < length; i++ {
		var j uint = 1
		for j <= length-i {
			r, err := getNBitFromRight(bs, j)
			if err != nil {
				return []byte{}, err
			}

			ob, ok := findOriginalByte(r, tNode, j, 0)
			if ok {
				original = append(original, ob)
				bs = shiftLeftBytes(bs, j)
				i += j - 1
				break
			}

			j++
		}
	}

	return original, nil
}
