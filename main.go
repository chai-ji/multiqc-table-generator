package main

import (
	"fmt"
	"flag"
	"log"
	"os"
	"encoding/csv"
	"strings"
	"html"
	"text/template"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// https://pkg.go.dev/encoding/csv
func tsvToHtml(path string)(string, error){
	var reader *csv.Reader
	fin, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fin.Close()
	reader = csv.NewReader(fin)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	var htmlStr strings.Builder
	htmlStr.WriteString("<table>\n")

	for {
		record, err := reader.Read()
		if err != nil { // io.EOF ends here too
			break
		}
		htmlStr.WriteString("  <tr>\n")
		for _, cell := range record {
			fmt.Fprintf(&htmlStr, "    <td>%s</td>\n", html.EscapeString(cell))
		}
		htmlStr.WriteString("  </tr>\n")
	}

	htmlStr.WriteString("</table>\n")

	return htmlStr.String(), nil
}

// indent adds n spaces to the start of each line (for YAML block scalars).
func indent(n int, s string) string {
	prefix := strings.Repeat(" ", n)
	// ensure the very first line is also indented
	return prefix + strings.ReplaceAll(s, "\n", "\n"+prefix)
}

// https://pkg.go.dev/text/template
func makeTableYAML(tableStr string)(string){
	tmplStr := `data: |4-
{{ indent 2 .Html }}
description: {{ .Description }}
plot_type: html
section_href: https://github.com/default-manifest-name
section_name: default-manifest-name {{ .SectionName }}
`
	type TemplateData struct {
		Html string
		Description string
		SectionName string
	}
	var templateData = TemplateData{
		Html: tableStr,
		Description: "This is the description",
		SectionName: "Section Name Goes Here",
	}

	tmpl := template.Must(template.New("yaml").
		Funcs(template.FuncMap{"indent": func(n int, s string) string { return indent(n, s) }}).
		Parse(tmplStr))

		var outputStr strings.Builder
	if err := tmpl.Execute(&outputStr, templateData); err != nil {
		panic(err)
	}


		return outputStr.String()
}

func main() {
	// command line args
	// outputFile := flag.String("output", "output.yml", "output filename")
	flag.Parse()

	// first positional arg passed
	inputFile := flag.Arg(0)
	if inputFile == "" {
		log.Fatalf("error: Input file path not provided")
	}

	// fmt.Printf("outputFile: %v\n", *outputFile)
	// fmt.Printf("inputFile: %v\n", inputFile)

	htmlStr, err := tsvToHtml(inputFile)
	if err != nil {
		log.Fatalf("error loading input table: %v", err)
	}

	fmt.Println(makeTableYAML(htmlStr))


}
