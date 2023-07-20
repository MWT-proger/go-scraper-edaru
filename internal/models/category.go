package models

type Category struct {
	Slug       string
	Name       string
	Href       string
	ParentSlug string
}

func (Category) GetType() string {
	return "Category"
}

func (c Category) GetArgsInsert() []any {
	if c.ParentSlug != "" {
		return []any{c.Slug, c.Name, c.Href, c.ParentSlug}
	} else {
		return []any{c.Slug, c.Name, c.Href, nil}
	}
}
