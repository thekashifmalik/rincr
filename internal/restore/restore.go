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
	fmt.Println(repository.PathExists(path))
	return nil
}
