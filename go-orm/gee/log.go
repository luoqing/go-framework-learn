package gee

import (
	"fmt"
	"sync"
	"time"
)

//errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
//logger.SetOutput(os.Stdout) 或者 logger.SetOutput(fd)

var (
	LEVEL_FLAGS = [...]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	TRACE = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

type Record struct {
	time  string
	code  string
	info  string
	level int
}

func (r *Record) String() string {
	return fmt.Sprintf("[%s][%s][%s] %s\n", LEVEL_FLAGS[r.level], r.time, r.code, r.info)
}

type Writer interface {
	Init() error
	Write(*Record) error
}

type Rotater interface {
	Rotate() error
	SetPathPattern(string) error
}

type Flusher interface {
	Flush() error
}

type Logger struct {
	writers     []Writer
	tunnel      chan *Record
	level       int
	lastTime    int64
	lastTimeStr string
	c           chan bool
	layout      string
	recordPool  *sync.Pool
}

func (l *Logger) Register(w Writer) {
	if err := w.Init(); err != nil {
		panic(err)
	}
	l.writers = append(l.writers, w)
}

func Info(format string, args ...interface{}) {
	r := &colorRecord{
		code:  "0",
		time:  time.Now().Format("2006-01-02 15:04:05"),
		info:  fmt.Sprintf(format, args...),
		level: INFO,
	}

	fmt.Println(r.String())
}

func Error(format string, args ...interface{}) {
	r := &colorRecord{
		code:  "999",
		time:  time.Now().Format("2006-01-02 15:04:05"),
		info:  fmt.Sprintf(format, args...),
		level: ERROR,
	}

	fmt.Println(r.String())

}

func Fatal(format string, args ...interface{}) {
	r := &colorRecord{
		code:  "-1",
		time:  time.Now().Format("2006-01-02 15:04:05"),
		info:  fmt.Sprintf(format, args...),
		level: FATAL,
	}

	fmt.Println(r.String())

}

type colorRecord Record

func (r *colorRecord) String() string {
	switch r.level {
	case TRACE:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[34m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL_FLAGS[r.level], r.code, r.info)
	case DEBUG:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[34m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL_FLAGS[r.level], r.code, r.info)

	case INFO:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[32m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL_FLAGS[r.level], r.code, r.info)

	case WARNING:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[33m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL_FLAGS[r.level], r.code, r.info)

	case ERROR:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[31m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL_FLAGS[r.level], r.code, r.info)

	case FATAL:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[35m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL_FLAGS[r.level], r.code, r.info)
	}

	return ""
}
