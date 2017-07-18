package tools

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"

	"github.com/fsouza/go-dockerclient"
)

var dockerLocal = "unix:///var/run/docker.sock"

// Docker client interface
type Docker interface {
	ListContainers(namePattern string, labelPattern string)
}

type dockerImpl struct {
	client *docker.Client
}

// NewDocker docker client constructor
func NewDocker() Docker {
	client, err := docker.NewClient(dockerLocal)
	if err != nil {
		panic(err)
	}

	return &dockerImpl{client: client}
}

func (clientImpl *dockerImpl) ListContainers(namePattern string, labelPattern string) {
	opts := docker.ListContainersOptions{All: true, Filters: map[string][]string{"label": {labelPattern}, "status": {"running"}}}
	containers, err := clientImpl.client.ListContainers(opts)
	if err != nil {
		panic(err)
	}
	log.Printf("list : %v", containers)
	if len(containers) == 0 {
		panic("no containers match")
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
