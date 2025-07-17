package repository

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/thekashifmalik/rincr/internal"
)

func (r *Repository) WriteLastFile(timestamp time.Time) error {
	formatted := timestamp.Format(internal.TIME_FORMAT)
	if !r.IsRemote() {
		f, err := os.Create(r.GetPath() + internal.LAST_FILE_PATH)
		if err != nil {
			return err
		}
		_, err = f.WriteString(formatted)
		if err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}
	} else {
		echoCmd := fmt.Sprintf(`"echo %v > %v"`, formatted, r.GetPath()+internal.LAST_FILE_PATH)
		cmd := exec.Command("ssh", r.GetHost(), "bash", "-c", echoCmd)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
