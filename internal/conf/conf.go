package conf

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

var (
	defaultConfig     = "/etc/alertfy/config.yaml"
	defaultListenAddr = ":8080"
)

// New creates a configuration using the provided arguments, environment
// variables, and config file. The order of precedence for configuration values
// is: flag arguments > environment variables > config file.
func New(args ...string) (*C, error) {
	k := koanf.New(".")

	err := k.Load(confmap.Provider(map[string]any{
		"hook.auth.enable":              false,
		"hook.auth.username":            "",
		"hook.auth.password":            "",
		"hook.log.level":                "info",
		"hook.log.format":               "text",
		"hook.terminationGracePeriod":   time.Second * 60,
		"ntfy.baseUrl":                  "",
		"ntfy.auth.enable":              false,
		"ntfy.auth.username":            "",
		"ntfy.auth.password":            "",
		"ntfy.notification.topic":       StringExpr{},
		"ntfy.notification.priority":    StringExpr{Text: "default"},
		"ntfy.notification.tags":        []Tag{},
		"ntfy.notification.title":       nil,
		"ntfy.notification.description": nil,
	}, "."), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load default configuration: %w", err)
	}

	f := parseFlags(args)
	confF, err := f.GetString("conf")
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	err = k.Load(file.Provider(confF), yaml.Parser())
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	err = k.Load(env.Provider("ALERTFY_", ".", func(s string) string {
		s = strings.TrimPrefix(s, "ALERTFY_")
		s = strings.ToLower(s)
		s = strings.ReplaceAll(s, "_", ".")
		return s
	}), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load env variables: %w", err)
	}

	conf := new(C)
	err = k.Unmarshal("", conf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return conf, nil
}

func parseFlags(args []string) *flag.FlagSet {
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Print(f.FlagUsages())
		os.Exit(0)
	}
	f.String("conf", defaultConfig, "path to config file")
	f.Parse(args)
	return f
}
