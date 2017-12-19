package cfg

import (
	"encoding/json"
	"log"
	//"os"
	"sync"

	"github.com/toolkits/file"
)

type HttpConfig struct {
	Addr    string `json:"addr"`
	Timeout int    `json:"timeout"`
}

type ClientVerifyConfig struct {
	Flag             bool     `json:"flag"`
	ClientVerifyPath string   `json:"client_verify_file"`
	VerifyVersion    []int    `json:"verify_version"`
	PassVersion      []int    `json:"pass_version"`
	PassPlatForm     []string `json:"pass_platform"`
	ErrTips          string   `json:"error"`
}

type SecurityConfig struct {
	Verify            bool                `json:"verify"`
	ServerCertPath    string              `json:"server_cert_file"`
	ServerKeyPath     string              `json:"server_key_file"`
	ServerKeyPassword string              `json:"server_key_pwd"`
	CaCertPath        string              `json:"ca_cert_file"`
	IssuerKeyPath     string              `json:"issuer_key_file"`
	IssuerKeyPassword string              `json:"issuer_key_pwd"`
	OldestVer         []int               `json:"oldest_supported_ver"`
	ClientVerify      *ClientVerifyConfig `json:"client_verify"`
}

type GlobalConfig struct {
	Debug          bool            `json:"debug"`
	Http           *HttpConfig     `json:"http"`
	Security       *SecurityConfig `json:"security"`
	RemoteExchange bool            `json:"remote_exchange"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
}
