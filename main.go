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
}

type Location struct {
	Id        int
	Locations []string
}

type Receive struct {
	Id string
}

type Artistsend struct {
	Block []ArtistAll
}

type Dates struct {
	Id    int
	Dates []string
}

var artistall []ArtistAll
var artist []Artist
var artistlocations map[string][]Location
var artistedates map[string][]Dates
var Recherche string
var Idartist int
var localisation []int

//----TRIER ARTISTES DANS L'ORDRE

// func sort() {
// 	var tab1 []rune
// 	var tab2 []rune
// 	for j := range artistall {
// 		tab1 = []rune(artistall[j].Name)
// 		for i := range artistall {
// 			tab2 = []rune(artistall[i].Name)
// 			if int(tab1[0]) < int(tab2[0]) {
// 				artistall[j], artistall[i] = artistall[i], artistall[j]
// 			}
// 		}
// 	}
// }

//----API LOCATION

func Getlocation() {
	res, err := http.Get(URL + "/locations")
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
	res, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
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
		if artistall[i].Id == identier {
			Idartist = i
		}

	}
}

//----PAGE ARTISTE PRECIS

func artistpage(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("artist.html"))

	y, _ := strconv.Atoi(r.URL.Path[8:])
	fmt.Println(y)
	Search()

	tmpl.Execute(w, artistall[y-1])
}

//----PAGE D'ACCEUIL

func mainpage(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("index.html"))

	blockart := Artistsend{Block: artistall}

	tmpl.Execute(w, blockart)
}

func alldata() {
	lenall := make([]ArtistAll, len(artist))
	artistall = lenall
	for i := range artist {
		artistall[i].Id = artist[i].Id
		artistall[i].Image = artist[i].Image
		artistall[i].Name = artist[i].Name
		artistall[i].Members = artist[i].Members
		artistall[i].CreationDate = artist[i].CreationDate
		artistall[i].FirstAlbum = artist[i].FirstAlbum
		artistall[i].Locations = artistlocations["index"][i].Locations
		artistall[i].Date = artistedates["index"][i].Dates

	}
}

func pageconcert(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("concert.html"))

	y, _ := strconv.Atoi(r.URL.Path[8:])
	tmpl.Execute(w, localisation[y-1])
}

//----MAIN

func main() {

	Getartist()
	Getlocation()
	GetDate()
	alldata()
	//sort()

	fs := http.FileServer(http.Dir("paul"))
	http.Handle("/paul/", http.StripPrefix("/paul/", fs))
	http.HandleFunc("/artist/", artistpage)
	http.HandleFunc("/", mainpage)
	http.HandleFunc("/concert", pageconcert)

	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}
