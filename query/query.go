package query

import (
	"net/url"
	"reflect"
	"regexp"
)

/**
 * Converts a complex query string into a nested structure
 * E.g. From "a[aa]=11&a[ab]=12&c=3"
 * To map[a:map[aa:11 ab:12] c:3]
 */
func StringToQuery(query string) (map[string]interface{}, error) {
	r := regexp.MustCompile("(.*)\\[(.*)\\]+")

	values, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}

	dst := map[string]interface{}{}

	for prefix, values := range values {
		if len(values) != 1 {
			continue
		}
		value := values[0]

		matches := r.FindAllStringSubmatch(prefix, -1)

		// Plain string, pass through
		if len(matches) == 0 {
			dst[prefix] = value
			continue
		}
		keys := matches[0][1:len(matches[0])]
		keysLen := len(keys)
		next := processKey(dst, value, 0, keysLen-1, keys)
		dst[keys[0]] = next

	}
	return dst, nil
}

func processKey(
	dst map[string]interface{},
	value string,
	depth,
	maxDepth int,
	keys []string,
) map[string]interface{} {

	key := keys[depth]

	// Are we at maximum depth?
	if depth == maxDepth {
		dst[key] = value
		return dst
	}

	// Make the empty key if it doesn't exist
	if _, ok := dst[key]; !ok {
		dst[key] = map[string]interface{}{}
	}

	// Continue up the chain if required
	dstMap, dstOk := mapify(dst[key])
	if dstOk {
		dst = processKey(dstMap, value, depth+1, maxDepth, keys)
	}
	return dst
}

func mapify(i interface{}) (map[string]interface{}, bool) {
	value := reflect.ValueOf(i)
	if value.Kind() == reflect.Map {
		m := map[string]interface{}{}
		for _, k := range value.MapKeys() {
			m[k.String()] = value.MapIndex(k).Interface()
		}
		return m, true
	}
	return map[string]interface{}{}, false
}
