package api

const (
    APIPrefix          string = "/service"
    URLPackageList     string = APIPrefix + "/v014/package/list"
    URLPackageSync     string = APIPrefix + "/v014/package/sync/:name"
    URLPackageMeta     string = APIPrefix + "/v014/package/meta/:name"
    URLUserAuth        string = APIPrefix + "/v014/user/auth"
)

const (
    URLHealthCheck     string = "/healthcheck"
    URLAppStats        string = "/stats"
)

const (
    FilePackageList    string = "list.json"

    FSPackageRootList  string = "/api-service/v014/package"
    FSPackageRootSync  string = "/api-service/v014/package/sync"
    FSPackageRootMeta  string = "/api-service/v014/package/meta"
)