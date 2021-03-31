# Background
Tpot is a simple tsh teleport wrapper. Currently, we are not able to get list of available node in tsh. Usually, to get the
list of nodes we need to go through teleport web ui then creating an alias to accessible through tsh. This make you're able to
get list of nodes from Terminal, pick one then login to the node by `tsh`.

# Feature
![tpot.gif](tpot.gif)
- Support multiple environment
- Browsing the node list and search it
- You're able to get the node list from a cache or fresh from the teleport server

# How does it work
this tool simply store the proxy environment under your `$HOME/.tpot/` directory.
whenever you try to get the node list it'll ask the teleport server to give the latest node list. Once, we got it, it'll store
 in the configuration file for caching purpose.

# Install
this tool requires `tsh` is installed in your machine.
if you don't have it yet, you can download & install from [this](https://gravitational.com/teleport/docs/user-manual/#installing-tsh).
to install this tool you can run this command.
```shell script
curl  https://raw.githubusercontent.com/adzimzf/tpot/master/download.sh | sh
```
or
```shell script
wget -O - https://raw.githubusercontent.com/adzimzf/tpot/master/download.sh | sh
```
It'll copy the binary to `/usr/bin`.

If you want to install into a specific directory you can add arguments `-s -- -b "directory"`, for example:
```shell
curl  https://raw.githubusercontent.com/adzimzf/tpot/master/download.sh | sh -s -- -b "/home/myuser/Desktop" 
```

If you are familiar with `Golang/Go` and your `Golang version` is `go1.13` you can install using `go install` by running this command:
```shell
go get github.com/adzimzf/tpot
go install github.com/adzimzf/tpot
```

# Usage
Before use this tools you need to add proxy configuration first by run this command
```shell script
tpot -c --add
```
It'll prompt your config editor, by default it'll use `nano`
- `Environment` is an identifier for your proxy config, eg. `staging` and `prod`
- `Proxy address` is a valid proxy address in http protocol, eg. `https://teleport.myport.com:3080`
- `User name` ia a user name used for teleport login. eg. `adzimzf`
- `Auth Connector` ia a 3rd party auth connector for SSO. eg. `gsuite`
- `Need 2Fa` does the proxy need 2FA or not. eg `true` or `false`


you can change the default editor by running this command
```shell script
tpot -c --edit
```
or
```shell script
tpot --config --edit
```

if the configuration installed successfully you can start use `tpot` by running this command
```shell script
tpot staging
```
or
```shell script
tpot staging -a
```
will get the node from the server then append to the cache,
or 
```shell script
tpot staging -r
```

When the list of node shows, you can navigate by `RIGHT`, `LEFT`, `UP` and `DOWN`. For searching the node, you can type the `node name` then hit `TAB`.
Hit `ENTER` to select the node and login. 


to get the node server instead of `cache`. if it gives you an error `Permision denied`, you can manually add `tpot` config dir by running this command
```shell script
mkdir $HOME/.tpot
```
with `775` permission, then you can re-add the configuration



That's all hope you find your need

