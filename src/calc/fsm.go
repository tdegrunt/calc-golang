// ## StateMachine

// A `StateMachine` is a simple, generic finite state machine.

package calc

import "fmt"

// `StateMachineDelegate` is basically a stateful callback function
// invoked after a state transition.
type StateMachineDelegate interface {
	StateChanged(fromState int, toState int, input interface{})
}

// A state machine rule specifies source & target states as well as a
// predicate used to determine if the input satisfies the rule.
type StateMachineRule struct {
	FromState int
	ToState int
	SatisfiedBy func(interface{}) bool
}

// We store a collection of rules, the current state and a delegate.
type StateMachine struct {
	rules []StateMachineRule
	state int
	delegate StateMachineDelegate
}

// ### Creation

// The initial state is the source state of the first rule.
func NewStateMachine(rules []StateMachineRule, delegate StateMachineDelegate) *StateMachine {
	return &StateMachine { 
		rules: rules, 
		state: rules[0].FromState, 
		delegate: delegate,
	}
}

// ### Usage

// The machine is fed input and if the input matches a rule the state is changed accordingly
// and the provided callback (if any) is invoked.
func (this *StateMachine) Process(input interface{}) (state int, err error) {
	if r := this.findRule(input); r != nil {
		this.state = r.ToState
		state = r.ToState
		if this.delegate != nil {
			this.delegate.StateChanged(r.FromState, r.ToState, input)
		}
	} else {
		err = fmt.Errorf("No matching rule for input %#v in state %d",
			input, this.state)
	}
	return
}

func (this *StateMachine) Reset() {
	this.state = this.rules[0].FromState
}

func (this StateMachine) findRule(input interface{}) *StateMachineRule {
	for _, r := range this.rules {
		if r.FromState == this.state && r.SatisfiedBy(input) {
			return &r
		}
	}
	return nil
}