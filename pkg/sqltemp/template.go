package sqltemp

import (
	"bytes"
	"reflect"
	"regexp"
	"strings"
	"text/template"
)

// генерация скрипта
func Template(name, sqlt string, data any) (string, error) {
	var re = regexp.MustCompile(`[ ]{2,}|[\t\n]+`)
	var sqlBuf bytes.Buffer

	tmp, err := template.New(name).Funcs(template.FuncMap{
		"isnnil": isnnil,
		"up":     up,
	}).Parse(sqlt)
	if err != nil {
		return "", err
	}
	err = tmp.Execute(&sqlBuf, data)
	if err != nil {
		return "", err
	}

	s := re.ReplaceAllString(sqlBuf.String(), ` `)

	return s, nil
}

// проверка аргумента на nil
func isnnil(obj ...any) bool {
	for _, c := range obj {
		if !(c == nil || (reflect.ValueOf(c).Kind() == reflect.Ptr && reflect.ValueOf(c).IsNil())) {
			return true
		}
	}
	return false
}

// проверка аргумента на nil
func up(s string) string {
	return strings.ToUpper(s)
}
