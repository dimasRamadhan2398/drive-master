package repositories

import (
	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB	
}

func NewRepositoryRegistry(db *gorm.DB) *Registry {
	return &Registry{db: db}
}

type IRepositoryRegistry interface {
	GetRegion() IRegionRepository
	GetCache() ICacheRepository
	GetEvent() IEventRepository
}

