package responsemanager

//Se encarga de manejar la logica de la respuesta que se va a dar ne Json
import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	persistence "../persistence"
	webReader "../webreader"

	"github.com/likexian/whois-go"
	"github.com/valyala/fasthttp"
)

var urlGlobal string

//DefinirURL define la url globla a la que se le va a hacer la busqueda
func DefinirURL(url string) {
	urlGlobal = url
}

//Server Estructura de un servidor visitado
type Server struct {
	Address  string `json:"address"`
	SslGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}

//Response : Estructura de respuesta para el json de la parte 1
type Response struct {
	ServerList       []Server `json:"servers"`
	ServersChanged   bool     `json:"servers_changed"`
	SslGrade         string   `json:"ssl_grade"`
	PreviousSslGrade string   `json:"previous_ssl_grade"`
	Logo             string   `json:"logo"`
	Title            string   `json:"title"`
	IsDown           bool     `json:"is_down"`
}

//Devuelve los tres atributos de interes del resutado de la funcion whois de la API importada
func whoIsProccesor(whoisRaw string) []string {
	//fmt.Println(whoisRaw)
	OrgName := ""
	Country := ""
	Updated := ""
	lineas := strings.Split(whoisRaw, "\n")
	//fmt.Println(whoisRaw) // Se puede imprimir la respuesta del who is
	for i := 0; i < len(lineas)-1; i++ {
		if strings.Contains(lineas[i], "OrgName:") == true {
			contain := strings.Split(lineas[i], ":")
			OrgName = contain[1]
		} else if strings.Contains(lineas[i], "org-name:") == true {
			contain := strings.Split(lineas[i], ":")
			OrgName = contain[1]
		} else if strings.Contains(lineas[i], "Country:") == true {
			contain := strings.Split(lineas[i], ":")
			Country = contain[1]
		} else if strings.Contains(lineas[i], "country:") == true {
			contain := strings.Split(lineas[i], ":")
			Country = contain[1]
		} else if strings.Contains(lineas[i], "Updated:") == true {
			contain := strings.Split(lineas[i], ":")
			Updated = contain[1]
			Updated = strings.TrimSpace(Updated)
			Updated += "T12:30:00.371Z"
		} else if strings.Contains(lineas[i], "last-modified:") == true {
			contain := strings.Split(lineas[i], " ")
			Updated = contain[2]
			Updated = strings.TrimSpace(Updated)
		}
	}
	OrgName = strings.TrimSpace(OrgName)
	Country = strings.TrimSpace(Country)
	rta := []string{OrgName, Country, Updated}
	return rta
}

//Usa las funciones de time par verificar si la ultima fecha en la 	que se actualizo el servidor
// fue en menos de una hora.
func didServerUpdate(entry string) bool {
	//strforma := entry + "T12:30:00.371Z"
	timeOfUpdate, err := time.Parse(time.RFC3339, entry)
	if err != nil {
		fmt.Println(err)
	}
	timeNow := time.Now()
	timeHaceUnaHora := timeNow.Add(time.Hour * -1)
	rta := false
	if timeOfUpdate.After(timeHaceUnaHora) {
		rta = true
	}
	return rta
}

//Devuelve el valor ssl en una escala numerica que puede ser posteriormente comparada para
//propositos de guardar la nota mas baja de una serie de servidores.
func sslValue(param string) int {
	rta := 0
	switch {
	case strings.Contains(param, "A+"):
		rta = 7
	case strings.Contains(param, "A"):
		rta = 6
	case strings.Contains(param, "B"):
		rta = 5
	case strings.Contains(param, "C"):
		rta = 4
	case strings.Contains(param, "D"):
		rta = 3
	case strings.Contains(param, "E"):
		rta = 2
	case strings.Contains(param, "F"):
		rta = 1
	}
	return rta
}

//Endpoint1: AllServers : funcion que procesa la informacion del dominio ingresado
func AllServers(ctx *fasthttp.RequestCtx) {
	fmt.Println("Endpoint Hit: Servidores ")
	//llama al metodo de la clase lectorApi que devulve el json con la info del dominio
	s, err := webReader.GetServidores(urlGlobal)
	if err != nil {
		fmt.Println(err)
	}
	var servers []Server
	var dominio Response
	var srvrChanged bool
	var lowestSslGrade string
	var prvSslGrade string
	down := false
	//Recorre los endpoints del dominio
	endpoints := s.Endpoints
	sslGradeComparator := 8 //PAra tener el referente de la ssl grade a comparar.

	for i := 0; i <= len(endpoints)-1; i++ {
		adrss := endpoints[i].IPAddress
		whoisRaw, err := whois.Whois(adrss)
		if err != nil {
			down = true
			fmt.Println(err)
		}
		ssgrade := endpoints[i].Grade
		numSslgrade := sslValue(ssgrade)
		if (sslGradeComparator > numSslgrade) == true {
			lowestSslGrade = ssgrade
			sslGradeComparator = numSslgrade
		}
		whoIsProcessed := whoIsProccesor(whoisRaw)
		country := whoIsProcessed[1]
		ownr := whoIsProcessed[0]
		lastUpdate := whoIsProcessed[2]
		fmt.Println("Last update info", lastUpdate)
		srvr := Server{Address: adrss, SslGrade: ssgrade, Country: country, Owner: ownr}
		servers = append(servers, srvr)
		srvrChanged = didServerUpdate(lastUpdate)
		prvSslGrade = persistence.PreviousGrade(urlGlobal)

	}
	timeHaceUnaHora := time.Now().Add(time.Hour * -1) // le manda el tiempo hace una hora para que considere el booleano d euna ohora o antes atras.
	tiempoencadena := timeHaceUnaHora.Format(time.RFC3339)

	persistence.InsertItems(urlGlobal, lowestSslGrade, tiempoencadena)
	logo := webReader.FindLogo(urlGlobal)
	ttl := webReader.FindTTL(urlGlobal)
	webReader.ResetTitle()
	webReader.ResetLogo()
	dominio = Response{ServerList: servers, ServersChanged: srvrChanged, SslGrade: lowestSslGrade, PreviousSslGrade: prvSslGrade, Logo: logo, Title: ttl, IsDown: down}

	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(ctx).Encode(dominio)
}

//Endpoint2: Devuelve el historial de busquedas que se han hecho
func HistorialBusquedas(ctx *fasthttp.RequestCtx){
	fmt.Println("Endpoint Hit: Historial ")
	rta := persistence.GetItems()
	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(ctx).Encode(rta)
}
