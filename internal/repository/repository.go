package repository

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/thekashifmalik/rincr/internal"
)

type Repository struct {
	path string
}

func NewRepository(path string) *Repository {
	return &Repository{
		path: path,
	}
}

func (r *Repository) IsRemote() bool {
	return len(r.parse()) == 2
}

func (r *Repository) GetHost() string {
	if r.IsRemote() {
		return r.parse()[0]
	} else {
		// TODO: Consider returning an error in this case; we should never be using the host when the repository is
		// local.
		return ""
	}
}

func (r *Repository) GetPath() string {
	if r.IsRemote() {
		return r.parse()[1]
	} else {
		return r.path
	}
}

func (r *Repository) parse() []string {
	return strings.SplitN(r.path, ":", 2)
}

func (r *Repository) PathExists(path string) bool {
	filePath := r.GetPath() + "/" + path
	if r.IsRemote() {
		quotedPath := `"` + filePath + `"`
		err := exec.Command("ssh", r.GetHost(), "stat", quotedPath).Run()
		return err == nil

	} else {
		_, err := os.Stat(filePath)
		if err == nil {
			return true
		}
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
		return false
	}
}

func (r *Repository) Exists() bool {
	return r.PathExists(internal.BACKUPS_DIR)
}

// func (r *Repository) PathExistsHistorical(path string, when time.Time) bool {
// 	// historicalPath := 0
// 	return r.PathExists(path)
// }
