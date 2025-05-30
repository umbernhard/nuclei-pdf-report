# nuclei-pdf-report
A golang tool for generating PDF reports from nuclei scans

## Requirements

`nuclei-pdf-report` calls `pdflatex`, so some version of that needs to be installed and accessible in your `PATH`. See here for instructions: https://www.tug.org/texlive/

## Usage

`nuclei-pdf-report` expects input on `stdin` of `jsonl` output from nuclei:

```
  nuclei -u example.com -jle | nuclei-pdf-report
```

You can also accomplish this by piping from a file:
```
  cat example.jsonl | nuclei-pdf-report
```

Currently this will put a report called `report.pdf` in your current directory, but this is a work in progress. 
