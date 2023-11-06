package main

import (
	"context"
	"time"
)

type Helix struct{}

const HELIX_IMAGE = "quay.io/lukemarsden/helix-runner:v0.0.2"

func (m *Helix) Service(ctx context.Context, outputPath *Directory) *Service {
	return dag.Container().
		From(HELIX_IMAGE).
		ExperimentalWithAllGPUs().
		WithEntrypoint([]string{"/app/helix/helix", "runner", "--timeout-seconds", "600", "--memory", "24GB"}).
		WithMountedDirectory("/app/sd-scripts/output_images", outputPath).
		WithExposedPort(8080).
		AsService()
}

func (m *Helix) Client(ctx context.Context) *Container {
	return dag.Container().
		From(HELIX_IMAGE).
		WithEntrypoint([]string{"/app/helix/helix", "run"})
}

// You don't want to have to load the model weights every time you use the AI
// model - you want to reuse the GPU memory so we run it as a service that
// persists across the lifetime of many dagger calls within a dag

func (m *Helix) Generate(ctx context.Context, outputPath *Directory, prompt string) (*Container, error) {
	// create HTTP service container with exposed port 8080
	helixRunner := m.Service(ctx, outputPath)

	container := m.Client(ctx).
		WithServiceBinding("helix-runner", helixRunner).
		WithEnvVariable("CACHE_BUSTER", time.Now().Format(time.RFC3339Nano)).
		WithExec([]string{"--api-host", "http://helix-runner:8080", "--type", "image", "--prompt", prompt})

	return container, nil
}

// example usage: "dagger call nvidia-smi"
func (m *Helix) NvidiaSmi(ctx context.Context) (string, error) {
	return dag.Container().
		From("nvidia/cuda:12.2.2-base-ubuntu22.04").
		ExperimentalWithAllGPUs().
		WithExec([]string{"nvidia-smi"}).
		Stdout(ctx)
}
