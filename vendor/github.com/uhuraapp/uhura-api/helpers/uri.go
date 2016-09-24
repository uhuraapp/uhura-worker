package helpers

import (
  "regexp"
  "strings"

  . "github.com/fiam/gounidecode/unidecode"
)

type Uriable struct {
}

func (u *Uriable) MakeUri(txt string) string {
  return MakeUri(txt)
}

func MakeUri(txt string) string {
  re := regexp.MustCompile(`\W`)
  uri := Unidecode(txt)
  uri = re.ReplaceAllString(uri, "")
  uri = strings.ToLower(uri)
  return uri
}
