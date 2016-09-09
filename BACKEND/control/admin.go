package control


import (
	"net/http"
	"log"
	"net"

	"github.com/zenazn/goji/web"
	"github.com/stkim1/BACKEND/util"
	"github.com/stkim1/BACKEND/model"
)

// Category route
func (controller *Controller) Admin(c web.C, r *http.Request) (string, int) {

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("userip: %q is not IP:port", r.RemoteAddr)
		return "", http.StatusNotFound
	}

	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		log.Printf("userip: %q is not IP:port", r.RemoteAddr)
		return "", http.StatusNotFound
	}
	forward := r.Header.Get("X-Forwarded-For")

	log.Printf("IP %, Forwarded for: %", ip, forward)

	param := c.URLParams["mode"]
	if len(param) == 0 {
	}

	var content map[string]interface{} = map[string]interface{} {
		"ISINDEX"			   : false,
		"SITENAME"			   : "PocketCluster Index",
		"DEFAULT_LANG"         : "utf-8",
		"SITEURL"              : "https://index.pocketcluster.io",
		"THEME_STATIC_DIR"     : "theme",
		"CATEGORIES"		   : model.GetDefaultCategory(),
	}

	return util.RenderPage("dashboard.html.mustache", content), http.StatusOK
}
