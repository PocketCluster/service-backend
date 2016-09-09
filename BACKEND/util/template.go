package util

import (
	"path"
	"log"

	"github.com/cbroglie/mustache"
)


func RenderPage(template string, data map[string]interface{}) string {
	filename := path.Join("views/", template)
	content, err := mustache.RenderFile(filename, data)
	if err != nil {
		log.Panic("Render failed : " + err.Error())
	}
	return content
}

func RenderLayout(bodyTemplate string, baseTemplate string, data map[string]interface{}) string {
	//filename := path.Join(path.Join(os.Getenv("PWD"), "tests"), "test1.mustache")
	basefile := path.Join("views/", baseTemplate)
	bodyfile := path.Join("views/", bodyTemplate)
	content, err := mustache.RenderFileInLayout(bodyfile, basefile, data)
	if err != nil {
		log.Panic("Render failed : " + err.Error())
	}
	return content
}