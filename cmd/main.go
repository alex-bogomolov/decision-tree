package main

import (
	"fmt"
	"github.com/alex-bogomolov/decision_tree"
	"os"
	"encoding/csv"
	"io"
	"strconv"
)

func main() {
	accuracy()
}

func accuracy() {
	dataset := decision_tree.LoadTitanic()

	tree := decision_tree.BuildTree(dataset, 300, 1)
	decision_tree.PrintTree(tree, 1)

	total := len(dataset)
	correct := 0

	for _, example := range dataset {
		last := len(example)-1
		input := example[:last]
		correctClass := example[last]
		prediction := decision_tree.Predict(tree, input)
		if prediction == correctClass {
			correct++
		}
	}

	accuracy := float64(correct) / float64(total)

	fmt.Printf("Accuracy: %.2f\n", accuracy)

	csvFile, err := os.Open("titanic_test.csv")
	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(csvFile)

	counter := 0

	outCsv := [][]string{{"PassengerId", "Survived"}}

	for {
		record, err := csvReader.Read()
		if counter == 0 {
			counter++
			continue
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		row := stringSliceToFloat(record)
		inRow := row[2:]
		prediction := decision_tree.Predict(tree, inRow)
		outCsv = append(outCsv, []string{fmt.Sprint(row[1]), fmt.Sprint(prediction)})
	}

	csvOutFile, err := os.OpenFile("/Users/admin/Desktop/Workspace/golang/src/github.com/alex-bogomolov/decision_tree/submission.csv", os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		panic(err)
	}

	writer := csv.NewWriter(csvOutFile)
	err = writer.WriteAll(outCsv)

	if err != nil {
		panic(err)
	}



	err = csvOutFile.Close()

	if err != nil {
		panic(err)
	}

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


