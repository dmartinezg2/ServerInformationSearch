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

	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("\t", "Buscador informacion de dominios", "\n")
	// fmt.Print("Ingrese el URL de la pagina que quiere investigar: ")
	// dominio, _ := reader.ReadString('\n')
	// dominio = strings.TrimSuffix(dominio, "\r\n")
	// dominio = strings.TrimSuffix(dominio, " ")
	// fmt.Println("Investigando la informacion del dominio: ", dominio)
	// s.AsignarDominio(dominio)

}
