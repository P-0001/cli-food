package utils

import "net/url"

func ToForm(data map[any]any) string {
	form := url.Values{}
	for y, z := range data {
		k := ToString(y)
		v := ToString(z)
		form.Set(k, v)
	}
	return form.Encode()
}
