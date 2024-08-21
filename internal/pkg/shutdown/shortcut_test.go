package shutdown //nolint:testpackage // test package

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGlobalStopper_Global(t *testing.T) { //nolint:paralleltest // cause global obj
	global := Global()
	require.NotNil(t, global)
}

func TestGlobalStopper_Register(t *testing.T) { //nolint:paralleltest // cause global obj
	mock := &mockStopFn{}
	Register(mock)

	require.Len(t, globalShutdown.stops, 1)
}

func TestGlobalStopper_RegisterFn(t *testing.T) { //nolint:paralleltest // cause global obj
	var called bool

	RegisterFn(func() { called = true })

	globalShutdown.Stop()

	require.True(t, called)
}

func TestGlobalStopper_Wait(t *testing.T) { //nolint:paralleltest // cause global obj
	mock := &mockStopFn{}
	Register(mock)

	channel := make(chan struct{})
	go func() {
		close(channel)
	}()

	Wait(channel)

	require.True(t, mock.called)
}

func TestGlobalStopper_WaitSignal(t *testing.T) { //nolint:paralleltest // cause global obj
	mock := &mockStopFn{}
	Register(mock)

	signals := make(chan os.Signal, 1)
	go func() {
		signals <- os.Interrupt
	}()

	WaitSignal(signals)

	require.True(t, mock.called)
}

func TestGlobalStopper_Stop(t *testing.T) { //nolint:paralleltest // cause global obj
	mock := &mockStopFn{}
	Register(mock)

	Stop()

	require.True(t, mock.called)
}
