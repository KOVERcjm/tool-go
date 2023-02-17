package init

import (
	kLogger "github.com/kovercjm/tool-go/logger"
	kRepository "github.com/kovercjm/tool-go/repository"
	kGORM "github.com/kovercjm/tool-go/repository/gorm"
)

func NewRepository(config *kRepository.Config, logger kLogger.Logger) (kRepository.Repository, error) {
	// TODO add support for other repository types
	return kGORM.Repository{}.Init(config, logger)
}
