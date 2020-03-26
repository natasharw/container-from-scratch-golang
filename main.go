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
		run()
	case "child":
		child()
	default:
		panic("incorrect arguments supplied. hint: must be followed by \"run <command>\"")
	}
}

func run() {
	fmt.Printf("Running %v as pid %d\n", os.Args[2:], os.Getpid())

	//take command line arguments as input
	//creates child process to set specifics of container
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//create new hostname, process ids, and mounts namespaces
	//CLONE_NEWUTS = new unix time sharing namespace
	//CLONE_NEWPID = new process ids namespaces
	//CLONE_NEWNS = new namespace for mounts
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		//do not share mounts of container back with the host
		Unshareflags: syscall.CLONE_NEWNS,
	}
	must(cmd.Run())
}

func child() {
	fmt.Printf("Running %v as pid %d in new namespace\n", os.Args[2:], os.Getpid())

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//set up a control group to specify what resources the container can use
	cg()

	//set a new hostname on creation of the container
	must(syscall.Sethostname([]byte("my-very-own-container")))

	//change container's root directory to be the desired file system
	must(syscall.Chroot("/path/to/your/my-new-fs"))
	must(syscall.Chdir("/"))

	//mount proc at proc, declaring it is a proc file system
	//this allows process ids to be isolated
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	must(cmd.Run())

	//clean up mounts
	must(syscall.Unmount("/proc", 0))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
