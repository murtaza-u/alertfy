package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateListenAddr(t *testing.T) {
	a := assert.New(t)
	inputs := map[string]bool{
		":8080":                true,
		"127.0.0.1:8080":       true,
		"0.0.0.0:8080":         true,
		"127.0.0.1":            false,
		"localhost:8080":       false,
		"foo:8080":             false,
		":70000":               false,
		"300.300.300.300:8080": false,
	}
	for input, isValid := range inputs {
		err := validateListenAddr(input)
		if isValid {
			a.NoErrorf(err, "INPUT=%s", input)
			continue
		}
		a.Errorf(err, "INPUT=%s", input)
	}
}

func TestValidateAuth(t *testing.T) {
	type input struct {
		enable   bool
		username string
		password string
		isValid  bool
	}
	a := assert.New(t)
	inputs := []input{
		{enable: true, username: "foo", password: "bar", isValid: true},
		{enable: false, username: "foo", password: "bar", isValid: true},
		{enable: false, username: "", password: "bar", isValid: true},
		{enable: false, username: "foo", password: "", isValid: true},
		{enable: true, username: "", password: "", isValid: false},
		{enable: true, username: "foo", password: "", isValid: false},
		{enable: true, username: "", password: "bar", isValid: false},
	}
	for _, i := range inputs {
		err := validateAuth(Auth{
			Enable:   i.enable,
			Username: i.username,
			Password: i.password,
		})
		if i.isValid {
			a.NoErrorf(err, "enable=%v,username=%s,password=%s",
				i.enable, i.username, i.password)
			continue
		}
		a.Errorf(err, "enable=%v,username=%s,password=%s",
			i.enable, i.username, i.password)
	}
}

func TestValidateLogLevel(t *testing.T) {
	a := assert.New(t)
	inputs := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"DEBUG": false,
		"foo":   false,
		"bar":   false,
	}
	for input, isValid := range inputs {
		err := validateLogLevel(input)
		if isValid {
			a.NoErrorf(err, "INPUT=%s", input)
			continue
		}
		a.Errorf(err, "INPUT=%s", input)
	}
}

func TestValidateLogFormat(t *testing.T) {
	a := assert.New(t)
	inputs := map[string]bool{
		"text": true,
		"json": true,
		"TEXT": false,
		"JSON": false,
		"foo":  false,
		"bar":  false,
	}
	for input, isValid := range inputs {
		err := validateLogFormat(input)
		if isValid {
			a.NoErrorf(err, "INPUT=%s", input)
			continue
		}
		a.Errorf(err, "INPUT=%s", input)
	}
}
