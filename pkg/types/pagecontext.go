package types

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type PageData struct {
	Title       string
	Description string
}

func NewPageData(title string, description string) *PageData {
	return &PageData{
		Title:       title,
		Description: description,
	}
}

type PageContext struct {
	Title         string
	Locale        string
	Localizer     *i18n.Localizer
	UniTranslator ut.Translator
	NavItems      []NavigationItem
	Pathname      string
}

func (p *PageContext) T(k string, args ...map[string]interface{}) string {
	if len(args) > 1 {
		panic("T(): too many arguments")
	}
	if len(args) == 0 {
		return p.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: k})
	}
	return p.Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: k, TemplateData: args[0]})
}
