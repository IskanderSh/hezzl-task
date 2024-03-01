package models

import "time"

type Good struct {
	ID          int           `db:"id"`
	ProjectID   int           `db:"project_id"`
	Name        string        `db:"name"`
	Description string        `db:"description"`
	Priority    int           `db:"priority"`
	Removed     bool          `db:"removed"`
	CreatedAt   time.Duration `db:"created_at"`
}

type CreateRequest struct {
	ProjectID int    `json:"project_id,omitempty" db:"project_id"`
	Name      string `json:"name" db:"name"`
	Priority  int    `json:"priority_id,omitempty" db:"priority"`
}

type UpdateRequest struct {
	ID          int    `json:"id,omitempty" db:"id"`
	ProjectID   int    `json:"project_id,omitempty" db:"project_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type DeleteRequest struct {
	ID        int `json:"id" db:"id"`
	ProjectID int `json:"project_id" db:"project_id"`
}

type DeleteResponse struct {
	ID        int  `json:"id" db:"id"`
	ProjectID int  `json:"project_id" db:"project_id"`
	Removed   bool `json:"removed" db:"removed"`
}

type ListGoodsResponse struct {
	Meta  Meta   `json:"meta"`
	Goods []Good `json:"goods"`
}

type Meta struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

type ReprioritizeRequest struct {
	ID          int `db:"id"`
	ProjectID   int `db:"project_id"`
	NewPriority int `json:"newPriority"`
}

type ReprioritizeResponse struct {
	Priorities []Priorities `json:"priorities"`
}

type Priorities struct {
	ID       int `json:"id" db:"id"`
	Priority int `json:"priority" db:"priority"`
}

type GoodCache struct {
	ProjectID   int
	Name        string
	Description string
	Priority    int
	Removed     bool
	CreatedAt   time.Duration
}

type GoodLog struct {
	ID          int           `db:"Id"`
	ProjectID   int           `db:"ProjectId"`
	Name        string        `db:"Name"`
	Description string        `db:"Description"`
	Priority    int           `db:"Priority"`
	Removed     bool          `db:"Removed"`
	EventTime   time.Duration `db:"EventTime"`
}
