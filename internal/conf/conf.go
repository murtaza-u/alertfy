package conf

import (
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

var (
	defaultConfig     = "/etc/amify/config.yaml"
	defaultListenAddr = ":8080"
)

// New creates a configuration using the provided arguments, environment
// variables, and config file. The order of precedence for configuration values
// is: flag arguments > environment variables > config file.
func New(args ...string) (*C, error) {
	k := koanf.New(".")

	f := parseFlags(args)
	confF, err := f.GetString("conf")
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	err = k.Load(file.Provider(confF), yaml.Parser())
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	err = k.Load(env.Provider("AMIFY_", ".", func(s string) string {
		s = strings.TrimPrefix(s, "AMIFY_")
		s = strings.ToLower(s)
		s = strings.ReplaceAll(s, "_", ".")
		return s
	}), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load env variables: %w", err)
	}

	err = k.Load(posflag.Provider(f, ".", k), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load post flags: %w", err)
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
	f.String("hook.listenAddr", defaultListenAddr, "address for the hook to listen on")
	f.String("hook.auth.username", "", "http basic auth username")
	f.String("hook.auth.password", "", "http basic auth password")
	f.Parse(args)
	return f
}
