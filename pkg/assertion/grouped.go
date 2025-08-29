package assertion

import (
	"github.com/guschin-goway/apitestkit/internal/adapter"
)

type GroupedAssertions struct {
	Required bool
	Checks   []Assertion
}

func Required(asserts ...Assertion) GroupedAssertions {
	return GroupedAssertions{Required: true, Checks: asserts}
}

func Optional(asserts ...Assertion) GroupedAssertions {
	return GroupedAssertions{Required: false, Checks: asserts}
}

// RunGrouped — выполняет группы проверок
func RunGrouped(t adapter.TestReporter, groups ...GroupedAssertions) {
	var optionalErrors []error

	for _, g := range groups {
		for _, a := range g.Checks {
			if err := a.Check(); err != nil {
				if g.Required {
					t.Fatal("Required assertion failed: %v", err)
				} else {
					optionalErrors = append(optionalErrors, err)
				}
			}
		}
	}

	if len(optionalErrors) > 0 {
		for _, e := range optionalErrors {
			t.Error("Optional assertion failed: %v", e)
		}
		t.FailNow()
	}
}
