package services

import "log/slog"

type GoodService struct {
	log             *slog.Logger
	storageProvider StorageProvider
}

type StorageProvider interface {
}

func NewGoodService(log *slog.Logger, provider StorageProvider) *GoodService {
	return &GoodService{log: log, storageProvider: provider}
}
