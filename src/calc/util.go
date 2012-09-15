package calc

import (
	"testing"
	"reflect"
	"strconv"
	"io"
)

func verify(t *testing.T, testCase string, input interface{}, output interface{}, expected interface{}) {
	if !reflect.DeepEqual(output, expected) {
		t.Errorf("\nTest case:\t%s\nInput:\t\t%v\nActual:\t\t%v\nExpected:\t%v\n", testCase, input, output, expected)
	}
}

const (
	stateStarted int = iota
	stateOperator
	stateOperand
	stateSentinel
	stateStopped
)

var expRules = []StateMachineRule {
	{ stateStarted, stateOperator, isOperator },
	{ stateOperator, stateOperand, isOperand },
	{ stateOperand, stateOperand, isOperand },
	{ stateOperand, stateSentinel, isSentinel },
	{ stateSentinel, stateOperator, isOperator },
	{ stateSentinel, stateStopped, isEOF },
}

func isEOF(input interface{}) bool {
	return reflect.DeepEqual(input, io.EOF)
}

func isSentinel(input interface{}) bool {
	if tok, ok := input.(string); ok {
		if tok == "\n" {
			return true
		}
	}
	return false
}

func isOperator(input interface{}) (b bool) {
	if tok, ok := input.(string); ok {
		_, b = operatorsByToken[tok]
	} else {
		b = false
	}
	return
}

func isOperand(input interface{}) bool {
	if tok, ok := input.(string); ok {
		if _, err := strconv.ParseFloat(tok, 64); err == nil {
			return true
		}
	}
	return false
}

var operatorsByToken = map[string]OpCode {
	"SUM" : OpSum,
	"MIN" : OpMin,
	"MAX" : OpMax,
	"AVG" : OpAvg,
}
