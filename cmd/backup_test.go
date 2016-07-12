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

func TestBackupReturnsErrorIfNoPathSpecified(t *testing.T) {
	err := backup(nil, []string{})

	require.Error(t, err, "Expected error when no path was specified")
}

func TestBackupReturnsErrorIfPathDoesNotExist(t *testing.T) {
	err := backup(nil, []string{"/does/not/exist"})

	require.Error(t, err, "Expected error when specified path does not exist")
}
