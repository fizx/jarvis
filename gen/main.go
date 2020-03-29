package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	fs := http.Dir("./template")
	err := vfsgen.Generate(fs, vfsgen.Options{
		Filename:     "generated/assets/assets.go",
		PackageName:  "assets",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
