package tagfast

import (
	"reflect"
	"strconv"
	"sync"
)

//[struct_name][field_name]
var CachedStructTags map[string]map[string]TagFast = make(map[string]map[string]TagFast)
var lock *sync.RWMutex = new(sync.RWMutex)

func CacheTag(struct_name string, field_name string, value TagFast) {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := CachedStructTags[struct_name]; !ok {
		CachedStructTags[struct_name] = make(map[string]TagFast)
	}
	CachedStructTags[struct_name][field_name] = value
}

func GetTag(struct_name string, field_name string) (r TagFast, ok bool) {
	lock.RLock()
	defer lock.RUnlock()
	var v map[string]TagFast
	v, ok = CachedStructTags[struct_name]
	if !ok {
		return
	}
	r, ok = v[field_name]
	return
}

//usage: Tag(typ.Name(),typ.Field(i).Name,typ.Field(i).Tag,"form")
func Tag(i_struct interface{}, field_no int, key string) (tag string) {
	t := reflect.TypeOf(i_struct)
	if t.Field(field_no).Tag == "" {
		return ""
	}
	if v, ok := GetTag(t.String(), t.Field(field_no).Name); ok {
		tag = v.Get(key)
	} else {
		v := TagFast{Tag: t.Field(field_no).Tag}
		tag = v.Get(key)
		CacheTag(t.String(), t.Field(field_no).Name, v)
	}
	return
}

func ClearTag() {
	CachedStructTags = make(map[string]map[string]TagFast)
}

func ParseStructTag(tag string) map[string]string {
	lock.Lock()
	defer lock.Unlock()
	var tagsArray map[string]string = make(map[string]string)
	for tag != "" {
		// skip leading space
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// scan to colon.
		// a space or a quote is a syntax error
		i = 0
		for i < len(tag) && tag[i] != ' ' && tag[i] != ':' && tag[i] != '"' {
			i++
		}
		if i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// scan quoted string to find value
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, _ := strconv.Unquote(qvalue)
		tagsArray[name] = value
	}
	return tagsArray
}

type TagFast struct {
	Tag    reflect.StructTag
	Cached map[string]string
}

func (a *TagFast) Get(key string) string {
	if a.Cached == nil {
		a.Cached = ParseStructTag(string(a.Tag))
	}
	lock.RLock()
	defer lock.RUnlock()
	if v, ok := a.Cached[key]; ok {
		return v
	}
	return ""
}
