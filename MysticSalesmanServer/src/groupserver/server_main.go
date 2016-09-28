package main

import "flag"

func main() {
	flag.Parse()
	listeningAddress := flag.String("H", ":6666", "Listening Address")
	app, err := NewServerApp(*listeningAddress)
	if err != nil {
		panic(err)
	}
	app.Run()
}
