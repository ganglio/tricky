package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/ganglio/memo"
)

type data struct {
	Classes []string            `json:"classes"`
	Words   map[string][]string `json:"words"`
}

var (
	words  data
	page   func() interface{}
	kitten []byte
)

func init() {
	js, err := os.Open("words.json")
	if err != nil {
		panic(err)
	}
	defer js.Close()
	body, err := ioutil.ReadAll(js)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &words)
	if err != nil {
		panic(err)
	}

	kt, err := os.Open("kitten.png")
	if err != nil {
		panic(err)
	}
	defer kt.Close()
	kitten, err = ioutil.ReadAll(kt)
	if err != nil {
		panic(err)
	}

	page = memo.MemoX(func() (interface{}, error) {
		return template.ParseFiles("page.gohtml")
	}, time.Second*5)
}

func main() {

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write(kitten)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rand.Seed(time.Now().UnixNano())

		class := words.Classes[rand.Uint32()%uint32(len(words.Classes))]
		word := words.Words[class][rand.Uint32()%uint32(len(words.Words[class]))]

		log.Printf("Rendering %s %s", class, word)

		data := struct {
			Class string
			Word  string
		}{
			Class: class,
			Word:  word,
		}

		w.Header().Set("Content-Type", "text/html")
		page().(*template.Template).Execute(w, data)
	})

	port := os.Getenv("PORT")
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
