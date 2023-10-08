package main

import "github.com/matthisholleville/mapsyncproxy/server"

func main() {
	s := server.New()
	s.Start(":8000")
}
