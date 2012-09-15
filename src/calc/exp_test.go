package calc

import (
	"testing"
	"errors"
)

type expResult struct {
	num float64
	err error
}

type expTest struct {
	msg string
	exp Expression
	res expResult
}

func TestExpression(t *testing.T) {
	tests := []expTest {
		{"The sum of 10, 20, 30 is 60", Expression{OpSum, []float64{10, 20, 30}}, expResult{60, nil}},
		{"The min of 30, 10, -1, 20 is -1", Expression{OpMin, []float64{30, 10, -1, 20}}, expResult{-1, nil}},
		{"The max of 0, 2, 4.5, -20 is 4.5", Expression{OpMax, []float64{0, 2, 4.5, -20}}, expResult{4.5, nil}},
		{"The avg of 5, 10, 15, 20 is 12.5", Expression{OpAvg, []float64{5, 10, 15, 20}}, expResult{12.5, nil}},
		{"Raise error when no operands provided", Expression{OpAvg, nil}, expResult{0, errors.New("No operands provided")}},
		{"Raise error when unknown opcode provided", Expression{-1, []float64{1, 2}}, expResult{0, errors.New("Unknown opcode")}},
	}

	for _, tt := range tests {
		num, err := tt.exp.Evaluate()
		res := expResult{num, err}
		verify(t, tt.msg, tt.exp, res, tt.res)
	}
}