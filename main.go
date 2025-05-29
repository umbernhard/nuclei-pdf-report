package main

import (
	"github.com/jojomi/go-latex"
	
	"encoding/json"
	"text/template"
	"log"
	"os"
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


	// TODO: configify this
	//WorkingDir := "/tmp/"
	TemplateName := "template.tex"
	OutputName := "report"


	// Parse input
	// TODO: stdin or from file
	var match Match

	err := json.NewDecoder(os.Stdin).Decode(&match)
	if err != nil {
		log.Fatal(err)
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

	err = tmpl.Execute(tempFile, match)
	if err != nil {
		log.Fatal(err)
	}
	/*
	err = ct.ExecuteTemplate(tmpl, match,TemplateName, "")
	if err != nil {
		log.Fatal(err)
	}
	*/

	// Run pdflatex
	ct.SetCompileFilename(tempFile.Name())
	log.Println("Generating PDF file:", OutputName +".pdf")
	err = ct.Pdflatex(tempFile.Name(), "-jobname=" + OutputName)

	if err != nil {
		log.Fatal(err)
	}

	//ct.MoveToDest(, ".")
	
	// Clean up
	ct.ClearLatexTempFiles(".")

}
