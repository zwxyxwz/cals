package taxcalculator

import (
	"fmt"
	"testing"

	"github.com/spf13/cast"
)

func TestCalOneMonth(t *testing.T) {
	cal, err := NewTaxCalculator(2024, RantingTag)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(cal.Cal(2000, 7))
	fmt.Println(cal.CalThisYear([]float64{5000, 5000, 5000, 5000, 5000, 5000}))
	fmt.Println(cal.CalThisYear([]float64{15000, 15000, 15000, 15000, 15000, 15000}))
}

func TestBackwords(t *testing.T) {
	cal, err := NewTaxCalculator(2025, RantingTag)
	if err != nil {
		t.Error(err)
		return
	}
	gap := float64(100)
	for base := float64(5000); base <= 50000; base += gap {
		var input []float64
		for c := 0; c <= 6; c++ {
			input = append(input, base)
		}
		res := cal.CalThisYear(input)
		if cast.ToFloat64(res[5])-cast.ToFloat64(res[6]) > 400 && cast.ToFloat64(res[5])-cast.ToFloat64(res[6]) < 600 {
			fmt.Println(res, "\n", cast.ToFloat64(res[5])-cast.ToFloat64(res[6]), base)
		}
	}
}
