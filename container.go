package container

import (
	"flag"
	"os"
	"os/exec"
	"syscall"
)

const SUB_PROCESS = "tiny-container"

var (
	HOSTNAME string
	COMMAND  string
)

func init() {
	flag.StringVar(&HOSTNAME, "hostname", "tiny-container", "")
	flag.StringVar(&COMMAND, "exec", "/bin/sh", "")

	if os.Args[0] == SUB_PROCESS {
		if err := Init(); err != nil {
			panic(err)
		}
		if err := Exec(COMMAND); err != nil {
			panic(err)
		}
		os.Exit(0)
	}
}

func Start() error {
	cmd := exec.Cmd{
		Path: "/proc/self/exe",
		Args: os.Args,
		SysProcAttr: &syscall.SysProcAttr{
			Pdeathsig:  syscall.SIGTERM,
			Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS,
		},
	}
	cmd.Args[0] = SUB_PROCESS

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Init() error {
	if err := syscall.Sethostname([]byte(HOSTNAME)); err != nil {
		return err
	}
	return nil
}

func Exec(command string) error {
	cmd := exec.Command(command)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{"PS1=-[" + HOSTNAME + "]- # "}

	return cmd.Run()
}
