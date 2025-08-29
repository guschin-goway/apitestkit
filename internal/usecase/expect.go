package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/guschin-goway/apitestkit/internal/adapter"
	"github.com/guschin-goway/apitestkit/pkg/assertion"

	"github.com/invopop/jsonschema"
	schema "github.com/santhosh-tekuri/jsonschema/v6"
)

type Expect struct {
	resp   *Response
	report adapter.TestReporter
}

func (e *Expect) Code(expected int) *Expect {
	if e.resp.status != expected {
		e.report.Fatal("unexpected code:", e.resp.status, "expected:", expected)
	}
	return e
}

func (e *Expect) JSONObject(out any) *Expect {
	if err := json.Unmarshal(e.resp.body, out); err != nil {
		e.report.Fatal("json unmarshal failed:", err)
	}
	return e
}

func (e *Expect) AssertAll(asserts ...assertion.Assertion) *Expect {
	var optErrors []error
	for _, a := range asserts {
		if err := a.Check(); err != nil {
			optErrors = append(optErrors, err)
		}
	}
	if len(optErrors) > 0 {
		for _, err := range optErrors {
			e.report.Error(err)
		}
		e.report.Fatal("optional assertion failed")
	}
	return e
}

func (e *Expect) Assert(groups ...assertion.GroupedAssertions) *Expect {
	assertion.RunGrouped(e.report, groups...)
	return e
}

func (e *Expect) JSONSchema(model any) *Expect {
	r := new(jsonschema.Reflector)
	schemaObj := r.ReflectFromType(reflect.TypeOf(model))

	schemaBytes, err := json.MarshalIndent(schemaObj, "", "  ")
	if err != nil {
		e.report.Fatal(err)
	}

	compiler := schema.NewCompiler()
	schemaJSON, err := schema.UnmarshalJSON(bytes.NewReader(schemaBytes))
	if err != nil {
		e.report.Fatal(err)
	}
	if err = compiler.AddResource("inline.json", schemaJSON); err != nil {
		e.report.Fatal(err)
	}

	sch, err := compiler.Compile("inline.json")
	if err != nil {
		e.report.Fatal(err)
	}

	if err = json.Unmarshal(e.resp.body, &model); err != nil {
		e.report.Fatal(err)
	}

	if err = sch.Validate(model); err != nil {
		e.report.Fatal(fmt.Sprintf("Schema validation failed: %v", err))
	}

	return e
}

func (e *Expect) Done() {}
