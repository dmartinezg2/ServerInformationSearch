package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	//router = fasthttprouter.New()
	//router.GET("/", homePage)
	// router.POST("/newservers", Endpoint1)
	// router.GET("/servidores", Endpoint1)
	// router.POST("/dominios", Endpoint2)
	//log.Fatal(fasthttp.ListenAndServe(":8000", router.Handler))

	log.Println("Servidor arrancando en puerto: " + s.port + " iniciando router handler: ")
	http.HandleFunc("/newservers", Endpoint3)
	http.HandleFunc("/dominios", Endpoint4)

	http.ListenAndServe(":8000", nil)

}

//homePage de la pagina
func homePage(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Homepage Endpoint Hit")
}

//Endpoint1 : muestra la informacion de un dominio
func Endpoint1(ctx *fasthttp.RequestCtx) {
	fmt.Println("Esta llegando al endpoint 1")
	var d = new(Dominio)
	err := json.Unmarshal(ctx.PostBody(), &d)
	if err != nil {
		fmt.Println(err)
	}
	var url = d.Domain
	responseManager.DefinirURL(url)
	responseManager.AllServers(ctx)
}

//Endpoint2 : descriptio
func Endpoint2(ctx *fasthttp.RequestCtx) {
	fmt.Println("Esta llegando la peticion al endpoint numero 2")
	persistence.GetItems(ctx)
}

//Endpoint3 : muestra la informacion de un dominio
func Endpoint3(w http.ResponseWriter, request *http.Request) {
	var d = new(Dominio)
	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&d)
	fmt.Println("Esta llegando al endpoint 1")

	var domain = d.Domain
	fmt.Println(domain + "dominio buscado")
	responseManager.DefinirURL(domain)
	dominio := responseManager.AllServers2()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(w).Encode(dominio); err != nil {
		panic(err)
	}
}

//Endpoint4 : muestra la informacion de un dominio
func Endpoint4(w http.ResponseWriter, request *http.Request) {
	fmt.Println("Esta llegando la peticion al endpoint numero 2")
	dominios := persistence.GetItemsMod()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(w).Encode(dominios); err != nil {
		panic(err)
	}
}

//AsignarDominio :
// func (s *Server) AsignarDominio(pDominio string) {
// 	responseManager.DefinirURL(pDominio)
// 	log.Println("Dominio asignado")
// 	dominio = pDominio

// }

// //HandleRequest : maneja las request http
// func handleRequest() {
// 	router.GET("/", homePage)
// 	router.GET("/servidores", Endpoint1)
// 	router.GET("/dominios", Endpoint2)

// }
