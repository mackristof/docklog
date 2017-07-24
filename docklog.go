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
	set    bool
	values []string
}

func (sf *stringFlag) Set(x string) error {
	sf.values = append(sf.values, x)
	sf.set = true
	return nil
}

func (sf *stringFlag) Strings() []string {
	return sf.values
}

func (sf *stringFlag) String() string {
	if len(sf.values) > 0 {
		return sf.values[0]
	}
	return ""
}

func main() {
	var label stringFlag
	var name stringFlag
	var certPath stringFlag
	flag.Var(&label, "label", "regexp of label container")
	flag.Var(&name, "name", "regexp of name container")
	flag.Var(&certPath, "certPath", "path where certificates are located")
	localFlag := flag.Bool("local", false, "access to local docker engine")
	remoteFlag := flag.Bool("remote", false, "access to remote docker engine")
	swarmFlag := flag.Bool("swarm", false, "access to docker swarm cluster")
	flag.Parse()
	fmt.Printf("url : %v\n", flag.Args())
	if label.set {
		fmt.Printf("label count: %d\n", len(label.Strings()))
		fmt.Printf("label: %v\n", label.values)
	}
	if name.set {
		fmt.Printf("name: %v\n", name.values)
	}
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		os.Exit(2)
	}()
	var param tools.DockerParam
	if *localFlag {
		fmt.Println("no host defined so use local docker engine")
		param = tools.DockerParam{URL: tools.DockerLocal}
	}

	if len(flag.Args()) == 1 && *remoteFlag {
		param = tools.DockerParam{
			URL:  fmt.Sprintf("tcp://%s:2376", flag.Args()[0]),
			Path: certPath.String(),
		}
	} else {
		if len(flag.Args()) == 1 && *swarmFlag {
			param = tools.DockerParam{
				URL:       fmt.Sprintf("tcp://%s:2376", flag.Args()[0]),
				Path:      certPath.String(),
				SwarmMode: true,
			}
		} else {
			fmt.Println("no argument defined to docker engine")
			os.Exit(2)
		}

	}
	fmt.Printf("connect to %s with path %s\n ", param.URL, param.Path)
	docker, err := tools.NewDocker(param)
	if err != nil {
		fmt.Printf("cant create docker client %s", err)
		os.Exit(2)
	}
	docker.ListContainers(name.String(), label.Strings())
	wait()
}

func wait() {
	done := make(chan (bool))
	time.AfterFunc(5*time.Minute, func() {
		fmt.Println("exiting")
		done <- true
	})
	<-done
	os.Exit(0)
}
