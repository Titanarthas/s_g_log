package main

import (
	"s_g_log/s_g_log"
	"sync"
	"time"
	"fmt"
	"os"
	"log"
)

const (
	g_r_max = 50000
	s_r_max = 15
)

func main(){
	time1 := time.Now()

	s_g_log_test()

	time2 := time.Since(time1)

	fmt.Println(time2)

	time3 := time.Now()

	log_test()

	time4 := time.Since(time3)

	fmt.Println(time4)
}

func s_g_log_test(){

	var wg sync.WaitGroup
	for i := 0; i < g_r_max; i++ {
		wg.Add(1)

		go func(j int) {
			for k := 0; k < s_r_max; k++ {
				s_g_log.Infof("test j %d k %d", j, k)
			}
			wg.Add(-1)
		}(i)
	}



	wg.Wait()
	s_g_log.Exit()  // 确保内存中所有日志被写入到文件
}


func log_test(){
	fileName := "./log/glog.log"
	logFile,err  := os.Create(fileName)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error")
	}
	debugLog := log.New(logFile,"[Info]",log.Llongfile)

	var wg sync.WaitGroup
	for i := 0; i < g_r_max; i++ {
		wg.Add(1)

		go func(j int) {
			for k := 0; k < s_r_max; k++ {
				debugLog.Println("test j %d k %d", j, k)
			}
			wg.Add(-1)
		}(i)

	}

	wg.Wait()
}
