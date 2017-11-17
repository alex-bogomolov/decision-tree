package decision_tree

import (
	"testing"
	"math"
)

func TestGiniIndex(t *testing.T) {
	inMatrix := []Matrix{
		{
			{1, 1},
			{1, 0},
		},
		{
			{1, 1},
			{1, 0},
		},
	}

	inVector := Vector{0, 1}

	expected := 0.5
	actual := GiniIndex(inMatrix, inVector)

	if expected != actual {
		t.Errorf("Actual: %v, expected: %v", actual, expected)
	}

	inMatrix = []Matrix{
		{
			{1, 0},
			{1, 0},
		},
		{
			{1, 1},
			{1, 1},
		},
	}

	expected = 0
	actual = GiniIndex(inMatrix, inVector)

	if expected != actual {
		t.Errorf("Actual: %v, expected: %v", actual, expected)
	}
}

func TestTestSplit(t *testing.T) {
	inMatrix := Matrix{
		{0, 28, 37, 59},
		{99, 98, 42, 59},
		{73, 64, 13, 44},
		{38, 97, 71, 4},
		{17, 46, 70, 68},
		{77, 91, 85, 12},
		{28, 6, 42, 83},
		{74, 93, 23, 83},
		{96, 76, 3, 68},
		{51, 37, 54, 7},
	}
	inIndex := 2
	inValue := 50.0

	expected1 := Matrix{
		{0, 28, 37, 59},
		{99, 98, 42, 59},
		{73, 64, 13, 44},
		{28, 6, 42, 83},
		{74, 93, 23, 83},
		{96, 76, 3, 68},
	}

	expected2 := Matrix{
		{38, 97, 71, 4},
		{17, 46, 70, 68},
		{77, 91, 85, 12},
		{51, 37, 54, 7},
	}

	actual1, actual2 := TestSplit(inIndex, inValue, inMatrix)

	if !expected1.Equals(actual1) {
		t.Errorf("Actual: %v, expected: %v", actual1, expected1)
	}

	if !expected2.Equals(actual2) {
		t.Errorf("Actual: %v, expected: %v", actual2, expected2)
	}
}

func TestVector_Unique(t *testing.T) {
	v := Vector{1, 1, 2, 2, 2, 3, 3, 3}
	expected := Vector{1, 2, 3}
	actual := v.Unique()

	if !expected.Equals(actual) {
		t.Errorf("Actual: %v, expected: %v", actual, expected)
	}
}

func TestGetSplit(t *testing.T) {
	root := GetSplit(toyDataset())
	root.groups[0] = root.groups[0][:3]
	root.groups[1] = root.groups[1][:3]

	expectedIndex := 0
	expectedValue := 6.64

	expectedGroups := []Matrix{
		{
			{ 2.77, 1.78, 0 },
			{ 1.73, 1.7, 0 },
			{ 3.68, 2.81, 0 },
		},
		{
			{ 7.5, 3.16, 1 },
			{ 9, 3.34, 1 },
			{ 7.44, 0.48, 1 },
		},
	}

	if expectedIndex != root.index {
		t.Errorf("Actual: %v, expected: %v", root.index, expectedIndex)
	}

	if math.Abs(expectedValue - root.value) > similarity {
		t.Errorf("Actual: %v, expected: %v", root.value, expectedValue)
	}

	if !expectedGroups[0].Equals(root.groups[0]) {
		t.Errorf("Actual: %v, expected: %v", root.groups[0], expectedGroups[0])
	}

	if !expectedGroups[1].Equals(root.groups[1]) {
		t.Errorf("Actual: %v, expected: %v", root.groups[1], expectedGroups[1])
	}
}

func TestToTerminal(t *testing.T) {
	in := Matrix{
		{1, 1},
		{2, 1},
		{3, 0},
		{4, 1},
	}

	expected := 1.0
	actual := ToTerminal(in)

	if actual != expected {
		t.Errorf("Actual: %v, expected: %v", actual, expected)
	}
}

func TestBuildTree(t *testing.T) {
	d := LoadIris()
	tree := BuildTree(d, 3, 10)
	PrintTree(tree, 0)
}

func TestPredict(t *testing.T) {
	tree := BuildTree(LoadIris(), 100, 1)
	in := Vector{1.4,0.2,5.1,3.5}

	expected := 0.0
	actual := Predict(tree, in)

	if expected != actual {
		t.Errorf("Expected: %v, actual: %v", expected, actual)
	}
}

func ExamplePredict() {

}

func BenchmarkBuildTree(b *testing.B) {
	d := LoadIris()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		BuildTree(d, 100, 1)
	}
}

func toyDataset() Matrix {
	return Matrix{
		{2.771244718, 1.784783929, 0},
		{1.728571309, 1.169761413, 0},
		{3.678319846, 2.81281357, 0},
		{3.961043357, 2.61995032, 0},
		{2.999208922, 2.209014212, 0},
		{7.497545867, 3.162953546, 1},
		{9.00220326, 3.339047188, 1},
		{7.444542326, 0.476683375, 1},
		{10.12493903, 3.234550982, 1},
		{6.642287351, 3.319983761, 1},
	}
}



