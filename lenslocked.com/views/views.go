package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

func layoutFiles(pattern string) []string {
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}
	return files
}

func NewView(layout string, files ...string) *View {
	addPathAndExt(files)

	files = append(files, layoutFiles(TemplateDir+"layouts/*"+TemplateExt)...)

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

// Render is used to render the view with the predefined layout
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

// addTemplatePath takes in a slice of strings and concatenates TemplateDir and TemplateExt
func addPathAndExt(files []string) {
	for i, s := range files {
		files[i] = TemplateDir + s + TemplateExt
	}
}
