package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

//----STRUCTURES

const URL = "https://groupietrackers.herokuapp.com/api"

type ArtistAll struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	Locations    []string
	Date         []string
}

type Artist struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	// Location     []string
}

type Profile struct {
	Artistpro      []Artist
	Artistlocation []Location
}

type Location struct {
	Id    map[string]int
	Dates map[string][]string
}

type Receive struct {
	Id string
}

type Artistsend struct {
	Block []Artist
}

type Dates struct {
	Id    map[string]int
	Dates map[string][]string
}

var artistall []ArtistAll
var artist []Artist
var artistlocations []Location
var artistedates []Dates
var Recherche string
var Idartist int

//----TRIER ARTISTES DANS L'ORDRE

func sort() {
	var tab1 []rune
	var tab2 []rune
	for j := range artist {
		tab1 = []rune(artist[j].Name)
		for i := range artist {
			tab2 = []rune(artist[i].Name)
			if int(tab1[0]) < int(tab2[0]) {
				artist[j], artist[i] = artist[i], artist[j]
			}
		}
	}

}

//----API LOCATION

func Getlocation() {
	res, err := http.Get(URL + "/relation")
	if err != nil {
		log.Fatal(err)
	}

	data, readErr := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if readErr != nil {
		log.Fatal(readErr)
	}
	Jsonerror := json.Unmarshal(data, &artistlocations)
	if Jsonerror != nil {
		log.Fatal(Jsonerror)
	}
}

//----API ARTISTE

func Getartist() {
	res, err := http.Get(URL + "/artists")
	if err != nil {
		log.Fatal(err)
	}

	data, readErr := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if readErr != nil {
		log.Fatal(readErr)
	}
	Jsonerror := json.Unmarshal(data, &artist)
	if Jsonerror != nil {
		log.Fatal(Jsonerror)
	}
}

//-- API DATE //
func GetDate() {
	res, err := http.Get(URL + "/dates")
	if err != nil {
		log.Fatal(err)
	}

	data, readErr := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if readErr != nil {
		log.Fatal(readErr)
	}
	Jsonerror := json.Unmarshal(data, &artistedates)
	if Jsonerror != nil {
		log.Fatal(Jsonerror)
	}
}

//-------Filter///

func Search() {
	tableid := []rune(Recherche)
	Recherche = ""
	for i := 0; i < len(tableid)-1; i++ {
		Recherche += string(tableid[i])
	}
	identier, _ := strconv.Atoi(Recherche)

	for i := range artist {
		if artist[i].Id == identier {
			Idartist = i
		}

	}
}

//----PAGE ARTISTE PRECIS

func artistpage(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("artist.html"))

	details := Receive{
		Id: r.FormValue("Id"),
	}
	Recherche = details.Id

	Search()

	tmpl.Execute(w, artist[Idartist])
}

//----PAGE D'ACCEUIL

func mainpage(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("index.html"))

	blockart := Artistsend{Block: artist}

	tmpl.Execute(w, blockart)
}

//----MAIN

func main() {

	Getartist()
	sort()
	//Getlocation()
	//GetDate()

	fs := http.FileServer(http.Dir("paul"))
	http.Handle("/paul/", http.StripPrefix("/paul/", fs))
	http.HandleFunc("/artist/", artistpage)
	http.HandleFunc("/", mainpage)

	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}
