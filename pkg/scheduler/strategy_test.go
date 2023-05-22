package scheduler

import (
	"fmt"
	"testing"
)

func Test_ExpandOrShrink_NotLimited(t *testing.T) {
	tables := [][]int64{}
	row1 := []int64{0, 0}
	row2 := []int64{100, 80}
	row3 := []int64{200, 140}
	row4 := []int64{1000, 1000}
	row5 := []int64{2000, 1998}
	row6 := []int64{10000, 5000}
	row7 := []int64{200000, 190000}
	row8 := []int64{10000, 80000}
	row9 := []int64{200000, 80000}
	row10 := []int64{1000, 800}
	row11 := []int64{100, 100}
	row12 := []int64{0, 0}
	row13 := []int64{0, 0}
	row14 := []int64{0, 0}
	row15 := []int64{0, 0}
	row16 := []int64{0, 0}
	row17 := []int64{0, 0}
	row18 := []int64{0, 0}
	row19 := []int64{0, 0}
	row20 := []int64{0, 0}

	tables = append(tables, row1)
	tables = append(tables, row2)
	tables = append(tables, row3)
	tables = append(tables, row4)
	tables = append(tables, row5)
	tables = append(tables, row6)
	tables = append(tables, row7)
	tables = append(tables, row8)
	tables = append(tables, row9)
	tables = append(tables, row10)
	tables = append(tables, row11)
	tables = append(tables, row12)
	tables = append(tables, row13)
	tables = append(tables, row14)
	tables = append(tables, row15)
	tables = append(tables, row16)
	tables = append(tables, row17)
	tables = append(tables, row18)
	tables = append(tables, row19)
	tables = append(tables, row20)

	gradienter := NewGradienter()
	var numActives int64
	var diff int64
	numActives = 100

	for i := 0; i < 20; i++ {
		diff = gradienter.ExpandOrShrink(tables[i][0], tables[i][1], numActives)
		numActives = numActives + diff
		fmt.Printf("%d, %d, %d, %d\n", tables[i][0], tables[i][1], diff, numActives)
	}
}

func Test_ExpandOrShrink_NumActivesLimited(t *testing.T) {
	tables := [][]int64{}
	row1 := []int64{0, 0}
	row2 := []int64{100, 80}
	row3 := []int64{200, 140}
	row4 := []int64{1000, 1000}
	row5 := []int64{2000, 1998}
	row6 := []int64{10000, 5000}
	row7 := []int64{200000, 190000}
	row8 := []int64{10000, 80000}
	row9 := []int64{200000, 80000}
	row10 := []int64{1000, 800}
	row11 := []int64{100, 100}
	row12 := []int64{0, 0}
	row13 := []int64{0, 0}
	row14 := []int64{0, 0}
	row15 := []int64{0, 0}
	row16 := []int64{0, 0}
	row17 := []int64{0, 0}
	row18 := []int64{0, 0}
	row19 := []int64{0, 0}
	row20 := []int64{0, 0}

	tables = append(tables, row1)
	tables = append(tables, row2)
	tables = append(tables, row3)
	tables = append(tables, row4)
	tables = append(tables, row5)
	tables = append(tables, row6)
	tables = append(tables, row7)
	tables = append(tables, row8)
	tables = append(tables, row9)
	tables = append(tables, row10)
	tables = append(tables, row11)
	tables = append(tables, row12)
	tables = append(tables, row13)
	tables = append(tables, row14)
	tables = append(tables, row15)
	tables = append(tables, row16)
	tables = append(tables, row17)
	tables = append(tables, row18)
	tables = append(tables, row19)
	tables = append(tables, row20)

	gradienter := NewGradienter()
	gradienter.SetMaxActives(5000)
	var numActives int64
	var diff int64
	numActives = 100

	for i := 0; i < 20; i++ {
		diff = gradienter.ExpandOrShrink(tables[i][0], tables[i][1], numActives)
		numActives = numActives + diff
		fmt.Printf("%d, %d, %d, %d\n", tables[i][0], tables[i][1], diff, numActives)
	}
}

func Test_ExpandOrShrink_MaxRateLimited(t *testing.T) {
	tables := [][]int64{}
	row1 := []int64{0, 0}
	row2 := []int64{100, 80}
	row3 := []int64{200, 140}
	row4 := []int64{1000, 1000}
	row5 := []int64{2000, 1998}
	row6 := []int64{10000, 5000}
	row7 := []int64{200000, 190000}
	row8 := []int64{10000, 80000}
	row9 := []int64{200000, 80000}
	row10 := []int64{1000, 800}
	row11 := []int64{100, 100}
	row12 := []int64{0, 0}
	row13 := []int64{0, 0}
	row14 := []int64{0, 0}
	row15 := []int64{0, 0}
	row16 := []int64{0, 0}
	row17 := []int64{0, 0}
	row18 := []int64{0, 0}
	row19 := []int64{0, 0}
	row20 := []int64{0, 0}

	tables = append(tables, row1)
	tables = append(tables, row2)
	tables = append(tables, row3)
	tables = append(tables, row4)
	tables = append(tables, row5)
	tables = append(tables, row6)
	tables = append(tables, row7)
	tables = append(tables, row8)
	tables = append(tables, row9)
	tables = append(tables, row10)
	tables = append(tables, row11)
	tables = append(tables, row12)
	tables = append(tables, row13)
	tables = append(tables, row14)
	tables = append(tables, row15)
	tables = append(tables, row16)
	tables = append(tables, row17)
	tables = append(tables, row18)
	tables = append(tables, row19)
	tables = append(tables, row20)

	gradienter := NewGradienter()
	gradienter.SetMaxRate(0.5)
	var numActives int64
	var diff int64
	numActives = 100

	for i := 0; i < 20; i++ {
		diff = gradienter.ExpandOrShrink(tables[i][0], tables[i][1], numActives)
		numActives = numActives + diff
		fmt.Printf("%d, %d, %d, %d\n", tables[i][0], tables[i][1], diff, numActives)
	}
}

func Test_ExpandOrShrink_MaxProcessedReqs(t *testing.T) {
	tables := [][]int64{}
	row1 := []int64{0, 0}
	row2 := []int64{100, 80}
	row3 := []int64{200, 140}
	row4 := []int64{1000, 1000}
	row5 := []int64{2000, 1998}
	row6 := []int64{10000, 5000}
	row7 := []int64{200000, 190000}
	row8 := []int64{10000, 80000}
	row9 := []int64{200000, 80000}
	row10 := []int64{1000, 800}
	row11 := []int64{100, 100}
	row12 := []int64{0, 0}
	row13 := []int64{0, 0}
	row14 := []int64{0, 0}
	row15 := []int64{0, 0}
	row16 := []int64{0, 0}
	row17 := []int64{0, 0}
	row18 := []int64{0, 0}
	row19 := []int64{0, 0}
	row20 := []int64{0, 0}

	tables = append(tables, row1)
	tables = append(tables, row2)
	tables = append(tables, row3)
	tables = append(tables, row4)
	tables = append(tables, row5)
	tables = append(tables, row6)
	tables = append(tables, row7)
	tables = append(tables, row8)
	tables = append(tables, row9)
	tables = append(tables, row10)
	tables = append(tables, row11)
	tables = append(tables, row12)
	tables = append(tables, row13)
	tables = append(tables, row14)
	tables = append(tables, row15)
	tables = append(tables, row16)
	tables = append(tables, row17)
	tables = append(tables, row18)
	tables = append(tables, row19)
	tables = append(tables, row20)

	gradienter := NewGradienter()
	gradienter.SetMaxProcessedReqs(5000)
	var numActives int64
	var diff int64
	numActives = 100

	for i := 0; i < 20; i++ {
		diff = gradienter.ExpandOrShrink(tables[i][0], tables[i][1], numActives)
		numActives = numActives + diff
		fmt.Printf("%d, %d, %d, %d\n", tables[i][0], tables[i][1], diff, numActives)
	}
}
