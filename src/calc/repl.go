package calc

import (
	"strconv"
	"io"
	"fmt"
	"unicode"
	"strings"
)

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

func NewRepl(in io.RuneReader, out io.Writer, err io.Writer) (r *Repl) {
	r = &Repl{in: in, out: out, err: err}
	return
}

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
			// eof: exit loop
			this.fsm.Process(io.EOF)
			return
		} else if r == unicode.ReplacementChar && count == 1 {
			fmt.Fprintln(this.err, "Invalid character in input.")
			this.skip = true
		} else if !this.handleRune(r) {
			return
		}
	}
}

func (this *Repl) handleRune(r rune) bool {
	if this.skip {
		if r == '\n' {
			this.skip = false
		}
	} else if _, ok := delimiters[r]; ok {
		// delimiter: flush buffer & feed the machine
		tok := this.flushBuffer()
		
		if tok == "QUIT" {
			fmt.Fprintln(this.err, "Goodbye!")
			return false
		}

		if _, err := this.fsm.Process(tok); err != nil {
			fmt.Fprintln(this.err, err.Error())
			this.skip = true
		} else if r == '\n' {
			// sentinel: end of expression
			this.fsm.Process(string(r))
		}
	} else {
		// append char to input buffer
		this.buf = append(this.buf, r)
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

