package conf

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/PaesslerAG/gval"
)

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
