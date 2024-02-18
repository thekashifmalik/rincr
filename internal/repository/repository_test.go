package repository_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thekashifmalik/rincr/internal/repository"
)

func TestNewRepository(t *testing.T) {
	repo := repository.NewRepository("test-path")
	require.NotNil(t, repo)
}

func TestIsRemoteLocal(t *testing.T) {
	repo := repository.NewRepository("test-path")
	require.False(t, repo.IsRemote())
}

func TestIsRemoteRemote(t *testing.T) {
	repo := repository.NewRepository("test-server:test-path")
	require.True(t, repo.IsRemote())
}

func TestGetHostLocal(t *testing.T) {
	repo := repository.NewRepository("test-path")
	require.Equal(t, repo.GetHost(), "")
}
func TestGetHostRemote(t *testing.T) {
	repo := repository.NewRepository("test-server:test-path")
	require.Equal(t, repo.GetHost(), "test-server")
}

func TestGetPathLocal(t *testing.T) {
	repo := repository.NewRepository("test-path")
	require.Equal(t, repo.GetPath(), "test-path")
}
func TestGetPathRemote(t *testing.T) {
	repo := repository.NewRepository("test-server:test-path")
	require.Equal(t, repo.GetPath(), "test-path")
}
