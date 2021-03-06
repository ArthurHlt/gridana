package converters

import (
	"bytes"
	"github.com/ArthurHlt/gridana/model"
	"github.com/russross/blackfriday/v2"
	tmplhtml "html/template"
	"regexp"
	"strings"
	tmpltext "text/template"
	"time"
)

func GenHTML(alert model.FormattedAlert, layout string) (string, error) {
	tmpl, err := tmpltext.New("").
		Option("missingkey=zero").
		Funcs(tmpltext.FuncMap(DefaultFuncs)).
		Parse(layout)
	if err != nil {
		return "", err
	}
	b := &bytes.Buffer{}
	err = tmpl.Execute(b, alert)
	if err != nil {
		return "", err
	}
	return string(blackfriday.Run(b.Bytes())), nil
}

func GenText(alert model.FormattedAlert, layout string) (string, error) {
	tmpl, err := tmpltext.New("").
		Option("missingkey=zero").
		Funcs(tmpltext.FuncMap(DefaultFuncs)).
		Parse(layout)
	if err != nil {
		return "", err
	}
	b := &bytes.Buffer{}
	err = tmpl.Execute(b, alert)
	if err != nil {
		return "", err
	}
	text := strings.Replace(b.String(), "\n", "", -1)
	text = strings.Replace(text, " ", "", -1)
	return text, nil
}

type FuncMap map[string]interface{}

var DefaultFuncs = FuncMap{
	"toUpper": strings.ToUpper,
	"toLower": strings.ToLower,
	"title":   strings.Title,
	// join is equal to strings.Join but inverts the argument order
	// for easier pipelining in templates.
	"join": func(sep string, s []string) string {
		return strings.Join(s, sep)
	},
	"timeFormat": func(format string, t time.Time) string {
		return t.Format(format)
	},
	"safeHtml": func(text string) tmplhtml.HTML {
		return tmplhtml.HTML(text)
	},
	"reReplaceAll": func(pattern, repl, text string) string {
		re := regexp.MustCompile(pattern)
		return re.ReplaceAllString(text, repl)
	},
}
