package main

import "flag"

func main() {
	listeningAddress := flag.String("H", ":6666", "Listening Address")
	managerAddress := flag.String("M", ":8080", "Manager HTTP Listening Address")
	flag.Parse()
	app, err := NewServerApp(*listeningAddress, *managerAddress)
	if err != nil {
		panic(err)
	}
	app.Run()
}
