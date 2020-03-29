package main

import (
	"fmt"

	"github.com/fizx/jarvis/generated/assets"
)

//go:generate go run gen/main.go

func main() {
	fmt.Println("hello world", assets.Assets)
}
