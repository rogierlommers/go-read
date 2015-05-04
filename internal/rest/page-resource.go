package rest

import (
	"encoding/base64"
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"github.com/rogierlommers/go-read/internal/model"
)

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	glog.Info("stats page")
}

func GenerateRSS(database model.ReadingListRecords) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//	sort.Sort(sort.Reverse(ById(database.Records)))

		now := time.Now()
		feed := &feeds.Feed{
			Title:       "Go-read",
			Link:        &feeds.Link{Href: "http://bla.com"},
			Description: "personal RSS feed with articles to be read",
			Author:      &feeds.Author{"Rogier Lommers", "rogier@lommers.org"},
			Created:     now,
		}

		for _, value := range database.Records {
			glog.Info("adding feed: ", value.URL)
			newItem := feeds.Item{Title: value.URL,
				Link: &feeds.Link{Href: value.URL},
				//Description: "A discussion on controlled parallelism in golang",
				//Author:      &feeds.Author{"Rogier Lommers", "jmoiron@jmoiron.net"},
				//Created:     now,
			}
			feed.Add(&newItem)
		}

		atom, err := feed.ToAtom()
		if err != nil {
			glog.Errorf("error creating RSS feed -> %s", err)
			return
		}
		w.Write([]byte(atom))
	}
}

func AddArticle(database model.ReadingListRecords) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		base64url := vars["base64url"]

		urlByteArray, err := base64.StdEncoding.DecodeString(base64url)
		if err != nil {
			glog.Errorf("error decoding url -> %s", err)
			return
		}

		url := string(urlByteArray[:])

		database = model.AddRecord(database, url)
		glog.Infof("add #%d url [%s] --> [%s]: ", len(database.Records), base64url, url)
		w.Write([]byte("url added..."))
	}
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	glog.Info("index page")
	glog.Info(r)
	fp := path.Join("static", "templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO extract serverlocation from header
	obj := map[string]string{"serverLocation": "http://localhost:8080"}

	if err := tmpl.Execute(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
