package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		panic("unknown comman")
	}
}
func parent() {
	fmt.Printf("running parent pid:%d %v\n", os.Getpid(), os.Args[2:])

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}

	wrapError(cmd.Run())
}

func child() {
	fmt.Printf("running child pid:%d %v\n", os.Getpid(), os.Args[2:])
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	wrapError(syscall.Chroot("/home/kubo/os/ubuntu"))
	wrapError(syscall.Mount("proc", "/proc", "proc", 0, ""))

	wrapError(cmd.Run())

	wrapError(syscall.Unmount("/proc", 0))

}

func wrapError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
