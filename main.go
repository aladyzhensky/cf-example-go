package main

import (
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"	

	"os"
	"log"
  "fmt"
)

type Book struct {
	Id 					int64 `db:"book_id"`
	Title       string
	Author      string
	Description string
}

func main() {
	dbmap := initDb()
  defer dbmap.Db.Close()

	m := martini.Classic()
	m.Map(dbmap)
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Get("/", ShowBooks)
	m.Post("/books", CreateBook)
	m.Get("/create", NewBooks)

  fmt.Println("listening...")
  err := http.ListenAndServe(":"+os.Getenv("PORT"), m)
  if err != nil {
    panic(err)
  }
  // log.Print("listening on " + "0.0.0.0:" + os.Getenv("PORT"))
	// m.RunOnAddr("0.0.0.0:" + os.Getenv("PORT"))
	// log.Fatal(http.ListenAndServe("0.0.0.0:" + os.Getenv("PORT"), m))
}

func NewBooks(r render.Render) {
	r.HTML(200, "create", nil)
}

func CreateBook(ren render.Render, r *http.Request, dbmap *gorp.DbMap) {
  new_book := Book{
  							Title: r.FormValue("title"), 
  							Author: r.FormValue("author"),
  						  Description: r.FormValue("description")}
  err := dbmap.Insert(&new_book)						

	PanicIf(err)
	ren.Redirect("/")
}

func ShowBooks(ren render.Render, r *http.Request, dbmap *gorp.DbMap) {
	var books_raws []Book
  _, err := dbmap.Select(&books_raws, "select * from books order by book_id")
  PanicIf(err)

	ren.HTML(200, "books", books_raws)
}

func initDb() *gorp.DbMap {
    // db, err := sql.Open("mysql", "root:@/go_sample")
	  db, err := sql.Open("mysql", os.Getenv("DB_URL"))
    PanicIf(err)

    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

    // Id property is an auto incrementing PK
    dbmap.AddTableWithName(Book{}, "books").SetKeys(true, "Id")

    err = dbmap.CreateTablesIfNotExists()
    PanicIf(err)
    dbmap.TraceOn("[gorp]", log.New(os.Stdout, "GO-sample:", log.Lmicroseconds)) 
    populateDb(dbmap)
    return dbmap
}

func populateDb(dbmap *gorp.DbMap) {	
    count, err := dbmap.SelectInt("select count(*) from books")
    PanicIf(err)
    if count == 0 {
    	book1 := Book{ 
    						Title: "JerBear goes to the City", 
  							Author: "Garnee Smashington",
  						  Description: "A young hipster bear seeks his fortune in the wild city of Irvine."}
  		book2 := Book{ 
    						Title: "Swarley''s Big Day", 
  							Author: "Barney Stinson",
  						  Description: "Putting his Playbook aside, one man seeks a lifetime of happiness."}				  
    	err = dbmap.Insert(&book1, &book2)
    	PanicIf(err)
    }

}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
