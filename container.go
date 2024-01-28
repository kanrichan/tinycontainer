package container

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const SUB_PROCESS = "tiny-container"

type Container struct {
	Root     string
	HostName string
	Command  string
}

func Start() error {
	cmd := exec.Cmd{
		Path: "/proc/self/exe",
		Args: os.Args,
		SysProcAttr: &syscall.SysProcAttr{
			Pdeathsig: syscall.SIGTERM,
			Cloneflags: syscall.CLONE_NEWNS |
				syscall.CLONE_NEWUTS |
				syscall.CLONE_NEWIPC |
				syscall.CLONE_NEWUSER |
				syscall.CLONE_NEWPID,
			UidMappings: []syscall.SysProcIDMap{
				{
					ContainerID: 0,
					HostID:      os.Getuid(),
					Size:        1,
				},
			},
			GidMappings: []syscall.SysProcIDMap{
				{
					ContainerID: 0,
					HostID:      os.Getgid(),
					Size:        1,
				},
			},
		},
	}
	cmd.Args[0] = SUB_PROCESS

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (c *Container) Init() error {
	proc := filepath.Join(c.Root, "/proc")
	os.MkdirAll(proc, 0755)
	if err := syscall.Mount("proc", proc, "proc", 0, ""); err != nil {
		return err
	}

	putold := filepath.Join(c.Root, "/.pivot_root")
	if err := os.Mkdir(putold, 0700); err != nil {
		return err
	}
	if err := syscall.Mount(c.Root, c.Root, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err
	}
	if err := syscall.PivotRoot(c.Root, putold); err != nil {
		return err
	}
	if err := os.Chdir("/"); err != nil {
		return err
	}
	putold = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(putold, syscall.MNT_DETACH); err != nil {
		return err
	}
	if err := os.Remove(putold); err != nil {
		return err
	}

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
