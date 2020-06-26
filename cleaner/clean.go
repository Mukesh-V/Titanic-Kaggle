package cleaner

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

//FindNaN helps in finding NaN entries
func FindNaN(file string) {
	f1, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(f1)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	r.FieldsPerRecord = 8

	count := make([]int, 8)
	for _, record := range records {
		for idx := range record {
			if record[idx] == "NaN" {
				count[idx] = count[idx] + 1
			}
			if record[idx] == "" {
				count[idx] = count[idx] + 1
			}
		}
	}

	for idx := range count {
		if count[idx] != 0 {
			fmt.Printf("[%d]%s : %d \n", idx, records[0][idx], count[idx])
		}
	}
}

//Encode helps in numerising categorical data
func Encode(file string, m map[string]int, idx int, isHeader bool) {
	f1, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	r := csv.NewReader(f1)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for i, record := range records {
		if isHeader == true && i == 0 {
			continue
		}
		f := fmt.Sprintf("%d", m[record[idx]])
		record[idx] = f
	}

	f2, err2 := os.Create(file)
	if err2 != nil {
		log.Fatal(err2)
	}
	defer f2.Close()

	w := csv.NewWriter(f2)
	w.WriteAll(records)
}
