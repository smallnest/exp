package sync

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSafeClose(t *testing.T) {
	var conn net.Conn
	justClosed, err := SafeClose(conn)
	assert.False(t, justClosed)
	assert.NoError(t, err)

	conn, err = net.Dial("tcp", "www.bing.com:80")
	require.NoError(t, err)
	require.NotNil(t, conn)

	err = conn.Close()
	require.NoError(t, err)

	justClosed, err = SafeClose(conn)
	assert.True(t, justClosed)
	assert.NoError(t, err)
}
