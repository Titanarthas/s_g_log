# s_g_log
看了下官方提供的log比较烂，加锁太多。按照c的思想撸了一个版本的log。

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
