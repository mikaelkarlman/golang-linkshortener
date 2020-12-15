package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

//Index server generic page
func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func faviconHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	http.ServeFile(w, r, "favicon.ico")
}

//MatchLinkID checks the ID, matches it and redirects if the key exists
func MatchLinkID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var linkID string = ps.ByName("linkID")

	if linkID != "favicon.ico" {

		fmt.Println("LinkID called: " + linkID)

		link := GetDatabaseData(linkID)

		if len(link) != 0 {

			http.Redirect(w, r, link, http.StatusMovedPermanently)

		} else {

			fmt.Fprintf(w, "LinkID: "+string(linkID)+" does not exist")

		}

	} else {

		faviconHandler(w, r, ps)

	}

}

//GetDatabaseData takes a linkID and looks it up in the database, if it exists it returns the link
func GetDatabaseData(linkID string) string {
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
	router.GET("/:linkID", MatchLinkID)

	log.Println(http.ListenAndServe(":8080", router))

}
