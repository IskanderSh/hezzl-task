package services

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/IskanderSh/hezzl-task/internal/lib/error/wrapper"
	"github.com/IskanderSh/hezzl-task/internal/models"
)

type GoodService struct {
	log             *slog.Logger
	storageProvider StorageProvider
	cacheProvider   CacheProvider
}

type StorageProvider interface {
	Create(projectID int, name string, priority int) (*models.Goods, error)
	GetAllGoods() (*[]models.Goods, error)
	UpdateGood(req *models.UpdateRequest) (*models.Goods, error)
	DeleteGood(req *models.DeleteRequest) (*models.DeleteResponse, error)
}

type CacheProvider interface {
	GetMaxPriority(ctx context.Context) (string, error)
	SetMaxPriority(ctx context.Context, priority int) error
}

func NewGoodService(
	log *slog.Logger,
	provider StorageProvider,
	cache CacheProvider,
) *GoodService {
	return &GoodService{
		log:             log,
		storageProvider: provider,
		cacheProvider:   cache,
	}
}

const (
	defaultPriority = 0
)

func (s *GoodService) CreateGood(ctx context.Context, req *models.CreateRequest) (*models.Goods, error) {
	const op = "services.CreateGood"

	log := s.log.With(slog.String("op", op))

	priorityID := defaultPriority

	priorityString, err := s.cacheProvider.GetMaxPriority(ctx)
	if err != nil {
		log.Warn("no priorityID in cache")
	} else {
		priorityID, err = strconv.Atoi(priorityString)
		if err != nil {
			log.Warn("couldn't convert priority id to int from string", priorityString)
		}
	}

	if priorityID == defaultPriority {
		priorityID, err = s.getMaxPriorityID(ctx)
		if err != nil {
			log.Warn("couldn't get maximum of priorityID from database")
			return nil, wrapper.Wrap(op, err)
		}
	}

	priorityID++ // new priorityID = maxPriorityID + 1

	good, err := s.storageProvider.Create(req.ProjectID, req.Name, priorityID)
	if err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	if err := s.cacheProvider.SetMaxPriority(ctx, priorityID); err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	return good, nil
}

func (s *GoodService) UpdateGood(ctx context.Context, req *models.UpdateRequest) (*models.Goods, error) {
	const op = "services.UpdateGood"

	//log := s.log.With(slog.String("op", op))

	good, err := s.storageProvider.UpdateGood(req)
	if err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	return good, nil
}

func (s *GoodService) DeleteGood(ctx context.Context, req *models.DeleteRequest) (*models.DeleteResponse, error) {
	const op = "services.DeleteGood"

	output, err := s.storageProvider.DeleteGood(req)
	if err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	return output, nil
}

func (s *GoodService) getMaxPriorityID(ctx context.Context) (int, error) {
	const op = "services.getMaxPriorityID"

	//log := s.log.With(slog.String("op", op))

	values, err := s.storageProvider.GetAllGoods()
	if err != nil {
		return 0, wrapper.Wrap(op, err)
	}

	maxPriority := defaultPriority
	for i := 0; i < len(*values); i++ {
		priority := (*values)[i].Priority

		if priority > maxPriority {
			maxPriority = priority
		}
	}

	return maxPriority, nil
}
