package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Rsvp struct {
	Name, Email, Phone string
	WillAttend         bool
}

var responses = make([]*Rsvp, 0, 10)
var templates = make(map[string]*template.Template, 3)

func loadTemplates() {
	templateNames := [5]string{"welcome", "form", "thanks", "sorry", "list"}
	for index, name := range templateNames {
		t, err := template.ParseFiles("layout.html", name+".html")
		if err == nil {
			templates[name] = t
			fmt.Println("Loaded template", index, name)
		} else {
			panic(err)
		}
	}
}
func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	err := templates["welcome"].Execute(writer, nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
func listHandler(writer http.ResponseWriter, request *http.Request) {
	err := templates["list"].Execute(writer, responses)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

type formData struct {
	*Rsvp
	Errors []string
}

func formHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		templates["form"].Execute(writer, formData{
			Rsvp: &Rsvp{}, Errors: []string{},
		})
	} else if request.Method == http.MethodPost {
		request.ParseForm()
		responseData := Rsvp{
			Name:       request.Form["name"][0],
			Email:      request.Form["email"][0],
			Phone:      request.Form["phone"][0],
			WillAttend: request.Form["willattend"][0] == "true",
		}
		responses = append(responses, &responseData)
		if responseData.WillAttend {
			templates["thanks"].Execute(writer, responseData.Name)
		} else {
			templates["sorry"].Execute(writer, responseData.Name)
		}
	}
}

func main() {
	loadTemplates()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/form", formHandler)

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
