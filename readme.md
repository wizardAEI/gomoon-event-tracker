# MAC

使用 gvm 安装不同架构的 go

arm64 架构：

安装 1.22.5 ( `gvm install go1.22.5` )

```bash
gvm use --default go1.22.5
go build
```

amd64 架构：

使用 Rosetta 启动 amd 终端安装 go 1.22.4 (安装amd版本的go `gvm install go1.22.4 -B`)

更改 go 版本 ( `gvm use go1.22.4` )

更改 go.mod 文件中的 go 版本为 1.22.4

```bash
go build -o eventTracker_x64 main.go
```

# Win

`go build -o eventTracker.exe main.go`
