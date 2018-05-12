# log_process
Go项目练习， 日志分析系统，能够分析在某个协议下的某个请求在某个请求方法的QPS&amp;响应时间&amp;流量


## 说明
首先，我们定义了一个`LogProcess` 结构体用以封装整个模块的基本功能。定义了两个`wc`和`rc` `channel` 用以在三个模块之间进行数据同步。定义了名叫`read`的读取器（Ｒeader）和名叫`write`的写入器(Writer)用以引入读写模块。然后，将读写模块抽离并单独实现。同时，实现了数据解析部分。在`main`函数中，我们实例化了读取器和写入器，并将这两个参数注入到`LogProcess`结构体中。最后，分别在三个协程执行三个模块。

我们利用监控器对象，向外暴露监控系统监控的一些信息： 例如 总处理日志行数，系统吞吐量， read&write channel，总运行时间和错误行数。
这些信息都是通过HTTP服务器暴露出来。


## 安装influxdb

```shell
$ docker pull influxdb
# 最简单的启动方式
$ docker run influxdb
```

[how to use influxdb docker image](https://hub.docker.com/_/influxdb/)


## 写入模块的实现
1. 初始化influxdb client
2. 从Write Channel中读取监控数据
3. 构造数据并写入influxdb


## how to use 
1. install `grafana` to monitor `influxdb` which will show you more information
2. Through `http://127.0.0.1:9193/monitor`, you will get the monitor infomation of this system.

```shell
$ go run log_process.go ./access.log http://127.0.0.1:8086@bruce@bruce@log_process@s

$ go run data_generate.go 
```

the shell scripts are used to pull and start `grafana` and `influxdb`

When starting the system firstly, you must create the database `log_process`, 
the user with username and password, `bruce`, and you could select infomation from table named `log_info` 

