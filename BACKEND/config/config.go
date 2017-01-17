package config

import (
    "io/ioutil"

    "gopkg.in/yaml.v2"
)

const (
    EnvConfigFilePath string = "CONFIG_PATH"
)

type General struct {
    TemplatePath    string    `yaml:"template_path"`
    ReadmePath      string    `yaml:"readme_path"`
    ServerPort      int       `yaml:"server_port"`
    MaxConcurrency  int       `yaml:"max_concurrency"`
}

type Site struct {
    SiteName        string    `yaml:"site_name"`
    SiteURL         string    `yaml:"site_url"`
    ThemeLink       string    `yaml:"theme_link"`
}

type Cookie struct {
    MacSecret       string    `yaml:"mac_secret"`
    Secure          bool      `yaml:"secure"`
}

type Database struct {
    DatabaseType    string    `yaml:"database_type"`
    DatabasePath    string    `yaml:"database_path"`
}

type Supplement struct {
    DatabasePath    string    `yaml:"database_path"`
}

type CSRF struct {
    Key             string    `yaml:"key"`
    Cookie          string    `yaml:"coockie"`
    Header          string    `yaml:"header"`
}

type Github struct {
    ClientID        string    `yaml:"client_id"`
    ClientSecret    string    `yaml:"client_secret"`
}

type VPN struct {
    ConnCheck       bool      `yaml:"conn_check"`
    VpnHost         string    `yaml:"vpn_host"`
}

type Update struct {
    ForceReadme        bool      `yaml:"force_readme"`
    // last updated record
    MetaUpdateRecord   string    `yaml:"meta_update_record"`
    // meta update loop period in minutes
    MetaUpdateInterval int64     `yaml:"meta_update_interval"`
    // individual repo meta update cycle in minutes
    MetaUpdateCycle    int64     `yaml:"meta_update_cycle"`
    // last updated record
    SuppUpdateRecord   string    `yaml:"supp_update_record"`
    // supp update loop period in minutes
    SuppUpdateInterval int64     `yaml:"supp_update_interval"`
    // individual repo supp update cycle in minutes
    SuppUpdateCycle    int64     `yaml:"supp_update_cycle"`
    // how many entities should we read from github
    MaxReleaseCollect  int       `yaml:"max_release_collect"`
    // how many entities should we rebuild for display
    MaxReleaseRebuild  int       `yaml:"max_release_rebuild"`
}

type Search struct {
    IndexStoragePath   string    `yaml:"index_storage_path"`
}

type Config struct {
    General         `yaml:"general",inline,flow`
    Site            `yaml:"site",inline,flow`
    Cookie          `yaml:"cookie",inline,flow`
    Database        `yaml:"database",inline,flow`
    Supplement      `yaml:"supplement",inline,flow`
    CSRF            `yaml:"csrf",inline,flow`
    Github          `yaml:"github",inline,flow`
    VPN             `yaml:"vpn",inline,flow`
    Update          `yaml:"update",inline,flow`
    Search          `yaml:"search",inline,flow`
}

func NewConfig(filepath string) (*Config, error) {
    data, err := ioutil.ReadFile(filepath)
    if err != nil {
        return nil, err
    }

    config := Config{}
    err = yaml.Unmarshal(data, &config)
    return &config, err
}

func SaveConfig(config *Config, filepath string) error {
    data, err := yaml.Marshal(config)
    if err != nil {
        return err
    }
    return ioutil.WriteFile(filepath, data, 0600)
}