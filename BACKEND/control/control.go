package control

import (
    "bytes"
    "fmt"
    "strings"
    "net"
    "net/http"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/gorilla/sessions"
    "github.com/zenazn/goji/web"
    "github.com/jinzhu/gorm"
    "github.com/google/go-github/github"

    "github.com/stkim1/BACKEND/model"
    "github.com/stkim1/BACKEND/config"
)

/* ------- GITHUG API CONTROL ------- */
const (
    githubWebURL string             = "https://github.com/"
)

func NewController(config *config.Config) *Controller {
    return &Controller{
        Config:         config,
        githubClient:   githubV3Client(config.Github.ClientID, config.Github.ClientSecret),
    }
}

// Model-View-Control
type Controller struct {
    githubClient        *github.Client
    *config.Config
}

func (ctrl *Controller) GetGithubRepoMeta(repoURL string) (*github.Repository, *github.Response, error) {
    // TODO : check if URL is in correct form
    if len(repoURL) == 0 {
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL address")
    }
    url := strings.Split(strings.Replace(repoURL , githubWebURL, "", -1), "/")
    owner, repo := url[0], url[1]
    if len(owner) == 0 || len(repo) == 0{
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL format")
    }
    return ctrl.githubClient.Repositories.Get(owner, repo)
}

func (ctrl *Controller) GetGithubContributors(repoURL string) ([]*github.Contributor, *github.Response, error) {
    // TODO : check if URL is in correct form
    if len(repoURL) == 0 {
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL address")
    }
    url := strings.Split(strings.Replace(repoURL , githubWebURL, "", -1), "/")
    owner, repo := url[0], url[1]
    if len(owner) == 0 || len(repo) == 0{
        return nil, nil, fmt.Errorf("[ERR] Invalid repository URL format")
    }
    opts := &github.ListContributorsOptions{Anon: "true"}
    return ctrl.githubClient.Repositories.ListContributors(owner, repo, opts)
}

func (ctrl *Controller) GetSession(c web.C) *sessions.Session {
    return c.Env["Session"].(*sessions.Session)
}

func (ctrl *Controller) GetGORM(c web.C) *gorm.DB {
    return c.Env["GORM"].(*gorm.DB)
}

func (ctrl *Controller) IsXhr(c web.C) bool {
    return c.Env["IsXhr"].(bool)
}

/*
    repo1, repo2, repo3 := GetAssignedRepoColumn(len(repositories))
    for index, _ := range repositories {
        subindex := index % SingleColumnCount
        switch int(index / SingleColumnCount) {
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
*/

// Assignment of repo to column within a page
const SingleColumnCount int = 10
const TotalRowCount int = 3

func GetAssignedRepoColumn(repoCount int) ([]*model.Repository, []*model.Repository, []*model.Repository) {

    remainCount := repoCount
    remainCheck := func() int {
        if remainCount <= 0 {
            return 0
        } else if 0 < remainCount && remainCount <= SingleColumnCount {
            count := remainCount
            remainCount = 0
            return count
        } else {
            remainCount -= SingleColumnCount
            return SingleColumnCount
        }
    }

    remain := remainCheck()
    var repo1, repo2, repo3 []*model.Repository
    if remain != 0 {
        repo1 = make([]*model.Repository, remain)
    }
    remain = remainCheck()
    if remain != 0 {
        repo2 = make([]*model.Repository, remain)
    }
    remain = remainCheck()
    if remain != 0 {
        repo3 = make([]*model.Repository, remain)
    }
    return repo1, repo2, repo3
}

/* ------- GETTING IP ADDRESS ------- */
//ipRange - a structure that holds the start and end of a range of ip addresses
type ipRange struct {
    start net.IP
    end net.IP
}

// inRange - check to see if a given ip address is within a range given
func inRange(r ipRange, ipAddress net.IP) bool {
    // strcmp type byte comparison
    if 0 < bytes.Compare(ipAddress, r.start) && 0 > bytes.Compare(ipAddress, r.end) {
        return true
    }
    return false
}

var privateRanges = []ipRange{
    ipRange{
        start: net.ParseIP("10.0.0.0"),
        end: net.ParseIP("10.255.255.255"),
    },
    ipRange{
        start: net.ParseIP("100.64.0.0"),
        end: net.ParseIP("100.127.255.255"),
    },
    ipRange{
        start: net.ParseIP("172.16.0.0"),
        end: net.ParseIP("172.31.255.255"),
    },
    ipRange{
        start: net.ParseIP("192.0.0.0"),
        end: net.ParseIP("192.0.0.255"),
    },
    ipRange{
        start: net.ParseIP("192.168.0.0"),
        end: net.ParseIP("192.168.255.255"),
    },
    ipRange{
        start: net.ParseIP("198.18.0.0"),
        end: net.ParseIP("198.19.255.255"),
    },
}

// isPrivateSubnet - check to see if this ip is in a private subnet
func isPrivateSubnet(ipAddress net.IP) bool {
    // my use case is only concerned with ipv4 atm
    if ipCheck := ipAddress.To4(); ipCheck != nil {
        // iterate over all our ranges
        for _, r := range privateRanges {
            // check if this ip is in a private range
            if inRange(r, ipAddress){
                return true
            }
        }
    }
    return false
}

func getIPAdress(r *http.Request) string {
    for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
        addresses := strings.Split(r.Header.Get(h), ",")
        // march from right to left until we get a public address
        // that will be the address right before our proxy.
        for i := len(addresses) -1 ; i >= 0; i-- {
            ip := addresses[i]
            // header can contain spaces too, strip those out.
            realIP := net.ParseIP(strings.Replace(ip, " ", "", -1))
            if !realIP.IsGlobalUnicast() && !isPrivateSubnet(realIP) {
                // bad address, go to next
                continue
            }
            return ip
        }
    }
    return ""
}

func githubV3Client(clientID, clientSecret string) *github.Client {
    tp := &github.UnauthenticatedRateLimitedTransport{
        ClientID:     clientID,
        ClientSecret: clientSecret,
    }
    return github.NewClient(tp.Client())
}

func (ctrl *Controller) IsSafeConnection(r *http.Request) bool {
    if !ctrl.VPN.ConnCheck {
        return true
    }
/*
    // access control based on IP
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
    forwarded := r.Header.Get("X-Forwarded-For")
    log.Print("Client IP " + string(clientIP) + " forwarded " + forwarded)
 */
    // access control based on IP
    ipAddress := getIPAdress(r)
    if ipAddress != ctrl.Config.VPN.VpnHost {
        log.Error(trace.Wrap(fmt.Errorf("Cannot display page without proper access from VPN : src[%s]", ipAddress)))
        return false
    }
    return true
}