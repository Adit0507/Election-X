package main

import (
	"log"

	"gopkg.in/mgo.v2"
)

var db *mgo.Session

// connect and disconnet from local running mongodb instance
func dialdb() error {
	var err error
	log.Println("dialing mongodb: localhost")

	db, err = mgo.Dial("localhost")
	return err
}

func closeDb() {
	db.Close()
	log.Println("closed database connection")
}

type poll struct {
	Options []string
}

// loads the poll objects and extract options from documents,
// which will be used to search X
func loadOptions() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()

	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}

	iter.Close()
	return options, iter.Err()
}

func main() {}
