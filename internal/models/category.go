package models

import "database/sql"

type Category struct {
	Slug       string         `db:"slug"`
	Name       string         `db:"name"`
	Href       string         `db:"href"`
	ParentSlug sql.NullString `db:"parent_slug"`
}

func (Category) GetType() string {
	return "Category"
}

func (c Category) GetArgsInsert() []any {
	if c.ParentSlug.String != "" {
		return []any{c.Slug, c.Name, c.Href, c.ParentSlug}
	} else {
		return []any{c.Slug, c.Name, c.Href, nil}
	}
}
