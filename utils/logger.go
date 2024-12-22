package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

func Format(message, color string) {
	_, file, line, ok := runtime.Caller(2)

	if !ok {
		log.Fatalln("Unable to get caller")
	}
	cwd, err := os.Getwd()

	if err != nil {
		log.Fatalln("Unable to get cwd")
	}

	fmt.Printf("%s%s\033[0m |  \033[38;5;177m%s:%d \033[0m| \033[0m%s\n", color, time.Now().Format("15:04:05"), file[len(cwd)+1:], line, message)
}

func Log(input string) {
	Format(input, "\033[38;2;30;215;96m")
}
func Warn(input string) {
	Format(input, "\033[38;2;255;215;95m")
}
func Error(input string, panic bool) {
	Format(input, "\033[38;2;219;60;66m")
	if panic {
		os.Exit(0)
	}
}
