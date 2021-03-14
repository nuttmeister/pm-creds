package file

import (
	"fmt"
	"os"
)

var stat = os.Stat

// CheckFilesExists returns true if the any of the files already exist and false
// if they don't. If there is an error stating the file false and error is returned
// but it can't be used to know if the file exists or not.
// Returns bool and error.
func CheckFilesExists(files []string) (bool, error) {
	exists := false
	for _, file := range files {
		if _, err := stat(file); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return false, fmt.Errorf("couldn't stat file %q. %w", file, err)
		}
		exists = true
	}
	return exists, nil
}
