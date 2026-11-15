package main

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
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
	{1, "Projet-Red-Hugo-Lior", "Tower of Duality.png", []string{"Go"}, "Développement d'un jeu vidéo de combat tour par tour en CLI.", "Jeu de combat tour par tourbasé sue la mythologie Grecque avec une tour à 2 sens : verrs le haut pour rejoindre la lumière des cieux, le bas pour rejoinre l'obscurité des ténèbres. Le jeu début par un choix de personnage : Sparte ou Athénien. Le jeu est composé d'une économie avec un marché où l'on retrouve un marchand, un forgeron. Également on peut y retrouver un système de fabrication d'équipements comme des armes et/ou des armures. Enfin le jeu est un jeu de combat en tour par tour, avec 10 niveaux possible : 5 vers les cieux et 5 vers les enfers.", []string{"combat tour par tour", "marché", "personnage"}, "https://github.com/li0r-ynov/Projet-Red-Hugo-Lior", "vidéo-tower-of-duality.mp4"},
	{2, "Harvey-TP-Web", "next lap - a new home.png", []string{"HTML", "CSS", "JAVASCRIPT"}, "Développement d'un jeu d'aventure en HTML/CSS, en se baladant de page en page à l'aide de boutons", "Jeu d'aventure en HTML/CSS basé sur le défilement de pages à l'aide de boutons. L'univers est la montagne et le snowboard, où l'on rencontre différents types de personnages : les locaux, les pros, les vidéastes, les parias. Le joueur avance de pages en pags jusqu'à arrivé à un contest de snowboard où il y participe avec une des équipes parmis les 4 groupes qu'il a rencontrer durant son aventure. CHaque fin est différencte en fonction de l'équipe que le jouer choisit.", []string{"Jeu d'aventure", "Pages Web", "Snowboard", "Équipes", "Contest"}, "https://github.com/li0r-ynov/Harvey-TP-Web", "vidéo-next-lap-a-new-home.mp4"},
}

func ListProjects(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request : %s %s", r.Method, r.URL.Path)

	renderTemplate(w, "liste-projets", allProjects)
}

func ShowProjectsById(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request : %s %s", r.Method, r.URL.Path)

	id, errId := strconv.Atoi(r.PathValue("id"))
	if errId != nil {
		http.Error(w, "Erreur données invalides...", http.StatusBadRequest)
		return
	}

	for _, p := range allProjects {
		if p.ID == id {
			renderTemplate(w, "projects-details", p)
			return
		}
	}
	http.Error(w, "Projet introuvable...", http.StatusNotFound)
}
