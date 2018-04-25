# s_g_log
（simple_go_log？本来命名是t_g_log的，t表示公司名，mmp，怕和谐）看了下官方提供的log比较烂，加锁太多，锁太多一方面是写日志慢，一方面在协程数量很多的情况下，大部分的协程会由于阻塞被挂起，整个程序效率大幅降低。

# 功能说明

1，日志分级和配置功能功能，调试的时候日志级别配低利于调试，线上运行的时候日志级别调高以减少不必要的日志

2，日志太大时分文件功能。

3，发生致命性错误时可以选择是否打印当前堆栈信息，调用ErrorStackf函数会打印堆栈信息，Errorf不会打印。

# 使用方法

1，下载包

go get github.com/Titanarthas/s_g_log/

2，import

在需要使用的文件中

import (
	"github.com/Titanarthas/s_g_log/s_g_log"
)

3，复制github.com/Titanarthas/s_g_log/目录下的conf目录到运行目录，其中的config文件为日志配置文件。

# 关于配置字段的说明

配置文件为json格式，例如

{

    "level": 1,  // 日志等级
  ，日志等级越高，打印的日志越少，debug为0，最低。
  
  "filename": "./log/long_connect_server_go.log",  // 日志文件名
  
  "log_msg_list_max": 16777216,  // 日志
  最大缓存  
  
  "file_count_max": 16777216  // 单文件最大日志条数
  
}

# 使用样例
package main

import (
	"github.com/Titanarthas/s_g_log/s_g_log"
)

func main(){

	defer s_g_log.Exit()
	
	s_g_log.Finef("aa")
	
	s_g_log.Warningf("bb")
	
	s_g_log.Errorf("err %d", 1)
	
}

# 测试对比
|   | 协程数量  |每个协程日志数量   |总耗时   |
| ------------ | ------------ | ------------ | ------------ |
| 官方log库  | 50000  |15   | 8.0730209s  |
|  s_g_log | 50000  | 15 | 3.359452s  |
