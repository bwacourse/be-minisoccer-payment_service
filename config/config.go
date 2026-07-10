package config

import (
	"os"
	"payment-service/common/util"

	"github.com/sirupsen/logrus"
)

var Config AppConfig

type AppConfig struct {
	// Add your configuration fields here
	Port                  int             `json:"port"`
	AppName               string          `json:"appName"`
	AppEnv                string          `json:"appEnv"`
	SignatureKey          string          `json:"signatureKey"`
	Database              Database        `json:"database"`
	RateLimiterMaxRequest float64         `json:"rateLimiterMaxRequest"`
	RateLimiterTimeSecond int             `json:"rateLimiterTimeSecond"`
	InternalService       InternalService `json:"internalService"`
	// Config for GCP
	GCSType                    string `json:"gcsType"`
	GCSProjectID               string `json:"gcsProjectID"`
	GCSPrivateKeyID            string `json:"gcsPrivateKeyID"`
	GCSPrivateKey              string `json:"gcsPrivateKey"`
	GCSClientEmail             string `json:"gcsClientEmail"`
	GCSClientID                string `json:"gcsClientID"`
	GCSAuthURI                 string `json:"gcsAuthURI"`
	GCSTokenURI                string `json:"gcsTokenURI"`
	GCSAuthProviderX509CertURL string `json:"gcsAuthProviderX509CertURL"`
	GCSClientX509CertURL       string `json:"gcsClientX509CertURL"`
	GCSUniverseDomain          string `json:"gcsUniverseDomain"`
	GCSBucketName              string `json:"gcsBucketName"`
}

type Database struct {
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	Name                  string `json:"name"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	MaxOpenConnections    int    `json:"maxOpenConnections"`
	MaxLifeTimeConnection int    `json:"maxLifeTimeConnection"`
	MaxIdleConnections    int    `json:"maxIdleConnections"`
	MaxIdleTime           int    `json:"maxIdleTime"`
}

type InternalService struct {
	User User `json:"user"`
}

type User struct {
	Host         string `json:"host"`
	SignatureKey string `json:"signatureKey"`
}

func Init() {
	err := util.BindFromJSON(&Config, "config.json", ".")
	if err != nil {
		logrus.Infof("failed to bind config from json: %v", err)

		err = util.BindFromConsul(&Config, os.Getenv("CONSUL_HTTP_URL"), os.Getenv("CONSUL_HTTP_KEY"))
		if err != nil {
			panic(err)
		}
	}
}
