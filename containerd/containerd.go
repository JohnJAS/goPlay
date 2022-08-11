package main

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"log"
)

// sandboxImage label defined here to avoid cyclic imports as they are used in both image store and cri server.
const (
	// sandboxImageLabelKey is the label value indicating the image is sandbox image
	sandboxImageLabelKey = "io.cri-containerd" + ".image.kind"
	// sandboxImageLabelValue is the label value indicating the image is sandbox image
	sandboxImageLabelValue = "sandbox"
)

func IsSandboxImage(labels map[string]string) bool {
	return labels[sandboxImageLabelKey] == sandboxImageLabelValue
}

func main() {
	ctx := context.Background()
	opts := containerd.WithDefaultNamespace("k8s.io")
	client, err := containerd.New("/run/containerd/containerd.sock", opts)
	if err != nil {
		log.Fatalf(err.Error())
	}
	imageList := []string{
		"localhost:5000/hpeswitomsandbox/pause:3.2",
		"k8s.gcr.io/pause:3.5",
	}

	for _, image := range imageList {
		image, err := client.GetImage(ctx, image)
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Println(image.Name())
		fmt.Println(image.Labels())
		fmt.Println(image.Metadata())
		fmt.Println(image.Config(ctx))
		fmt.Println(image.Platform())
		fmt.Println(image.RootFS(ctx))

		fmt.Println(IsSandboxImage(image.Labels()))
	}

	defer client.Close()
}

