package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

//------------STRUCTURES------------

type Artist []struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Location struct {
	Index []struct {
		Id       int      `json:"id"`
		Location []string `json:"locations"`
	} `json:"index"`
}

type Relation struct {
	Index []struct {
		ID            int                 `json:"id"`
		DatesLocation map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

type Dates struct {
	Index []struct {
		Id    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

type API struct {
	ID       int
	Artist   Artist
	Location Location
	Dates    Dates
	Relation Relation
}

type OneArtist struct {
	Artist   Artist
	Location []string
	Relation map[string][]string
	Dates    []string
}

type Page struct {
	Valeur string
}

type Pagemain struct {
	Data API
}

type Artistmain struct {
	Data OneArtist
}

var sel string
var APITracker API

func UnmarshalJason(url string, api interface{}) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	data, readErr := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if readErr != nil {
		log.Fatal(readErr)
	}
	Jsonerror := json.Unmarshal(data, &api)
	if Jsonerror != nil {
		log.Fatal(Jsonerror)
	}

}
func Erreur404(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("Erreur.html"))
	p := Page{Valeur: "essaie"}
	tmplate := tmpl.ExecuteTemplate(w, "Erreur.html", p)
	if tmplate != nil {
		http.Error(w, tmplate.Error(), http.StatusInternalServerError)
	}

}

func contact(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("contact.html"))
	p := Page{Valeur: "essaie"}
	tmplate := tmpl.ExecuteTemplate(w, "contact.html", p)
	if tmplate != nil {
		http.Error(w, tmplate.Error(), http.StatusInternalServerError)
	}

}

func concertPage(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("concerts.html"))

	y, _ := strconv.Atoi(r.URL.Path[8:])
	fmt.Println(y)

	tmpl.Execute(w, nil)

}

func artistpage(w http.ResponseWriter, r *http.Request) {
	var dates []string
	var relation map[string][]string
	var location []string
	var artist Artist
	var OneArtist OneArtist
	tmpl := template.Must(template.ParseFiles("artist.html"))

	y, _ := strconv.Atoi(r.URL.Path[8:])
	location = APITracker.Location.Index[y-1].Location
	OneArtist.Location = location
	fmt.Println(OneArtist.Location)
	artist = append(artist, APITracker.Artist[y-1])
	OneArtist.Artist = artist
	relation = APITracker.Relation.Index[y-1].DatesLocation
	OneArtist.Relation = relation
	dates = APITracker.Dates.Index[y-1].Dates
	OneArtist.Dates = dates
	P := Artistmain{Data: OneArtist}

	fmt.Println(y)

	fmt.Println(APITracker.Artist[y-1])

	// Search()

	tmpl.Execute(w, P)
}

//----PAGE D'ACCEUIL

func mainpage(w http.ResponseWriter, r *http.Request) {
	P := Pagemain{Data: APITracker}
	tmpl := template.Must(template.ParseFiles("index.html"))

	tmpl.Execute(w, P)

}

func jsonranger() {
	urlArtists := "https://groupietrackers.herokuapp.com/api/artists"
	urlLocation := "https://groupietrackers.herokuapp.com/api/locations"
	urlDates := "https://groupietrackers.herokuapp.com/api/dates"
	urlRelation := "https://groupietrackers.herokuapp.com/api/relation"

	UnmarshalJason(urlArtists, &APITracker.Artist)
	UnmarshalJason(urlLocation, &APITracker.Location)
	UnmarshalJason(urlDates, &APITracker.Dates)
	UnmarshalJason(urlRelation, &APITracker.Relation)

	fmt.Println(APITracker)

}

func filterM(filterMembers string, Members []string) bool {
	IFM, _ := strconv.Atoi(filterMembers)
	if len(Members) == IFM {
		return true
	}
	return false
}

func filtreA(filtreAlbum, Artistalbum string) bool {
	PlageAlbum, _ := strconv.Atoi(filtreAlbum)
	Datestring := Artistalbum
	date, _ := strconv.Atoi(Datestring)

	if date == PlageAlbum {
		return true
	}
	return false
}

func FiltrCreation(creationDate int, filtreCreation string) bool {
	PlageCreation, _ := strconv.Atoi(filtreCreation)
	if creationDate == PlageCreation {
		return true
	}
	return false
}

//----MAIN

func main() {
	jsonranger()

	fs := http.FileServer(http.Dir("paul"))
	http.Handle("/paul/", http.StripPrefix("/paul/", fs))
	http.HandleFunc("/concerts/", concertPage)
	http.HandleFunc("/artist/", artistpage)
	http.HandleFunc("/contact/", contact)
	http.HandleFunc("/", mainpage)
	http.HandleFunc("/error/", Erreur404)
	// http.HandleFunc("/search", Search)

	if err := http.ListenAndServe(":8010", nil); err != nil {
		log.Fatal(err)
	}
}
