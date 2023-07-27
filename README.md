# Structure

<pre>
<code>
logagent
├─ README.md
├─ common
│    └─ common.go
├─ conf
│    └─ config.ini
├─ etcd
│    ├─ etcd.go
│    └─ example
│           └─ main.go
├─ go.mod
├─ go.sum
├─ kafka
│    ├─ consumerExample
│    │    └─ main.go
│    └─ kafka.go
├─ logger
│    └─ logger.go
├─ main.go
├─ setting
│    └─ setting.go
├─ sysinfo
└─ tailfile
     ├─ tailfile.go
     └─ tailfile_manager.go
</code>
</pre>

# Skill List

*   go channel 消息传输
*   logrus日志库
*   Kafka接收日志收集的消息

* etcd存储日志文件配置
* tail日志文件跟踪并收集



1. 初始化日志设置，读取配置信息，启动kafka和etcd；
2. 根据本机ip从etcd中读取对应的日志文件信息`path, topic`；
3. 解析日志信息并交给tail跟踪，启动goroutine监控etcd日志配置改动；
   1. 当配置信息发生改动时，启动task监听新日志，关闭被删除日志项的task；
4. tail启动多个task分别监控对应的日志文件，并将读取到的条目发送给kafka；
5. (TODO) log-transfer 将收集到的日志信息发送给ElasticSearch，Kibana等。