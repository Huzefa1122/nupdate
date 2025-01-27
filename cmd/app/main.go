package main

import (
	"nupdate/internal/generator"
)

func main() {
	err := generator.RunNpmOutDated()
	if err != nil {
		return
	}
	err = generator.GenerateNewPackageJSON()
	if err != nil {
		return
	}

}
