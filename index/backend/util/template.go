package util

import (
    "path"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/cbroglie/mustache"
)

func RenderPage(templatePath, templateFile string, data map[string]interface{}) string {
    filename := path.Join(templatePath, templateFile)
    content, err := mustache.RenderFile(filename, data)
    if err != nil {
        log.Error(errors.WithMessage(err, "RenderPage failed "))
        return ""
    }
    return content
}

func RenderLayout(templatePath, bodyTemplate, baseTemplate string, data map[string]interface{}) string {
    //filename := path.Join(path.Join(os.Getenv("PWD"), "tests"), "test1.mustache")
    basefile := path.Join(templatePath, baseTemplate)
    bodyfile := path.Join(templatePath, bodyTemplate)
    content, err := mustache.RenderFileInLayout(bodyfile, basefile, data)
    if err != nil {
        log.Error(errors.WithMessage(err, "RenderLayout failed "))
    }
    return content
}