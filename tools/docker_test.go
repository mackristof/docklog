package tools

import (
	"testing"

	"github.com/fsouza/go-dockerclient"
)

func Test_dockerImpl_ListContainers(t *testing.T) {
	client, err := docker.NewClient(DockerLocal)
	if err != nil {
		t.Fatal("dockerClient init problem")
	}
	param := DockerParam{URL: DockerLocal, SwarmMode: false}
	dockerImpl, err := NewDocker(param)
	if err != nil {
		t.Fatal("dockerImpl client init problem")
	}
	type fields struct {
		client *docker.Client
		swarm  bool
	}
	type args struct {
		namePattern  string
		labelPattern []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "should found containers where name contains 'toto'",
			fields: fields{client: client, swarm: false},
			args:   args{namePattern: "toto"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			dockerImpl.ListContainers(tt.args.namePattern, tt.args.labelPattern)
		})
	}
}
