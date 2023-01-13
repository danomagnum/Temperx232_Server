package main

import (
	"log"
	"os"

	"github.com/danomagnum/gologix"
)

type InStr struct {
	Temperature float32
	Humidity    float32
	Spare       [32 - 8]byte
}
type OutStr struct {
	Spare [32]byte
}

var inInstance InStr
var outInstance OutStr

func class1serve() {
	r := gologix.PathRouter{}

	// define the Input and Output instances.  (Input and output here is from the plc's perspective)

	// an IO handler in slot 2
	//p3 := gologix.IOProvider[InStr, OutStr]{}
	p3 := gologix.IOProvider[InStr, OutStr]{
		In:  &inInstance,
		Out: &outInstance,
	}
	path3, err := gologix.ParsePath("1,0")
	if err != nil {
		log.Printf("problem parsing path. %v", err)
		os.Exit(1)
	}
	r.AddHandler(path3.Bytes(), &p3)

	s := gologix.NewServer(&r)
	s.Serve()
}
