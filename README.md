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

### Install from the latest release (https://github.com/telkomdev/tob/releases)
choose the binary from the release according to your platform, for example for the Linux platform

#### Download binary

```shell
$  wget https://github.com/telkomdev/tob/releases/download/1.1.0/tob-1.1.0.linux-amd64.tar.gz
```

#### Important !!!, always check the SHA256 Checksum before using it

Download `sha256sum.txt` according to the binary version you downloaded https://github.com/telkomdev/tob/releases/download/1.1.0/sha256sums.txt

```shell
$ wget https://github.com/telkomdev/tob/releases/download/1.1.0/sha256sums.txt
```

#### Verify SHA256 Checksum

Linux

```shell
$ sha256sum tob-1.1.0.linux-amd64.tar.gz -c sha256sums.txt
tob-1.1.0.linux-amd64.tar.gz: OK
```

Mac OSX

```shell
$ shasum -a 256 tob-1.1.0.darwin-amd64.tar.gz -c sha256sums.txt
tob-1.1.0.darwin-amd64.tar.gz: OK
```

You should be able to see that the checksum value for the file is valid, `tob-1.1.0.linux-amd64.tar.gz: OK` and `tob-1.1.0.darwin-amd64.tar.gz: OK`. 
Indicates the file is not damaged, not modified and safe to use.

#### Extract

```shell
$ tar -xvzf tob-1.1.0.linux-amd64.tar.gz
```

#### Run

```shell
$ ./tob -c config.json
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

### Service and Kind
currently tob supports below `KIND` of services
- **airflow**
- **airflowflower**
- **elasticsearch**
- **mongodb**
- **mysql**
- **postgresql**
- **redis**
- **web**

`KIND` represents one or many services. So you can monitor more than one service with the same `KIND`. For example, you can monitor multiple PostgreSQL instances. Or you can monitor multiple web applications.

`checkInterval: in Seconds` is how often your service is called by tob.

`enable` you set `true` when you want to monitor the service. Set it to `false`, if you don't want to monitor it.

`config.json`

```json
"postgresql_one": {
    "kind": "postgresql",
    "url": "postgres://demo:12345@localhost:5432/demo?sslmode=disable",
    "checkInterval": 10,
    "enable": false
},

"postgresql_two": {
    "kind": "postgresql",
    "url": "postgres://demo:12345@localhost:5433/demo?sslmode=disable",
    "checkInterval": 10,
    "enable": false
},

"web_internal": {
    "kind": "web",
    "url": "https://portal.mycompany.com/health-check",
    "checkInterval": 5,
    "enable": true
},

"web_main_1": {
    "kind": "web",
    "url": "https://mycompany.com/health-check",
    "checkInterval": 5,
    "enable": true
}
```

### Notificator

Currently tob supports the following types of `Notificator`. `Notificator` is where the tob will send notifications when one or more of the services you're monitoring have problems.

- **Discord**
- **Email with SMTP**
- **Slack (webhook)** https://api.slack.com/messaging/webhooks
- **Telegram**
- **Webhook** | For security reasons, your `webhook endpoint` must verify the HTTP header: `x-tob-token` that is in every incoming http request.

Example of `x-tob-token` webhook verification in nodejs application

```javascript
const express = require('express');
const app = express();
const bodyParser = require('body-parser');

const PORT = 3001;

const tobToken = "461b919e-1bf4-42db-a8ff-4f21633bbf10";

app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());

app.post('/webhook/tob', (req, res) => {
    const headers = req.headers;
    const token = headers["x-tob-token"];
    
    const verifiedRequest = token === tobToken;
    if (!verifiedRequest) {
        return res.status(401).send({'message': 'token is not valid'});
    }

    console.log(req.body);

    return res.status(200).send({'message': 'webbhook received'});
});

app.listen(PORT, () => console.log(`Server listening on port: ${PORT}`));
```

### TODO

- add Kafka service