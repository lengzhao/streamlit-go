package ptemplate

import (
	"embed"
	"html/template"
)

//go:embed page.html
var pageTemplateFS embed.FS

// GetPageTemplate 获取页面模板
func GetPageTemplate() (*template.Template, error) {
	return template.ParseFS(pageTemplateFS, "page.html")
}
