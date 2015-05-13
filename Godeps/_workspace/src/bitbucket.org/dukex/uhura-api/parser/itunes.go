package parser

import (
	rss "github.com/jteeuwen/go-pkg-rss"
)

type iTunes struct{}

type iTunesExtensiable interface {
	GetExtensions(string) map[string][]rss.Extension
}

func (c iTunes) value(f iTunesExtensiable, k string) string {
	if i := c.get(f, k); len(i) > 0 {
		return i[0].Value
	}
	return ""
}

func (c iTunes) attr(f iTunesExtensiable, k, attr string) string {
	if i := c.get(f, k); len(i) > 0 {
		return i[0].Attrs[attr]
	}
	return ""
}

func (c iTunes) get(f iTunesExtensiable, k string) []rss.Extension {
	if f != nil &&
		f.GetExtensions(ITUNES_EXT) != nil &&
		f.GetExtensions(ITUNES_EXT)[k] != nil &&
		len(f.GetExtensions(ITUNES_EXT)[k]) > 0 {
		return f.GetExtensions(ITUNES_EXT)[k]
	}
	return nil
}
