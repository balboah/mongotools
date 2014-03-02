package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/balboah/mongotools/bson"
	"io"
	"labix.org/v2/mgo"
	mbson "labix.org/v2/mgo/bson"
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

func cmdDump() <-chan mbson.M {
	if *in == "" {
		log.Fatal("No server specified")
	}

	c := make(chan mbson.M)
	go func() {
		defer close(c)
		s, err := mgo.Dial(fmt.Sprintf("%s?connect=direct", *in))
		if err != nil {
			log.Fatal(err)
		}
		defer s.Close()

		db := s.DB("")
		if *col == "" {
			// TODO: Dump all collections by default
			// cols, _ := db.CollectionNames()
			log.Fatal("No collection specified")
		}
		iter := db.C(*col).Find(nil).Iter()
		var result interface{}
		for iter.Next(&result) {
			fmt.Println(result)
			c <- result.(mbson.M)
		}
		if iter.Timeout() {
			log.Fatal("Cursor timed out")
		}
		if err := iter.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	return c
}

func main() {
	flag.Parse()

	cmd := flag.Arg(0)
	switch cmd {
	case "restore":
		cmdRestore()
	case "dump":
		if *out == "" {
			log.Fatal("No output specified")
		}
		dumpFile, err := os.Create(*out)
		defer dumpFile.Close()
		if err != nil {
			log.Fatal(err)
		}

		for o := range cmdDump() {
			id, ok := o["_id"].(mbson.ObjectId)
			if !ok {
				log.Println("Could not find the object id")
			}
			log.Println("Save to", id.Hex())
			if err := bson.MarshalToStream(dumpFile, o); err != nil {
				log.Fatal(err)
			}
		}
	default:
		log.Fatal("Missing argument, must be one of: dump, restore")
	}
}
