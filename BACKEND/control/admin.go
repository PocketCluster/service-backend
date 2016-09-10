package control

import (
	"net/http"
	"log"
	"net"

	"github.com/zenazn/goji/web"
	"github.com/stkim1/BACKEND/util"
)

// Category route
func (controller *Controller) DashboardOverview(c web.C, r *http.Request) (string, int) {

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
	//forward := r.Header.Get("X-Forwarded-For")

	var content map[string]interface{} = map[string]interface{} {
		"ISINDEX"			   : false,
		"SITENAME"			   : "PocketCluster Index",
		"DEFAULT_LANG"         : "utf-8",
		"SITEURL"              : "https://index.pocketcluster.io",
		"THEME_STATIC_DIR"     : "theme",
	}

	return util.RenderLayout("dashboard/overview.html.mustache", "dashboard/base.html.mustache", content), http.StatusOK
}

// Category route
func (controller *Controller) DashboardRepository(c web.C, r *http.Request) (string, int) {

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
	//forward := r.Header.Get("X-Forwarded-For")

	mode := c.URLParams["mode"]
	if len(mode) == 0 {
	}

	var content map[string]interface{} = map[string]interface{} {
		"ISINDEX"			   : false,
		"SITENAME"			   : "PocketCluster Index",
		"DEFAULT_LANG"         : "utf-8",
		"SITEURL"              : "https://index.pocketcluster.io",
		"THEME_STATIC_DIR"     : "theme",
	}

	return util.RenderLayout("dashboard/repository.html.mustache", "dashboard/base.html.mustache", content), http.StatusOK
}