package assertion

import "fmt"

// Assertion — интерфейс для проверки
type Assertion interface {
	Check() error
}

// Equal — проверка на равенство
func Equal(expected, actual any) Assertion {
	return assertionFunc(func() error {
		if expected != actual {
			return fmt.Errorf("equal failed: expected %v, got %v", expected, actual)
		}
		return nil
	})
}

// NotEmpty — проверка, что значение не пустое
func NotEmpty(actual any) Assertion {
	return assertionFunc(func() error {
		if actual == "" || actual == nil {
			return fmt.Errorf("NotEmpty failed: got %v", actual)
		}
		return nil
	})
}

// Contains — проверка на вхождение подстроки
func Contains(actual, substr string) Assertion {
	return assertionFunc(func() error {
		if !contains(actual, substr) {
			return fmt.Errorf("contains failed: '%s' does not contain '%s'", actual, substr)
		}
		return nil
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (substr == "" || (len(s) > 0 &&
		(len([]rune(s)) >= len([]rune(substr))) &&
		(s[:len(substr)] == substr || contains(s[1:], substr))))
}

type assertionFunc func() error

func (f assertionFunc) Check() error { return f() }
