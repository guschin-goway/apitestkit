package adapter

import "testing"

type TReporter struct {
	t *testing.T
}

type TestReporter interface {
	Fatal(args ...any)
	Error(args ...any)
	Helper()
	FailNow()
}

func NewTReporter(t *testing.T) *TReporter {
	return &TReporter{t: t}
}

func (r *TReporter) Fatal(args ...any) {
	r.t.Fatal(args...)
}

func (r *TReporter) Error(args ...any) {
	r.t.Error(args...)
}

func (r *TReporter) Helper() {
	r.t.Helper()
}

func (r *TReporter) FailNow() {
	r.t.FailNow()
}
