package postgres

import (
	"errors"
	"fmt"
	"strings"

	"github.com/IskanderSh/hezzl-task/internal/lib/error/wrapper"
	"github.com/IskanderSh/hezzl-task/internal/models"
)

const (
	defaultRemoved = false
)

var (
	ErrGoodNotFound = errors.New("good with such id in project not found")
)

func (s *Storage) Create(req *models.CreateRequest) (*models.Good, error) {
	const op = "storage.goods.Create"

	var good models.Good

	row := s.db.QueryRow(createGoodQuery, req.ProjectID, req.Name, req.Priority, defaultRemoved)

	if err := row.Scan(&good); err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	return &good, nil
}

func (s *Storage) GetAllGoods() (*[]models.Good, error) {
	const op = "storage.goods.GetAllGoods"

	rows, err := s.db.Query(getAllGoods)
	if err != nil {
		return nil, wrapper.Wrap(op, err)
	}
	defer rows.Close()

	var goods []models.Good

	for rows.Next() {
		value := models.Good{}

		if err := rows.Scan(&value); err != nil {
			return nil, wrapper.Wrap(op, err)
		}

		goods = append(goods, value)
	}

	return &goods, nil
}

func (s *Storage) UpdateGood(req *models.UpdateRequest) (*models.Good, error) {
	const op = "storage.goods.UpdateGood"

	value := models.Good{}
	row := s.db.QueryRow(getGood, req.ID, req.ProjectID)
	if err := row.Scan(&value); err != nil {
		return nil, wrapper.Wrap(op, ErrGoodNotFound)
	}

	row = s.db.QueryRow(updateGood, req.Name, req.Description, req.ID, req.ProjectID)
	if err := row.Scan(&value); err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	return &value, nil
}

func (s *Storage) DeleteGood(req *models.DeleteRequest) (*models.DeleteResponse, error) {
	const op = "storage.goods.DeleteGood"

	value := models.DeleteResponse{}
	row := s.db.QueryRow(getGood, req.ID, req.ProjectID)
	if err := row.Scan(&value); err != nil {
		return nil, wrapper.Wrap(op, ErrGoodNotFound)
	}

	row = s.db.QueryRow(deleteGood, req.ID, req.ProjectID)
	if err := row.Scan(&value); err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	return &value, nil
}

func (s *Storage) ListGoods(ids *[]int) (*[]models.Good, error) {
	const op = "storage.goods.ListGoods"

	idsConstraint := strings.Builder{}

	for i, idx := range *ids {
		if i != 0 {
			idsConstraint.WriteString(", ")
		}
		idsConstraint.WriteString(fmt.Sprintf("%d", idx))
	}

	query := fmt.Sprintf(listGoodsWithIds, idsConstraint.String())

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, wrapper.Wrap(op, err)
	}
	defer rows.Close()

	var goods []models.Good

	for rows.Next() {
		value := models.Good{}

		if err := rows.Scan(&value); err != nil {
			return nil, wrapper.Wrap(op, err)
		}

		goods = append(goods, value)
	}

	return &goods, nil
}

func (s *Storage) ReprioritizeGoods(req *models.ReprioritizeRequest) (*[]models.Priorities, error) {
	const op = "storage.goods.ReprioritizeGoods"

	value := models.Good{}
	row := s.db.QueryRow(getGood, req.ID, req.ProjectID)
	if err := row.Scan(&value); err != nil {
		return nil, wrapper.Wrap(op, errors.New("no good with this params"))
	}

	rows, err := s.db.Query(reprioritizeGood, req.NewPriority)
	if err != nil {
		return nil, wrapper.Wrap(op, err)
	}
	defer rows.Close()

	var priorities []models.Priorities

	for rows.Next() {
		value := models.Priorities{}

		if err := rows.Scan(&value); err != nil {
			return nil, wrapper.Wrap(op, err)
		}

		priorities = append(priorities, value)
	}

	return &priorities, nil
}
