package types

import (
	"html/template"
	"net/http"
)

type XPSource interface {
	TotalXPAmount() float64
}

type GoldSource interface {
	TotalValueInGold() int
}

type Loot interface {
	XPSource
	GoldSource
}

type SQLObject interface {
	DiffWithExisting(SQLObject) (SQLObject, error)
}

type HTMLAble interface {
	HTML(*http.Request) template.Template
}

type Blurbable interface {
	HTMLBlurb(*http.Request) template.Template
}
