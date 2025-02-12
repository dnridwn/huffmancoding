package huffmancoding

import (
	"fmt"
	"sort"
	"testing"
)

func TestBuildByteFreqMap(t *testing.T) {
	testCases := []struct {
		bs       []byte
		expected byteFreqMap
	}{
		{
			bs: []byte{0xFF, 0x80, 0x01},
			expected: byteFreqMap{
				0xFF: 1,
				0x80: 1,
				0x01: 1,
			},
		},
		{
			bs: []byte{0xFF, 0xFF, 0x01},
			expected: byteFreqMap{
				0xFF: 2,
				0x01: 1,
			},
		},
	}

	for k, tc := range testCases {
		t.Run(fmt.Sprintf("Test Build Freq Map #%d", (k+1)), func(t *testing.T) {
			bfm := buildByteFreqMap(tc.bs)
			for ek, e := range tc.expected {
				f, ok := bfm[ek]
				if !ok {
					t.Errorf("key %b does not exists in byte freq map", ek)
				} else if e != f {
					t.Errorf("frequent does not match. expected %d, got %d", e, f)
				}
			}
		})
	}
}

func TestBuildPriorityQueue(t *testing.T) {
	testCases := []struct {
		bfm      byteFreqMap
		expected priorityQueue
	}{
		{
			bfm: byteFreqMap{
				0xFF: 1,
				0x80: 2,
				0x01: 3,
			},
			expected: priorityQueue{
				&Node{Byte: 0xFF, Freq: 1},
				&Node{Byte: 0x80, Freq: 2},
				&Node{Byte: 0x01, Freq: 3},
			},
		},
		{
			bfm: byteFreqMap{
				0x01: 1,
				0xFF: 2,
			},
			expected: priorityQueue{
				&Node{Byte: 0x01, Freq: 1},
				&Node{Byte: 0xFF, Freq: 2},
			},
		},
	}

	for k, tc := range testCases {
		t.Run(fmt.Sprintf("Test Build Priority Queue #%d", k+1), func(t *testing.T) {
			pq := buildPriorityQueue(tc.bfm)
			sort.Sort(pq)
			for i := 0; i < len(pq); i++ {
				if pq[i].Byte != tc.expected[i].Byte {
					t.Errorf("byte not match. expected %b, got %b", tc.expected[i].Byte, pq[i].Byte)
				} else if pq[i].Freq != tc.expected[i].Freq {
					t.Errorf("freq not match. expected %d, got %d", tc.expected[i].Freq, pq[i].Freq)
				}
			}
		})
	}
}

func isNodesMatch(n1, n2 *Node) error {
	if n1 == nil && n2 != nil {
		return fmt.Errorf("first node is nil")
	}

	if n1 != nil && n2 == nil {
		return fmt.Errorf("second node is nil")
	}

	if n1.Byte != n2.Byte {
		return fmt.Errorf("byte not match. n1: %b, n2: %b", n1.Byte, n2.Byte)
	}

	if n1.Freq != n2.Freq {
		return fmt.Errorf("freq not match. n1: %d, n2: %d", n1.Freq, n2.Freq)
	}

	if n1.Left != nil || n2.Left != nil {
		if err := isNodesMatch(n1.Left, n2.Left); err != nil {
			return err
		}
	}

	if n1.Right != nil || n2.Right != nil {
		if err := isNodesMatch(n1.Right, n2.Right); err != nil {
			return err
		}
	}

	return nil
}

func TestBuildTreeNode(t *testing.T) {
	testCases := []struct {
		pq       priorityQueue
		expected *Node
	}{
		{
			pq: priorityQueue{
				&Node{Byte: 'a', Freq: 1},
				&Node{Byte: 'b', Freq: 1},
				&Node{Byte: 'c', Freq: 1},
			},
			expected: &Node{
				Byte: '$',
				Freq: 3,
				Left: &Node{
					Byte: 'c',
					Freq: 1,
				},
				Right: &Node{
					Byte: '$',
					Freq: 2,
					Left: &Node{
						Byte: 'a',
						Freq: 1,
					},
					Right: &Node{
						Byte: 'b',
						Freq: 1,
					},
				},
			},
		},
	}

	for k, tc := range testCases {
		t.Run(fmt.Sprintf("Test Build Tree Node #%d", k+1), func(t *testing.T) {
			if err := isNodesMatch(tc.expected, buildTreeNode(tc.pq)); err != nil {
				t.Errorf("nodes do not match: %v", err)
			}
		})
	}
}

func TestGetBitFromTreeNode(t *testing.T) {
	testCases := []struct {
		tNode     *Node
		c         byte
		r         *uint64
		l         *uint
		expectedR *uint64
		expectedL *uint
	}{
		{
			tNode: &Node{
				Byte: '$',
				Freq: 3,
				Left: &Node{
					Byte: 'b',
					Freq: 1,
				},
				Right: &Node{
					Byte: '$',
					Freq: 2,
					Left: &Node{
						Byte: 'a',
						Freq: 1,
					},
					Right: &Node{
						Byte: 'c',
						Freq: 1,
					},
				},
			},
			c: 'a',
			r: func() *uint64 {
				var b uint64 = 0
				return &b
			}(),
			l: func() *uint {
				var l uint = 0
				return &l
			}(),
			expectedR: func() *uint64 {
				var b uint64 = 0b00000010
				return &b
			}(),
			expectedL: func() *uint {
				var l uint = 2
				return &l
			}(),
		},
	}

	for k, tc := range testCases {
		t.Run(fmt.Sprintf("Test Get Bit From Tree Node #%d", k+1), func(t *testing.T) {
			b := getBitFromTreeNode(tc.tNode, tc.c, tc.r, tc.l)
			if !b {
				t.Errorf("byte %b not found in tree node", tc.c)
			} else {
				if *tc.r != *tc.expectedR {
					t.Errorf("result not match. expected: %b, got: %b", *tc.expectedR, *tc.r)
				} else if *tc.l != *tc.expectedL {
					t.Errorf("length not match. expected: %d, got: %d", *tc.expectedL, *tc.l)
				}
			}
		})
	}
}

func TestCompress(t *testing.T) {
	testCases := []struct {
		bs       []byte
		expected []byte
	}{
		{
			bs:       []byte("lossless"),
			expected: []byte{0b01001110, 0b10001100},
		},
	}

	for k, tc := range testCases {
		t.Run(fmt.Sprintf("Test Compress #%d", k+1), func(t *testing.T) {
			compressed, _, _, err := Compress(tc.bs)
			if err != nil {
				t.Errorf("Error occurred: %v", err)
			} else {
				if len(compressed) != len(tc.expected) {
					t.Errorf("bytes length not match. expected: %d, got: %d", len(tc.expected), len(compressed))
				} else {
					for bk, b := range tc.expected {
						if b != compressed[bk] {
							t.Errorf("byte not match. expected: %08b, got: %08b", b, compressed[bk])
						}
					}
				}
			}
		})
	}
}

func TestCompressDecompress(t *testing.T) {
	testCases := []struct {
		text []byte
	}{
		{
			text: []byte("Golang is great"),
		},
		{
			text: []byte("H0r@r3Y0u"),
		},
		{
			text: []byte("Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."),
		},
		{
			text: []byte("The Go Awakening. Joko was a seasoned developer drowning in a sea of complex, bloated codebases. One day, he discovered Go—simple, elegant, and fast. Intrigued, he rewrote an old project in Go. The code was lean, the performance skyrocketed, and debugging felt like a breeze. Concurrency? Goroutines made it effortless. Soon, Joko became the \"Go Guru\" at work. His peers, once skeptical, embraced the language. With Go, coding felt fresh again. Joko smiled—he had found his perfect tool."),
		},
	}

	for k, tc := range testCases {
		t.Run(fmt.Sprintf("Test Compress Decompress #%d", k+1), func(t *testing.T) {
			compressed, length, tNode, err := Compress(tc.text)
			if err != nil {
				t.Errorf("error occurred when compressing: %v", err)
			} else {
				decompressed, err := Decompress(compressed, length, tNode)
				if err != nil {
					t.Errorf("error occurred when decompressing: %v", err)
				} else {
					if string(decompressed) != string(tc.text) {
						t.Errorf("result not match. expected: %08b, got: %08b", tc.text, decompressed)
					}
				}
			}
		})
	}
}
