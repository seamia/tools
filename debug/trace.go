package debug

import (
	"os"
	"fmt"
	"bufio"
	"time"
)

var (
	enabled = false
	logFileName = ""
)

func Trace(format string, a ...interface{}) {
	name := os.Getenv("LogFile")
	if len(name) == 0 {
		name = os.Args[0] + ".log"
	}
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		defer f.Close()

		w := bufio.NewWriter(f)
		fmt.Fprintf(w, format, a...)
		w.Flush()
	}
}


func init() {
	enabled = true
	name := os.Getenv("LogFile")
	if len(name) == 0 {
		name = os.Args[0] + ".log"
	} else if name == "<none>" {
		enabled = false
	}

	if enabled {
		logFileName = name
		Trace("\n\n------------ app: %s:%d ------------ %s\n",
			os.Args[0],
			os.Getpid(),
			time.Now().Format(time.RFC850))
	}
}