package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

const hostName string = "my-container"
const maxProcesses string = "30"

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("incorrect arguments supplied. hint: follow program with \"run <command>\"")
	}
}

func run() {
	fmt.Printf("Running %v as pid %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// CLONE_NEWUTS: new unix time sharing namespace (hostname)
	// CLONE_NEWPID: new processes namespace
	// CLONE_NEWNS: new mounts namespace
	// Unshareflags: do not share mounts of container back with the host
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
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

	// set up control group against container resources
	setCg()

	// set hostname for container
	must(syscall.Sethostname([]byte(hostName)))

	// change container's root directory
	must(syscall.Chroot("./root-file-system/"))
	must(syscall.Chdir("/"))

	// mount proc at proc, declaring it is a proc file system
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	must(cmd.Run())

	// clean up mounts
	must(syscall.Unmount("/proc", 0))
}

func setCg() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")

	os.Mkdir(filepath.Join(pids, "newcg"), 0755)

	// state rule to set max number of system processes within control group
	must(ioutil.WriteFile(filepath.Join(pids, "newcg/pids.max"), []byte(maxProcesses), 0700))

	// remove when container exits
	must(ioutil.WriteFile(filepath.Join(pids, "newcg/notify_on_release"), []byte("1"), 0700))

	// write current process id into file to identify it belongs to control group
	must(ioutil.WriteFile(filepath.Join(pids, "newcg/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
