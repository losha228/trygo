package main

import (
	"context"
	"fmt"
	"hash/fnv"

	"github.com/cespare/xxhash/v2"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/docker/docker/api/types"
	dockerapi "github.com/docker/docker/client"
)

func main() {
	containerdTest()
	dockerClientTest()
}

func dockerClientTest() error {
	fmt.Printf("dockerClientTest:\n")
	runtimeURI := "unix://var/run/docker.sock"
	//cli, err := dockerapi.NewClientWithOpts(dockerapi.FromEnv)

	cli, err := dockerapi.NewClient(runtimeURI, "1.41", nil, nil)

	if err != nil {
		fmt.Printf("fail to create docker client: %v", err.Error())
		return err
	}
	fmt.Printf("cli %v", cli)
	defer cli.Close()
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}

	return nil
}

func containerdTest() error {
	ctx := namespaces.WithNamespace(context.Background(), "docker")
	//	cri = namespaces.WithNamespace(ctx, "cri")
	runtimeURI := "/run/containerd/containerd.sock"
	client, err := containerd.New(runtimeURI, containerd.WithDefaultNamespace("docker"))
	if err != nil {
		fmt.Printf("error to init containerd:%v\n", err.Error())
		return err
	}
	defer client.Close()
	images, err := client.ImageService().List(ctx)
	if err != nil {
		fmt.Printf("error to get image service:%v\n", err.Error())
		return err
	}
	for _, m := range images {
		fmt.Printf("%s\n", m.Name)
	}
	return nil
}

func makeHash(input string) string {
	h := xxhash.New()
	h.Write([]byte(input))
	return fmt.Sprintf("%d", h.Sum64())
}

func HashObject(input string) string { //nolint:revive
	objHash := fnv.New64a()
	//objHash.WriteString(input)
	objHash.Write([]byte(input))
	//WriteHashObject(objHash, object)
	return fmt.Sprintf("%d", objHash.Sum64())
}
