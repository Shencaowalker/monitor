#维护好服务生产者列表
#airserviceproducersum  currentcount跟normalcount相同时，指标值为1； currentcount不等于normalcount的值时，指标值为0,告警；
#aironserviceproducerhealthycount  healthycount跟normalcount相同时，指标值为2； healthycount小于normalcount的值时但是healthycount不为1时，指标值为1,轻度警告要处理；healthycount为0或者当前生产者不存在于配置文件中时，指标值为0,严重报警；

groups:
- name: serviceproducer
  rules:
  - alert: serviceproducer_not_available
    expr: airserviceproducersum == 0
    for: 2m
    labels:
      severity: warning
    annotations:
      description: "{{ $labels.name }} 2分钟内生产者种类于配置文件不符，存在生产者丢>失已不能提供部分服务,需及时处理。"
      summary: "{{ $labels.name }} 生产者种类恢复正常。"
  - alert: serviceproducercount_less_than_normal
    expr: aironserviceproducerhealthycount != 2
    for: 2m
    labels:
      severity: warning
    annotations:
      description: "{{ $labels.name }} 2分钟内生产者数量为{{ $labels.healthytcount }}，正常数量为{{ $labels.normalcount }}，请及时处理。"
      summary: "{{ $labels.name }} 生产者数量已恢复."




#####config.ini 最原始配置,配置丢失后可以使用下面来替换。

[global]
cmbalarminterface=http://10.5.100.4:5002/get?message=
consulipport=10.5.101.3:8500
delayupdateseconds=300
nacosip=nacos.local.aloudata.work
nacosport=8848
namespaceid=test050
pageno=1
pagesize=100
pushgatewayipport=10.5.101.3:9091
servicelist=dws cip das afp asp bde tse arctic jobserver
serviceport=8501
    #说明
        cmbalarminterface   佳泺提供的往招呼群发送信息的接口  接收get请求(暂时可以不配置)
        consulipport        consul 的ip:port
        delayupdateseconds  等待多长秒把运行中的nacos配置更新到本地，主要解决发版完毕之后，生产者无法立刻启动。
        nacosip             nacos ip
        nacosport           nacos port
        namespaceid         nacos 命名空间名称
        pageno              nacos 页数 一般不用改
        pagesize            nacos 一页有多少生产者、消费者列表，一般100差不多
        pushgatewayipport   pushgateway的ip:port
        servicelist         服务列表（nacos上没有air）
        serviceport         服务监听端口



####接口说明
## 读取当前nacos运行的生产者状态到本地基础配置 get请求,上线完毕之后需要手动触发一次。
/updatenacosstandardconf

## 向consul注册任务 post请求 post接口接收json数据
    /registered
        #请求体   注册jvm服务  需要修改  id  env["dev","test"]  app_type["air","as","big"]
        {
            "id":"dws",
            "group":"jvm_info",
            "address":"10.5.101.3",
            "port":"13100",
            "env":"all",
            "m_type":"tcp",
            "app_type":"air"
        }

        #请求体   注册node服务  需要修改  id  env["dev","test"]  app_type["air","as","big"]
        {
            "id":"dws",
            "group":"node",
            "address":"10.5.101.3",
            "port":"9100",
            "env":"all",
            "m_type":"tcp",
            "app_type":"air"
        }

        #请求体   注mid服务  需要修改:id port env["dev","test"]  app_type["air","as","big"]
        {
            "id":"dws",
            "group":"mid",
            "address":"10.5.101.3",
            "port":"5005",
            "env":"mid",
            "m_type":"tcp",
            "app_type":"air"
        }


## 下线consul注册的服务  接收x-www-form-urlencoded类型的post请求或者普通get请求
    /downlin
        #请求例子  删除consul_exporter_all_10.5.101.3_13100
        http://127.0.0.1:8501/downline?id=consul_exporter_all_10.5.101.3_13100




## 更新生产者指标
curl -X get http://127.0.0.1:8502/updateNacosproducerMonitor

## 更新生产者基线（每次上线需要手动触发一次）
curl -X get http://127.0.0.1:8502/updatenacosstandardconf

## 获取日志报错指标
curl -X get http://127.0.0.1:8502/updatetargetlogMetrics


### 日志相关配置
[big_logs]
label_list={"project": "as","environment":"joinsight-testing"}
label_name=starttime traceid status api abstractaaa
lokiexclre=
lokire=ERROR
recordslimit=100000
regex=(\d+-\d+-\d+\s\S+)\s\[(.*?)\]\s\[.*?\]\s(\w+)\s+(\S+)\s-\s(.*)

[can_logs]
label_list={"project": "as","environment":"joinsight-testing"}
label_name=starttime traceid status api abstractaaa
lokiexclre=success=true
lokire=ERROR
recordslimit=100000
regex=(\d+-\d+-\d+\s\S+)\s\[(.*?)\]\s\[.*?\]\s(\w+)\s+(\S+)\s-\s(.*)

[dingdingwebhook]
listrebots={"air":{"test":{"accessToken":"dd0dfc3e5f598a4c94b6bcfebbfc1d88281898b9658642294d6ff60f11a4f149","secret":"SECf20d8b02a0dee2df7dd02d2ff55498d0f7508fb234e0cadc4e5adace2e74acfe"},"prod":{"accessToken":"","secret":""}},"joinsight":{"test":{"accessToken":"dd0dfc3e5f598a4c94b6bcfebbfc1d88281898b9658642294d6ff60f11a4f149","secret":"SECf20d8b02a0dee2df7dd02d2ff55498d0f7508fb234e0cadc4e5adace2e74acfe"},"prod":{"accessToken":"","secret":""}},"bigmeta":{"test":{"accessToken":"dd0dfc3e5f598a4c94b6bcfebbfc1d88281898b9658642294d6ff60f11a4f149","secret":"SECf20d8b02a0dee2df7dd02d2ff55498d0f7508fb234e0cadc4e5adace2e74acfe"},"prod":{"accessToken":"","secret":""}},"default":{"default":{"accessToken":"8f3a1b5cd395675557bbd7a5fb80c666b66866f49f48a6d165f9ae3fa147ef28","secret":""}}}

[global]
cmbalarminterface=http://10.5.100.4:5003/get?message=
consulipport=10.5.101.3:8500
delayupdateseconds=15
nacosip=10.16.0.106
nacosport=8848
namespaceid=can-demo
pageno=1
pagesize=100
pushgatewayipport=10.5.101.3:9091
servicelist=jscip semantic
serviceport=8502

[mixedformat]
collectionscopeseconds=3600
label_len=128
latencycollectionseconds=2000
logs=can_logs big_logs
lokiipport=10.5.20.35:3100