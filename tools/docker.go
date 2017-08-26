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
	"github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
)

var DockerLocal = "unix:///var/run/docker.sock"

// Docker client interface
type Docker interface {
	ListContainers(namePattern string, labelPattern []string) []Container
	GetLogs(containers []Container)
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

type Container struct {
	ID      string
	Service bool
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

func (clientImpl *dockerImpl) ListContainers(namePattern string, labelPattern []string) []Container {
	result := make([]Container, 0)
	filterMap := map[string][]string{"status": {"running"}}
	if len(labelPattern) > 0 {
		// "label": {labelPattern}
		filterMap["label"] = labelPattern
	}
	if !clientImpl.swarm {
		opts := docker.ListContainersOptions{All: true, Filters: filterMap}
		fmt.Printf("search opts: %+v\n", opts)
		containers, err := clientImpl.client.ListContainers(opts)
		if err != nil {
			panic("can't list containers with filter")
		}
		filteredContainers := make([]docker.APIContainers, 0)
		for i, container := range containers {
			for _, name := range container.Names {
				if strings.Contains(name, namePattern) {
					fmt.Printf("container %s found\n", name)
					filteredContainers = append(filteredContainers, containers[i])
				}
			}
		}
		log.Printf("list : %+v", filteredContainers)
		if len(filteredContainers) == 0 {
			fmt.Println("no containers match")
			os.Exit(0)
		}
		for _, container := range containers {
			result = append(result, Container{ID: container.ID, Service: false})
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
		log.Printf("list : %+v", filteredServices)
		if len(filteredServices) == 0 {
			fmt.Println("no service match")
			os.Exit(0)
		}
		for _, service := range services {
			result = append(result, Container{ID: service.ID, Service: true})
		}
	}
	return result

}

func (clientImpl *dockerImpl) GetLogs(containers []Container) {
	var stream io.Writer = bufio.NewWriterSize(os.Stdout, 1)
	for _, container := range containers {
		var b bytes.Buffer
		writer := bufio.NewWriter(&b)
		w := io.MultiWriter(writer, stream)
		go clientImpl.getLog(container, w)
	}
}

func (clientImpl *dockerImpl) getLog(container Container, stream io.Writer) {
	if container.Service {
		opts := docker.LogsServiceOptions{
			Service:      container.ID,
			OutputStream: stream,
			RawTerminal:  true,
			Follow:       true,
			Stdout:       true,
			Stderr:       true,
			Timestamps:   false,
		}
		err := clientImpl.client.GetServiceLogs(opts)
		if err != nil {
			panic(err)
		}
	} else {
		opts := docker.LogsOptions{
			Container:    container.ID,
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
}
