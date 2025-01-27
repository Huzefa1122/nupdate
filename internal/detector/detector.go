package detector

import (
	"fmt"
	"os"
	"path"
)

func DetectPackageJSON() (*os.File, error) {
	currentDir, err := os.Getwd()
	var file *os.File
	if err != nil {
		return nil, err
	}
    fmt.Println(currentDir)
	file, err = os.Open(path.Join(currentDir, "package.json"))
	if err != nil {
		return nil, fmt.Errorf("Failed to detect package.json file, %w",err)
	}
	return file, nil
}
func DetectOutDatedJSON() (*os.File, error) {
	currentDir, err := os.Getwd()
	var file *os.File
	if err != nil {
		return nil, err
	}
	file, err = os.Open(path.Join(currentDir, "outdated.json"))
	if err != nil {
		return nil, fmt.Errorf("Failed to detect package.json file, Are you root directory of the project")
	}
	return file, nil
}
