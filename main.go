package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	Schedule     map[string]string
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
	Block        []ArtistAll
	Num          int
	ResultNumber string
}

type ArtistSendSearch struct {
	BlockSearch  []ArtistAll
	ResultNumber string
}

type Dates struct {
	Id    int
	Dates []string
}

var artistall []ArtistAll

var artist []Artist
var artistlocations map[string][]Location
var artistedates map[string][]Dates

var Idartist int

//----TRIER ARTISTES DANS L'ORDRE

// func sort() {
// 	var tab1 []rune
// 	var tab2 []rune
// 	for j := range artist {
// 		tab1 = []rune(artist[j].Name)
// 		for i := range artist {
// 			tab2 = []rune(artist[i].Name)
// 			if int(tab1[0]) < int(tab2[0]) {
// 				artist[j], artist[i] = artist[i], artist[j]
// 			}
// 			if int(tab1[0]) == int(tab2[0]) {
// 				if int(tab1[1]) < int(tab2[1]) {
// 					artist[j], artist[i] = artist[i], artist[j]
// 				}
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

// func Search(w http.ResponseWriter, r *http.Request) {
// 	// fmt.Println(len(artistText))
// 	// fmt.Println(artistText)

// 	if recherche == "" {
// 		return
// 	}

// 	// blockArt := appendSearchResults(artistsearch, tmp)
// 	// blockArtSearch := Artistsend{Block: blockArt}
// 	// fmt.Printf("%+v\n", blockArtSearch)
// 	// fmt.Println(len(blockArtSearch.Block))

// 	// blockart := Artistsend{Block: artistall}

// 	// tmpl := template.Must(template.ParseFiles("search.html"))
// 	// tmpl.Execute(w, blockArtSearch)
// }

// ---- Create Schedule ----//
func CreateSchedule(y int) {
	schedule := make(map[string]string)
	locations := artistall[y-1].Locations
	dates := artistall[y-1].Date
	for i := range locations {
		schedule[locations[i]] = dates[i]
	}
	artistall[y-1].Schedule = schedule
}

//----PAGE ARTISTE PRECIS

func artistpage(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("artist.html"))

	y, _ := strconv.Atoi(r.URL.Path[8:])
	fmt.Println(y)
	// Search()
	CreateSchedule(y)
	tmpl.Execute(w, artistall[y-1])
}

//----PAGE D'ACCEUIL

func mainpage(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("index.html"))
	mMain := " Results Available"

	artistText := r.FormValue("searchNames")

	if artistText == "" {
		resultsMain := 0
		for i := range artistall {
			resultsMain = i
		}
		blockart := Artistsend{
			Block:        artistall,
			Num:          resultsMain + 1,
			ResultNumber: mMain,
		}
		tmpl.Execute(w, blockart)
	} else {
		recherche := artistText
		var artistSearch []ArtistAll
		var result []int
		for i := range artistall {
			// fmt.Println(artistall[i].Name)
			if strings.Contains(strings.ToLower(artistall[i].Name), strings.ToLower(recherche)) {
				result = append(result, artistall[i].Id)
				artistSearch = append(artistSearch, artistall[i])
			}
		}

		recherche = ""
		// nRes := string(len(result))
		mNoResultsFound := " Results Found: Try to search another artist!"
		// mResultsFound := nRes + " Results Available!"
		if len(result) != 0 {
			searchSend := Artistsend{
				Block:        artistSearch,
				Num:          len(result),
				ResultNumber: mMain,
			}
			tmpl.Execute(w, searchSend)
		} else {
			blockart := Artistsend{
				Block:        nil,
				Num:          len(result),
				ResultNumber: mNoResultsFound,
			}
			tmpl.Execute(w, blockart)
		}

		// fmt.Println(len(artistSearch[0].Name))

	}
	// recherche = artistText

}

// func appendSearchResults(ar []ArtistAll) []ArtistAll {
// 	var bl []ArtistAll
// 	lenall := make([]ArtistAll, len(ar))
// 	bl = lenall
// 	for i := range ar {
// 		bl[i].Id = ar[i].Id
// 		bl[i].Image = ar[i].Image
// 		bl[i].Name = ar[i].Name
// 		bl[i].Members = ar[i].Members
// 		bl[i].CreationDate = ar[i].CreationDate
// 		bl[i].FirstAlbum = ar[i].FirstAlbum
// 		bl[i].Locations = artistlocations["index"][i].Locations
// 		bl[i].Date = artistedates["index"][i].Dates

// 	}
// 	return bl
// }

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

//----MAIN

func main() {

	Getartist()
	Getlocation()
	GetDate()
	alldata()
	// sort()

	fs := http.FileServer(http.Dir("paul"))
	http.Handle("/paul/", http.StripPrefix("/paul/", fs))
	http.HandleFunc("/artist/", artistpage)
	http.HandleFunc("/", mainpage)
	// http.HandleFunc("/search", Search)

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}
}
