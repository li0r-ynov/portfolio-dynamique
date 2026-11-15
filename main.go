package main

import (
	"bytes"
	"log"
	"net/http"
	"text/template"
)

// Chargement de tous les templates HTML du dossier "templates" au démarrage.
// template.Must fait planter le programme si un template est invalide.
var _temp = template.Must(template.ParseGlob("./templates/*.html"))

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", Home)
	router.HandleFunc("GET /projects", ListProjects)
	router.HandleFunc("GET /projects/{id}", ShowProjectsById)
	router.HandleFunc("GET /", NotFound)

	fileServer := http.FileServer(http.Dir("./asset"))
	router.Handle("/asset/", http.StripPrefix("/asset/", fileServer))

	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Fatal(err)
	}
}

// renderTemplate exécute un template dans un buffer puis l'envoie au client.
// Passer par un buffer évite d'envoyer une page à moitié écrite en cas d'erreur.
func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	// variable temp d'écriture du template
	buffer := bytes.Buffer{}

	errRender := _temp.ExecuteTemplate(&buffer, name, data)
	if errRender != nil {
		log.Printf("Erreur (Template) : Impossible de rendre le template %s : %v", name, errRender)
		http.Error(w, "Erreur chargement de la page", http.StatusInternalServerError)
		return
	}

	// Les en-têtes doivent être définis avant WriteHeader, sinon ils sont ignorés
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	buffer.WriteTo(w)
}

func Home(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request : %s %s", r.Method, r.URL.Path)

	renderTemplate(w, "index", nil)
}

type Projects struct {
	ID          int
	Title       string
	CoverImage  string
	TechStack   []string
	ShortIntro  string
	Description string
	Features    []string
	RepoLink    string
	DemoVideo   string
}

var allProjects = []Projects{
	{1,
	}

}

func ListProjects(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request : %s %s", r.Method, r.URL.Path)

	renderTemplate(w, "liste-projets", allProjects)
}
