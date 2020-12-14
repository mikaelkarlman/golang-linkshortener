package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var links = map[string]string{
	"troll":  "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
	"jebait": "https://www.youtube.com/watch?v=d1YBv2mWll0",
	"google": "https://google.com",
}

//Index server generic page
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

//MatchLinkID checks the ID, matches it and redirects if the key exists
func MatchLinkID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var linkID string = ps.ByName("linkID")

	if val, ok := links[linkID]; ok {

		http.Redirect(w, r, val, http.StatusMovedPermanently)

	} else {

		fmt.Fprintf(w, "LinkID: "+string(linkID)+" does not exist")

	}
}

func main() {

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/:linkID", MatchLinkID)

	log.Println(http.ListenAndServe(":8080", router))

}
