package api

const (
    APIPrefix          string = "/service"
    URLPackageList     string = APIPrefix + "/v014/package/list"
    URLPackageMeta     string = APIPrefix + "/v014/package/meta/:name"
)

const (
    URLHealthCheck     string = "/healthcheck"
    URLAppStats        string = "/stats"
)

const (
    FSPackageRoot      string = "/api-service/v014/package"
    FilePackageList    string = "list.json"
)