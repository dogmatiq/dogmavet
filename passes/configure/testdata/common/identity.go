package application

import (
	"github.com/dogmatiq/dogma"
)

// someCondition is a placeholder used when some branching condition is needed,
// but the truthiness of the condition is not relevant.
var someCondition bool

type (
	unidentified  struct{}
	identified1   struct{}
	identified2   struct{}
	someBranches1 struct{}
	someBranches2 struct{}
	someBranches3 struct{}
	duplicate1    struct{}
	duplicate2    struct{}
)

// Configure calls c.Identity() exactly once for every invocation, even though
// it has branches.
func (identified1) Configure(c dogma.ApplicationConfigurer) {
	if someCondition {
		c.Identity("name", "6c8faf1c-3748-43d5-bf58-d10663ee7c1d")
	} else {
		c.Identity("name", "6c8faf1c-3748-43d5-bf58-d10663ee7c1d")
	}
}

// Configure calls c.Identity() exactly once for every invocation, even though
// it has branches.
func (identified2) Configure(c dogma.ApplicationConfigurer) {
	if someCondition {
		c.Identity("name", "aea1fb69-49cd-402e-b6da-e259285b69b3")
		return
	}

	c.Identity("name", "aea1fb69-49cd-402e-b6da-e259285b69b3")
}

// Configure does not call c.Identity() at all.
func (unidentified) Configure(c dogma.ApplicationConfigurer) {} // want `Configure\(\) must call c.Identity\(\)`

// Configure only calls c.Identity() on some of its control-flow paths.
func (someBranches1) Configure(c dogma.ApplicationConfigurer) {
	if someCondition { // want `this control-flow statement causes c.Identity\(\) to remain uncalled on some execution paths`
		c.Identity("name", "b7f78aeb-84d4-404f-a677-3dc179e9d656")
	}
}

// Configure only calls c.Identity() on some of its control-flow paths.
func (someBranches2) Configure(c dogma.ApplicationConfigurer) {
	if someCondition { // want `this control-flow statement causes c.Identity\(\) to remain uncalled on some execution paths`
		return
	}

	c.Identity("name", "9bd2101e-0356-4896-979c-dc025273d5a0")
}

// Configure only calls c.Identity() on some of its control-flow paths.
//
// This test verifies that the diagnostic is shown on the most-specific
// flow-control statement.
func (someBranches3) Configure(c dogma.ApplicationConfigurer) {
	if someCondition {
		c.Identity("name", "cba5363b-b560-4c2b-8a53-825a6fb2770e")
	} else {
		if someCondition { // want `this control-flow statement causes c.Identity\(\) to remain uncalled on some execution paths`
			c.Identity("name", "cba5363b-b560-4c2b-8a53-825a6fb2770e")
		}
	}
}

// Configure calls c.Identity() more than once within the same block.
func (duplicate1) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("name", "05d1e6ec-41a0-43b4-a866-4b5cce3bf2fd")
	c.Identity("name", "05d1e6ec-41a0-43b4-a866-4b5cce3bf2fd") // want `c.Identity\(\) must be called exactly once`
	c.Identity("name", "05d1e6ec-41a0-43b4-a866-4b5cce3bf2fd") // want `c.Identity\(\) must be called exactly once`
}

// Configure calls c.Identity() more than once on some execution paths.
func (duplicate2) Configure(c dogma.ApplicationConfigurer) {
	if someCondition {
		c.Identity("name", "c454bbac-1a93-49d0-8d8c-a4671e3e3025")
	}

	c.Identity("name", "c454bbac-1a93-49d0-8d8c-a4671e3e3025") // want `c.Identity\(\) has already been called on at least one execution path`
}
