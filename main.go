package main

import (
	"context"
	"log"
	"time"
)

type Helix struct{}

const HELIX_IMAGE = "quay.io/lukemarsden/helix-runner:v0.0.8"

// TODO: need to make client download file from runner via the API

func (m *Helix) Service(ctx context.Context) *Service {
	return dag.Container().
		From(HELIX_IMAGE).
		ExperimentalWithAllGPUs().
		WithEntrypoint([]string{"/app/helix/helix", "runner", "--timeout-seconds", "600", "--memory", "24GB"}).
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

func (m *Helix) Generate(ctx context.Context, prompt string) (*Container, error) {
	// create HTTP service container with exposed port 8080
	helixRunner := m.Service(ctx)

	container := m.Client(ctx).
		WithServiceBinding("helix-runner", helixRunner).
		WithEnvVariable("CACHE_BUSTER", time.Now().Format(time.RFC3339Nano)).
		WithExec([]string{"--api-host", "http://helix-runner:8080", "--type", "image", "--prompt", prompt})

	return container, nil
}

func (m *Helix) GenerateFile(ctx context.Context, prompt string) (*File, error) {
	container, err := m.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}
	stdout, err := container.Stdout(ctx)
	if err != nil {
		return nil, err
	}
	stderr, err := container.Stderr(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("Got stdout from generate: %s", stdout)
	log.Printf("Got stderr from generate: %s", stderr)
	return container.File("/app/helix/output.png"), nil
}

// example usage: "dagger call nvidia-smi"
func (m *Helix) NvidiaSmi(ctx context.Context) (string, error) {
	return dag.Container().
		From("nvidia/cuda:12.2.2-base-ubuntu22.04").
		ExperimentalWithAllGPUs().
		WithExec([]string{"nvidia-smi"}).
		Stdout(ctx)
}
