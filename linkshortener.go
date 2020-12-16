package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

//Index serves generic page
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	http.ServeFile(w, r, "public/index.html")

}

//faviconHandler serves the favicon
func faviconHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	http.ServeFile(w, r, "favicon.ico")

}

func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

//addLink imports data from form to the database
func addLink(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "<h1>Error: %s</h1>\n", err)
	}

	fmt.Fprintf(w, "<h1>Submitted message!</h1>")
	fmt.Println(r.PostFormValue("linkID"))
	fmt.Println(r.PostFormValue("link"))

	fmt.Println(isValidURL(r.PostFormValue("link")))
}

//matchLinkID checks the ID, matches it and redirects if the key exists
func matchLinkID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var linkID string = strings.ToLower(ps.ByName("linkID"))

	if linkID != "favicon.ico" {

		fmt.Println("LinkID called: " + linkID)

		link := getDatabaseData(linkID)

		if len(link) != 0 {

			http.Redirect(w, r, link, http.StatusMovedPermanently)

		} else {

			fmt.Fprintf(w, "LinkID: "+string(linkID)+" does not exist")

		}

	} else {

		faviconHandler(w, r, ps)

	}

}

//getDatabaseData takes a linkID and looks it up in the database, if it exists it returns the link
func getDatabaseData(linkID string) string {
	//Postgres
	connStr := "user=" + os.Getenv("DB_USER") + " dbname=" + os.Getenv("DB_NAME") + " password=" + os.Getenv("DB_PASSWORD") + " host=" + os.Getenv("DB_HOST") + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var link string
	err = db.QueryRow("SELECT link FROM links WHERE linkID = $1 LIMIT 1;", linkID).Scan(&link)
	fmt.Println(err)
	fmt.Println(link)

	return link
}

func main() {
	//dotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//HTTPRouter
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/:linkID", matchLinkID)
	router.POST("/add_link", addLink)

	log.Println(http.ListenAndServe(":8080", router))

}
