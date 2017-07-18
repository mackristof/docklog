package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mackristof/docklog/tools"
)

type stringFlag struct {
	set   bool
	value string
}

func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *stringFlag) String() string {
	return sf.value
}

func main() {
	var label stringFlag
	var name stringFlag
	flag.Var(&label, "label", "regexp of label container")
	flag.Var(&name, "name", "regexp of name container")
	flag.Parse()
	fmt.Printf("args : %v\n", flag.Args())
	if label.set {
		fmt.Printf("label: %s\n", label.value)
	}
	if name.set {
		fmt.Printf("name: %s\n", name.value)
	}
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		os.Exit(2)
	}()
	docker := tools.NewDocker()
	docker.ListContainers(name.value, label.value)
	wait()
}

func wait() {
	done := make(chan (bool))
	time.AfterFunc(5*time.Second, func() {
		fmt.Println("exiting")
		done <- true
	})
	<-done
	os.Exit(0)
}
