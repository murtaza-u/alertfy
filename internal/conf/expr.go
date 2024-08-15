package conf

import (
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"github.com/PaesslerAG/gval"
)

// Expr consists of an evaluable gval expression. It implements the
// encoding.TextUnmarshaler interface.
type Expr struct {
	Text      string
	Evaluable gval.Evaluable
}

func (e *Expr) UnmarshalText(text []byte) error {
	if text == nil {
		return nil
	}
	s := strings.TrimSpace(string(text))
	ev, err := gval.Full().NewEvaluable(s)
	if err != nil {
		return fmt.Errorf("invalid expression %q: %w", s, err)
	}
	e.Text = s
	e.Evaluable = ev
	return nil
}

// StringExpr consists of an expression that can either be evaluated to a
// string or an evaluable. It implements the encoding.TextUnmarshaler
// interface.
type StringExpr struct {
	Text string
	Expr *Expr
}

func (se *StringExpr) UnmarshalText(text []byte) error {
	if text == nil {
		return nil
	}

	s := strings.TrimSpace(string(text))
	se.Text = s

	// if s is alphanumeric => given expression needs to be treated as text
	// else => given expression needs to be treated as an evaluable
	isAlphanum := true
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			isAlphanum = false
			break
		}
	}
	if isAlphanum {
		return nil
	}

	var expr Expr
	if err := expr.UnmarshalText([]byte(s)); err != nil {
		return err
	}
	se.Expr = &expr

	return nil
}

// Template consists of a parsed text template that can be executed at runtime.
// It implements the encoding.TextUnmarshaler interface.
type Template struct {
	template.Template
}

func (t *Template) UnmarshalText(text []byte) error {
	if text == nil {
		return nil
	}

	s := strings.TrimSpace(string(text))
	tmpl, err := template.New("").Parse(s)
	if err != nil {
		return fmt.Errorf("failed to parse template `%s`: %w", s, err)
	}

	t.Template = *tmpl
	return nil
}
