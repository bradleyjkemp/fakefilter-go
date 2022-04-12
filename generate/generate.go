package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"text/template"
)

type Fakelist struct {
	UpdatedTime int `json:"t"`
	Domains     map[string]struct {
		Provider        string `json:"provider"`
		Firstseen       int    `json:"firstseen"`
		Lastseen        int    `json:"lastseen"`
		RandomSubdomain bool   `json:"randomSubdomain"`
	} `json:"domains"`
}

func main() {
	resp, err := http.Get("https://raw.githubusercontent.com/7c/fakefilter/main/json/data.json")
	if err != nil {
		log.Fatal(err)
	}

	list := Fakelist{}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("fakefilter.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := tpl.Execute(f, list); err != nil {
		log.Fatal(err)
	}
}

var tpl = template.Must(template.New("").Parse(`package fakefilter

var domains = map[string]struct{}{
{{range $domain, $meta := .Domains}}    "{{$domain}}": {},
{{end}}}

func IsFakeDomain(domain string) bool {
	_, fake := domains[domain]
	return fake
}
`))
