package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/lukemarsden/helix/api/pkg/system"
	"github.com/lukemarsden/helix/api/pkg/types"
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

func (m *Helix) Generate(ctx context.Context, outputPath *Directory, prompt string) (string, error) {
	// create HTTP service container with exposed port 8080
	httpSrv, err := m.Service(ctx, outputPath)
	if err != nil {
		return "", err
	}

	// get endpoint
	val, err := httpSrv.Endpoint(ctx)
	if err != nil {
		return "", err
	}

	interaction := types.Interaction{
		ID:       "cli-user",
		Created:  time.Now(),
		Creator:  "user",
		Message:  prompt,
		Finished: true,
	}
	interactionSystem := types.Interaction{
		ID:       "cli-system",
		Created:  time.Now(),
		Creator:  "system",
		Finished: false,
	}

	id := system.GenerateUUID()
	session := types.Session{
		ID:           "cli-" + id,
		Name:         "cli",
		Created:      time.Now(),
		Updated:      time.Now(),
		Mode:         "inference",
		Type:         types.SessionType("image"),
		ModelName:    types.Model_SDXL,
		FinetuneFile: "",
		Interactions: []types.Interaction{interaction, interactionSystem},
		Owner:        "cli-user",
		OwnerType:    "user",
	}

	bs, err := json.Marshal(session)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", val+"/api/v1/worker/session", bytes.NewBuffer(bs))
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	log.Printf("Response: %+v", resp)

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
