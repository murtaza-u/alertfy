package conf_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/murtaza-u/alertfy/internal/conf"

	"github.com/stretchr/testify/assert"
)

type Params struct {
	X   int
	Y   bool
	Z   string
	Foo foo
	Bar map[string]string
}

type foo struct {
	Bar string
}

type InputExpr struct {
	Typ    string
	Expr   string
	Output any
}

type InputTemplate struct {
	Template string
	Output   string
}

func TestExpression(t *testing.T) {
	params := Params{X: 10, Y: false, Z: "foo"}
	inputs := []InputExpr{
		{
			Typ:    "bool",
			Expr:   `X == 10`,
			Output: true,
		},
		{
			Typ:    "string",
			Expr:   `Y ? "bob" : "alice"`,
			Output: "alice",
		},
		{
			Typ:    "string",
			Expr:   `Z == "foo" ? "bar" : "blah"`,
			Output: "bar",
		},
	}
	for _, i := range inputs {
		var expr conf.Expr
		err := expr.UnmarshalText([]byte(i.Expr))
		isNil := assert.Nilf(t, err, "unmarshalling expression: %s", i.Expr)
		if !isNil {
			continue
		}

		ctx := context.Background()
		var out any

		switch i.Typ {
		case "bool":
			out, err = expr.Evaluable.EvalBool(ctx, params)
		case "string":
			out, err = expr.Evaluable.EvalString(ctx, params)
		case "int":
			out, err = expr.Evaluable.EvalInt(ctx, params)
		}

		isNil = assert.Nilf(t, err, "evaluating expression: `%s`", i.Expr)
		if isNil {
			assert.Equalf(t, i.Output, out, "expression: `%s`", i.Expr)
		}
	}
}

func TestStringExpression(t *testing.T) {
	params := Params{X: 10, Y: false, Z: "foo", Foo: foo{Bar: "bar"}}
	inputs := []InputExpr{
		{
			Typ:    "string",
			Expr:   `foobar`,
			Output: "foobar",
		},
		{
			Typ:    "string",
			Expr:   `Y ? "bob" : "alice"`,
			Output: "alice",
		},
		{
			Typ:    "string",
			Expr:   `Z == "foo" ? "bar" : "blah"`,
			Output: "bar",
		},
		{
			Typ:    "string",
			Expr:   `Foo.Bar`,
			Output: "bar",
		},
	}
	for _, i := range inputs {
		var sexpr conf.StringExpr
		err := sexpr.UnmarshalText([]byte(i.Expr))
		isNil := assert.Nilf(t, err, "unmarshalling expression: %s", i.Expr)
		if !isNil {
			continue
		}

		ctx := context.Background()
		var out any

		switch i.Typ {
		case "bool":
			if sexpr.Expr == nil {
				out = sexpr.Text
			} else {
				out, err = sexpr.Expr.Evaluable.EvalBool(ctx, params)
			}
		case "string":
			if sexpr.Expr == nil {
				out = sexpr.Text
			} else {
				out, err = sexpr.Expr.Evaluable.EvalString(ctx, params)
			}
		case "int":
			if sexpr.Expr == nil {
				out = sexpr.Text
			} else {
				out, err = sexpr.Expr.Evaluable.EvalInt(ctx, params)
			}
		}

		isNil = assert.Nilf(t, err, "evaluating expression: `%s`", i.Expr)
		if isNil {
			assert.Equalf(t, i.Output, out, "expression: `%s`", i.Expr)
		}
	}
}

func TestTemplate(t *testing.T) {
	params := Params{
		X:   10,
		Y:   false,
		Z:   "foo",
		Foo: foo{Bar: "bar"},
		Bar: map[string]string{
			"a": "b",
		},
	}
	inputs := []InputTemplate{
		{
			Template: `
			{{- if eq .X 10 -}}
				It is 10
			{{- else -}}
				It is something else
			{{- end -}}
			`,
			Output: "It is 10",
		},
		{
			Template: `{{ index .Bar "a" }}`,
			Output:   "b",
		},
		{
			Template: `{{ index .Bar "c" }}`,
			Output:   "",
		},
	}
	for idx, i := range inputs {
		var tmpl conf.Template
		err := tmpl.UnmarshalText([]byte(i.Template))
		isNil := assert.Nilf(t, err, "unmarshalling template: %d", idx)
		if !isNil {
			continue
		}

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, params)
		isNil = assert.Nilf(t, err, "evaluating template: %d", idx)
		if isNil {
			assert.Equalf(t, i.Output, buf.String(), "template: %d", idx)
		}
	}
}
