package main

import (
	"github.com/jojomi/go-latex"
	
	"encoding/json"
	"log"
	"os"
	"text/template"
)

type Info struct {
	Name string `json:"name"`
	Severity string `json:"severity"`
}

type Match struct {
	Info Info `json:"info"`
	TemplateType string `json:"type"`
	Host string `json:"host"`
	Port string `json:"port"`
	URL string `json:"url"`
	MatchedAt string `json:"matched-at"`
	Request string `json:"request"`
	Response string `json:"response"`
	IP string `json:"ip"`
	Timestamp string `json:"timestamp"`
	CurlCommand string `json:"curl-command"`
}

func main() {

	// Parse input
	// TODO: stdin or from file
	//
	var match Match

	err := json.NewDecoder(os.Stdin).Decode(&match)
	if err != nil {
		log.Fatal(err)
	}

	// populate LaTex template
	// TODO: should I try to use text/template here?
	template, err := template.ParseFiles("./template.tex")
	if err != nil {
		log.Fatal(err)
	}

	// Output tex file? Do we need to?
	f, err := os.Create("/tmp/populated.tex")
	if err != nil {
		log.Fatal(err)
	}

	err = template.Execute(f, match)
	if err != nil {
		log.Fatal(err)
	}


	
	// Run pdflatex
	var ct latex.CompileTask = latex.NewCompileTask()
	ct.SetSourceDir("/tmp")
	ct.SetVerbosity(latex.VerbosityAll)
	ct.SetCompileFilename("/tmp/populated.tex")
	log.Println(ct.CompileFilenamePdf())
	err = ct.Pdflatex("/tmp/populated.tex", "")

	if err != nil {
		log.Fatal(err)
	}

	ct.MoveToDest(ct.CompileFilenamePdf(), ".")
	
	//
	// Clean up

}
