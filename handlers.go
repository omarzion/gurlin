package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Bounce is http server root, redirects if path is found 404 otherwise
func Bounce(w http.ResponseWriter, r *http.Request) {
	src := mux.Vars(r)["src"]
	redirect, err := GetRedirect(src)
	if err {
		fmt.Fprintf(w, "¯\\_(ツ)_/¯")
	} else {
		http.Redirect(w, r, redirect.Dest, 301)
	}
}

// Index is the index of the api path
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

// Available is /{apiRoot}/available/{src} and returns Message indicating whether redirect path is used or not
func Available(w http.ResponseWriter, r *http.Request) {
	src := mux.Vars(r)["src"]
	err := SrcAvailable(src)
	msg := "available"
	if !err {
		msg = "taken"
	}

	if err := SendResponse(w, MakeMsg(fmt.Sprintf("%v %v", src, msg))); err != nil {
		panic(err)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	form := r.Form
	to := form.Get("to")
	from := form.Get("from")
	var msg Message

	switch {
	case err != nil:
		msg = MakeErr("Unknown error")
	case to == "":
		msg = MakeErr("Destination address required, did you forget to add 'to:'?")
	case from != "":
		to = FixDestination(to)
		from = strings.Replace(from, "/", "", -1)
		if re, ok := GetRedirect(from); !ok && to != re.Dest {
			msg = MakeErr(fmt.Sprintf("%v taken", from))
		} else {

			defer AddRedirect(Redirect{Src: from, Dest: to})
			msg = MakeMsg(fmt.Sprintf("%v -> %v", from, to))
		}
	default:
		to = FixDestination(to)
		from = GetAvailableSrc(to)
		defer AddRedirect(Redirect{Src: from, Dest: to})
		msg = MakeMsg(fmt.Sprintf("%v -> %v", from, to))
	}
	if err = SendResponse(w, msg); err != nil {
		panic(err)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
}

func Disable(w http.ResponseWriter, r *http.Request) {
}

func SendResponse(w http.ResponseWriter, msg interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(msg)
}
