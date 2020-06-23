package app

import (
	"encoding/json"
	"fmt"
	"log"

	persistence "../persistence"
	responseManager "../responsemanager"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

//Server : estrucutrua del servidor
type Server struct {
	port string
}

//Dominio : la esteructra para guadar el url de entrad
type Dominio struct {
	Domain string `json:"domain"`
}

var router *fasthttprouter.Router

var dominio string

//NewServer :Inicia un nuevo servidor
func NewServer() Server {
	return Server{}
}

// Init : Asigna el puerto de conexion
func (s *Server) Init(port string) {
	log.Println("Inicializando servidor...")
	s.port = ":" + port

}

//Start : comienza el serviodr en el puerto asignado
func (s *Server) Start() {
	//Iniciar el router
	router = fasthttprouter.New()

	router.POST("/buscar", Endpoint1)
	router.GET("/historial", Endpoint2)
	log.Println("Servidor corriendo y escuchando en el puerto: " + s.port + "\n")
	log.Fatal(fasthttp.ListenAndServe(":8000", router.Handler))

}

//Endpoint1 : Hace la busqueda de la informacion de un dominio y de sus servidores.
func Endpoint1(ctx *fasthttp.RequestCtx) {
	fmt.Println("Endpoint 1 HIT")
	var d = new(Dominio)
	err := json.Unmarshal(ctx.PostBody(), &d)
	if err != nil {
		fmt.Println(err)
	}
	var url = d.Domain
	responseManager.DefinirURL(url)
	responseManager.AllServers(ctx)
}

//Endpoint2 : Hace la busqueda del historial de busquedas en la base de datos.
func Endpoint2(ctx *fasthttp.RequestCtx) {
	fmt.Println("Endpoint 2 HIT")

	persistence.GetItems(ctx)
}
