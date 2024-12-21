package utils

import (
	"fmt"
	"os"
	"time"
)

func Format(message, color string) {
	fmt.Printf("\033[1;37m%s\033[0m | %s%s\033[0m\n", time.Now().Format("15:04:05"), color, message)
}

func Log(str string) {
	Format(str, "\033[0m")
}
func Warn(str string) {
	Format(str, "\033[38;2;241;196;15m")
}
func Error(input string, panic bool) {
	Format(input, "\033[38;2;219;60;66m")
	if panic {
		os.Exit(0)
	}
}
