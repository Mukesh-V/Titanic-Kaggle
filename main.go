package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"titanic/cleaner"
	knnimpute "titanic/imputeKNN"

	"github.com/kniren/gota/dataframe"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/trees"
)

func main() {

	if os.Args[1] == "clean" {

		//for train.csv

		train, err := os.Open("train.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer train.Close()

		trainDF := dataframe.ReadCSV(train)
		records := trainDF.Records()
		var y []string
		for idx := range records {
			y = append(y, records[idx][1])
		}

		removed := trainDF.Drop([]int{0, 1, 3, 8, 9, 10, 11})
		newRecords := removed.Records()
		var final [][]string

		for idx, record := range newRecords {
			record = append(record, y[idx])
			final = append(final, record)
		}

		train2, err := os.Create("train_init.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer train2.Close()

		w := csv.NewWriter(train2)
		w.WriteAll(final)

		//for test.csv

		test, err := os.Open("test.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer test.Close()

		testDF := dataframe.ReadCSV(test)
		removedt := testDF.Drop([]int{0, 2, 7, 8, 9, 10})

		test2, err := os.Create("test_init.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer test2.Close()

		removedt.WriteCSV(test2)
	}

	if os.Args[1] == "encode" {
		g := make(map[string]int)
		g["male"] = 0
		g["female"] = 1
		cleaner.Encode("train_init.csv", g, 1, true)
		cleaner.Encode("test_init.csv", g, 1, true)
	}

	if os.Args[1] == "impute" {
		knnimpute.ImputeKNN("train_init.csv", true, 2)
		knnimpute.ImputeKNN("test_init.csv", true, 2)
	}

	if os.Args[1] == "train" {

		train, err := base.ParseCSVToInstances("train_init.csv", true)
		if err != nil {
			log.Fatal(err)
		}
		// 44111378

		rand.Seed(44111330)
		tree := trees.NewID3DecisionTree(0.3)

		cv, err := evaluation.GenerateCrossFoldValidationConfusionMatrices(train, tree, 5)
		if err != nil {
			log.Fatal(err)
		}
		mean, variance := evaluation.GetCrossValidatedMetric(cv, evaluation.GetAccuracy)
		stdev := math.Sqrt(variance)

		fmt.Printf("\nAccuracy\n%.2f (+/- %.2f)\n\n", mean, stdev*2)

		test, err := base.ParseCSVToInstances("test_init.csv", true)
		if err != nil {
			log.Fatal(err)
		}

		output, _ := tree.Predict(test)
		if err != nil {
			log.Fatal(err)
		}
		printCSV(output)
	}

}

func printCSV(op base.FixedDataGrid) {
	var finalOP [][]string

	headers := make([]string, 2)
	headers[0] = "PassengerId"
	headers[1] = "Survived"
	finalOP = append(finalOP, headers)

	for idx := 0; idx < 418; idx++ {
		record := make([]string, 2)
		record[0] = fmt.Sprintf("%d", idx+892)
		item, _ := strconv.ParseFloat(op.RowString(idx), 64)
		record[1] = fmt.Sprintf("%d", int(item))
		finalOP = append(finalOP, record)
	}

	final, err := os.Create("final.csv")
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(final)
	w.WriteAll(finalOP)
}
