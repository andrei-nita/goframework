package framework

import (
	"log"
	"os"
)

var (
	file    *os.File
	logfile bool
)

func CreateLog(logToFile bool) error {
	var err error
	logfile = logToFile

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if logfile {
		file, err = os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		if err != nil {
			return err
		}
		log.SetOutput(file)
	} else {
		log.SetOutput(os.Stdout)
	}
	return err
}

func CloseLog() {
	if logfile {
		_ = file.Close()
	}
}
