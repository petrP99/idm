package info

import (
	"idm/inner/common"
	"idm/inner/database"
)

type Service struct {
}

func NewConnectionService() *Service {
	return &Service{}
}

func (serv *Service) CheckDbConnection(cfg common.Config) bool {
	return database.CheckDbConnection(cfg)
}
