package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
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
	flag.Var(&label, "label", "regexp of label container")
	flag.Parse()
	logrus.Infof("args : %v", flag.Args())
	if label.set {
		logrus.Infof("label: %s", label.value)
	}
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		os.Exit(2)
	}()
	wait()
}

func wait() {
	done := make(chan (bool))
	time.AfterFunc(5*time.Second, func() {
		logrus.Info("exiting")
		done <- true
	})
	<-done
	os.Exit(0)
}
