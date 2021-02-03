# Gallop

Go gin 快速开发脚手架，本仓库的实现思路和部分代码多数来源于[https://github.com/shenyisyn/goft-gin][goft-gin]

[goft-gin]: https://github.com/shenyisyn/goft-gin

Gallop主要是使用inject去贯传，这个项目并实现mvc架构，对gin进行封装。用于快速开发。

### 使用到的开源库,排名不分先后

1. gin-gonic/gin  https://gin-gonic.com/docs/
2. spf13/viper   https://github.com/spf13/viper
3. facebookarchive/inject https://github.com/facebookarchive/inject
4. jinzhu/copier https://github.com/jinzhu/copier.git
5. goframe https://goframe.org/display/gf


### 系统配置说明

1.集成了日志zap组件,会以json对象输出
```yaml
logger:
    type: file
    path: ./logs
    level: debug # 支持debug< info< warn< error< dpanic< panic< fatal
    stack: false # 当设置为true只有>=warn的等级才会显示stack
```
为了方便继承，我们增加了LogFiled属性，用户自定义日志输出如增加req_id
```go
filed:=logger.LogField{}
filed["req_id"]="xxxxxxxxxxxxx"
logger.Info("1233",filed)
//输出： {"level":"INFO","ts":"2020-08-28 14:26:32","func":"controller/hello.go:20","msg":"1233","req_id":"xxxxxxxxxxxxx"}
```