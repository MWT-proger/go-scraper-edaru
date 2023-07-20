package models

type BaseModeler interface {
	GetType() string
	GetArgsInsert() []any
}
