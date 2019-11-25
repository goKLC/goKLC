package goKLC

import (
	"fmt"
	"github.com/flosch/pongo2"
)

var fileExtention = ".html"
var fileBase = "./View/"

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
		fmt.Println(err.Error())
	}

	return content
}

func getFilePath(name string) string {

	return fmt.Sprintf("%s%s%s", fileBase, name, fileExtention)
}