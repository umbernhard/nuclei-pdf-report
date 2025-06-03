# nuclei-pdf-report
A golang tool for generating PDF reports from nuclei scans

## Requirements

`nuclei-pdf-report` calls `pdflatex`, so some version of that needs to be installed and accessible in your `PATH`. See here for instructions: https://www.tug.org/texlive/

If you're on Ubuntu LTS, the version of LaTeX in `texlive-full` is really out of date. Consider using homebrew or nix to install it, or do a manual install.

## Usage

`nuclei-pdf-report` expects input on `stdin` of `jsonl` output from nuclei:

```bash
  nuclei -u example.com -j | nuclei-pdf-report
```

You can also accomplish this by piping from a file:
```bash
  cat example.jsonl | nuclei-pdf-report
```

This will put a report called `report.pdf` in your current directory. To specify a different name for this file, you can use the `o` flag:
```bash
cat example.jsonl | nuclei-pdf-report -o another-name
```

This documentation is still a work in progress. 
