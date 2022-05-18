package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/josh3021/nomadcoin/blockchain"
)

const templateDir string = "explorer/templates/"

var strPort string
var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(w http.ResponseWriter, r *http.Request) {
	// data := homeData{"Home", blockchain.BlockChain().GetAllBlocks()}
	// templates.ExecuteTemplate(w, "home", data)
}

func add(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		templates.ExecuteTemplate(w, "add", nil)
	case http.MethodPost:
		r.ParseForm()
		// data := r.Form.Get("data")
		blockchain.Blockchain().AddBlock()
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
}

// Start Explorer Server
func Start(port int) {
	handler := http.NewServeMux()
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	fmt.Printf("🚀 Explorer is Listening on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
