package models

type ReceptCategory struct {
	ReceptID     int    `db:"recept_id"`
	CategorySlug string `db:"category_slug"`
}

func (ReceptCategory) GetType() string {
	return "ReceptCategory"
}
func (c ReceptCategory) GetArgsInsert() []any {
	return []any{c.ReceptID, c.CategorySlug}
}
