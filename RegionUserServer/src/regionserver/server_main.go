package main

import "flag"

func main() {
	listeningAddress := flag.String("H", ":60000", "Listening Address")
	flag.Parse()

	app, err := NewServerApp(*listeningAddress)
	if err != nil {
		panic(err)
	}

	app.Run()
}
