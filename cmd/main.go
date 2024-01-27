package main

import (
	container "github.com/kanrichan/tinycontainer"
)

func main() {
	if err := container.Start(); err != nil {
		panic(err)
	}
}
