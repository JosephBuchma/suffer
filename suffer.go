package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const (
	WAVPath = "/tmp/__suffer.wav"

	PrintDelay = 30 * time.Millisecond
	BitchTime  = time.Second*2 - PrintDelay
	BufferSize = 1024
)

var sufferCommand = func() *exec.Cmd {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("aplay", WAVPath)
	case "darwin":
		return exec.Command("afplay", WAVPath)
	}
	panic("unsupported OS")
}()

func Suffer() (bitchTime <-chan time.Time) {
	go func() {
		if _, err := os.Stat(WAVPath); os.IsNotExist(err) {
			b, _ := Asset("Suffer-Bitch.wav")
			err = ioutil.WriteFile(WAVPath, b, 777)
			if err != nil {
				fmt.Println("Can't suffer :( failed to write file to /tmp")
			}
		}

		if err := sufferCommand.Run(); err != nil {
			fmt.Printf("Can't suffer :( failed to play audio\n%s\n", err)
		}
	}()
	return time.After(BitchTime)
}

func main() {
	runtime.GOMAXPROCS(1)

	var interval <-chan time.Time

	if len(os.Args) < 2 {
		fmt.Println("\nsuffer usage: \n\n\tsuffer <command>")
		os.Exit(1)
	}
	args := os.Args[1:]
	cmd := exec.Command(args[0], args[1:]...)

	outw := bufio.NewWriter(os.Stdout)
	errw := bufio.NewWriter(os.Stderr)

	cmd.Stdout = outw
	cmd.Stderr = errw
	cmd.Stdin = os.Stdin
	done := make(chan error)
	go func() {
		err := cmd.Run()
		done <- err
	}()

	for {
		interval = time.After(PrintDelay)
		select {
		case e := <-done:
			if e != nil {
				<-Suffer()
				defer fmt.Println(e)
			}
			outw.Flush()
			errw.Flush()
			return

		case <-interval:
			outw.Flush()
		}
	}
}
