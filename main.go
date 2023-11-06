package main

import (
	"context"
	"fmt"
	"log"
	"time"

	helix_types "github.com/helixml/helix/api/pkg/types"
)

type Helix struct{}

// example usage: "dagger call nvidia-smi"
func (m *Helix) NvidiaSmi(ctx context.Context) (string, error) {
	return dag.Container().
		From("nvidia/cuda:12.2.2-base-ubuntu22.04").
		ExperimentalWithAllGPUs().
		WithExec([]string{"nvidia-smi"}).
		Stdout(ctx)
}

func (m *Helix) Service(ctx context.Context, outputPath *Directory) (*Service, error) {
	return dag.Container().
		From("quay.io/lukemarsden/helix-runner:v0.0.2").
		ExperimentalWithAllGPUs().
		WithExec([]string{"/app/helix/helix", "runner", "--timeout-seconds", "600", "--memory", "24GB"}).
		WithMountedDirectory("/app/sd-scripts/output_images", outputPath).
		WithExposedPort(8080).
		AsService(), nil
}

// You don't want to have to load the model weights every time you use the AI
// model - you want to reuse the GPU memory so we run it as a service that
// persists across the lifetime of many dagger calls within a dag

func (m *Helix) Test(ctx context.Context, outputPath *Directory) (string, error) {
	// create HTTP service container with exposed port 8080
	httpSrv, err := m.Service(ctx, outputPath)
	if err != nil {
		return "", err
	}

	helix_types.Session{}

	// get endpoint
	val, err := httpSrv.Endpoint(ctx)
	if err != nil {
		return "", err
	}

	fmt.Println(val)
	log.Printf("one")

	time.Sleep(10 * time.Second)

	val, err = dag.Container().
		From("alpine").
		WithServiceBinding("www", httpSrv).
		WithExec([]string{"wget", "-O-", "http://www:8080"}).
		Stdout(ctx)

	log.Printf("two")
	time.Sleep(10 * time.Second)

	if err != nil {
		return "", err
	}

	return val, nil
}

///////////////

// example usage: "dagger call container-echo --string-arg yo"
func (m *Helix) ContainerEcho(stringArg string) *Container {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg})
}

// example usage: "dagger call grep-dir --directory-arg . --pattern GrepDir"
func (m *Helix) GrepDir(ctx context.Context, directoryArg *Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}
