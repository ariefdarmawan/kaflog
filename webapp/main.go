package main

import (
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"git.eaciitapp.com/sebar/knot"
	"github.com/ariefdarmawan/kaflog"
	"github.com/eaciit/toolkit"
)

func main() {
	brokers, topic := func() ([]string, string) {
		bs := strings.Split(os.Getenv("KAFLOG_BROKERS"), ",")
		topic := os.Getenv("KAFLOG_TOPIC")
		return bs, topic
	}()

	if len(brokers) != 0 {
		toolkit.Printfn("Enable kafka logging. Topic %s, brokers: %v", topic, brokers)
	}
	host := ":9123"
	log, _ := toolkit.NewLog(true, true, "./logs", "$LOGTYPE_$DATE.log", "yyyyMMdd")
	log.AddHook(kaflog.Hook(host, topic, brokers...), "ERROR", "WARNING")

	s := knot.NewServer().SetLogger(log)
	s.AddPlugin(knot.NewPlugin("recoverer", func(next knot.KnotHandler) knot.KnotHandler {
		return func(ctx *knot.WebContext) {
			defer func() {
				if rec := recover(); rec != nil {
					s.Logger().Errorf("panic detected. %v", rec)
				}
			}()

			next(ctx)
		}
	}))

	s.Route("/action", func(ctx *knot.WebContext) {
		value := strings.ToLower(ctx.Request.URL.Query().Get("name"))

		if value == "error" {
			ctx.Write([]byte("You should see error from your kafka"), http.StatusOK)
			panic("Panic from /action")
		} else if value == "warning" {
			ctx.Write([]byte("You should see warning from your kafka"), http.StatusOK)
			s.Logger().Warning("warning from /action")
			return
		}

		ctx.Write([]byte("All is good"), http.StatusOK)
	})

	err := s.Start(host)
	if err != nil {
		toolkit.Printf("Error: %s \n", err.Error())
		return
	}

	cExit := make(chan os.Signal)
	signal.Notify(cExit, syscall.SIGINT)
	<-cExit

	s.Stop()
}
