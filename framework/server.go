package framework

import (
	"encoding/json"
	"log"
	"os"
)

func init() {
	readDeploy()
}

var Server server

const (
	ModeDev  = "dev"
	ModeProd = "prod"
)

// server represents the memory structure of deploy.json
type server struct {
	Dev                   *Dev
	Prod                  *Prod
	Mode                  string `json:"mode"`
	CSRF                  string `json:"csrf"`
	Port, PortSSL, Domain string
	UseSSL                bool
	CacheStaticFiles      bool
	CacheTempls           bool
}

type Dev struct {
	Port             string `json:"port"`
	PortSSL          string `json:"port_ssl"`
	UseSSL           bool   `json:"use_ssl"`
	CacheStaticFiles bool   `json:"cache_static_files"`
	CacheTempls      bool   `json:"cache_templates"`
}

type Prod struct {
	User             string `json:"user"`
	Domain           string `json:"domain"`
	IP               string `json:"ip"`
	Password         string `json:"password"`
	Port             string `json:"port"`
	PortSSL          string `json:"port_ssl"`
	UseSSL           bool   `json:"use_ssl"`
	CacheStaticFiles bool   `json:"cache_static_files"`
	CacheTempls      bool   `json:"cache_templates"`
}

func readDeploy() {
	file, err := os.Open("server.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if err = json.NewDecoder(file).Decode(&Server); err != nil {
		log.Fatalln("cannot get configuration from file", err)
	}

	if Server.Mode == ModeProd {
		Server.Port = Server.Prod.Port
		Server.PortSSL = Server.Prod.PortSSL
		Server.Domain = Server.Prod.Domain
		Server.UseSSL = Server.Prod.UseSSL
		Server.CacheStaticFiles = Server.Prod.CacheStaticFiles
		Server.CacheTempls = Server.Prod.CacheTempls
	} else {
		Server.Port = Server.Dev.Port
		Server.PortSSL = Server.Dev.PortSSL
		Server.Domain = "localhost"
		Server.UseSSL = Server.Dev.UseSSL
		Server.CacheStaticFiles = Server.Dev.CacheStaticFiles
		Server.CacheTempls = Server.Dev.CacheTempls
	}
}

func IsDev() bool {
	return Server.Mode == ModeDev
}
