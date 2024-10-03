package logs

import "log"

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

func ErrorLog(message string) {
	log.Println(Red + message + Reset)
}

func WarnLog(message string) {
	log.Println(Yellow + message + Reset)
}

func OkLog(message string) {
	log.Println(Green + message + Reset)
}

func MessageLog(message string) {
	log.Println(Blue + message + Reset)
}
