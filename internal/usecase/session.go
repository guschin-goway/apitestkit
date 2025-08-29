package usecase

import "github.com/guschin-goway/apitestkit/internal/adapter"

type Response struct {
	status int
	body   []byte
	report adapter.TestReporter
}

func (r *Response) Expect() *Expect {
	return &Expect{resp: r, report: r.report}
}
