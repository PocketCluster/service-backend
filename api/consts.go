package api

const (
    APIPrefix          string = "/service"
    URLPackageList     string = APIPrefix + "/v014/package/list"
    URLPackageMeta     string = APIPrefix + "/v014/package/meta/:name"
    URLUserAuth        string = APIPrefix + "/v014/user/auth"
)

const (
    URLHealthCheck     string = "/healthcheck"
    URLAppStats        string = "/stats"
)

const (
    FSPackageListRoot  string = "/api-service/v014/package"
    FilePackageList    string = "list.json"

    FSPackageMetaRoot  string = "/api-service/v014/package/meta"
)