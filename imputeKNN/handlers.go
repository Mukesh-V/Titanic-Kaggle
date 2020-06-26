package knnimpute

import (
	"log"
	"math"
	"sort"
	"strconv"
)

type distanceMap struct {
	idx      int
	distance float64
}

func findMissing(records [][]string, impute int) []int {
	var missing []int
	for i := range records {
		if records[i][impute] == "" || records[i][impute] == "NaN" {
			missing = append(missing, i)
		}
	}
	return missing
}

func euclid(row1 []string, row2 []string, impute int) float64 {
	dist := 0.0
	// var v1 float64
	// var v2 float64
	for idx := range row1 {
		if idx == impute {
			continue
		}
		v1, err1 := strconv.ParseFloat(row1[idx], 64)
		if err1 != nil {
			log.Fatal(err1)
		}
		v2, _ := strconv.ParseFloat(row2[idx], 64)
		dist += math.Pow(v1-v2, 2)
	}
	return math.Sqrt(dist)
}

func neighbours(records [][]string, row []string, impute int, num int) []int {
	var dists []distanceMap
	var neighbourIndices []int

	for idx, record := range records {
		dist := distanceMap{idx, euclid(row, record, impute)}
		dists = append(dists, dist)
	}

	sort.Slice(dists, func(p, q int) bool {
		return dists[p].distance < dists[q].distance
	})
	dists = dists[:num]

	for i := 0; i < num; i++ {
		neighbourIndices = append(neighbourIndices, dists[i].idx)
	}

	return neighbourIndices
}
