package restore

import (
	"fmt"

	"github.com/thekashifmalik/rincr/internal/repository"
)

func Restore(repository *repository.Repository, path string, output string) error {
	if !repository.Exists() {
		return fmt.Errorf("No repository found")
	}
	fmt.Println("repository.IsRemote()")
	fmt.Println(repository.IsRemote())
	fmt.Println("repository.PathExists(path)")
	existsLatest := repository.PathExists(path)
	fmt.Println(existsLatest)
	if !existsLatest {
		fmt.Println("repository.GetBackupTimes()")
		fmt.Println(repository.GetBackupTimes())
	}
	return nil
}
