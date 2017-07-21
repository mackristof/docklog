package tools

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strings"

	"fmt"

	"github.com/fsouza/go-dockerclient"
)

var DockerLocal = "unix:///var/run/docker.sock"

// Docker client interface
type Docker interface {
	ListContainers(namePattern string, labelPattern []string)
}

type DockerParam struct {
	Url  string
	Path string
}

type dockerImpl struct {
	client *docker.Client
}

// NewDocker docker client constructor
func NewDocker(param DockerParam) Docker {
	var err error
	var client *docker.Client
	if strings.HasPrefix(param.Url, "tcp://") {
		if len(param.Path) == 0 {
			fmt.Println("Path where certificates are located must be set when url start with tcp://")
			os.Exit(2)
		}
		ca := fmt.Sprintf("%s/ca.pem", param.Path)
		cert := fmt.Sprintf("%s/cert.pem", param.Path)
		key := fmt.Sprintf("%s/key.pem", param.Path)
		client, err = docker.NewTLSClient(param.Url, cert, key, ca)
	} else {
		if strings.HasPrefix(param.Url, "unix://") {
			client, err = docker.NewClient(param.Url)
		} else {
			client, err = docker.NewClient(param.Url)
		}
	}

	if err != nil {
		panic(err)
	}

	return &dockerImpl{client: client}
}

func (clientImpl *dockerImpl) ListContainers(namePattern string, labelPattern []string) {
	filterMap := map[string][]string{"status": {"running"}}
	if len(labelPattern) > 0 {
		// "label": {labelPattern}
		filterMap["label"] = labelPattern
	}
	opts := docker.ListContainersOptions{All: true, Filters: filterMap}
	fmt.Printf("search opts: %v\n", opts)
	containers, err := clientImpl.client.ListContainers(opts)
	if err != nil {
		panic(err)
	}
	filteredContainers := make([]docker.APIContainers, 0)
	for i, container := range containers {
		for _, name := range container.Names {
			if strings.Contains(name, namePattern) {
				filteredContainers = append(filteredContainers, containers[i])
			}
		}
	}
	log.Printf("list : %v", containers)
	if len(containers) == 0 {
		fmt.Println("no containers match")
		os.Exit(0)
	}
	var stream io.Writer = bufio.NewWriterSize(os.Stdout, 1)
	for _, container := range containers {
		var b bytes.Buffer
		writer := bufio.NewWriter(&b)
		w := io.MultiWriter(writer, stream)
		go clientImpl.getLogContainer(container.ID, w)
	}

}

func (clientImpl *dockerImpl) getLogContainer(containerID string, stream io.Writer) {

	opts := docker.LogsOptions{
		Container:    containerID,
		OutputStream: stream,
		RawTerminal:  true,
		Follow:       true,
		Stdout:       true,
		Stderr:       true,
		Timestamps:   false,
	}
	err := clientImpl.client.Logs(opts)
	if err != nil {
		panic(err)
	}
}
