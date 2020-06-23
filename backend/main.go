package main

import (
	"flag"

	"./src/system/app"
)

var port string

//var dominio string

func init() {
	flag.StringVar(&port, "port", "8000", "Asignando el puerto que el servidor deberia escuchar.")
	flag.Parse()
}

func main() {
	s := app.NewServer()
	s.Init(port) //El backend corre en el puerto 8000
	s.Start()
}
