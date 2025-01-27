package generator

import (
	"encoding/json"
	"fmt"
	"io"
	"nupdate/internal/detector"
	"os"
	"os/exec"
)

type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description,omitempty"`
	Scripts         map[string]string `json:"scripts,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	DevDependencies map[string]string `json:"devDependencies,omitempty"`
}

func RunNpmOutDated() error {
	// Create or truncate the outdated.json file
	outFile, err := os.Create("outdated.json")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return err
	}
	defer outFile.Close()

	// Run the command and redirect output
	cmd := exec.Command("npm", "outdated", "--json")
	cmd.Stdout = outFile   // Redirect stdout to the file
	cmd.Stderr = os.Stderr // Optionally, redirect stderr to the terminal

	_ = cmd.Run()

	fmt.Println("Outdated packages saved to outdated.json")
	return nil
}

func GenerateNewPackageJSON() error {
	packageFile, err := detector.DetectPackageJSON()
	if err != nil {
		fmt.Printf("Error opening package.json: %v\n", err)
		return err
	}
	defer packageFile.Close()

	var pkg PackageJSON
	packageData, _ := io.ReadAll(packageFile)
	err = json.Unmarshal(packageData, &pkg)
	if err != nil {
		fmt.Printf("Error parsing package.json: %v\n", err)
		return err
	}

	// Step 2: Read npm outdated --json output
	outdatedFile, err := detector.DetectOutDatedJSON()
	if err != nil {
		fmt.Printf("Error opening outdated.json: %v\n", err)
		return err
	}
	defer outdatedFile.Close()

	var outdated map[string]struct {
		Current string `json:"current"`
		Wanted  string `json:"wanted"`
		Latest  string `json:"latest"`
	}

	outdatedData, _ := io.ReadAll(outdatedFile)
	err = json.Unmarshal(outdatedData, &outdated)
	if err != nil {
		fmt.Printf("Error parsing outdated.json: %v\n", err)
		return err
	}

	// Step 3: Update versions in dependencies and devDependencies
	updateVersions := func(deps map[string]string) {
		for pkgName := range deps {
			if info, exists := outdated[pkgName]; exists {
				deps[pkgName] = fmt.Sprintf("^%s", info.Wanted)
			}
		}
	}

	updateVersions(pkg.Dependencies)
	updateVersions(pkg.DevDependencies)

	// Step 4: Write updated package.json to a new file
	newPackageFile, err := os.Create("new-package.json")
	if err != nil {
		fmt.Printf("Error creating new package.json: %v\n", err)
		return err
	}
	defer newPackageFile.Close()

	encoder := json.NewEncoder(newPackageFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(pkg)
	if err != nil {
		fmt.Printf("Error writing new package.json: %v\n", err)
		return err
	}

	fmt.Println("Updated package.json written to new-package.json")
	return nil
}
