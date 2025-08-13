package logger

import (
	"context"
	"fmt"
	"log"
	"os"
)

const (
	Debug = "DEBUG"
	Info  = "INFO"
	Warn  = "WARN"
	Err   = "ERROR"
)

func InitLogger(ctx context.Context, logCh chan string) {
	log := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	for {
		select {
		case msg, ok := <-logCh:
			if !ok {
				log.Println(FormatLog(Warn, "channel closed"))
				return
			}
			log.Println(msg)
		case <-ctx.Done():
			for {
				select {
				case msg := <-logCh:
					log.Println(msg)
				default:
					log.Println(FormatLog(Info, "logger stopped"))
					return
				}
			}
		}
	}
}

func FormatLog(level, msg string) string {
	return fmt.Sprintf("[%s] %s", level, msg)
}
