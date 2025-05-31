package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/jojomi/go-latex"
	"github.com/op/go-logging"

	"encoding/json"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

// Max characters before seqsplit
const MAX = 110

// Setup log for output
var log = logging.MustGetLogger("")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}% %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

// Structs for the json we're parsing.
// TODO: maybe move this to a separate file?
type Info struct {
	Name        string `json:"name"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
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

type Summary struct {
	NumCritical int
	NumHigh int
	NumMedium int
	NumLow int
	NumInfo int
	Total int
}

type Report struct {
	Summary Summary
	Matches []Match
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

// Process a response to make it nicer in the PDF
func processHttp(httpobject string) string {

	lines := strings.Split(httpobject, "\r\n")

	var output []string

	for _, line := range lines {
		// Remove characters LaTeX doesn't like
		outline := sanitize(line)
		// If the line has a Cookie or a link in it, wrap the offending content
		// in seqsplit so it doesn't overflow.
		needsSeqSplit := false

		// These are characters in things like URLs that might need seqsplit
		re := regexp.MustCompile("\\?|&")
		needsSeqSplit = (re.MatchString(outline) && len(outline) > MAX) || strings.Contains(outline, "Cookie")
		if needsSeqSplit {
			log.Debug("Found cookie")
			outline = "\\seqsplit{" + outline + "}"
		} else if strings.Contains(outline, "href") {
			re := regexp.MustCompile(`href='(.*?)'`)
			outline = re.ReplaceAllString(outline, "\\seqsplit{href='$1'}")
		}

		output = append(output, outline)
	}

	return strings.Join(output, "\r\n")
}

// Fix LaTeX unfriendly characters
func sanitize(input string) string {

	replace := strings.NewReplacer("{", "\\{", "}", "\\}", "_", "\\_", "&", "\\&")
	return replace.Replace(input)

}

func preprocess(matches []Match) []Match {

	// First sort by criticality
	severityOrder := []string{"critical", "high", "medium", "low", "info"}
	sort.Slice(matches, func(i, j int) bool {
		indexI := findIndex(severityOrder, matches[i].Info.Severity)
		indexJ := findIndex(severityOrder, matches[j].Info.Severity)
		return indexI < indexJ
	})

	for i, match := range matches {

		match.Info.Name = sanitize(match.Info.Name)
		match.Info.Description = sanitize(match.Info.Description)

		// Rewrite the severity for pretier printing
		match.Info.Severity = cases.Title(language.English, cases.NoLower).String(match.Info.Severity)

		// Process requests(?)
		match.Request = processHttp(match.Request)

		// Process response
		match.Response = processHttp(match.Response)

		matches[i] = match
	}

	return matches
}

var opts struct {
	LatexVerbosity string `long:"latex-verbosity" description:"Set the verbosity for LaTeX output" default:"default" choice:"none" choice:"default" choice:"more" choice:"all"`
	Verbosity      bool   `short:"v" long:"verbose" description:"Show verbose debugging information"`
	Output         string `short:"o" long:"output" description:"Name of the PDF to output" default:"report"`
	SaveTexFile    bool   `short:"s" long:"save-tex" description:"Save a copy of the .tex file (useful for debugging)"`
	SaveAllTexFile bool   `long:"save-all-tex" description:"Save all files produced by LaTeX (.aux, .log, etc.)"`
	Template       string `short:"t" long:"template" description:"The name of a template to use" default:"template.tex"`
}

func main() {

	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)

	_, err := flags.Parse(&opts)
	if flags.WroteHelp(err) {
		return
	} else if err != nil {
		log.Critical(err)
	}

	backendLeveled.SetLevel(logging.ERROR, "")
	if opts.Verbosity {
		backendLeveled.SetLevel(logging.DEBUG, "")
	}
	logging.SetBackend(backendLeveled)

	// TODO: configify this
	//WorkingDir := "/tmp/"
	TemplateName := opts.Template
	OutputName := opts.Output

	var summary Summary

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
			log.Critical(err)
		}

		switch match.Info.Severity {
		case "critical":
			summary.NumCritical++
		case "high":
			summary.NumHigh++
		case "medium":
			summary.NumMedium++
		case "low":
	 		summary.NumLow++
		case "info":
			summary.NumInfo++
		}

		matches = append(matches, match)
	}
	summary.Total = len(matches)

	// TODO: we need to preprocess, esepecially the HTTP requests
	// because LaTeX/templates really struggle by themselves
	processed := preprocess(matches)

	var ct latex.CompileTask = latex.NewCompileTask()
	ct.SetSourceDir(".") // Is this needed?

	var latexVerbosity latex.VerbosityLevel
	switch opts.LatexVerbosity {
	case "all":
		latexVerbosity = latex.VerbosityAll
	case "more":
		latexVerbosity = latex.VerbosityMore
	case "none":
		latexVerbosity = latex.VerbosityNone
	case "default":
	default:
		latexVerbosity = latex.VerbosityDefault
	}
	ct.SetVerbosity(latexVerbosity)

	// populate LaTex template
	tmpl, err := template.ParseFiles(TemplateName)
	if err != nil {
		log.Critical(err)
	}

	// Output tex file? Do we need to?
	log.Debug("Executing template", TemplateName)

	// create temp file
	tempFile, err := os.CreateTemp("", "*.tex")
	if err != nil {
		log.Critical(err)
	}


	var report Report
	report.Summary = summary
	report.Matches = processed

	err = tmpl.Execute(tempFile, report)
	if err != nil {
		log.Critical(err)
	}

	// Run pdflatex
	ct.SetCompileFilename(tempFile.Name())
	log.Debug("Generating PDF file:", ct.CompileFilenamePdf())
	err = ct.Pdflatex(tempFile.Name(), "")

	if err != nil {
		log.Critical(err)
	}

	// Clean up
	// TODO: parameterize this for debugging
	ct.ClearLatexTempFiles(".")
	finalName := filepath.Base(ct.CompileFilenamePdf())
	log.Debug("Moving compiled file from", finalName, "to", OutputName+".pdf")

	err = os.Rename(finalName, OutputName+".pdf")
	if err != nil {
		log.Critical(err)
	}

	if opts.SaveTexFile {
		err = os.Rename(ct.CompileFilename(), "./"+OutputName+".tex")
	}

}
