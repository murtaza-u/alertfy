package conf_test

import (
	"context"
	"testing"

	"github.com/murtaza-u/amify/internal/conf"

	"github.com/stretchr/testify/assert"
)

type Params struct {
	X   int
	Y   bool
	Z   string
	Foo foo
}

type foo struct {
	Bar string
}

type Input struct {
	Typ    string
	Expr   string
	Output any
}

func TestExpression(t *testing.T) {
	params := Params{X: 10, Y: false, Z: "foo"}
	inputs := []Input{
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
	inputs := []Input{
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
