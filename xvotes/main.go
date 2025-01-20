package main

import (
	"log"

	"github.com/nsqio/go-nsq"
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
	iter := db.DB("ballots").C("polls").Find(nil).Iter() //allows us to access each poll one by one
	// much mpre memory efficient coz it only uses single poll object

	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}

	iter.Close()
	return options, iter.Err()
}

// takes the votes channel and publish each string that is received from it
func publishVotes(votes <-chan string) <-chan struct{} {
	stopChan := make(chan struct{}, 1)

	pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	go func() {
		for vote := range votes {
			pub.Publish("votes", []byte(vote))	//publish vote
		}

		log.Println("Publisher: Stopping")
		pub.Stop()

		log.Println("Publisher: Stopped")
		stopChan <- struct{}{}
	}()

	return stopChan
}

func main() {}
