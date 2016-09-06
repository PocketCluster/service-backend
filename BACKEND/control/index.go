package control

import (
	"net/http"

	"github.com/zenazn/goji/web"
	"github.com/jinzhu/gorm"
	"github.com/stkim1/BACKEND/util"
	"github.com/stkim1/BACKEND/model"
)

// Home page route
func (controller *Controller) Index(c web.C, r *http.Request) (string, int) {
	const singleColumnCount int = 8
	const totalRowCount int = 3

	var repositories []model.Repository
	repo1 := make([]*model.Repository, singleColumnCount)
	repo2 := make([]*model.Repository, singleColumnCount)
	repo3 := make([]*model.Repository, singleColumnCount)

	var db *gorm.DB = controller.GetGORM(c)
	db.Order("updated desc").Limit(singleColumnCount * totalRowCount).Find(&repositories)
	if len(repositories) == 0 {
		return "", http.StatusNotFound
	}

	for index, _ := range repositories {
		subindex := index % singleColumnCount
		switch int(index / singleColumnCount ) {
			case 0: {
				repo1[subindex] = &repositories[index]
				break
			}
			case 1: {
				repo2[subindex] = &repositories[index]
				break
			}
			case 2: {
				repo3[subindex] = &repositories[index]
				break
			}
		}
	}

	var content map[string]interface{} = map[string]interface{} {
		"ISINDEX"			   : true,
		"DEFAULT_LANG"         : "utf-8",
		"SITEURL"              : "https://index.pocketcluster.io",
		"THEME_STATIC_DIR"     : "theme",
		"repo1"				   : &repo1,
		"repo2"				   : &repo2,
		"repo3"				   : &repo3,
	}
	return util.Render("index.html.mustache", "base.html.mustache", content), http.StatusOK
}
