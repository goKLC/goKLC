package goKLC

import (
	"fmt"
	"github.com/flosch/pongo2"
)

type Context map[string]interface{}

type View struct {
	Filename string
	Context  Context
}

func NewView(fileName string, context Context) *View {
	return &View{Filename: fileName, Context: context}
}

func (v *View) Render() string {
	template := pongo2.Must(pongo2.FromFile(getFilePath(v.Filename)))
	content, err := template.Execute(pongo2.Context(v.Context))

	if err != nil {
		_app.Log().Error(err.Error(), nil)
	}

	return content
}

func getFilePath(name string) string {
	var fileExtension = _config.Get("ViewFileExtension", ".html")
	var fileBase = _config.Get("ViewFolder", "./View/")

	return fmt.Sprintf("%s%s%s", fileBase, name, fileExtension)
}
