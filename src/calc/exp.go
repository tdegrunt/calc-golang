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
	operator OpCode
	operands []float64
}

type ExpressionFactory int

func (this ExpressionFactory) NewEvaluator() Evaluator {
	return new(Expression)
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
	if this.operands == nil || len(this.operands) == 0 {
		return 0, errors.New("No operands provided")
	}

	if fn, ok := operators[this.operator]; ok {
		return fn(this.operands), nil
	}

	return 0, errors.New("Unknown opcode")
}

func (this *Expression) Operator(operator *OpCode) OpCode {
	if (operator != nil) {
		this.operator = *operator
	}
	return this.operator
}

func (this *Expression) Operands(operands []float64) []float64 {
	if (operands != nil) {
		this.operands = operands
	}
	return this.operands
}