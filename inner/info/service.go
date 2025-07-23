package info

import (
	"idm/inner/common"
	"idm/inner/database"
)

type Service struct {
}

func (serv *Service) CheckDbConnection(cfg common.Config) bool {
	return database.CheckDbConnection(cfg)
}
