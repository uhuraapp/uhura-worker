package helpers

import "reflect"

func Contains(slice interface{}, val interface{}) bool {
	sv := reflect.ValueOf(slice)

	for i := 0; i < sv.Len(); i++ {
		if sv.Index(i).Interface() == val {
			return true
		}
	}
	return false
}

type mapf func(interface{}) int64

func MapInt(in interface{}, fn mapf) []int64 {
	val := reflect.ValueOf(in)
	out := make([]int64, val.Len())

	for i := 0; i < val.Len(); i++ {
		out[i] = fn(val.Index(i).Interface())
	}
	return out
}

func Returns(slice interface{}, name string, val interface{}) interface{} {
	sv := reflect.ValueOf(slice)

	for i := 0; i < sv.Len(); i++ {
		if sv.Index(i).FieldByName(name).Interface() == val {
			return sv.Index(i).Interface()
		}
	}
	return false
}

func AppendIfMissing(slice []string, i string) []string {
	if i == "" {
		return slice
	}
	for _, e := range slice {
		if e == i {
			return slice
		}
	}
	return append(slice, i)
}

func RemoveDuplicates(s []int64) []int64 {
	m := map[int64]bool{}

	for _, v := range s {
		if _, seen := m[v]; !seen {
			s[len(m)] = v
			m[v] = true
		}
	}

	r := s[:len(m)]

	return r
}
