package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type Testartist []struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type TestLocation struct {
	Index []struct {
		Id       int      `json:"id"`
		Location []string `json:"locations"`
	} `json:"index"`
}

type Testrelation struct {
	Index []struct {
		ID            int                 `json:"id"`
		DatesLocation map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

type TestDates struct {
	Index []struct {
		Id    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

type Apimain struct {
	ID       int
	Artist   Testartist
	Location TestLocation
	Dates    TestDates
	Relation Testrelation
}

var APITracker Apimain

func test1(url string, api interface{}) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	data, readErr := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if readErr != nil {
		log.Fatal(readErr)
	}
	json.Unmarshal(data, &api)

}

func test(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("./test.html"))

	tmpl.Execute(w, APITracker)
}
func jsonranger() {
	urlArtists := "https://groupietrackers.herokuapp.com/api/artists"
	urlLocation := "https://groupietrackers.herokuapp.com/api/locations"
	urlDates := "https://groupietrackers.herokuapp.com/api/dates"
	urlRelation := "https://groupietrackers.herokuapp.com/api/relation"

	test1(urlArtists, &APITracker.Artist)
	test1(urlLocation, &APITracker.Location)
	test1(urlDates, &APITracker.Dates)
	test1(urlRelation, &APITracker.Relation)

	fmt.Println(APITracker)

}

func main() {
	jsonranger()

	http.HandleFunc("/test", test)
	// http.HandleFunc("/search", Search)
	fs := http.FileServer(http.Dir("./server6"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	if err := http.ListenAndServe(":8070", nil); err != nil {
		log.Fatal(err)
	}
}
