package postgres

import (
	"errors"

	"github.com/IskanderSh/hezzl-task/internal/lib/error/wrapper"
	"github.com/IskanderSh/hezzl-task/internal/models"
)

const (
	removedBase = false
)

func (s *Storage) Create(projectID int, name string, priority int) (*models.Goods, error) {
	const op = "storage.goods.Create"

	var good models.Goods

	row := s.db.QueryRow(createGoodQuery, projectID, name, priority, removedBase)

	if err := row.Scan(&good); err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	return &good, nil
}

func (s *Storage) GetAllGoods() (*[]models.Goods, error) {
	const op = "storage.goods.GetAllGoods"

	rows, err := s.db.Query(getAllGoods)
	if err != nil {
		return nil, wrapper.Wrap(op, err)
	}
	defer rows.Close()

	var goods []models.Goods

	for rows.Next() {
		value := models.Goods{}

		if err := rows.Scan(&value); err != nil {
			return nil, wrapper.Wrap(op, err)
		}

		goods = append(goods, value)
	}

	return &goods, nil
}

func (s *Storage) UpdateGood(req *models.UpdateRequest) (*models.Goods, error) {
	const op = "storage.goods.UpdateGood"

	value := models.Goods{}
	row := s.db.QueryRow(getGood, req.ID, req.ProjectID)
	if err := row.Scan(&value); err != nil {
		return nil, wrapper.Wrap(op, errors.New("no good with this params"))
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
	row := s.db.QueryRow(deleteGood, req.ID, req.ProjectID)
	if err := row.Scan(&value); err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	return &value, nil
}
