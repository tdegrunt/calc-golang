package calc

import (
	"testing"
	"io"
)

type fsmResult struct {
	state int
	err error
}

type fsmTest struct {
	msg string
	input interface{}
	res fsmResult
}

func TestStateMachine(t *testing.T) {
	tests := []fsmTest {
		{ "Transition from 'started' to 'operator'", "SUM", fsmResult{stateOperator, nil} },
		{ "Transition from 'operator' to 'operand'", "10", fsmResult{stateOperand, nil} },
		{ "Transition from 'operand' to 'operand'", "20", fsmResult{stateOperand, nil} },
		{ "Transition from 'operand' to 'sentinel'", "\n", fsmResult{stateSentinel, nil} },
		{ "Transtion from 'sentinel' to 'stopped'", io.EOF, fsmResult{stateStopped, nil} },
	}

	fsm := NewStateMachine(expRules, nil)

	for _, tt := range tests {
		state, err := fsm.Process(tt.input)
		res := fsmResult{state, err}
		verify(t, tt.msg, tt.input, res, tt.res)
	}
}