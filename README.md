<div align="center">

[<img src="./assets/tob.PNG" width="100">](https://github.com/telkomdev/tob)
<h3>Tob => Bot</h3>
A Notification Bot written in Go
</div>

### Architecture

[<img src="./assets/tob_arch.PNG" width="600">](https://github.com/telkomdev/tob)

### Screenshot
<h4>Discord</h5>

[<img src="./assets/discord_n.PNG" width="400">](https://github.com/telkomdev/tob)

<h4>Email</h5>

[<img src="./assets/email_n.PNG" width="400">](https://github.com/telkomdev/tob)

<h4>Slack</h5>

[<img src="./assets/slack_n.PNG" width="400">](https://github.com/telkomdev/tob)

<h4>Telegram</h5>

[<img src="./assets/telegram_n.PNG" width="400">](https://github.com/telkomdev/tob)

## Getting Started

### Install from release (https://github.com/telkomdev/tob/releases)
choose the binary from the release according to your platform, for example for the Linux platform

#### Download binary

```shell
$  wget https://github.com/telkomdev/tob/releases/download/1.0.0/tob-1.0.0.linux-amd64.tar.gz
```

#### Important !!!, always check the SHA256 Checksum before using it

Download `sha256sum.txt` according to the binary version you downloaded https://github.com/telkomdev/tob/releases/download/1.0.0/sha256sums.txt

```shell
$ wget https://github.com/telkomdev/tob/releases/download/1.0.0/sha256sums.txt
```

#### Verify SHA256 Checksum

Linux

```shell
$ sha256sum tob-1.0.0.linux-amd64.tar.gz -c sha256sums.txt
tob-1.0.0.linux-amd64.tar.gz: OK
```

Mac OSX

```shell
$ shasum -a 256 tob-1.0.0.darwin-amd64.tar.gz -c sha256sums.txt
tob-1.0.0.darwin-amd64.tar.gz: OK
```

You should be able to see that the checksum value for the file is valid, `tob-1.0.0.linux-amd64.tar.gz: OK` and `tob-1.0.0.darwin-amd64.tar.gz: OK`. 
Indicates the file is not damaged, not modified and safe to use.

#### Extract

```shell
$ tar -xvzf tob-1.0.0.linux-amd64.tar.gz
```

#### Run

```shell
$ ./tob -h
```

### Build from source

Requirements
- Go version 1.16 or higher

Clone tob to your Machine
```shell
$ git clone https://github.com/telkomdev/tob.git
$ cd tob/
```

```shell
$ make build
```

`tob` options
```shell
$ ./tob -h
```

Running `tob` with config file
```shell
$ ./tob -c config.json
```

### TODO

- add Kafka service