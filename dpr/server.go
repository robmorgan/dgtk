package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strings"
)

type Server struct {
	Address            string
	AwsAccessKeyId     string
	AwsSecretAccessKey string
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.Address, s)
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method + " http://" + r.Host + r.URL.String())
	defer r.Body.Close()

	normalizedPath := root + r.URL.Path
	if strings.HasSuffix(normalizedPath, "/") {
		normalizedPath += "index"
	}

	s := &Resource{r}
	w.Header().Set("X-Docker-Registry-Version", "0.0.1")
	w.Header().Add("X-Docker-Endpoints", r.Host)
	switch r.Method {
	case "PUT":
		e := s.store()
		if e != nil {
			log.Println(e.Error())
			http.Error(w, e.Error(), 500)
			return
		}
		w.Header().Add("WWW-Authenticate", `Token signature=123abc,repository="dynport/test",access=write`)
		w.Header().Add("X-Docker-Token", "token")
		w.Header().Add("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.String(), "/images") {
			w.WriteHeader(204)
		} else {
			w.WriteHeader(200)
		}
		return
	case "GET":
		if strings.HasSuffix(r.URL.String(), "/tags") {
			writeTags(s.localPath(), w)
		} else if strings.HasSuffix(r.URL.String(), "/ancestry") {
			p := root + path.Dir(r.URL.String())
			list := []string{path.Base(p)}
			for {
				img, e := loadImage(p + "/json")
				if e != nil {
					log.Print(e.Error())
					break
				}
				if img.Parent == "" {
					break
				}
				list = append(list, img.Parent)
				log.Println(list)
				p = path.Dir(p) + "/" + img.Parent
			}
			w.Header().Set("Content-Type", "application/json")
			if e := json.NewEncoder(w).Encode(list); e != nil {
				log.Print("ERROR: " + e.Error())
			}
		} else if s.Exists() {
			_, e := s.Write(w)
			if e != nil {
				log.Print("ERROR: " + e.Error())
			}
			return
		}
	}
	w.WriteHeader(404)
}