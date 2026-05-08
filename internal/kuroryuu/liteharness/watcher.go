package liteharness

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	EnvAgentID = "LITEHARNESS_AGENT_ID"
	EnvRoot    = "LITEHARNESS_ROOT"

	DefaultInterval = 2 * time.Second
)

type Message struct {
	From string `json:"from"`
	Type string `json:"type"`
	Body string `json:"body"`
}

type Deliver func(context.Context, Message) error

type Watcher struct {
	AgentID  string
	Root     string
	Interval time.Duration
	Log      *slog.Logger
}

func FromEnv() (Watcher, bool, error) {
	agentID := strings.TrimSpace(os.Getenv(EnvAgentID))
	if agentID == "" {
		return Watcher{}, false, nil
	}

	root, err := Root()
	if err != nil {
		return Watcher{}, false, err
	}

	return Watcher{
		AgentID:  agentID,
		Root:     root,
		Interval: DefaultInterval,
	}, true, nil
}

func Root() (string, error) {
	if root := strings.TrimSpace(os.Getenv(EnvRoot)); root != "" {
		return root, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve LiteHarness home: %w", err)
	}
	return filepath.Join(home, ".liteharness"), nil
}

func (w Watcher) Run(ctx context.Context, deliver Deliver) error {
	if err := w.validate(deliver); err != nil {
		return err
	}
	if err := w.Poll(ctx, deliver); err != nil {
		w.logger().Warn("LiteHarness inbox poll failed", "error", err)
	}

	tick := time.NewTicker(w.interval())
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
			if err := w.Poll(ctx, deliver); err != nil {
				w.logger().Warn("LiteHarness inbox poll failed", "error", err)
			}
		}
	}
}

func (w Watcher) Poll(ctx context.Context, deliver Deliver) error {
	if err := w.validate(deliver); err != nil {
		return err
	}
	if err := w.ensure(); err != nil {
		return err
	}

	entries, err := os.ReadDir(w.newDir())
	if err != nil {
		return fmt.Errorf("read LiteHarness inbox: %w", err)
	}

	var result error
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(w.newDir(), entry.Name())
		msg, err := read(path)
		if err != nil {
			w.logger().Warn("Skipping malformed LiteHarness message", "path", path, "error", err)
			result = errors.Join(result, err, w.moveToCur(path))
			continue
		}
		if strings.TrimSpace(msg.Body) == "" {
			w.logger().Debug("Skipping empty LiteHarness message", "path", path)
			result = errors.Join(result, w.moveToCur(path))
			continue
		}
		if err := deliver(ctx, msg); err != nil {
			result = errors.Join(result, fmt.Errorf("deliver LiteHarness message %s: %w", entry.Name(), err))
			continue
		}
		result = errors.Join(result, w.moveToCur(path))
	}
	return result
}

func read(path string) (Message, error) {
	bts, err := os.ReadFile(path)
	if err != nil {
		return Message{}, fmt.Errorf("read LiteHarness message: %w", err)
	}

	var msg Message
	if err := json.Unmarshal(bts, &msg); err != nil {
		return Message{}, fmt.Errorf("parse LiteHarness message: %w", err)
	}
	return msg, nil
}

func (w Watcher) validate(deliver Deliver) error {
	if strings.TrimSpace(w.AgentID) == "" {
		return errors.New("LiteHarness agent id is required")
	}
	if strings.TrimSpace(w.Root) == "" {
		return errors.New("LiteHarness root is required")
	}
	if deliver == nil {
		return errors.New("LiteHarness deliver callback is required")
	}
	return nil
}

func (w Watcher) ensure() error {
	for _, dir := range []string{w.newDir(), w.curDir()} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return fmt.Errorf("create LiteHarness mailbox directory %s: %w", dir, err)
		}
	}
	return nil
}

func (w Watcher) moveToCur(path string) error {
	target := filepath.Join(w.curDir(), filepath.Base(path))
	if err := os.Rename(path, target); err == nil {
		return nil
	}

	fallback := filepath.Join(w.curDir(), fmt.Sprintf("%s.%d", filepath.Base(path), time.Now().UnixNano()))
	if err := os.Rename(path, fallback); err != nil {
		return fmt.Errorf("move LiteHarness message to cur: %w", err)
	}
	return nil
}

func (w Watcher) newDir() string {
	return filepath.Join(w.mailboxDir(), "new")
}

func (w Watcher) curDir() string {
	return filepath.Join(w.mailboxDir(), "cur")
}

func (w Watcher) mailboxDir() string {
	return filepath.Join(w.Root, "mailboxes", w.AgentID)
}

func (w Watcher) interval() time.Duration {
	if w.Interval > 0 {
		return w.Interval
	}
	return DefaultInterval
}

func (w Watcher) logger() *slog.Logger {
	if w.Log != nil {
		return w.Log
	}
	return slog.Default()
}
