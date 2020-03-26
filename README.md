# Go container
A sandbox project to implement a simple container from scratch in Go

Guided by Liz Rice’s [GOTO 2018 talk](https://www.youtube.com/watch?v=8fi7uSYlOdc)

---

## What is the purpose of this?
To demonstrate the basics of how a container can be constructed inside a host operating system


## Prerequisites
* Must be run on a Linux machine
* Must have Go installed and PATH variable set up in `~/.profile`
```
export GOPATH=$HOME/go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```
* Ensure that you are operating with root privileges on the host. Hint: you may need to use `sudo bash`  


## How to run
#### 1. Decide what distribution you want for the container. For example Ubuntu or Alpine. Download a copy of its root file system into a new directory.
Example using Ubuntu Minimal:
```
mkdir my-new-fs
```
```
cd my-new-fs
```
```
curl -o ubuntu-minimal.tar.gz http://cloud-images.ubuntu.com/minimal/releases/bionic/release/ubuntu-18.04-minimal-cloudimg-amd64-root.tar.xz
```
```
tar xvf ubuntu-minimal.tar.qz
```
```
rm ubuntu-minimal.tar.qz
```
#### 2. Change the path in the Chroot command in `main.go` to the new directory:
Example: 
`must(syscall.Chroot("/path/to/your/my-new-fs"))` —> `must(syscall.Chroot(“/home/natasha/my-new-fs”))`
#### 3. Run the container, specifying a desired process
For example, to run a bash session:
```
go run main.go run /bin/bash
```
#### 4. Success! You are running your very own container

#### 5. Inside a bash session, test out commands inside the container
  - `ps` - check that the process ids are different to those visible on the host
  - `ls /` - check that the root file system is different to that of the host
  - `hostname` - check that the hostname is set is different to that of the host
  - `sleep 100`
  
#### 6. Use `exit` to exit out of container

---
## What is happening?

* A container with its own namespaces for hostname, process ids and mounts is set up by `syscall.CLONE_NEWUTS`, `syscall.CLONE_NEWPID` and `syscall.CLONE_NEWNS`
* A different hostname for the container is set with `syscall.Sethostname`
* The container is set to point towards a new root file system (whatever you decided to base it on) through `syscall.Chroot()` and `syscall.Chdir()`
* A new `proc/` folder is mounted with `syscall.Mount()` which allows process ids to be isolated from those of the host operating system. This allows `ps` command to only show the container’s processes
* A control group is set up with the custom function `cg()`. The arbitrary rule set up is to limit the max number of processes that the container can run to 30. This could be used to limit memory or CPU usage instead.
