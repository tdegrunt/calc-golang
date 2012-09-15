package calc

import (
	"testing"
	"io"
)

type fsmProcessResult struct {
	state int
	err error
}

type fsmTest struct {
	msg string
	input interface{}
	res fsmProcessResult
}

func TestStateMachine(t *testing.T) {
	tests := []fsmTest {
		{ "Transition from 'started' to 'operator'", "SUM", fsmProcessResult{stateOperator, nil} },
		{ "Transition from 'operator' to 'operand'", "10", fsmProcessResult{stateOperand, nil} },
		{ "Transition from 'operand' to 'operand'", "20", fsmProcessResult{stateOperand, nil} },
		{ "Transition from 'operand' to 'sentinel'", "\n", fsmProcessResult{stateSentinel, nil} },
		{ "Transtion from 'sentinel' to 'stopped'", io.EOF, fsmProcessResult{stateStopped, nil} },
	}

	fsm := NewStateMachine(expRules, nil)

	for _, tt := range tests {
		state, err := fsm.Process(tt.input)
		res := fsmProcessResult{state, err}
		verify(t, tt.msg, tt.input, res, tt.res)
	}
}