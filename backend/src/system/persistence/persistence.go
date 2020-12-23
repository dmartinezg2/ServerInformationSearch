package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/valyala/fasthttp"
)

//informacion para desplegar el historial de busquedas
type Historial struct {
	Dominio  string `json:"url"`
	DateVisited string `json:"date_visited"`
}
//Rta la estrucutra de respuesta del endpoint 2 de la pagina.
type Rta struct {
	Items []Historial `json:"items"`
}



//GetItems : busca todos los dominios buscados en la base de datos.
func GetItems(ctx *fasthttp.RequestCtx) {
	// Connect to the "bank" database.
	db, err := sql.Open("postgres", "postgresql://daviddb@localhost:26257/retobd?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	rows, err := db.Query("SELECT dominio, dateVisited FROM busquedas")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("Initial urls:")
	var items []Historial
	for rows.Next() {
		var dominio, dateVisited string 
		if err := rows.Scan(&dominio,&dateVisited); err != nil {
			log.Fatal(err)
		}
		items = append(items, Historial{dominio, dateVisited})
		fmt.Printf("%s\n", dominio)
	}
	rta := Rta{items}
	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(ctx).Encode(rta)
}

// InsertItems inserta una busqueda en la tabla busquedas de la base de datos.
func InsertItems(url string, sslgrade string, lasVisit string) {
	// Connect to the "bank" database.
	db, err := sql.Open("postgres", "postgresql://daviddb@localhost:26257/retobd?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	if _, err := db.Exec(
		"INSERT INTO busquedas (dominio, grade, dateVisited) VALUES ('" + url + "', '" + sslgrade + "', '" + lasVisit + "')"); err != nil {
		log.Fatal(err)
	}
}

// PreviousGrade : busca la ssl grade previa
func PreviousGrade(url string) string {
	var respuesta = "Inexistente"
	db, err := sql.Open("postgres", "postgresql://daviddb@localhost:26257/retobd?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	rows, err := db.Query("SELECT  grade, dateVisited FROM busquedas WHERE dominio = '" + url + "'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {

		var grade, dateVisited string
		if err := rows.Scan(&grade, &dateVisited); err != nil {
			log.Fatal(err)
		}
		fechaVisita, err := time.Parse(time.RFC3339, dateVisited)
		if err != nil {
			fmt.Println("No se pudo parsear el tiempo de visita", err)
		}
		if fechaVisita.Before(time.Now().Add(time.Hour * -1)) {
			respuesta = grade
		}

	}
	return respuesta
}
