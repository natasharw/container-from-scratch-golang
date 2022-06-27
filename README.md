# Container From Scratch
Implementing a simple container from scratch in Go


---

## Overview
This small program demonstrates a simplified version of how a container can be constructed inside a host operating system.

It was guided by Liz Rice’s amazing [GOTO 2018 talk](https://www.youtube.com/watch?v=8fi7uSYlOdc).


## Prerequisites
* Linux machine for testing purpose only (e.g. temporary VM)
* Go installed
* Operating with root privileges. Hint: you may need to use `sudo bash`


## How to run
#### 1. Clone repository
```
$ git clone https://github.com/natasharw/container-from-scratch-golang.git
```
#### 2. Extract an image of a root file system for the container into the empty placeholder directory
Good choices for distribution would be something lightweight. Here is an example using a minimal Ubuntu distribution:
```
$ cd container-from-scratch-golang/root-file-system
```
```
$ curl -o ubuntu-minimal.tar.gz http://cloud-images.ubuntu.com/minimal/releases/bionic/release/ubuntu-18.04-minimal-cloudimg-amd64-root.tar.xz
```
```
$ tar xvf ubuntu-minimal.tar.qz
```
```
$ rm ubuntu-minimal.tar.qz
```
#### 3. Move back to main directory
```
$ cd ..
```
#### 4. Run the container, specifying a desired process
For example, to run a bash session:
```
$ go run main.go run /bin/bash
```
#### 5. Success! You are now running a simple container

#### 6. From inside a container bash session, test some commands:
  - `$ ps` - check that the process ids are different to those visible on the host
  - `$ ls /` - check that the root file system is different to that of the host
  - `$ hostname` - check that the hostname is set is different to that of the host
  
#### 7. Stop running the container
```
$ exit
```

## What is happening?

* A container with its own isolated namespaces for hostname, process IDs and mounts are set up by `syscall.CLONE_NEWUTS`, `syscall.CLONE_NEWPID` and `syscall.CLONE_NEWNS`
* `syscall.Sethostname()` sets a different hostname for the container
* The container points towards the new root file system with `syscall.Chroot()` and `syscall.Chdir()`
* A new `proc/` folder is mounted with `syscall.Mount()`, allowing process IDs to be isolated from those of the host operating system
  * `$ ps` will now only show the container’s processes
* A control group for the container is created with the custom function `setCg()`
  * The arbitrary rule used is to limit the maximum number of processes that the container can run. A control group could be used to limit memory or CPU usage instead if you wanted
