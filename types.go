package decision_tree

import (
	"math"
)

type Matrix []Vector

const similarity = 1e2

func (m Matrix) Equals(other Matrix) bool {
	mShape1, mShape2 := m.Shape()
	otherShape1, otherShape2 := other.Shape()

	if mShape1 != otherShape1 || mShape2 != otherShape2 {
		return false
	}

	for i, v := range m {
		if !v.Equals(other[i]) {
			return false
		}
	}

	return true
}

func (m Matrix) Shape() (int, int) {
	if len(m) == 0 {
		return 0, 0
	}
	return len(m), len(m[0])
}

func (m Matrix) GetCol(index int) Vector {
	var out Vector

	for _, row := range m {
		out = append(out, row[index])
	}

	return out
}

type Vector []float64

func (v Vector) Equals(other Vector) bool {
	if len(v) != len(other) {
		return false
	}

	for i, e := range v {
		if math.Abs(e-other[i]) > similarity {
			return false
		}
	}

	return true
}

func (v Vector) Unique() Vector {
	found := make(map[float64]bool)

	var out Vector

	for _, el := range v {
		if _, ok := found[el]; !ok {
			out = append(out, el)
			found[el] = true
		}
	}

	return out
}

type Node struct {
	index         int
	value         float64
	groups        []Matrix
	terminal      bool
	terminalValue float64
	left          *Node
	right         *Node
}

type String string

func (s String) Multiply(a int) String {
	var out String = ""

	for i := 0; i < a; i++ {
		out += s
	}

	return out
}