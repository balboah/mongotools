package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/balboah/mongotools/restore"
	"io"
	"log"
	"os"
)

var dump = flag.String("dump", "", "filename of dump")

func main() {
	flag.Parse()

	if *dump == "" {
		log.Fatal("no dump specified")
	}
	input, err := os.Open(*dump)
	if err != nil {
		log.Fatal(err)
	}

	o := new(interface{})
	for err == nil {
		if *o != nil {
			b, err := json.MarshalIndent(o, "", "\t")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s,\n", string(b))
		}
		err = restore.UnmarshalFromStream(input, o)
	}
	if err != io.EOF {
		log.Println(err)
	}
}
