package lnd

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

// Note: These are integration tests. They only work if an lnd instance is running.
func TestAddInvoice(t *testing.T) {
	client, err := New()
	require.NoError(t, err)

	_, err = client.AddInvoice(context.Background(), 10, 3600, "test")
	require.NoError(t, err)
}

func TestGetInfo(t *testing.T) {
	client, err := New()
	require.NoError(t, err)

	_, err = client.GetInfo(context.Background())
	require.NoError(t, err)
}