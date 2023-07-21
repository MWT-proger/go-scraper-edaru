package models

import (
	"database/sql"
	"time"
)

type Ingredient struct {
	ID          int            `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	Href        sql.NullString `db:"href"`
	ParentId    sql.NullInt64  `db:"parent_id"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

func (Ingredient) GetType() string {
	return "Ingredient"
}

func (c Ingredient) GetArgsInsert() []any {
	if c.ParentId.Int64 != 0 {
		return []any{c.ID, c.Name, c.Description, c.Href, c.ParentId, c.UpdatedAt}
	} else {
		return []any{c.ID, c.Name, c.Description, c.Href, nil, c.UpdatedAt}
	}
}
