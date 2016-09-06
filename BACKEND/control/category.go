package control

import (
	"net/http"
	"strings"

	"github.com/zenazn/goji/web"
	"github.com/jinzhu/gorm"
	"github.com/stkim1/BACKEND/util"
	"github.com/stkim1/BACKEND/model"
)

// Home page route
func (controller *Controller) Category(c web.C, r *http.Request) (string, int) {
	const singleColumnCount int = 8
	const totalRowCount int = 3

	var repositories []model.Repository
	var repo1, repo2, repo3 []*model.Repository
	var category string = strings.ToLower(c.URLParams["cat"])
	//FIXME : Titalize
	var title string = strings.TrimSpace(c.URLParams["cat"])
	if !model.IsCategoryPresent(category) {
		return "", http.StatusNotFound
	}

	var db *gorm.DB = controller.GetGORM(c)
	db.Where("category = ?", category).Order("updated desc").Limit(singleColumnCount * totalRowCount).Find(&repositories)
	if len(repositories) == 0 {
		return "", http.StatusNotFound
	}

	// if queried repositories are smaller than required for a page
	// FIXME : use recursive func
	var repoCount int = len(repositories)
	if repoCount <= singleColumnCount {
		repo1 = make([]*model.Repository, repoCount)
	} else {
		repo1 = make([]*model.Repository, singleColumnCount)

		var repoRemain int = repoCount - singleColumnCount
		if repoRemain <= singleColumnCount {
			repo2 = make([]*model.Repository, repoRemain)
		} else {
			repo2 = make([]*model.Repository, singleColumnCount)

			repo3 = make([]*model.Repository, repoRemain - singleColumnCount)
		}
	}

	for index, _ := range repositories {
		var subindex int = index % singleColumnCount
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
		"ISINDEX"			   : false,
		"SITENAME"			   : "PocketCluster Index",
		"DEFAULT_LANG"         : "utf-8",
		"SITEURL"              : "https://index.pocketcluster.io",
		"THEME_STATIC_DIR"     : "theme",
		"CATEGORIES"		   : model.GetActivatedCategory(category),
		"title"				   : title,
		"repo1"				   : &repo1,
		"repo2"				   : &repo2,
		"repo3"				   : &repo3,
	}
	return util.Render("index.html.mustache", "base.html.mustache", content), http.StatusOK
}
