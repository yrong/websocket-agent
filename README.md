opsagent
==========

`opsagent` is customized from open source project [websockted](https://github.com/joewalnes/websocketd) with advanced searching ability.

## features

- Server side scripts can access details about the WebSocket HTTP request (e.g. remote host, query parameters, cookies, path, etc) via standard CGI environment variables.
- As well as serving websocket daemons it also includes a static file server and classic CGI server for convenience.
- Command line help available via websocketd --help.
- searching ability provided by integration with elasticsearch

## installation

```
go get -u github.com/golang/dep/cmd/dep
cd $GOPATH/src;git clone https://e.coding.net/scirichon/OpsAgent.git
cd OpsAgent;dep ensure
go install
cd examples/windows-vbscript
OpsAgent console.cmd
start console.html
```