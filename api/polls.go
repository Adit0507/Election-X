package main

import (
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

	case "OPTIONS":
		w.Header().Add("Access-Control-Allow-Methods", "DELETE")
		respond(w, r, http.StatusOK, nil)
		return
	}

	respondHTTPErr(w, r, http.StatusNotFound)
}

// reading polls
func (s *Server) handlePollsGet(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy() //creatin copy of DB session that will allow us to interact with Mongo
	defer session.Close()

	c := session.DB("ballots").C("polls")

	var q *mgo.Query
	p := NewPath(r.URL.Path)
	if p.HasID() {
		// gets specific poll
		q = c.FindId(bson.ObjectIdHex(p.ID))
	} else {
		q = c.Find(nil) //get all polls
	}

	var res []*poll
	if err := q.All(&res); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, &res)
}

// creating polls
func (s *Server) handlePollsPost(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("ballots").C("polls")
	var p poll
	if err := decodeBody(r, &p); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed toread poll from request", err)
		return
	}
	apikey, ok := APIKey(r.Context())
	if ok {
		p.APIKey = []string{apikey}
	}
	p.ID = bson.NewObjectId()
	if err := c.Insert(p); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to insertpoll", err)
		return
	}
	w.Header().Set("Location", "polls/"+p.ID.Hex())
	respond(w, r, http.StatusCreated, nil)
}

func (s *Server) handlePollsDelete(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()

	c := session.DB("ballots").C("polls")
	p := NewPath(r.URL.Path)
	if !p.HasID() {
		respondErr(w, r, http.StatusMethodNotAllowed, "Cannot delete all polls")
		return
	}
	if err := c.RemoveId(bson.ObjectIdHex(p.ID)); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to delete poll", err)
		return
	}

	respond(w, r, http.StatusOK, nil)
}
