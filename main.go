package main

import (
	"github.com/jojomi/go-latex"

	"encoding/json"
	"log"
	"os"
	"io"
	"text/template"
	"path/filepath"
	"sort"
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


// Utils for sorting resultsby severity
func findIndex(slice []string, value string) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return len(slice) // Return length if not found (treat as lowest priority)
}



func preprocess(matches []Match) []Match {

	// First sort by criticality
	severityOrder := []string{"critical", "high", "medium", "low", "info"}
	sort.Slice(matches, func(i, j int) bool {
		indexI := findIndex(severityOrder, matches[i].Info.Severity)
		indexJ := findIndex(severityOrder, matches[j].Info.Severity)
		return indexI < indexJ
	})

	
	var processed []Match



	processed = matches

	return processed
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


	// TODO: we need to preprocess, esepecially the HTTP requests
	// because LaTeX/templates really struggle by themselves
	processed := preprocess(matches)

	
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

	err = tmpl.Execute(tempFile, processed)
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


	// Clean up
	ct.ClearLatexTempFiles(".")
	finalName := filepath.Base(ct.CompileFilenamePdf())
	log.Println("Moving compiled file from", finalName, "to", OutputName + ".pdf")

	err = os.Rename(finalName, OutputName + ".pdf")
	if err != nil {
		log.Fatal(err)
	}

}
