package main

import (
	"flag"
	"fmt"

	"urllookupservice/common"
	"urllookupservice/routers"
	"urllookupservice/server"
)

func main() {
	fmt.Println("Starting URL lookup service...")

	port := flag.String("port", common.DefaultUrlLookUpServicePort, "port at which server will listen to")
	flag.Parse()

	s := server.NewServer(*port)
	routers.CreateRouter(s)
	s.StartServer()
}
