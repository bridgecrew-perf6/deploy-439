# deploy
项目计划实现一个能够将本地项目压缩上传然后构建容器的工具  
## 使用方法
在config.json配置项目的基本信息  
在server.json配置上传服务器的账户信息  
``` JSON
"demo" : {
    "folder" : "D:/programming/web",
    "hostname" : "192.168.33.12",
    "ssh_port" : 22,
    "run_command": "docker run -p 8080:8080 -d web:test"
}
```
执行    
``` go
  go run main.go -p demo
```
即可将项目推送到hostname对应的服务器上，并自动构建容器  
这里run_command 是构建容器后执行的运行容器的命令  
