package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/balboah/mongotools/bson"
	"io"
	"labix.org/v2/mgo"
	"log"
	"os"
)

var (
	in  = flag.String("in", "", "where to read from")
	out = flag.String("out", "", "what to write to")
	col = flag.String("collection", "", "targeted collection")
)

func cmdRestore() {
	if *in == "" {
		log.Fatal("no dump specified")
	}
	input, err := os.Open(*in)
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
		err = bson.UnmarshalFromStream(input, o)
	}
	if err != io.EOF {
		log.Println(err)
	}
}

func cmdDump() {
	if *in == "" {
		log.Fatal("No server specified")
	}
	if *out == "" {
		log.Fatal("No output specified")
	}
	s, err := mgo.Dial(fmt.Sprintf("%s?connect=direct", *in))
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	db := s.DB("")
	if *col == "" {
		log.Fatal("No collection specified")
	}
	// TODO: Dump all collections by default
	// cols, _ := db.CollectionNames()

	iter := db.C(*col).Find(nil).Iter()
	result := new(interface{})
	dumpFile, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}
	for iter.Next(result) {
		fmt.Println(*result)
		if err := bson.MarshalToStream(dumpFile, *result); err != nil {
			log.Fatal(err)
		}
	}
	if iter.Timeout() {
		log.Fatal("Cursor timed out")
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	cmd := flag.Arg(0)
	switch cmd {
	case "restore":
		cmdRestore()
	case "dump":
		cmdDump()
	default:
		log.Fatal("Missing argument, must be one of: dump, restore")
	}
}
