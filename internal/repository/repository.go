package repository

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

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

func (r *Repository) GetFullPath() string {
	return r.path
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

func (r *Repository) GetBackupTimes() ([]time.Time, error) {
	backupsDirPath := r.GetPath() + internal.BACKUPS_DIR_PATH
	files := []string{}
	if r.IsRemote() {
		_filesRaw, err := exec.Command("ssh", r.GetHost(), "ls", "-A", backupsDirPath).Output()
		if err != nil {
			return nil, err
		}
		_files := strings.Split(string(_filesRaw), "\n")
		for _, file := range _files {
			if file != "" {
				files = append(files, file)
			}
		}
	} else {
		entries, err := os.ReadDir(backupsDirPath)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			files = append(files, entry.Name())
		}
	}
	backupTimes := []time.Time{}
	for _, file := range files {
		if file != "last" {
			backupTime, err := time.ParseInLocation(internal.TIME_FORMAT, file, time.Local)
			if err != nil {
				return nil, err
			}
			backupTimes = append(backupTimes, backupTime)
		}
	}
	return backupTimes, nil
}

func (r *Repository) DeleteBackupsByTime(backupTimes []time.Time) {
	for _, backupTime := range backupTimes {
		fmt.Printf("deleting backup: %v\n", backupTime)
		// TODO: Handle any errors here
		backupPath := fmt.Sprintf("%v/%v/%v", r.GetPath(), internal.BACKUPS_DIR, backupTime.Format(internal.TIME_FORMAT))
		if r.IsRemote() {
			exec.Command("ssh", r.GetHost(), "rm", "-rf", backupPath).Run()
		} else {
			os.RemoveAll(backupPath)
		}
	}
}
