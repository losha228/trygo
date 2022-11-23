package main

import (
	"context"
	"fmt"
	"hash/fnv"

	"github.com/cespare/xxhash/v2"
	"github.com/containerd/containerd"
)

func main() {
	containerdTest()
}

func containerdTest() error {
	runtimeURI := "/run/containerd/containerd.sock"
	client, err := containerd.New(runtimeURI)
	if err != nil {
		return err
	}
	images, _ := client.ImageService().List(context.TODO())
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
