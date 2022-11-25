package main

import (
	"bufio"
	"context"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cespare/xxhash/v2"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/docker/docker/api/types"
	dockerapi "github.com/docker/docker/client"
	"k8s.io/utils/path"
)

func main() {
	readCert("/tmp")
	//containerdTest()
	//dockerClientTest()
}

func readCert(folder string) (string, error) {
	// sonick8sclient.pfx.notify
	libRegEx, e := regexp.Compile("^.+\\.(pfx.notify)$")
	if e != nil {
		fmt.Printf("can not build regex: %s", e.Error())
		return "", e
	}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		fmt.Printf("can not build regex: %s", err.Error())
		return "", err
	}

	notifyFile := ""
	for _, f := range files {
		if libRegEx.MatchString(f.Name()) {
			notifyFile = filepath.Join(folder, f.Name())
		}
	}

	if notifyFile == "" {
		return "", fmt.Errorf("Invalid notify file.")
	}
	fmt.Printf("Found notify file: %s", notifyFile)

	file, err := os.Open(notifyFile)
	if err != nil {
		fmt.Printf("File to open file: %s", err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	cert := ""
	for scanner.Scan() {
		p := strings.Trim(scanner.Text(), " ")
		if ok, err := path.Exists(path.CheckFollowSymlink, p); err == nil && ok {
			cert = p
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("scanner err: %s\n", err.Error())
		return "", err
	}

	if cert == "" {
		return "", fmt.Errorf("Invalid notify file.")
	}
	fmt.Printf("cert is %s\n", cert)
	return cert, nil
}

func dockerClientTest() error {
	fmt.Printf("dockerClientTest:\n")
	runtimeURI := "unix:///var/run/docker.sock"
	//cli, err := dockerapi.NewClientWithOpts(dockerapi.FromEnv)

	cli, err := dockerapi.NewClient(runtimeURI, "1.41", nil, nil)

	if err != nil {
		fmt.Printf("fail to create docker client: %v", err.Error())
		return err
	}
	defer cli.Close()
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All:    true,
		Latest: false,
		Limit:  100,
		Size:   true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("List images:")
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
