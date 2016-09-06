package util

import (
	"html/template"
	"bytes"
	"path"

	"github.com/cbroglie/mustache"
	"github.com/Sirupsen/logrus"
)

func Parse(t *template.Template, name string, data interface{}) string {
	var doc bytes.Buffer
	t.ExecuteTemplate(&doc, name, data)
	return doc.String()
}

func Render(bodyTemplate string, baseTemplate string, data map[string]interface{}) string {
	//filename := path.Join(path.Join(os.Getenv("PWD"), "tests"), "test1.mustache")
	basefile := path.Join("views/", baseTemplate)
	bodyfile := path.Join("views/", bodyTemplate)
	content, err := mustache.RenderFileInLayout(bodyfile, basefile, data)
	if err != nil {
		logrus.Panic("Render failed : " + err.Error())
	}
	return content
}