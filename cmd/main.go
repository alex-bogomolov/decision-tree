package main

import (
	"fmt"
	dt "github.com/alex-bogomolov/decision_tree"
	_ "os"
	_ "encoding/csv"
	_ "io"
	"strconv"
)

func main() {
	accuracy()
}

func accuracy() {
	dataset, labels := dt.LoadTitanic()

	tree := dt.BuildTree(dataset, 3, 10)
	tree.SetLabels(labels)
	dt.VisualizeTree(tree)

	total := len(dataset)
	correct := 0

	for _, example := range dataset {
		last := len(example)-1
		input := example[:last]
		correctClass := example[last]
		prediction := dt.Predict(tree, input)
		if prediction == correctClass {
			correct++
		}
	}

	accuracy := float64(correct) / float64(total)

	fmt.Printf("Accuracy: %.2f\n", accuracy)
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


