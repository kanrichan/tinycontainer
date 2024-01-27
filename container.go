package container

import (
	"os"
	"os/exec"
	"syscall"
)

const SUB_PROCESS = "tiny-container"

type Container struct {
	HostName string
	Command  string
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

func (c *Container) Init() error {
	if err := syscall.Sethostname([]byte(c.HostName)); err != nil {
		return err
	}
	return nil
}

func (c *Container) Exec() error {
	cmd := exec.Command(c.Command)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{"PS1=-[" + c.HostName + "]- # "}

	return cmd.Run()
}
