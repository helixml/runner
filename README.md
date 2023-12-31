# Easy Generative AI with Dagger and Helix

Requires:

* Linux
* NVIDIA GPU

## Setup

Install dagger CLI >= v0.9.3
```
curl -sfL https://releases.dagger.io/dagger/install.sh | sudo BIN_DIR=/usr/local/bin sh
```

Enable GPU support:
```
export _EXPERIMENTAL_DAGGER_GPU_SUPPORT=1
```

Test nvidia-smi:

```
dagger call -m github.com/helixml/runner nvidia-smi
```
```
✔ dagger call nvidia-smi [9.06s]
┃ Mon Nov  6 13:38:57 2023
┃ +---------------------------------------------------------------------------------------+
┃ | NVIDIA-SMI 545.23.06              Driver Version: 545.23.06    CUDA Version: 12.3     |
┃ |-----------------------------------------+----------------------+----------------------+
┃ | GPU  Name                 Persistence-M | Bus-Id        Disp.A | Volatile Uncorr. ECC |
┃ | Fan  Temp   Perf          Pwr:Usage/Cap |         Memory-Usage | GPU-Util  Compute M. |
┃ |                                         |                      |               MIG M. |
┃ |=========================================+======================+======================|
┃ |   0  NVIDIA GeForce RTX 3090        On  | 00000000:01:00.0 Off |                  N/A |
┃ |  0%   44C    P8              25W / 350W |      6MiB / 24576MiB |      0%      Default |
┃ |                                         |                      |                  N/A |
┃ +-----------------------------------------+----------------------+----------------------+
┃
┃ +---------------------------------------------------------------------------------------+
┃ | Processes:                                                                            |
┃ |  GPU   GI   CI        PID   Type   Process name                            GPU Memory |
┃ |        ID   ID                                                             Usage      |
┃ |=======================================================================================|
┃ |  No running processes found                                                           |
┃ +---------------------------------------------------------------------------------------+
• Engine: 25fcbf2f3c88 (version v0.9.3)
⧗ 47.50s ✔ 75
```

## Create an image (SDXL)

```
dagger download -m github.com/helixml/runner generate-file --prompt fish
```

The resulting image will be written to the current working directory.

## Chat with daggerbot (Mistral-7B)

Coming soon.

```
dagger shell -m github.com/helixml/runner chat
```

## View status of helix runner

Coming soon.

```
dagger shell -m github.com/helixml/runner status
```

## Roadmap

* Chat with daggerbot (Mistral-7B) with a CLI chatbot (demoing reusing GPU memory)
* View GPU memory status of helix runner while it's running with a nice CLI tool
* Demo of a DAG chaining models together e.g. use Mistral-7B to come up with a prompt for SDXL
* CLI support for fine-tuning SDXL
* CLI support for fine-tuning Mistral-7B
* Connect your runner to helix.ml so you can manage and serve it through a web interface

## More info

Learn more about helix at [https://github.com/helixml/helix](https://github.com/helixml/helix). Note the license of helix itself may differ from the license of this Dagger module.
