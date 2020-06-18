package webreader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//Endpoint : Modela la lista de Enpoints que se dan en la respuesta del API de SSL
type Endpoint struct {
	IPAddress         string `json:"ipAddress"`
	ServerName        string `json:"serverName"`
	StatusMessage     string `json:"statusMessage"`
	Grade             string `json:"grade"`
	GradeTrustIgnored string `json:"gradeTrustIgnored"`
	HasWarnings       bool   `json:"hasWarnings"`
	IsExceptional     bool   `json:"isExceptional"`
	Progress          int64  `json:"progress"`
	Duration          int64  `json:"duration"`
	Delegation        int64  `json:"delegation"`
}

//SslAPIResponse : Modela la estructura de respuesta del API de SSL
type SslAPIResponse struct {
	Host            string     `json:"host"`
	Port            int64      `json:"port"`
	Protocol        string     `json:"protocol"`
	IsPublic        bool       `json:"isPublic"`
	Status          string     `json:"status"`
	StartTime       int64      `json:"startTime"`
	TestTime        int64      `json:"testTime"`
	EngineVersion   string     `json:"engineVersion"`
	CriteriaVersion string     `json:"criteriaVersion"`
	Endpoints       []Endpoint `json:"endpoints"`
}

//API : bla bla
type API struct {
}

//HacerRequest Pide el Json asociado
func HacerRequest(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return []byte(body), err
}

//GetServidores : obtiene todos los servidores
func GetServidores(url string) (*SslAPIResponse, error) {

	body, err := HacerRequest("https://api.ssllabs.com/api/v3/analyze?host=" + url)
	if err != nil {
		return nil, err
	}
	s, err := ParseServidores(body)
	return s, err
}

//ParseServidores : hace el parse del json a una estructura local al programa
func ParseServidores(body []byte) (*SslAPIResponse, error) {
	var s = new(SslAPIResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return s, err
}

var rtaTitulo = ""

//FindTTL : Encuentra el valor del titulo y del logo de una pagina.
// Basado en los algoritmos de https://www.devdungeon.com/content/web-scraping-go
// Utiliza una libreria que permite hacer jqueries en go. El link del repositorio esta en imports
func FindTTL(url string) string {
	urlCorrecta := "https://www."
	urlCorrecta += url
	res, err := http.Get(urlCorrecta)
	if err != nil {
		log.Fatal(err)
	}
	// para cerrar la peticion htttp cuando se acabe la ejecucion del metodo.
	defer res.Body.Close()
	jDoc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}
	jDoc.Find("title").Each(procesarTitulo)
	return rtaTitulo
}

//sacaa el valor del html del titulo
func procesarTitulo(index int, element *goquery.Selection) {
	hrefTitulo, error := element.Html()
	if error != nil {
		fmt.Println(error)
	}
	rtaTitulo = hrefTitulo
}

// ResetTitle : hace el resetdel titulo
func ResetTitle() {
	rtaTitulo = "No hay titulo"
}

var rtaLogo = ""

// FindLogo ecuentra el logo de la pagina analizando los links que se encuentran en el HTML
// Basado en los algoritmos de https://www.devdungeon.com/content/web-scraping-go
// Utiliza una libreria que permite hacer jqueries en go. El link del repositorio esta en imports
func FindLogo(url string) string {
	urlCorrecta := "https://www."
	urlCorrecta += url
	res, err := http.Get(urlCorrecta)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	jDoc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}
	jDoc.Find("link").Each(procesarLink)
	return rtaLogo
}

//Proccesa las lineas de html que contengan link como tag y evalua cual de ellas es el icono
//Considerando si contiene icon en su link
func procesarLink(index int, element *goquery.Selection) {
	// See if the href attribute exists on the element
	hrefLogo, exists := element.Attr("href")
	if exists {
		if strings.Contains(hrefLogo, "icon") {
			rtaLogo = hrefLogo
		}
	}

}

func ResetLogo() {
	rtaLogo = "No hay logo"
}
