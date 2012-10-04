// ## Expression

// The `Expression` type represents simple arithmetic expressions consisting of an operator
// and one or more operands and yielding a single result.

package calc

import "errors"

// An *opcode* is used to represent the operator.
type OpCode int

const (
	OpSum OpCode = iota
	OpMin
	OpMax
	OpAvg
)

// Operands are stored as a slice of type `float64`.
type Expression struct {
	operator OpCode
	operands []float64
}

// ### Creation

// We use a factory to create `Expression` objects in a generic way.
type ExpressionFactory int

func (this ExpressionFactory) NewEvaluator() Evaluator {
	return new(Expression)
}

// ### Usage

// We use a map to lookup operations by opcode. Each operation is simply a function that 
// applies the operator to the operands.
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

// ### Getters / Setters

// Get or set the operator or operands.
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