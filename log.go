package log

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// TODO: Включена или нет отладка берем из файла настроек
//

type Log struct {
	ServiceName string
	DebugMode   bool
	StartTime   time.Time
}

// Error - запись в [service]_err.log - критичные для выполнения ошибки
func (l *Log) Error(msg string) {
	fmt.Println("error", msg)
	l.write("error", msg)
}

// Info - запись в [service]_info.log - протокол работы (+ некритичные worning)
func (l *Log) Info(msg string) {
	fmt.Println(msg)
	l.write("info", msg)
}

// Debug - запись в [service]_debug.log - все что нужно на этапе отладки
func (l *Log) Debug(msg string) {
	if l.DebugMode {
		fmt.Println(msg)
		l.write("debug", msg)
	}
}

//Start - initialize ServiceName and start timer
func (l *Log) Start(serviceName string) {
	l.ServiceName = serviceName
	l.write("info", "Zzz")
	l.StartTime = time.Now()
}

// End add log run duration
func (l *Log) End() {
	duration := time.Since(l.StartTime)
	l.Info(duration.String())
}

//to do: сделать для многопоточной работы через каналы
// Error, Info, Debug - пишут в свои каналы
// обрботчики каналов их слушают и в отдельных горутинах пишут в файлы
// go Write - слушает канал и записывает в файлы

func (l *Log) write(logType string, msg string) {
	fileName := l.ServiceName + "_" + logType + ".log"
	//fileName := "C:\\" + l.ServiceName + "_" + logType + ".log"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	curTime := time.Now()
	const formatDate = "2006-01-02 15:04:05"
	msg = curTime.Format(formatDate) + ";" + sToCSV(msg) + ";" + "\n"
	if _, err = file.WriteString(msg); err != nil {
		fmt.Println(err)
	}

}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// сделать CSV-валидным

func sToCSV(st string) string {

	st = strings.Replace(st, "\"", "\"\"", -1)
	if strings.ContainsAny(st, "\"") || strings.ContainsAny(st, ";") || strings.ContainsAny(st, "\n") {
		st = "\"" + st + "\""
	}
	return st
}
