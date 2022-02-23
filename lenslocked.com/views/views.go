package views

import (
	"bytes"
	"github.com/gorilla/csrf"
	"github.com/pkg/errors"
	"github.com/username/project-name/context"
	"html/template"
	"io"
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

	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", errors.New("csrfField is not implemented")
		},
	}).ParseFiles(files...)
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
	v.Render(w, r, nil)
}

// Render is used to render the view with the predefined layout
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	csrfField := csrf.TemplateField(r)
	tmpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})
	if err := tmpl.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		http.Error(w, "Something went wrong. "+
			"If the problem persists please email support@lenslocked.com",
			http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

// addTemplatePath takes in a slice of strings and concatenates TemplateDir and TemplateExt
func addPathAndExt(files []string) {
	for i, s := range files {
		files[i] = TemplateDir + s + TemplateExt
	}
}
