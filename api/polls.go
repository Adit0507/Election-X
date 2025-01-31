package main

import (
	"errors"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// descirbes the polls we have created
type poll struct {
	ID      bson.ObjectId     `bson:"_id" json:"id"`
	Title   string            `json:"title"`
	Options []string          `json:"options"`
	Results map[string]string `json:"results,omitempty"`
	APIKey  []string          `json:"apikey"`
}

func (s *Server) handlePolls(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handlePollsGet(w, r)
		return

	case "POST":
		s.handlePollsPost(w, r)
		return

	case "DELETE":
		s.handlePollsDelete(w, r)
		return
	}

	respondHTTPErr(w, r, http.StatusNotFound)
}

func (s *Server) handlePollsGet(w http.ResponseWriter, r*http.Request) {
	session := s.db.Copy()	//creatin copy of DB session that will allow us to interact with Mongo
	defer session.Close()
	
	c := session.DB("ballots").C("polls")
	
	var q *mgo.Query	
	p := NewPath(r.URL.Path)
	if p.HasID() {
		// gets specific poll
		q = c.FindId(bson.ObjectIdHex(p.ID))
	} else {
		q = c.Find(nil)	//get all polls
	}

	var res []*poll
	if err := q.All(&res); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, &res)
}


func (s *Server) handlePollsPost(w http.ResponseWriter, r*http.Request) {
	respondErr(w, r, http.StatusInternalServerError)
	errors.New("not implemented")
}
func (s *Server) handlePollsDelete(w http.ResponseWriter, r*http.Request) {
	respondErr(w, r, http.StatusInternalServerError)
	errors.New("not implemented")
}

