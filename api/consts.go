package api

const (
    APIPrefix          string = "/service"
    URLPackageList     string = APIPrefix + "/v014/package/list"
    URLPackageRepo     string = APIPrefix + "/v014/package/repo"
    URLPackageSync     string = APIPrefix + "/v014/package/sync/:name"
    URLPackageMeta     string = APIPrefix + "/v014/package/meta/:name"
    URLAuthCheck       string = APIPrefix + "/v014/auth/check"


    URLHealthCheck     string = "/healthcheck"
    URLAppStats        string = "/stats"


    FilePackageList    string = "list.json"
    FilePackageRepo    string = "srcs.json"

    FSPackageRootList  string = "/api-service/v014/package"
    FSPackageRootRepo  string = "/api-service/v014/package/repo"
    FSPackageRootSync  string = "/api-service/v014/package/sync"
    FSPackageRootMeta  string = "/api-service/v014/package/meta"
)