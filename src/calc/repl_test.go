package calc

import (
	"testing"
	"bytes"
)

type replResult struct {
	out string
	err string
}

type testProcessor struct { 
	states []int
	errors []error
	delegate StateMachineDelegate
	i int
}

func (this testProcessor) Process(input interface{}) (state int, err error) {
	state, err = this.states[this.i], this.errors[this.i]
	this.i += 1
	this.delegate.StateChanged(state-1, state, input)
	return
}

type testEvaluator struct {
	op OpCode
	res float64
	err error
}

func (this testEvaluator) Evaluate() (float64, error) {
	return this.res, this.err
}

func (this testEvaluator) Operator(operator *OpCode) OpCode {
	return this.op
}

func (this testEvaluator) Operands(operands []float64) []float64 {
	return nil
}

type testEvaluatorFactory struct { 
	exp *testEvaluator
}

func (this testEvaluatorFactory) NewEvaluator() Evaluator {
	return this.exp
} 

func TestRepl(t *testing.T) {
	var in, out, err bytes.Buffer
	repl := NewRepl(&in, &out, &err)

	tests := [] struct {
		msg string
		in string
		res replResult
		fsm testProcessor
		fac testEvaluatorFactory
	}{
		{ "The sum of 10, 20, 30 is 60", "sum:10,20,30\n", replResult{"SUM:60\n", ""}, 
			testProcessor { 
				states: []int { stateOperator, stateOperand, stateOperand, stateOperand, stateSentinel, stateStopped },
				errors: []error { nil, nil, nil, nil, nil, nil },
				delegate: repl,
			},
			testEvaluatorFactory { &testEvaluator{OpSum, 60, nil} },
		},
	}
	
	for _, tt := range tests {
		in.WriteString(tt.in)
		repl.Read(tt.fsm, tt.fac)
		res := replResult {
			out.String(), 
			err.String(),
		}

		verify(t, tt.msg, tt.in, res, tt.res)

		in.Reset()
		out.Reset()
		err.Reset()
	}
}