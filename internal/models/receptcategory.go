package models

type ReceptCategory struct {
	ReceptID   int    `db:"recept_id"`
	CategoryID string `db:"category_id"`
}

func (ReceptCategory) GetType() string {
	return "ReceptCategory"
}
func (c ReceptCategory) GetArgsInsert() []any {
	return []any{c.ReceptID, c.CategoryID}
}
