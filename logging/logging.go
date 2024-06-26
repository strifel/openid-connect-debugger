package logging

import (
	"fmt"
	color "github.com/fatih/color"
	"log"
	"os"
)

var (
	debugLogger  = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	infoLogger  = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	successLogger  = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	warnLogger  = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "", log.Ldate|log.Ltime)
)

func Debug(format string, v ...interface{}){
	message := fmt.Sprintf(format, v...)
	debugLogger.Println(fmt.Sprintf("%s %s", color.MagentaString("[DEBUG]"), message))
}

func Info(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	infoLogger.Println(fmt.Sprintf("%s %s", color.CyanString("[INFO]"), message))
}
func Success(format string, v ...interface{}){
	message := fmt.Sprintf(format, v...)
	debugLogger.Println(fmt.Sprintf("%s %s", color.GreenString("[OK]"), message))
}

func Warn(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	warnLogger.Println(fmt.Sprintf("%s %s", color.YellowString("[WARN]"), message))
}

func Error(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	errorLogger.Println(fmt.Sprintf("%s %s", color.RedString("[ERROR]"), message))
}