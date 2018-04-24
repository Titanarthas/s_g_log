package s_g_log

import (
	"log"
	"os"
	"time"
	"fmt"
	"runtime"
	"sync"
	"encoding/json"
	"strings"
)

const (
	LOG_DEBUG int = iota
	LOG_FINE
	LOG_INFO
	LOG_WARNING
	LOG_ERROR
)
var msgChan (chan []byte)
var logLevel int
var wg sync.WaitGroup
var ioFile *os.File

type configuration struct {
	Level int    // 这个地方必须大写，要不然解析不到
	Filename    string
	Log_msg_list_max int
	File_count_max int
}

var conf configuration
func init() {

	file, _ := os.Open("conf/logconf.json")
	defer file.Close()

	if bExist,_ := PathExists("conf/logconf.json"); !bExist {
		fmt.Printf("conf/logconf.json is not exist")
	}

	decoder := json.NewDecoder(file)
	conf = configuration{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Error:", err)
		conf.Filename = "s_g_log.log"
		conf.Level = 0
	}

	if conf.Log_msg_list_max < 1024 {
		conf.Log_msg_list_max = 1024
	}

	if conf.File_count_max < 1024 {
		conf.File_count_max = 1024
	}

	msgChan = make(chan []byte, conf.Log_msg_list_max)

	openLogFile()

	logLevel = conf.Level

	wg.Add(1)
	go DisPathLogMsg()
}

func DisPathLogMsg() {
	defer wg.Add(-1)
	i := 0
	for msg := range msgChan {
		ioFile.Write(msg)

		i++
		if i > conf.File_count_max {
			ioFile.Close()
			openLogFile()
			i = 0
		}
	}
}

func openLogFile(){
	fileName := conf.Filename
	if bExist,_ := PathExists(fileName); bExist {
		strNewPath := fileName + GetMillisecondTime()
		os.Rename(fileName, strNewPath)
	}

	var err error
	ioFile,err  = CreateFile(fileName)

	if err != nil {
		log.Fatalln("open file error")
	}
}

func CreateFile(name string) (*os.File, error) {
	// make sure dir is exist
	parentDir := substr(name, 0, strings.LastIndex(name, "/"))

	if len(parentDir) > 0 {
		os.MkdirAll(parentDir, 777)
	}
	// create file
	return os.Create(name)
}

func substr(s string, pos, length int) string {
	if length < 0 {
		return ""
	}
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func Exit() {
	close(msgChan)
	wg.Wait()
}

func Printf(format string, v ...interface{}) {
	output(2, fmt.Sprintf(format, v...))
}

func Debugf(format string, v ...interface{}) {
	if logLevel <= LOG_DEBUG {
		output(2, fmt.Sprintf("%s %s", "LOG_DEBUG", fmt.Sprintf(format, v...)))
	}
}

func Finef(format string, v ...interface{}) {
	if logLevel <= LOG_FINE {
		output(2, fmt.Sprintf("%s %s", "LOG_FINE", fmt.Sprintf(format, v...)))
	}
}

func Infof(format string, v ...interface{}) {
	if logLevel <= LOG_INFO {
		output(2, fmt.Sprintf("%s %s", "LOG_INFO", fmt.Sprintf(format, v...)))
	}
}

func Warningf(format string, v ...interface{}) {
	if logLevel <= LOG_WARNING {
		output(2, fmt.Sprintf("%s %s", "LOG_WARNING", fmt.Sprintf(format, v...)))
	}
}

func Errorf(format string, v ...interface{}) {
	if logLevel <= LOG_ERROR {
		output(2, fmt.Sprintf("%s %s", "LOG_ERROR", fmt.Sprintf(format, v...)))
	}
}


func ErrorStackf(format string, v ...interface{}) {
	// if logLevel <= LOG_ERROR
	{
		var bufStack []byte
		buf := make([]byte, 1024)
		for {
			n := runtime.Stack(buf, false)
			if n < len(buf) {
				bufStack = buf[:n]
				break
			}
			buf = make([]byte, 2*len(buf))
		}
		output(2, fmt.Sprintf("%s %s \nstack is \n%s", "LOG_ERROR", fmt.Sprintf(format, v...), bufStack))
	}
}

func output(calldepth int, s string)  {
	now := time.Now() // get this early.
	var file string
	var line int
	var pc uintptr
	var funcname string

	{
		var ok bool
		pc, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}

		funcname = runtime.FuncForPC(pc).Name()
	}
	buf := make([]byte, 0)
	formatHeader(&buf, now, file, line, funcname)
	buf = append(buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		buf = append(buf, '\n')
	}
	// _, err := l.out.Write(l.buf)
	msgChan <- buf
	// return err
}

func formatHeader(buf *[]byte, t time.Time, file string, line int, funcname string) {
	// if l.flag&Ldate != 0
	{
		year, month, day := t.Date()
		itoa(buf, year, 4)
		*buf = append(*buf, '/')
		itoa(buf, int(month), 2)
		*buf = append(*buf, '/')
		itoa(buf, day, 2)
		*buf = append(*buf, ' ')
	}
	// if l.flag&(Ltime|Lmicroseconds) != 0
	{
		hour, min, sec := t.Clock()
		itoa(buf, hour, 2)
		*buf = append(*buf, ':')
		itoa(buf, min, 2)
		*buf = append(*buf, ':')
		itoa(buf, sec, 2)
		// if l.flag&Lmicroseconds != 0
		{
			*buf = append(*buf, '.')
			itoa(buf, t.Nanosecond()/1e3, 6)
		}
		*buf = append(*buf, ' ')
	}

	// if l.flag&(Lshortfile|Llongfile) != 0
	{
		// if l.flag&Lshortfile != 0
		{
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ':')
		*buf = append(*buf, funcname...)
		*buf = append(*buf, ": "...)
	}
}

func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetMillisecondTime() string {
	now := time.Now()
	year, mon, day := now.Date()
	hour, min, sec := now.Clock()
	zone := now.Nanosecond()

	return fmt.Sprintf("%d-%d-%d-%02d-%02d-%02d-%d", year, mon, day, hour, min, sec, zone / 1e6)
}
