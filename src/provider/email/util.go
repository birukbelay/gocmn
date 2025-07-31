package email

import (
	"bytes"
	"html/template"
	"io"
	"io/fs"
)

// ParseEmailTemplate parses temaplate from a given path
func ParseEmailTemplate(templatePath string, templateStruct any) (string, error) {
	// Parse the HTML template file
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	// Create a buffer to store the rendered template
	var buf bytes.Buffer

	// Execute the template with the provided data
	err = tmpl.Execute(&buf, templateStruct)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ParseEmailTemplate parses from an embeded file
func ParseEmbededTemplate(templatePath fs.File, templateStruct any) (string, error) {
	b, err := io.ReadAll(templatePath)
	if err != nil {
		panic(err)
	}
	// Parse the HTML template file
	tmpl, err := template.New("tmpl").Parse(string(b))
	if err != nil {
		return "", err
	}

	// Create a buffer to store the rendered template
	var buf bytes.Buffer

	// Execute the template with the provided data
	err = tmpl.Execute(&buf, templateStruct)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
