package main

import (
	"flag"
)

func main() {
	listenAddress := flag.String("H", ":50203", "Listening Address")
	flag.Parse()
	app, err := NewServerApp(*listenAddress)
	if err != nil {
		panic(err)
	}
	app.Run()
}
