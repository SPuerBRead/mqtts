	███╗   ███╗ ██████╗ ████████╗████████╗███████╗
	████╗ ████║██╔═══██╗╚══██╔══╝╚══██╔══╝██╔════╝
	██╔████╔██║██║   ██║   ██║      ██║   ███████╗
	██║╚██╔╝██║██║▄▄ ██║   ██║      ██║   ╚════██║
	██║ ╚═╝ ██║╚██████╔╝   ██║      ██║   ███████║
	╚═╝     ╚═╝ ╚══▀▀═╝    ╚═╝      ╚═╝   ╚══════╝


# MQTTS
[![GoV](https://img.shields.io/badge/golang-1.16.4-brightgreen.svg)]()

![](./img/render1623551568329.gif)

支持安全检查类型
-----------
* 匿名登陆 (批量)
* emqx embox_plugin_template任意用户名密码登陆 (批量)
* 用户名密码爆破 (批量)
* 获取服务端信息
* 尽可能获取所有topic信息
* 获取证书信息

支持协议类型
-----------
* TCP
* SSL
* WS
* WSS

使用说明
-----------
自动探测（包含匿名登陆、任意用户名密码登陆、用户名密码爆破）

`./mqtts -t 127.0.0.1 -p 1883 -au`

获取服务端信息

`./mqtts -t 127.0.0.1 -p 1883 -s`

获取topic信息

`./mqtts -t 127.0.0.1 -p 1883 -ts -w 60`

批量测试

`./mqtts -tf ./target.txt -au`

其他参数见 `./mqtts -h`

批量扫描文件格式(空格分割，*必填项)

*host *port protocol clientId username password


编译源代码
-----------

### mac

`CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o mqtts_darwin_amd64 main.go `

### linux

`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mqtts_linux_amd64 main.go`

### win64

`CGO_ENABLED=0 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o mqtts_windows_amd64.exe main.go`

### win32

`CGO_ENABLED=0 GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc go build -o mqtts_windows_386.exe main.go`







