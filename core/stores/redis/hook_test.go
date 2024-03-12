package redis

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"testing"
	"time"

	red "github.com/go-redis/redis/v8"
	"github.com/sjclijie/go-zero/core/breaker"
	"github.com/sjclijie/go-zero/core/logx/logtest"
	"github.com/stretchr/testify/assert"
)

func TestHookProcessCase3(t *testing.T) {
	writer := log.Writer()
	var buf strings.Builder
	log.SetOutput(&buf)
	defer log.SetOutput(writer)

	assert.Nil(t, durationHook.AfterProcess(context.Background(), red.NewCmd(context.Background())))
	assert.True(t, buf.Len() == 0)
}

func TestHookProcessCase4(t *testing.T) {
	writer := log.Writer()
	var buf strings.Builder
	log.SetOutput(&buf)
	defer log.SetOutput(writer)

	ctx := context.WithValue(context.Background(), startTimeKey, "foo")
	assert.Nil(t, durationHook.AfterProcess(ctx, red.NewCmd(context.Background())))
	assert.True(t, buf.Len() == 0)
}

func TestHookProcessPipelineCase3(t *testing.T) {
	w := logtest.NewCollector(t)

	assert.Nil(t, durationHook.AfterProcessPipeline(context.Background(), []red.Cmder{
		red.NewCmd(context.Background()),
	}))
	assert.True(t, len(w.String()) == 0)
}

func TestHookProcessPipelineCase4(t *testing.T) {
	w := logtest.NewCollector(t)

	ctx := context.WithValue(context.Background(), startTimeKey, "foo")
	assert.Nil(t, durationHook.AfterProcessPipeline(ctx, []red.Cmder{
		red.NewCmd(context.Background()),
	}))
	assert.True(t, len(w.String()) == 0)
}

func TestHookProcessPipelineCase5(t *testing.T) {
	writer := log.Writer()
	var buf strings.Builder
	log.SetOutput(&buf)
	defer log.SetOutput(writer)

	ctx := context.WithValue(context.Background(), startTimeKey, "foo")
	assert.Nil(t, durationHook.AfterProcessPipeline(ctx, []red.Cmder{
		red.NewCmd(context.Background()),
	}))
	assert.True(t, buf.Len() == 0)
}

func TestLogDuration(t *testing.T) {
	w := logtest.NewCollector(t)

	logDuration(context.Background(), []red.Cmder{
		red.NewCmd(context.Background(), "get", "foo"),
	}, 1*time.Second)
	assert.True(t, strings.Contains(w.String(), "get foo"))

	logDuration(context.Background(), []red.Cmder{
		red.NewCmd(context.Background(), "get", "foo"),
		red.NewCmd(context.Background(), "set", "bar", 0),
	}, 1*time.Second)
	assert.True(t, strings.Contains(w.String(), `get foo\nset bar 0`))
}

func TestFormatError(t *testing.T) {
	// Test case: err is OpError
	err := &net.OpError{
		Err: mockOpError{},
	}
	assert.Equal(t, "timeout", formatError(err))

	// Test case: err is nil
	assert.Equal(t, "", formatError(nil))

	// Test case: err is red.Nil
	assert.Equal(t, "", formatError(red.Nil))

	// Test case: err is io.EOF
	assert.Equal(t, "eof", formatError(io.EOF))

	// Test case: err is context.DeadlineExceeded
	assert.Equal(t, "context deadline", formatError(context.DeadlineExceeded))

	// Test case: err is breaker.ErrServiceUnavailable
	assert.Equal(t, "breaker", formatError(breaker.ErrServiceUnavailable))

	// Test case: err is unknown
	assert.Equal(t, "unexpected error", formatError(errors.New("some error")))
}

type mockOpError struct {
}

func (mockOpError) Error() string {
	return "mock error"
}

func (mockOpError) Timeout() bool {
	return true
}
