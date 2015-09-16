package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"

	. "github.com/etsy/mixer/config"
	. "github.com/etsy/mixer/db"

	"github.com/etsy/mixer/Godeps/_workspace/src/github.com/gorilla/mux"
)

type Handlers struct {
}

func NewHandlers() (*Handlers, error) {

	/*err := config.Config.Load()*/
	/*if err != nil {*/
	/*fmt.Println("error reading or parsing config config:", err)*/
	/*}*/

	h := &Handlers{}
	return h, nil
}

func (h *Handlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Mixer"}
	t, _ := template.ParseFiles("static/index.html")
	t.Execute(w, p)
}

func (h *Handlers) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	fmt.Fprint(w, "custom 404")
}

func (h *Handlers) AuthHandler(w http.ResponseWriter, r *http.Request) {
	username := parseAuth(r)
	p := GetPersonDataFromUsername(username)
	myjson, _ := json.Marshal(p)
	w.Write([]byte(myjson))
}

func parseAuth(r *http.Request) (userName string) {
	extractHeader := func(key string) (value string) {
		if val, exists := r.Header[key]; exists {
			value = val[0]
		}
		return
	}
	userName = extractHeader(Config.Userauth.Header)
	return
}

func (h *Handlers) PersonPutHandler(w http.ResponseWriter, r *http.Request) {
	var p Person
	vars := mux.Vars(r)
	person_id, _ := strconv.ParseInt(vars["id"], 10, 64)
	p.Id = person_id

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(body, &p)
	if err != nil {
		fmt.Println(err)
	}

	UpdatePerson(p)
}

func (h *Handlers) PersonPostHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	var p Person
	err = json.Unmarshal(body, &p)
	if err != nil {
		fmt.Println(err)
	}

	p = InsertPerson(p)
	myjson, _ := json.Marshal(p)
	fmt.Println(string(myjson))
	w.Write([]byte(myjson))
}

func (h *Handlers) PersonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	person_id, _ := strconv.ParseInt(vars["id"], 10, 64)
	p, err := GetPersonData(person_id)
	if err != nil {
		// return some 404 status to the app
		fmt.Println(err)
		w.WriteHeader(404)
		return
	}
	myjson, _ := json.Marshal(p)
	w.Write([]byte(myjson))
}

func (h *Handlers) AllPeopleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupname := vars["group"]
	p := GetPeopleData(groupname)
	myjson, _ := json.Marshal(p)
	w.Write([]byte(myjson))
}

func (h *Handlers) StaffRedirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	url := fmt.Sprintf("%s/%s", Config.Staff.DirectoryUrl, username)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handlers) MixerHandler(w http.ResponseWriter, r *http.Request) {
	g, _ := GetMixers()
	myjson, _ := json.Marshal(g)
	w.Write([]byte(myjson))
}
