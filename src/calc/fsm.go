package calc

import "fmt"

type StateMachineDelegate interface {
	StateChanged(fromState int, toState int, input interface{})
}

type StateMachineRule struct {
	FromState int
	ToState int
	SatisfiedBy func(interface{}) bool
}

type StateMachine struct {
	rules []StateMachineRule
	state int
	delegate StateMachineDelegate
}

func NewStateMachine(rules []StateMachineRule, delegate StateMachineDelegate) *StateMachine {
	return &StateMachine { 
		rules: rules, 
		state: rules[0].FromState, 
		delegate: delegate,
	}
}

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

func (this StateMachine) findRule(input interface{}) *StateMachineRule {
	for _, r := range this.rules {
		if r.FromState == this.state && r.SatisfiedBy(input) {
			return &r
		}
	}
	return nil
}