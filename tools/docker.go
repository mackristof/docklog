package tools

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strings"

	"fmt"

	"github.com/docker/docker/api/types/swarm"
	"github.com/mackristof/go-dockerclient"
	"github.com/pkg/errors"
)

var DockerLocal = "unix:///var/run/docker.sock"

// Docker client interface
type Docker interface {
	ListContainers(namePattern string, labelPattern []string)
}

type DockerParam struct {
	URL       string
	Path      string
	SwarmMode bool
}

type dockerImpl struct {
	client *docker.Client
	swarm  bool
}

// NewDocker docker client constructor
func NewDocker(param DockerParam) (Docker, error) {
	var err error
	var client *docker.Client
	var isSwarm bool
	if strings.HasPrefix(param.URL, "tcp://") { //|| strings.HasPrefix(param.URL, "https://") {
		if len(param.Path) == 0 {
			fmt.Println("Path where certificates are located must be set when url start with tcp://")
			os.Exit(2)
		}
		ca := fmt.Sprintf("%s/ca.pem", param.Path)
		cert := fmt.Sprintf("%s/cert.pem", param.Path)
		key := fmt.Sprintf("%s/key.pem", param.Path)
		if !param.SwarmMode {
			client, err = docker.NewTLSClient(param.URL, cert, key, ca)
		} else {
			fmt.Println("connect to swarm")
			isSwarm = true
			client, err = docker.NewTLSClient(param.URL, cert, key, ca)
		}

	} else {
		if strings.HasPrefix(param.URL, "unix://") {
			client, err = docker.NewClient(param.URL)
		} else {
			panic(fmt.Sprintf("can't connect to %s", param.URL))
		}
	}

	if err != nil {
		return nil, errors.Wrap(err, "can't connect to docker engine")
	}

	return &dockerImpl{client: client, swarm: isSwarm}, nil
}

func (clientImpl *dockerImpl) ListContainers(namePattern string, labelPattern []string) {
	filterMap := map[string][]string{"status": {"running"}}
	if len(labelPattern) > 0 {
		// "label": {labelPattern}
		filterMap["label"] = labelPattern
	}
	if !clientImpl.swarm {
		opts := docker.ListContainersOptions{All: true, Filters: filterMap}
		fmt.Printf("search opts: %v\n", opts)
		containers, err := clientImpl.client.ListContainers(opts)
		if err != nil {
			panic("can't list containers with filter")
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
	} else { //swarm mode
		services, err := clientImpl.client.ListServices(docker.ListServicesOptions{})
		if err != nil {
			panic(err)
		}
		filteredServices := make([]swarm.Service, 0)
		for i, service := range services {

			if strings.Contains(service.Spec.Name, namePattern) {
				filteredServices = append(filteredServices, services[i])
			}

		}
		log.Printf("list : %v", services)
		if len(services) == 0 {
			fmt.Println("no service match")
			os.Exit(0)
		}
		var stream io.Writer = bufio.NewWriterSize(os.Stdout, 1)
		for _, service := range services {
			var b bytes.Buffer
			writer := bufio.NewWriter(&b)
			w := io.MultiWriter(writer, stream)
			go clientImpl.getLogService(service.ID, w)
		}
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

func (clientImpl *dockerImpl) getLogService(serviceID string, stream io.Writer) {

	opts := docker.LogsServiceOptions{
		Service:      serviceID,
		OutputStream: stream,
		RawTerminal:  true,
		Follow:       true,
		Stdout:       true,
		Stderr:       true,
		Timestamps:   false,
	}
	err := clientImpl.client.LogsService(opts)
	if err != nil {
		panic(err)
	}
}
