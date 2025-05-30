package main

import (
	"github.com/jojomi/go-latex"

	"encoding/json"
	"log"
	"os"
	"io"
	"text/template"
	"path/filepath"
)

type Info struct {
	Name     string `json:"name"`
	Severity string `json:"severity"`
}

type Match struct {
	Info         Info   `json:"info"`
	TemplateType string `json:"type"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	URL          string `json:"url"`
	MatchedAt    string `json:"matched-at"`
	Request      string `json:"request"`
	Response     string `json:"response"`
	IP           string `json:"ip"`
	Timestamp    string `json:"timestamp"`
	CurlCommand  string `json:"curl-command"`
}

func main() {

	// TODO: configify this
	//WorkingDir := "/tmp/"
	TemplateName := "template.tex"
	OutputName := "report"

	// Parse input
	// TODO: order this by severity
	var matches []Match


	// TODO: top-level stats?

	// TODO: stdin or from file
	decoder := json.NewDecoder(os.Stdin)

	for {
		var match Match
		if err := decoder.Decode(&match); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		matches = append(matches, match)
	}

	
	var ct latex.CompileTask = latex.NewCompileTask()
	ct.SetSourceDir(".") // Is this needed?
	ct.SetVerbosity(latex.VerbosityAll)

	// populate LaTex template
	tmpl, err := template.ParseFiles(TemplateName)
	if err != nil {
		log.Fatal(err)
	}

	// Output tex file? Do we need to?
	log.Println("Executing template", TemplateName)

	// create temp file
	tempFile, err := os.CreateTemp("", "*.tex")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(tempFile, matches)
	if err != nil {
		log.Fatal(err)
	}

	// Run pdflatex
	ct.SetCompileFilename(tempFile.Name())
	log.Println("Generating PDF file:", ct.CompileFilenamePdf())
	err = ct.Pdflatex(tempFile.Name(), "")

	if err != nil {
		log.Fatal(err)
	}

//	ct.MoveToDest(ct.CompileFilenamePdf(), OutputName + ".pdf")


	// Clean up
	ct.ClearLatexTempFiles(".")
	finalName := filepath.Base(ct.CompileFilenamePdf())
	log.Println("Moving compiled file from", finalName, "to", OutputName + ".pdf")

	err = os.Rename(finalName, OutputName + ".pdf")
	if err != nil {
		log.Fatal(err)
	}

}
