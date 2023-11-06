# Easy Generative AI with Dagger and Helix

Requires:

* Linux
* NVIDIA GPU

## Setup

Install dagger CLI >= v0.9.3
```
cd /usr/local
curl -L https://dl.dagger.io/dagger/install.sh | sudo sh
```

Enable GPU support:
```
export _EXPERIMENTAL_DAGGER_GPU_SUPPORT=1
```

Test nvidia-smi:

```
dagger call nvidia-smi
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

## Start the helix runner

The runner will automatically load and unload AI models based on demand to make optimal use of your GPU memory and minimize latency for requests.

```
dagger call start
```

## Create an image!

```
dagger call generate --type image --prompt "A pig in space"
```

The resulting image will be written to the current working directory.

Note the first time will be slower, but the second time the long-running server will have the model cached intelligently in GPU memory, and the response will be much quicker.

## Chat with daggerbot!

Coming soon.

```
dagger shell chat
```

## View status of helix runner

Coming soon.

```
dagger shell status
```

## Roadmap

* Support for fine-tuning
* Connect your runner to helix.ml and so you can manage and serve it through a web interface

## More info

Learn more about helix at [https://github.com/helixml/helix](https://github.com/helixml/helix). Note the license of helix itself may differ from the license of this Dagger module.