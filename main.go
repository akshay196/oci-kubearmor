package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
)

const (
	policyConfigMediaType = "application/vnd.cncf.kubearmor.config.v1+json"
	policyLayerMediaType  = "application/vnd.cncf.kubearmor.policy.layer.v1.yaml"
)

func pushFiles() error {
	if len(os.Args) < 3 {
		panic(errors.New("you should specify policy path as a first argument and an image as a second argument"))
	}
	policyRef := os.Args[1]
	image := os.Args[2]

	// 0. Create a file store
	curr, _ := os.Getwd()
	fs, err := file.New(curr)
	if err != nil {
		return err
	}
	defer fs.Close()
	ctx := context.Background()

	// 1. Add files to a file store
	fileNames := []string{policyRef}
	fileDescriptors := make([]v1.Descriptor, 0, len(fileNames))
	for _, name := range fileNames {
		fileDescriptor, err := fs.Add(ctx, name, policyLayerMediaType, "")
		if err != nil {
			return err
		}
		fileDescriptors = append(fileDescriptors, fileDescriptor)
		fmt.Printf("file descriptor for %s: %v\n", name, fileDescriptor)
	}

	// 2. Pack the files and tag the packed manifest
	manifestDescriptor, err := oras.Pack(ctx, fs, policyConfigMediaType, fileDescriptors, oras.PackOptions{
		PackImageManifest: true,
	})
	if err != nil {
		return err
	}
	fmt.Println("manifest descriptor:", manifestDescriptor)

	var reg, tag string
	idx := strings.LastIndex(image, ":")
	if idx != -1 {
		reg = image[0:idx]
		tag = image[idx+1:]
	}

	if err = fs.Tag(ctx, manifestDescriptor, tag); err != nil {
		return err
	}

	// 3. Connect to a remote repository
	repo, err := remote.NewRepository(reg)
	repo.PlainHTTP = true
	if err != nil {
		panic(err)
	}

	// 4. Copy from the file store to the remote repository
	_, err = oras.Copy(ctx, fs, tag, repo, tag, oras.DefaultCopyOptions)
	return err
}

func main() {
	if err := pushFiles(); err != nil {
		panic(err)
	}
}
