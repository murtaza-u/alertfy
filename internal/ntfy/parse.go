package ntfy

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/murtaza-u/amify/internal/alert"
	"github.com/murtaza-u/amify/internal/conf"
)

// Parser is defines a Parse method to process the alert and extract relevant
// data.
type Parser interface {
	Parse(context.Context, alert.Alert) *Data
}

// NewParser creates a new instance of a parser. The returned parser will use
// the provided configuration to parse alerts.
func NewParser(conf conf.Ntfy) Parser {
	return parser{
		conf: conf,
	}
}

// Parse processes the provided alert and extracts various pieces of data. If
// any step in the process fails, appropriate error messages are logged, and
// the method returns nil.
func (p parser) Parse(ctx context.Context, alert alert.Alert) *Data {
	title, err := p.Title(alert)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"failed to parse title for notification. Aborting",
			slog.String("fingerprint", alert.Fingerprint),
			slog.String("error", err.Error()),
		)
		return nil
	}

	desc, err := p.Description(alert)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"failed to parse description for notification. Aborting",
			slog.String("fingerprint", alert.Fingerprint),
			slog.String("error", err.Error()),
		)
		return nil
	}

	topic, err := p.Topic(ctx, alert)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"failed to parse notification topic. Aborting",
			slog.String("fingerprint", alert.Fingerprint),
			slog.String("error", err.Error()),
		)
	}

	priority, err := p.Priority(ctx, alert)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"failed to parse notification priority. Defaulting to `default`",
			slog.String("fingerprint", alert.Fingerprint),
			slog.String("error", err.Error()),
		)
		priority = defaultPriority
	}

	tags := p.Tags(ctx, alert)

	url, err := p.URL(topic)
	if err != nil {
		slog.LogAttrs(
			ctx,
			slog.LevelError,
			"failed to get ntfy url. Aborting",
			slog.String("fingerprint", alert.Fingerprint),
			slog.String("error", err.Error()),
		)
		return nil
	}

	// If the description is empty, send the title as the description so that
	// the ntfy app doesn't fall back to setting "triggered" as the
	// description.
	if desc == "" {
		desc = title
		title = ""
	}

	return &Data{
		URL:         url,
		Title:       title,
		Description: desc,
		Tags:        tags,
		Priority:    priority,
	}
}

// parser is the default Parser implememtation.
type parser struct {
	conf conf.Ntfy
}

// URL constructs the URL for a given topic by appending the topic to the base
// URL.
func (p parser) URL(topic string) (string, error) {
	base := p.conf.BaseURL
	url, err := url.JoinPath(base, topic)
	if err != nil {
		return "", fmt.Errorf("appending %q to %q: %w", topic, base, err)
	}
	return url, nil
}

// Title generates the title for the alert by executing the template stored in
// the configuration.
func (p parser) Title(alert alert.Alert) (string, error) {
	buf := new(bytes.Buffer)
	err := p.conf.Notification.Title.Execute(buf, alert)
	if err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}
	return buf.String(), nil
}

// Description generates the description for the alert by executing the
// template stored in the configuration.
func (p parser) Description(alert alert.Alert) (string, error) {
	buf := new(bytes.Buffer)
	err := p.conf.Notification.Description.Execute(buf, alert)
	if err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}
	return buf.String(), nil
}

// Topic extracts the topic from the alert by evaluating the topic expression
// defined in the configuration. If the topic expression is nil, it simply
// returns the topic text. Otherwise, it evaluates the expression and returns
// the result.
func (p parser) Topic(c context.Context, alert alert.Alert) (string, error) {
	topic := p.conf.Notification.Topic
	if topic.Expr == nil {
		return topic.Text, nil
	}
	out, err := topic.Expr.Evaluable.EvalString(c, alert)
	if err != nil {
		return "", fmt.Errorf("evaluating expression: %w", err)
	}
	return out, nil
}

// Priority extracts the priority from the alert by evaluating the priority
// expression defined in the configuration. If the priority text is not set, it
// defaults to a `defaultPriority`. Otherwise, it evaluates the expression and
// returns the result.
func (p parser) Priority(c context.Context, alert alert.Alert) (string, error) {
	priority := p.conf.Notification.Priority
	if priority.Text == "" {
		slog.LogAttrs(
			c,
			slog.LevelDebug,
			"ntfy.notification.priority not set. Defaulting to `default`",
			slog.String("fingerprint", alert.Fingerprint),
		)
		return defaultPriority, nil
	}
	out, err := priority.Expr.Evaluable.EvalString(c, alert)
	if err != nil {
		return "", fmt.Errorf("evaluating expression: %w", err)
	}
	return out, nil
}

// Tags constructs a comma-separated list of tags for the alert based on the
// tags defined in the configuration. Each tag is included if its condition
// evaluates to true or if no condition is specified.
func (p parser) Tags(c context.Context, alert alert.Alert) string {
	tags := p.conf.Notification.Tags
	if len(tags) == 0 {
		return ""
	}

	var stitched string
	for _, tag := range tags {
		if tag.Tag == "" {
			continue
		}
		if tag.Condition.Text == "" {
			stitched += fmt.Sprintf("%s,", tag.Tag)
			continue
		}
		include, err := tag.Condition.Evaluable.EvalBool(c, alert)
		if err != nil {
			slog.LogAttrs(
				c,
				slog.LevelError,
				"evaluating tag condition failed. Skipping tag",
				slog.String("error", err.Error()),
				slog.String("condition", tag.Condition.Text),
				slog.String("fingerprint", alert.Fingerprint),
			)
			continue
		}
		if !include {
			slog.LogAttrs(
				c,
				slog.LevelDebug,
				"tag condition evaluated to false",
				slog.String("condition", tag.Condition.Text),
				slog.String("fingerprint", alert.Fingerprint),
			)
			continue
		}
		stitched += fmt.Sprintf("%s,", tag.Tag)
	}

	return strings.TrimRight(stitched, ",")
}
