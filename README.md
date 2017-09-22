 Go Client for Apollo
================

* Source Code: https://github.com/zouyx/agollo/
* Fork 主要改了时间的配置，过期时间为7天，每次`StartRefreshConfig`都会更新过期时间。修正了timeout，此值需要大于30s，因为为长连接，需保持30s。

方便Golang接入配置中心框架 [Apollo](https://github.com/ctripcorp/apollo) 所开发的Golang版本客户端。

Installation
------------

如果还没有安装Go开发环境，请参考以下文档[Getting Started](http://golang.org/doc/install.html) ，安装完成后，请执行以下命令：

``` shell
gopm get github.com/cihub/seelog -v -g
gopm get github.com/coocood/freecache -v -g
gopm get github.com/zouyx/agollo -v -g
```

或者

``` shell
go get -u github.com/cihub/seelog
go get -u github.com/coocood/freecache
go get -u github.com/zouyx/agollo
```


*请注意*: 最好使用Go 1.8进行开发

# Features
* 实时同步配置
* 灰度配置
* 客户端容灾

# Usage

- 修改配置
> app.properties:
```json
{
    "appId": "SampleApp",
    "cluster": "default",
    "namespaceName": "application",
    "ip": "localhost:8080"
}
```  
> seelog.xml
```xml
<seelog type="sync" mininterval="2000000" maxinterval="100000000" critmsgcount="500" minlevel="debug">

    <outputs formatid="all">
        <console formatid="fmterror"/>
    </outputs>
    <formats>
        <format id="fmtinfo" format="[%Level]  [%Time] %Msg%n"/>
        <format id="fmterror" format="[%LEVEL] [%Time] [%FuncShort @ %File.%Line] %Msg%n"/>
        <format id="all" format="[%Level] [%Time]  [@ %File.%Line] %Msg%n"/>
        <format id="criticalemail" format="Critical error on our server!\n    %Time %Date %RelFile %Func %Msg \nSent by Seelog"/>
    </formats>
</seelog>
```

- 启动agollo

``` go
func main() {
  go agollo.Start()
  for {
    fmt.Println("apollo is :", agollo.GetIntValue("timeout", 100), agollo.GetStringValue("str", "100"))  
    #服务读取配置需要等待1s重，等agollo.Start()初始化完成
    time.Sleep(1 * time.Second)
    ttl, _ := agollo.GetApolloConfigCache().TTL([]byte("timeout"))
    fmt.Println("ttl is:",ttl)
  }
  time.Sleep(100 * time.Minute)
}
```

- 获取Apollo的配置
  - String
  
  ```
  agollo.GetStringValue(Key,DefaultValue)
  ```
  - Int
  
  ```
  agollo.GetIntValue(Key,DefaultValue)
  ```

  - Float
  
  ```
  agollo.GetFloatValue(Key,DefaultValue)
  ```

  - Bool
  
  ```
  agollo.GetBoolValue(Key,DefaultValue)
  ```
  



