package knnimpute

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

//ImputeKNN fills holes in data
func ImputeKNN(file string, hasHeader bool, col int) {
	f1, err1 := os.Open(file)
	if err1 != nil {
		log.Fatal(err1)
	}
	defer f1.Close()

	r := csv.NewReader(f1)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var df [][]string

	if hasHeader {
		df = append(df, records[0])
		records = records[1:]
	}

	missing := findMissing(records, col)
	for _, idx := range missing {
		nei := neighbours(records, records[idx], col, 4)
		inter := fmt.Sprintf("%f", getVal(records, idx, nei, col))
		records[idx][col] = inter
	}

	for _, record := range records {
		df = append(df, record)
	}

	f2, err2 := os.Create(file)
	if err2 != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	w := csv.NewWriter(f2)
	w.WriteAll(df)

}

func getVal(records [][]string, row int, neighbours []int, impute int) float64 {
	max := 0.0
	for _, nei := range neighbours {
		item, err := strconv.ParseFloat(records[nei][impute], 64)
		if err != nil {
			log.Fatal(err)
		}
		if max < item {
			max = item
		}
	}
	return max
}
