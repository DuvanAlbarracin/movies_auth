package config

import (
	"os"

	"github.com/DuvanAlbarracin/movies_auth/pkg/utils"
)

type Config struct {
	Port         string
	DBUrl        string
	JWTSecretKey string
}

func LoadConfig() (c Config, err error) {
	portData, err := os.ReadFile("/run/secrets/auth_port")
	if err != nil {
		return
	}

	dbUrlData, err := os.ReadFile("/run/secrets/db_url")
	if err != nil {
		return
	}

	jwtData, err := os.ReadFile("/run/secrets/jwt_key")
	if err != nil {
		return
	}

	c.Port = utils.TrimString(string(portData))
	c.DBUrl = utils.TrimString(string(dbUrlData))
	c.JWTSecretKey = utils.TrimString(string(jwtData))
	return
}
