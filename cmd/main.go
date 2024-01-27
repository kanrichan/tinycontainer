package main

import (
	"flag"
	"os"

	container "github.com/kanrichan/tinycontainer"
)

var (
	HOSTNAME string
	COMMAND  string
)

func init() {
	flag.StringVar(&HOSTNAME, "hostname", "tiny-container", "")
	flag.StringVar(&COMMAND, "exec", "/bin/sh", "")

	c := container.Container{
		HostName: HOSTNAME,
		Command:  COMMAND,
	}

	if os.Args[0] == container.SUB_PROCESS {
		if err := c.Init(); err != nil {
			panic(err)
		}
		if err := c.Exec(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}
}

func main() {
	if err := container.Start(); err != nil {
		panic(err)
	}
}
