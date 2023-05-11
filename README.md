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
$  wget https://github.com/telkomdev/tob/releases/download/1.4.1/tob-1.4.1.linux-amd64.tar.gz
```

#### Important !!!, always check the SHA256 Checksum before using it

Download `sha256sum.txt` according to the binary version you downloaded https://github.com/telkomdev/tob/releases/download/1.4.1/sha256sums.txt

```shell
$ wget https://github.com/telkomdev/tob/releases/download/1.4.1/sha256sums.txt
```

#### Verify SHA256 Checksum

Linux

```shell
$ sha256sum tob-1.4.1.linux-amd64.tar.gz -c sha256sums.txt
tob-1.4.1.linux-amd64.tar.gz: OK
```

Mac OSX

```shell
$ shasum -a 256 tob-1.4.1.darwin-amd64.tar.gz -c sha256sums.txt
tob-1.4.1.darwin-amd64.tar.gz: OK
```

You should be able to see that the checksum value for the file is valid, `tob-1.4.1.linux-amd64.tar.gz: OK` and `tob-1.4.1.darwin-amd64.tar.gz: OK`. 
Indicates the file is not damaged, not modified and safe to use.

#### Extract

```shell
$ tar -xvzf tob-1.4.1.linux-amd64.tar.gz
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
- **kafka**
- **mongodb**
- **mysql**
- **postgresql**
- **redis**
- **web**
- **diskstatus**

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

### Disk Status Monitoring

To monitor `Disk Status` on a Server Computer, `tob` requires a special `agent` that can be called by `tob`. 
So we need to deploy an `agent`, in this case `tob-http-agent` to the Server Computer whose `Disk Status` we need to monitor.

#### Download `tob-http-agent` binary

```shell
$  wget https://github.com/telkomdev/tob/releases/download/1.4.1/tob-http-agent-1.0.0.linux-amd64.tar.gz
```

#### Important !!!, always check the SHA256 Checksum before using it

Download `tob-http-agent-sha256sums.txt` according to the binary version you downloaded https://github.com/telkomdev/tob/releases/download/1.4.1/tob-http-agent-sha256sums.txt

```shell
$ wget https://github.com/telkomdev/tob/releases/download/1.4.1/tob-http-agent-sha256sums.txt
```

#### Verify `tob-http-agent` SHA256 Checksum

Linux

```shell
$ sha256sum tob-http-agent-1.0.0.linux-amd64.tar.gz -c tob-http-agent-sha256sums.txt
tob-http-agent-1.0.0.linux-amd64.tar.gz: OK
```

#### Extract `tob-http-agent`

```shell
$ tar -xvzf tob-http-agent-1.0.0.linux-amd64.tar.gz
```

#### Run `tob-http-agent` as a daemon with `systemd`

Create `tob-http-agent.service` `systemd` unit service
```shell
$ sudo vi /etc/systemd/system/tob-http-agent.service
```

Copy content from this file to the `/etc/systemd/system/tob-http-agent.service` and save
```
https://github.com/telkomdev/tob/blob/master/deployments/systemd/tob-http-agent.service
```

Reload `systemd daemon`
```shell
$ sudo systemctl daemon-reload
```

Start `tob-http-agent` service
```shell
$ sudo systemctl start tob-http-agent
```

Check if its running
```shell
$ sudo systemctl status tob-http-agent
```

#### Expose `tob-http-agent` service with `nginx`

Create `tob-http-agent.conf`
```shell
$ sudo vi /etc/nginx/sites-available/tob-http-agent.conf
```

Copy content from this file to the `/etc/nginx/sites-available/tob-http-agent.conf` and save
```
https://github.com/telkomdev/tob/blob/master/deployments/nginx/tob-http-agent-nginx.conf
```

Create `/etc/nginx/sites-available/tob-http-agent.conf` symlink
```shell
$ sudo ln -s /etc/nginx/sites-available/tob-http-agent.conf /etc/nginx/sites-enabled/
```

Make sure the `nginx` configuration is not error
```shell
$ sudo nginx -t
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

Restart `nginx`
```shell
$  sudo systemctl restart nginx
```

#### Add `diskstatus` config to the `tob` service config

```json
"ubuntu_1_storage_status": {
    "kind": "diskstatus",
    "url": "http://tob-http-agent.yourdomain.com",
    "checkInterval": 5,
    "thresholdDiskUsage": 90,
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

The `tobToken` variable must be the same as the `tobToken` config located in the `config.json` file

```json
"webhook": {
    "url": "https://api.yourcompany.com/webhook/tob",
    "tobToken": "461b919e-1bf4-42db-a8ff-4f21633bbf10",
    "enable": true
}
```


`tob` will send a message/payload in the following form to the webhook endpoint that you have specified in the config above.

```json
{ 
    "message": "mysql_cluster_1 is DOWN" 
}
```