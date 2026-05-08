package liteharness

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromEnvDisabledWithoutAgentID(t *testing.T) {
	t.Setenv(EnvAgentID, "")

	_, enabled, err := FromEnv()
	require.NoError(t, err)
	require.False(t, enabled)
}

func TestFromEnvUsesConfiguredRoot(t *testing.T) {
	root := t.TempDir()
	t.Setenv(EnvAgentID, "crush-agent")
	t.Setenv(EnvRoot, root)

	watcher, enabled, err := FromEnv()
	require.NoError(t, err)
	require.True(t, enabled)
	require.Equal(t, "crush-agent", watcher.AgentID)
	require.Equal(t, root, watcher.Root)
}

func TestPollDeliversAndMovesMessage(t *testing.T) {
	t.Parallel()

	watcher := testWatcher(t)
	path := writeMessage(t, watcher, "001.json", `{"from":"sentinel","type":"notification","body":"hello crush","extra":true}`)

	var got []Message
	err := watcher.Poll(context.Background(), func(_ context.Context, msg Message) error {
		got = append(got, msg)
		return nil
	})

	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "sentinel", got[0].From)
	require.Equal(t, "notification", got[0].Type)
	require.Equal(t, "hello crush", got[0].Body)
	require.NoFileExists(t, path)
	require.FileExists(t, filepath.Join(watcher.curDir(), "001.json"))
}

func TestPollKeepsMessageWhenDeliveryFails(t *testing.T) {
	t.Parallel()

	watcher := testWatcher(t)
	path := writeMessage(t, watcher, "001.json", `{"body":"retry me"}`)
	want := errors.New("not ready")

	err := watcher.Poll(context.Background(), func(context.Context, Message) error {
		return want
	})

	require.ErrorIs(t, err, want)
	require.FileExists(t, path)
	require.NoFileExists(t, filepath.Join(watcher.curDir(), "001.json"))
}

func TestPollQuarantinesMalformedMessage(t *testing.T) {
	t.Parallel()

	watcher := testWatcher(t)
	path := writeMessage(t, watcher, "bad.json", `{bad json`)

	err := watcher.Poll(context.Background(), func(context.Context, Message) error {
		t.Fatal("malformed messages should not be delivered")
		return nil
	})

	require.Error(t, err)
	require.NoFileExists(t, path)
	require.FileExists(t, filepath.Join(watcher.curDir(), "bad.json"))
}

func testWatcher(t *testing.T) Watcher {
	t.Helper()

	return Watcher{
		AgentID: "agent",
		Root:    t.TempDir(),
	}
}

func writeMessage(t *testing.T, watcher Watcher, name string, body string) string {
	t.Helper()

	require.NoError(t, os.MkdirAll(watcher.newDir(), 0o700))
	require.NoError(t, os.MkdirAll(watcher.curDir(), 0o700))

	path := filepath.Join(watcher.newDir(), name)
	require.NoError(t, os.WriteFile(path, []byte(body), 0o600))
	return path
}
