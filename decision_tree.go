package decision_tree

import (
	"fmt"
	"os"
	"encoding/csv"
	"io"
	"strconv"
	"encoding/json"
	"html/template"
)

// Calculate the Gini index for a split dataset
func GiniIndex(groups []Matrix, classes Vector) float64 {
	// count all samples at split point
	nInstances := float64(numInstances(groups))
	// sum weighted Gini index for each group
	gini := 0.0

	for _, group := range groups {
		size := float64(len(group))
		// avoid divide by zero
		if size == 0 {
			continue
		}
		score := 0.0
		// score the group based on the score for each class

		for _, classVal := range classes {
			p := 0.0

			for _, row := range group {
				rowClass := row[len(row) - 1]
				if rowClass == classVal {
					p += 1
				}
			}

			p /= size

			score += p * p
		}

		// weight the group score by its relative size
		gini += (1.0 - score) * (size / nInstances)
	}

	return gini
}

// Split a dataset based on an attribute and an attribute value
func TestSplit(index int, value float64, dataset Matrix) (Matrix, Matrix) {
	left, right := Matrix{}, Matrix{}

	for _, row := range dataset {
		if row[index] < value {
			left = append(left, row)
		} else {
			right = append(right, row)
		}
	}

	return left, right
}

// Select the best split point for a dataset
func GetSplit(dataset Matrix) *Node {
	_, datasetCols := dataset.Shape()
	classValues := dataset.GetCol(datasetCols - 1).Unique()

	var bIndex int
	var bValue, bScore float64
	var bGroups []Matrix

	bScore = 100

	for i := 0; i < datasetCols - 1; i++ {
		for _, row := range dataset {
			leftGroup, rightGroup := TestSplit(i, row[i], dataset)
			gini := GiniIndex([]Matrix{leftGroup, rightGroup}, classValues)

			if gini < bScore {
				bIndex, bValue, bScore, bGroups = i, row[i], gini, []Matrix{leftGroup, rightGroup}
			}
		}
	}

	out := Node{
		index:bIndex,
		value:bValue,
		groups:bGroups,
	}

	return &out
}

func ToTerminal(group Matrix) float64 {
	_, cols := group.Shape()
	outcomes := group.GetCol(cols - 1)

	frequencies := make(map[float64] int)

	for _, outcome := range outcomes {
		if _, ok := frequencies[outcome]; ok {
			frequencies[outcome]++
		} else {
			frequencies[outcome] = 1;
		}
	}

	counter := 0
	value := -1.0

	for outcome, frequency := range frequencies {
		if frequency > counter {
			counter = frequency
			value = outcome
		}
	}

	return value
}

func Split(node *Node, maxDepth, minSize, depth int) {
	left, right := node.groups[0], node.groups[1]
	node.groups = nil

	if len(left) == 0 || len(right) == 0 {
		t := ToTerminal(append(left, right...))
		n := Node{
			terminal:true,
			terminalValue:t,
		}
		node.left = &n
		node.right = &n
		return
	}

	if depth >= maxDepth {
		node.left = &Node{
			terminal:true,
			terminalValue:ToTerminal(left),
		}

		node.right = &Node{
			terminal:true,
			terminalValue:ToTerminal(right),
		}
		return
	}
	if len(left) <= minSize {
		node.left = &Node{
			terminal:true,
			terminalValue:ToTerminal(left),
		}
	} else {
		node.left = GetSplit(left)
		Split(node.left, maxDepth, minSize, depth+1)
	}

	if len(right) <= minSize {
		node.right = &Node {
			terminal:true,
			terminalValue:ToTerminal(right),
		}
	} else {
		node.right = GetSplit(right)
		Split(node.right, maxDepth, minSize, depth+1)
	}
}

func BuildTree(train Matrix, maxDepth int, minSize int) *Node {
	root := GetSplit(train)
	Split(root, maxDepth, minSize, 1)
	return root
}

func PrintTree(node *Node, depth int) {
	if node.terminal {
		fmt.Printf("%s[%.3f]\n", String(" ").Multiply(depth), node.terminalValue)
	} else {
		fmt.Printf("%s[X%d < %.3f]\n", String(" ").Multiply(depth), node.index + 1, node.value)
		PrintTree(node.left, depth + 1)
		PrintTree(node.right, depth + 1)
	}
}

func (n Node) ToNodeJSON() *NodeJSON {
	return nodeJSON(&n, n.labels)
}

func nodeJSON(node *Node, labels []string) *NodeJSON {
	if node == nil {
		return nil
	}

	if node.terminal {
		n := &NodeJSON{
			Text:NodeText{
				Name:fmt.Sprintf("%.0f\n", node.terminalValue),
			},
		}
		return n
	}

	n := &NodeJSON{
		Text:NodeText{
			Name:fmt.Sprintf("%s < %.3f\n", labels[node.index], node.value),
		},
	}

	n.Children = make([]*NodeJSON, 2, 2)
	n.Children[0] = nodeJSON(node.left, labels)
	n.Children[1] = nodeJSON(node.right, labels)
	return n
}

func Predict(node *Node, row Vector) float64  {
	return predict(node, row).terminalValue
}

func predict(node *Node, row Vector) *Node {
	if node.terminal {
		return node
	} else {
		if row[node.index] < node.value {
			return predict(node.left, row)
		} else {
			return predict(node.right, row)
		}
	}
}

func LoadIris() (Matrix, []string) {
	csvFile, err := os.Open("iris.csv")
	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(csvFile)

	var out Matrix
	var labels []string

	counter := 0

	for {
		record, err := csvReader.Read()
		if counter == 0 {
			labels = record[:len(record) - 1]
			counter++
			continue
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		out = append(out, stringSliceToFloat(record))
	}

	return out, labels
}

func LoadTitanic() (Matrix, []string) {
	csvFile, err := os.Open("titanic.csv")
	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(csvFile)

	var out Matrix
	var labels []string

	counter := 0

	for {
		record, err := csvReader.Read()
		if counter == 0 {
			counter++
			labels = record[1:len(record)-1]
			continue
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		out = append(out, stringSliceToFloat(record[1:]))
	}

	return out, labels
}

func VisualizeTree(node *Node) {
	root := nodeJSON(node, node.labels)

	b, _ := json.Marshal(root)

	t, err := template.ParseFiles("template/index.html")

	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile("tree.html", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0600)

	if err != nil {
		panic(err)
	}

	t.Lookup("index.html").Execute(f, template.JS(b))
}


func numInstances(groups []Matrix) int {
	out := 0

	for _, group := range groups {
		out += len(group)
	}

	return out
}

func stringSliceToFloat(in []string) []float64 {
	out := make([]float64, len(in), len(in))

	for i, v := range in {
		f, err := strconv.ParseFloat(v, 64)

		if err != nil {
			panic(err)
		}

		out[i] = f
	}

	return out
}
