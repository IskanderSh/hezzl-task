package services

import (
	"context"
	"encoding/json"
	"fmt"
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
	Create(req *models.CreateRequest) (*models.Good, error)
	GetAllGoods() (*[]models.Good, error)
	UpdateGood(req *models.UpdateRequest) (*models.Good, error)
	DeleteGood(req *models.DeleteRequest) (*models.DeleteResponse, error)
	ListGoods(ids *[]int) (*[]models.Good, error)
	ReprioritizeGoods(req *models.ReprioritizeRequest) (*[]models.Priorities, error)
}

type CacheProvider interface {
	GetMaxPriority(ctx context.Context) (string, error)
	SetMaxPriority(ctx context.Context, priority int) error
	SaveGood(ctx context.Context, key string, value *models.GoodCache) error
	GetGood(ctx context.Context, key string) (string, error)
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

func (s *GoodService) CreateGood(ctx context.Context, req *models.CreateRequest) (*models.Good, error) {
	const op = "services.CreateGood"

	log := s.log.With(slog.String("op", op))

	priority := defaultPriority

	priorityString, err := s.cacheProvider.GetMaxPriority(ctx)
	if err != nil {
		log.Warn("no priorityID in cache")
	} else {
		priority, err = strconv.Atoi(priorityString)
		if err != nil {
			log.Warn("couldn't convert priority id to int from string", priorityString)
		}
	}

	if priority == defaultPriority {
		priority, err = s.getMaxPriorityID(ctx)
		if err != nil {
			log.Warn("couldn't get maximum of priorityID from database")
			return nil, wrapper.Wrap(op, err)
		}
	}

	priority++ // new priorityID = maxPriorityID + 1

	req.Priority = priority

	good, err := s.storageProvider.Create(req)
	if err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	if err := s.cacheProvider.SetMaxPriority(ctx, priority); err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	return good, nil
}

func (s *GoodService) UpdateGood(ctx context.Context, req *models.UpdateRequest) (*models.Good, error) {
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

func (s *GoodService) GetGoods(ctx context.Context, limit, offset int) (*models.ListGoodsResponse, error) {
	const op = "services.GetGoods"

	log := s.log.With(slog.String("op", op))

	idsNotInCache := make([]int, 0, limit)
	output := make([]models.Good, 0, limit)

	for id := offset; id <= offset+limit; id++ {
		key := fmt.Sprintf("%d", id)
		value, err := s.cacheProvider.GetGood(ctx, key)
		if err != nil {
			log.Info("there is no value in cache with key", key)
			idsNotInCache = append(idsNotInCache, id)
		} else {
			good := models.Good{}

			if err := json.Unmarshal([]byte(value), &good); err != nil {
				log.Warn("couldn't unmarshal")
				idsNotInCache = append(idsNotInCache, id)
				continue
			}

			output = append(output, good)
		}
	}

	goodsNotInCache, err := s.storageProvider.ListGoods(&idsNotInCache)
	if err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	for _, good := range *goodsNotInCache {
		output = append(output, good)

		key := fmt.Sprintf("%d", good.ID)
		value := &models.GoodCache{
			ProjectID:   good.ProjectID,
			Name:        good.Name,
			Description: good.Description,
			Priority:    good.Priority,
			Removed:     good.Removed,
			CreatedAt:   good.CreatedAt,
		}

		if err := s.cacheProvider.SaveGood(ctx, key, value); err != nil {
			log.Warn("couldn't save good to cache", key)
		}
	}

	total := 0
	removed := 0
	for _, value := range output {
		if value.Removed {
			removed++
		}
		total++
	}

	return &models.ListGoodsResponse{
		Meta: models.Meta{
			Total:   total,
			Removed: removed,
			Limit:   limit,
			Offset:  offset,
		},
		Goods: output,
	}, nil
}

func (s *GoodService) ReprioritizeGood(ctx context.Context, req *models.ReprioritizeRequest) (*models.ReprioritizeResponse, error) {
	const op = "services.ReprioritizeGood"

	output, err := s.storageProvider.ReprioritizeGoods(req)
	if err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	return &models.ReprioritizeResponse{
		Priorities: *output,
	}, nil
}

func (s *GoodService) getMaxPriorityID(ctx context.Context) (int, error) {
	const op = "services.getMaxPriorityID"

	//log := s.log.With(slog.String("op", op))

	values, err := s.storageProvider.GetAllGoods()
	if err != nil {
		return 0, wrapper.Wrap(op, err)
	}

	maxPriority := defaultPriority
	for _, value := range *values {
		priority := value.Priority

		if priority > maxPriority {
			maxPriority = priority
		}
	}

	return maxPriority, nil
}
