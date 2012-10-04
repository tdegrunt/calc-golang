// ## Repl

// A simple, interactive expression evaluation environment.

package calc

import (
	"strconv"
	"io"
	"fmt"
	"unicode"
	"strings"
)


// In order to keep `Repl` generic we depend on intefaces rather than concretions.
type Processor interface {
	Process(interface{}) (int, error)
}

type Evaluator interface {
	Evaluate() (float64, error)
	Operator(*OpCode) OpCode
	Operands([]float64) []float64
}

type EvaluatorFactory interface {
	NewEvaluator() Evaluator
}

// We need to maintain a buffer containing the characters read thus far but not yet
// evaluated, our state machine, a factory used to create new expressions, the current
// expression, IO streams, as well as a boolean indicating whether to skip characters
// from the input stream.
type Repl struct {
	buf []rune
	fsm Processor
	fac EvaluatorFactory
	exp Evaluator
	out io.Writer
	err io.Writer
	in io.RuneReader
	skip bool
}

// ### Creation

// To instantiate a repl we need three IO streams corresponding to stdin, stdout & stderr.
func NewRepl(in io.RuneReader, out io.Writer, err io.Writer) (r *Repl) {
	r = &Repl{in: in, out: out, err: err}
	return
}

// ### Usage

// `ReadDefault` uses the default concretions for `Processor` and `EvaluatorFactory`.
func (this *Repl) ReadDefault() {
	fsm := NewStateMachine(expRules, this)
	fac := new(ExpressionFactory)
	this.Read(fsm, fac)
}

func (this *Repl) Read(fsm Processor, fac EvaluatorFactory) {
	this.fsm = fsm
	this.fac = fac
	this.exp = fac.NewEvaluator()

	fmt.Fprintln(this.err, "Enter expressions to evaluate followed by a newline.",
		"Type \"QUIT\" to exit.")

	for {
		if r, count, err := this.in.ReadRune(); err != nil {
			this.fsm.Process(io.EOF) 
			return // eof: exit loop
		} else if r == unicode.ReplacementChar && count == 1 {
			fmt.Fprintln(this.err, "Invalid character in input.")
			this.skip = true
		} else if !this.handleRune(r) {
			return
		}
	}
}

// ### Internals

// Here we determine how to action a character:

// * If we encounter a delimiter we flush the input buffer to a string and feed it into the state machine. 
// * If we've reached the end of the expression and notify the state machine. 
// * Otherwise, we append the character to the input buffer.
func (this *Repl) handleRune(r rune) bool {
	if this.skip {
		if r == '\n' {
			this.skip = false
		}
	} else if _, ok := delimiters[r]; ok {
		
		tok := this.flushBuffer() // delimiter: flush buffer & feed the machine
		
		if tok == "QUIT" {
			fmt.Fprintln(this.err, "Goodbye!")
			return false
		}

		if _, err := this.fsm.Process(tok); err != nil {
			fmt.Fprintln(this.err, err.Error())
			this.skip = true
		} else if r == '\n' {
			this.fsm.Process(string(r)) // sentinel: end of expression
		}
	} else {
		this.buf = append(this.buf, r) // append char to input buffer
	}

	return true
}

func (this *Repl) flushBuffer() (str string) {
	str = strings.TrimSpace(string(this.buf))
	str = strings.ToUpper(str)
	this.buf = this.buf[:0]
	return
}

func (this Repl) evalAndPrint() {
	if num, err := this.exp.Evaluate(); err != nil {
		fmt.Fprintln(this.err, err)
	} else {
		fmt.Fprintf(this.out, "%s: %f\n", 
			tokensByOperator[this.exp.Operator(nil)], num)
	}
	return
}

// This is our state machine callback. Depending on the current state we either

// * evaluate the current expression
// * append the input to the collection of operands
// * evaluate and print the result
// * do nothing
func (this *Repl) StateChanged(fromState int, toState int, input interface{}) {
	if tok, ok := input.(string); ok {
		switch toState {
		case stateOperator:
			op := operatorsByToken[tok]
			this.exp = this.fac.NewEvaluator()
			this.exp.Operator(&op)
		case stateOperand:
			num, _ := strconv.ParseFloat(tok, 64)
			this.exp.Operands(append(this.exp.Operands(nil), num))
		case stateSentinel:
			this.evalAndPrint()
			this.exp = nil
		}
	}
	return
}

var delimiters = map[rune]bool {
	':' : true,
	',' : true,
	'\n' : true,
}

var tokensByOperator = map[OpCode]string {
	OpSum : "SUM",
	OpMin : "MIN",
	OpMax : "MAX",
	OpAvg : "AVG",
}