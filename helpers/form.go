package helpers

import (
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func ParseQueryParams(r *http.Request, dst interface{}) error {
	v := reflect.ValueOf(dst).Elem()
	t := v.Type()

	query := r.URL.Query()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("query")
		if tag == "" {
			continue
		}

		value := query.Get(tag)
		// Check if we have multiple values for this parameter
		values := query[tag]

		if len(values) == 0 {
			if HandleDefaultTag(field, v.Field(i)) {
				continue
			}
			continue
		}

		switch field.Type.Kind() {
		case reflect.String:
			v.Field(i).SetString(value)
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.String {
				// Handle array notation parameters (tag[])
				bracketValues := query[tag+"[]"]
				if len(bracketValues) > 0 {
					v.Field(i).Set(reflect.ValueOf(bracketValues))
					continue
				}

				// Handle repeated parameters (tag=val1&tag=val2)
				if len(values) > 1 {
					v.Field(i).Set(reflect.ValueOf(values))
					continue
				}

				// Handle comma-separated values (tag=val1,val2,val3)
				if value != "" && (field.Tag.Get("split") == "comma" || field.Tag.Get("split") == "") {
					splitValues := strings.Split(value, ",")
					// Trim spaces from each value
					for j := range splitValues {
						splitValues[j] = strings.TrimSpace(splitValues[j])
					}
					v.Field(i).Set(reflect.ValueOf(splitValues))
					continue
				}

				// If only one value and no comma, set as single-item slice
				if value != "" {
					v.Field(i).Set(reflect.ValueOf([]string{value}))
				}
			}
		case reflect.Int:
			if intValue, err := strconv.Atoi(query.Get(tag)); err == nil {
				v.Field(i).SetInt(int64(intValue))
			}
		case reflect.Ptr:
			setupPointerType(v, i, query, tag)
		case reflect.Bool:
			if boolValue, err := strconv.ParseBool(query.Get(tag)); err == nil {
				v.Field(i).SetBool(boolValue)
			}
		case reflect.Struct:
			if field.Type == reflect.TypeOf(time.Time{}) {
				if tValue, err := time.Parse(time.RFC3339, value); err == nil {
					v.Field(i).Set(reflect.ValueOf(tValue))
				}
			}
			// Add other types as needed
		}
	}

	return nil
}

func setupPointerType(v reflect.Value, i int, query url.Values, tag string) {
	if v.Field(i).IsNil() {
		v.Field(i).Set(reflect.New(v.Field(i).Type().Elem()))
	}

	switch v.Field(i).Elem().Kind() {
	case reflect.Bool:
		if boolValue, err := strconv.ParseBool(query.Get(tag)); err == nil {
			v.Field(i).Elem().SetBool(boolValue)
		}
		// Add additional cases here as needed
	}
}

func HandleDefaultTag(field reflect.StructField, fieldValue reflect.Value) bool {
	defaultTag := field.Tag.Get("default")
	if defaultTag == "" {
		return false
	}

	switch fieldValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if defaultInt, err := strconv.Atoi(defaultTag); err == nil {
			fieldValue.SetInt(int64(defaultInt))
			return true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if defaultUint, err := strconv.ParseUint(defaultTag, 10, 64); err == nil {
			fieldValue.SetUint(defaultUint)
			return true
		}
	case reflect.Float32, reflect.Float64:
		if defaultFloat, err := strconv.ParseFloat(defaultTag, 64); err == nil {
			fieldValue.SetFloat(defaultFloat)
			return true
		}
	case reflect.Bool:
		if defaultBool, err := strconv.ParseBool(defaultTag); err == nil {
			fieldValue.SetBool(defaultBool)
			return true
		}
	case reflect.String:
		fieldValue.SetString(defaultTag)
		return true
	}

	return false
}
