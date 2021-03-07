package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
)

//------------STRUCTURES------------
//----
type Artist []struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Albumyears   int
}

type Location struct {
	Index []struct {
		Id            int      `json:"id"`
		Location      []string `json:"locations"`
		LocationClear []string
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
type APIFilter struct {
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
	Data           API
	FiltreLocation []string
	FiltreCreation []string
	FiltreAlbum    []string
}

type Artistmain struct {
	Data OneArtist
}
type ReceiveFilter struct {
	RecLocation string
	RecCreation string
	RecAlbum    string
	RecMembre   string
}

var filtrecreation []string
var filtrelocation []string
var filtreAlbum []string
var APITracker API

func GetFilter() {
	for i := range APITracker.Artist {
		add := ""
		count := 0
		for _, j := range APITracker.Artist[i].FirstAlbum {
			if count == 2 {
				add += string(j)
			}
			if j == '-' {
				count++
			}
		}
		intyears, _ := strconv.Atoi(add)
		APITracker.Artist[i].Albumyears = intyears
		valide := true
		for _, j := range filtreAlbum {
			if j == add {
				valide = false
			}
		}
		if valide == true {
			filtreAlbum = append(filtreAlbum, add)
		}
		strcreation := strconv.Itoa(APITracker.Artist[i].CreationDate)
		valide = true
		for _, j := range filtrecreation {
			if j == strcreation {
				valide = false
			}
		}
		if valide == true {
			filtrecreation = append(filtrecreation, strcreation)
		}

	}
	sort.Strings(filtreAlbum)
	sort.Strings(filtrecreation)
	for i := range APITracker.Location.Index {
		for j := range APITracker.Location.Index[i].LocationClear {
			valide := true
			add := APITracker.Location.Index[i].LocationClear[j]
			for h := range filtrelocation {
				if filtrelocation[h] == add {
					valide = false
				}
			}
			if valide == true {
				filtrelocation = append(filtrelocation, add)
			}
		}
	}
	sort.Strings(filtrelocation)
}

func GetArtistFilter(filtercountry string, filtercreation string, filternbrmember string, filteralbum string) API {
	TRLocation := false
	TRCreation := false
	TRMember := false
	TRAlbum := false

	creation, _ := strconv.Atoi(filtercreation)
	nbrmember, _ := strconv.Atoi(filternbrmember)
	album, _ := strconv.Atoi(filteralbum)
	var Apifiltre API
	for i := range APITracker.Artist {
		TRLocation = false
		TRCreation = false
		TRMember = false
		TRAlbum = false
		if filtercountry == "" {
			TRLocation = true
		} else {
			for _, j := range APITracker.Location.Index[i].LocationClear {
				if j == filtercountry {
					TRLocation = true
				}
			}
		}
		if filtercreation == "" {
			TRCreation = true
		} else if creation == APITracker.Artist[i].CreationDate {
			TRCreation = true
		}
		if filternbrmember == "0" {
			TRMember = true
		} else if nbrmember == len(APITracker.Artist[i].Members) {
			TRMember = true
		}
		if filteralbum == "" {
			TRAlbum = true
		} else if album == APITracker.Artist[i].Albumyears {
			TRAlbum = true
		}
		if TRLocation == true && TRCreation == true && TRMember == true && TRAlbum == true {
			Apifiltre.Artist = append(Apifiltre.Artist, APITracker.Artist[i])
			Apifiltre.Location.Index = append(Apifiltre.Location.Index, APITracker.Location.Index[i])
			Apifiltre.Dates.Index = append(Apifiltre.Dates.Index, APITracker.Dates.Index[i])
			Apifiltre.Relation.Index = append(Apifiltre.Relation.Index, APITracker.Relation.Index[i])
		}
	}
	return Apifiltre
}

//enleve
func Getbeautiful() {
	for i := range APITracker.Location.Index {
		for j := range APITracker.Location.Index[i].Location {
			valid := true
			tablocate := []rune(APITracker.Location.Index[i].Location[j])
			for l := range tablocate {
				if valid == true {
					if tablocate[l] == ':' {
						valid = false
					}
					if l == 0 {
						tablocate[l] -= 32
					}
					if tablocate[l] == '-' || tablocate[l] == '_' {
						tablocate[l] = ' '
						tablocate[l+1] -= 32
					}
				}
			}
			APITracker.Location.Index[i].LocationClear = append(APITracker.Location.Index[i].LocationClear, string(tablocate))
		}
	}
}

//on sort les informations de l'api
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

//page erreur lorsque l'on trouve pas la reponse avec la searchbar
func Erreur404(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("Erreur.html"))
	p := Page{Valeur: "essaie"}
	tmplate := tmpl.ExecuteTemplate(w, "Erreur.html", p)
	if tmplate != nil {
		http.Error(w, tmplate.Error(), http.StatusInternalServerError)
	}

}

//page contact avec info developper
func contact(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("contact.html"))
	p := Page{Valeur: "essaie"}
	tmplate := tmpl.ExecuteTemplate(w, "contact.html", p)
	if tmplate != nil {
		http.Error(w, tmplate.Error(), http.StatusInternalServerError)
	}

}

//
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
	// fmt.Println(OneArtist.Location)
	artist = append(artist, APITracker.Artist[y-1])
	OneArtist.Artist = artist
	relation = APITracker.Relation.Index[y-1].DatesLocation
	OneArtist.Relation = relation
	dates = APITracker.Dates.Index[y-1].Dates
	OneArtist.Dates = dates
	P := Artistmain{Data: OneArtist}
	//on met tout dans la variable one artist, on vient chercher les informations dans les strucures nommé,
	//on les stock dans one artiste avec l'index de l'artiste demander avec Url.Path
	//tout est stocké dans la structures api, on append le tout pour que toutes les informations apparaise dans la page Artist
	// et on renvoi le tout sur la page artiste avec le template execute qui contient toute les informations avec OneArtist
	tmpl.Execute(w, P)
}

//----PAGE D'ACCEUIL
//Apitracker est la variable qui contient la structure api qui contient elle meme tout les autres structures qui permet de grab toute les informations presente dans l'api
func mainpage(w http.ResponseWriter, r *http.Request) {
	FiltreApi := APITracker
	tmpl := template.Must(template.ParseFiles("index.html"))

	Filtre := ReceiveFilter{
		RecLocation: r.FormValue("Location"),
		RecCreation: r.FormValue("Creationdate"),
		RecAlbum:    r.FormValue("Albumdate"),
		RecMembre:   r.FormValue("Nombredemembre"),
	}
	fmt.Println("---------------------------")
	fmt.Println(Filtre.RecLocation)
	fmt.Println(Filtre.RecCreation)
	fmt.Println(Filtre.RecMembre)
	fmt.Println(Filtre.RecAlbum)
	fmt.Println("---------------------------")
	if Filtre.RecLocation != "" || Filtre.RecCreation != "" || (Filtre.RecMembre != "0" && Filtre.RecMembre != "") || Filtre.RecAlbum != "" {
		FiltreApi = GetArtistFilter(Filtre.RecLocation, Filtre.RecCreation, Filtre.RecMembre, Filtre.RecAlbum)
	}
	if len(FiltreApi.Artist) == 0 {
		tmpl1 := template.Must(template.ParseFiles("Erreur.html"))
		tmpl1.Execute(w, nil)
		fmt.Println("salut")
		return

	}
	P := Pagemain{Data: FiltreApi, FiltreAlbum: filtreAlbum, FiltreCreation: filtrecreation, FiltreLocation: filtrelocation}
	fmt.Println(FiltreApi)

	tmpl.Execute(w, P)

}

//on sort tout les information de l'url api et les stockons dans APItracker et nous utiliser la fonction unmashall pour les stocker dans la variable Apitracker
func jsonranger() {
	urlArtists := "https://groupietrackers.herokuapp.com/api/artists"
	urlLocation := "https://groupietrackers.herokuapp.com/api/locations"
	urlDates := "https://groupietrackers.herokuapp.com/api/dates"
	urlRelation := "https://groupietrackers.herokuapp.com/api/relation"

	UnmarshalJason(urlArtists, &APITracker.Artist)
	UnmarshalJason(urlLocation, &APITracker.Location)
	UnmarshalJason(urlDates, &APITracker.Dates)
	UnmarshalJason(urlRelation, &APITracker.Relation)

	// fmt.Println(APITracker)

}

//----MAIN

func main() {
	jsonranger()

	Getbeautiful()
	GetFilter()

	fmt.Println(filtreAlbum)
	fmt.Println(filtrecreation)

	fs := http.FileServer(http.Dir("paul"))
	http.Handle("/paul/", http.StripPrefix("/paul/", fs))
	http.HandleFunc("/artist/", artistpage)
	http.HandleFunc("/contact/", contact)
	http.HandleFunc("/", mainpage)
	http.HandleFunc("/error/", Erreur404)
	// http.HandleFunc("/search", Search)

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}
}
