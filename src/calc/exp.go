package calc

import "errors"

type OpCode int

const (
	OpSum OpCode = iota
	OpMin
	OpMax
	OpAvg
)

type Expression struct {
	Operator OpCode
	Operands []float64
}

var operators = map[OpCode]func([]float64) float64{
	OpSum: func(operands []float64) float64 {
		sum := float64(0)
		for _, n := range operands {
			sum += n
		}
		return sum
	},
	OpMin: func(operands []float64) float64 {
		min := operands[0]
		for _, n := range operands {
			if n < min {
				min = n
			}
		}
		return min
	},
	OpMax: func(operands []float64) float64 {
		max := operands[0]
		for _, n := range operands {
			if n > max {
				max = n
			}
		}
		return max
	},
	OpAvg: func(operands []float64) float64 {
		total := float64(0)
		for _, n := range operands {
			total += n
		}
		return total / float64(len(operands))
	},
}

func (this Expression) Evaluate() (float64, error) {
	if this.Operands == nil || len(this.Operands) == 0 {
		return 0, errors.New("No operands provided")
	}

	if fn, ok := operators[this.Operator]; ok {
		return fn(this.Operands), nil
	}

	return 0, errors.New("Unknown opcode")
}
