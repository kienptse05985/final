package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/rs/cors"
)

var (
	configPrefix string
	configSource string
	config       Config
)

func main() {
	flag.Parse()
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			main()
		}
	}()

	err := LoadEnv(&config, configPrefix, configSource)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		fmt.Printf("Listen and serve proxy API %s\n", config.Binding)
		fmt.Println(http.ListenAndServe(config.Binding, NewApiV1()))
	}()

	Running()
}

func NewApiV1() http.Handler {

	router := gin.Default()
	router.POST("/api/v1/scan", ScanURL)

	return cors.Default().Handler(router)
}

func Running() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	for {
		select {
		case <-sig:
			log.Println("\nSignal received, stopping")
			return
		}
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.StringVar(&configPrefix, "configPrefix", "API", "config prefix")
	flag.StringVar(&configSource, "configSource", ".env", "config source")
}