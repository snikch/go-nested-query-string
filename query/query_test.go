package query

import (
	"reflect"
	"testing"
)

func TestStringToQuery(t *testing.T) {
	rawString := "a[aa]=11&a[ab]=12&c=3"

	query, err := StringToQuery(rawString)
	if err != nil {
		t.Errorf("stringToQuery failed: %s", err.Error())
	}
	expectedQuery := map[string]interface{}{
		"a": map[string]interface{}{
			"aa": "11",
			"ab": "12",
		},
		"c": "3",
	}

	if !reflect.DeepEqual(query, expectedQuery) {
		t.Errorf("Expected %s, but got %s", expectedQuery, query)
	}
}
