package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackupCommandWasAdded(t *testing.T) {
	cmd, _, err := TarwsCmd.Find([]string{"backup"})

	require.NoError(t, err, "Expected no error when finding install command but got %+v", err)

	assert.NotNil(t, cmd, "Expected a real command but got nil")
}
