package web

import (
	"github.com/caixr9527/go-cloud/web/render"
	"html/template"
)

type TemplateOps struct {
	TemplatePattern string
	FuncMap         template.FuncMap
	HTMLRender      render.HTMLRender
}
type Options struct {
	TemplateOps
}
