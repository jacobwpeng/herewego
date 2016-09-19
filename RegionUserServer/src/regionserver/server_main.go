package main

import "flag"

func main() {
	listeningAddress := flag.String("H", ":60000", "Listening Address")
	flag.Parse()

	//serverState := newServerState()
	//serverState.UpdateUser(2191195, 100, 1)
	//serverState.UpdateUser(2191195, 100, 2)
	//fmt.Println("PickUser", serverState.PickUser(0, 9, 2))
	//status := serverState.GetRegionStatus()
	//fmt.Println(status)
	app, err := NewServerApp(*listeningAddress)
	if err != nil {
		panic(err)
	}

	app.Run()
}
