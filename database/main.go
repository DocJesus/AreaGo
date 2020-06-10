package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/handlers"
	_ "github.com/lib/pq"
	"github.com/mux"
)

// User est en struct ce que la table SQL est
type User struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Passwd string `json:"passwd"`
}

var db *sql.DB
var err error

func init() {
	db, err = sql.Open("postgres", "postgres://postgres:toto@localhost/userarea")
	if err != nil {
		fmt.Printf("swot")
		log.Fatal(err)
	}
}

// AllUser retourne tout les user de la talbe utilsiateur
var AllUser = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM utilisateur")
	if err != nil {
		fmt.Printf("swot2")
		log.Fatal(err)
	}
	defer rows.Close()

	usrs := make([]*User, 0)
	for rows.Next() {
		usr := new(User)
		err := rows.Scan(&usr.Id, &usr.Name, &usr.Passwd)
		if err != nil {
			log.Fatal(err)
		}
		usrs = append(usrs, usr)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, usr := range usrs {
		fmt.Fprintf(w, "%s %s %s\n", usr.Id, usr.Name, usr.Passwd)
	}
})

// UsrUser retourne les infos de l'user à partir d'un certain Username
var UsrUser = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	name := r.Form.Get("user")

	rows, err := db.Query("SELECT * FROM utilisateur where username='" + name + "'")
	if err != nil {
		fmt.Printf("swot2")
		log.Fatal(err)
	}
	defer rows.Close()

	usrs := make([]*User, 0)
	for rows.Next() {
		usr := new(User)
		err := rows.Scan(&usr.Id, &usr.Name, &usr.Passwd)
		if err != nil {
			log.Fatal(err)
		}
		usrs = append(usrs, usr)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, usr := range usrs {
		fmt.Fprintf(w, "%s %s %s | %s\n", usr.Id, usr.Name, usr.Passwd, name)
	}
})

// newUser va créer un user dans la base de donnée
var newUser = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	//r.ParseForm()
	//log.Fatal(r.Body)

	var usr User
	err := json.NewDecoder(r.Body).Decode(&usr)

	if err != nil {
		log.Fatal(err)
	}
	//log.Fatal("id = " + usr.Name + " " + usr.Passwd + " " + usr.Id)

	//name := r.Form.Get("newUser")
	//passwd := r.Form.Get("newPasswd")
	var id string

	/*type User struct {
		id     string
		name   string
		passwd string
	}*/

	//récupère tout les users pour pouvoir avoir l'ID du nouvel user (dernier ID + 1)
	rows, _ := db.Query("SELECT count(*) FROM utilisateur")
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&usr.Id)
	}

	tmpID, _ := strconv.Atoi(usr.Id)
	tmpID++
	usr.Id = strconv.Itoa(tmpID)

	_, secErr := db.Query("insert into utilisateur (id, username, passwd) values ('" + usr.Id + "', '" + usr.Name + "', '" + usr.Passwd + "')")

	if secErr != nil {
		log.Fatal(secErr)
	}

	fmt.Fprintf(w, "nouvel utilsiateur créé %s %s %s", id, usr.Name, usr.Passwd)
})

func main() {

	r := mux.NewRouter()

	r.Handle("/users", AllUser).Methods("GET")
	r.Handle("/users/user", UsrUser).Methods("GET")
	r.Handle("/users/newUser", newUser).Methods("POST")

	http.ListenAndServe(":4242", handlers.LoggingHandler(os.Stdout, r))

}
