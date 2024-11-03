## Tob Plugin

https://pkg.go.dev/plugin

You can create custom functionality for services. This is needed when you want to create a service that is not currently available on Tob, or when you want to create a custom message that will appear on the monitoring Dashboard or a message that will appear on the Notificator.

There are several limitations currently at https://pkg.go.dev/plugin. 
- Binary Plugins must be built with the same code as the code you use to build the main binary.

- The version of the operating system used to run the plugin binary and main binary must be the same as the one you used to build the plugin binary and main binary. For example, if you want to run the plugin binary and main binary on Ubuntu 20.04, you have to build it using Ubuntu 20.04.

### Getting started

- Clone the latest Tob code
```shell
https://github.com/telkomdev/tob.git
```

- create your plugin folder, for example `dummyplugin`
```shell
mkdir dummyplugin
```

- copy the `templateplugin.go` and `Makefile` into your `dummyplugin` folder
```shell
cp docs/plugin/templateplugin.go dummyplugin/
cp docs/plugin/Makefile dummyplugin/
```

- build the binary
```shell
cd dummyplugin/
make build
```

- add configuration to your `json config`, by adjusting 2 fields. The `kind` field is filled with `plugin`, and `pluginPath` is filled with plugin files with the extension `.so`.

```json
"dummy_plugin_one": {
    "kind": "plugin",
    "url": "https://www.google.com",
    "checkInterval": 5,
    "enable": false,
    "tags": ["product 1", "product 2"],
    "pics": ["bob", "john"],
    "pluginPath": "/home/john/tob/dummyplugin/dummyplugin.so",
    .....
```