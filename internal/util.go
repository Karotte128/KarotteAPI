package internal

import "reflect"

func GetNestedValue[Type any](m map[string]any, path ...string) (Type, bool) {
	var current any = m
	var zero Type

	for _, key := range path {
		next, ok := current.(map[string]any)
		if !ok {
			return zero, false
		}

		current, ok = next[key]
		if !ok {
			return zero, false
		}
	}

	targetType := reflect.TypeOf((*Type)(nil)).Elem()
	value := reflect.ValueOf(current)

	// Handle invalid (nil interface)
	if !value.IsValid() {
		return zero, false
	}

	// Exact or assignable match
	if value.Type().AssignableTo(targetType) {
		return value.Interface().(Type), true
	}

	// Handle slice conversion: []any -> []TElement
	if targetType.Kind() == reflect.Slice && value.Kind() == reflect.Slice {
		elemType := targetType.Elem()
		result := reflect.MakeSlice(targetType, value.Len(), value.Len())

		for i := 0; i < value.Len(); i++ {
			elem := value.Index(i)

			// Unwrap interface values safely
			if elem.Kind() == reflect.Interface {
				if elem.IsNil() {
					return zero, false
				}
				elem = elem.Elem()
			}

			if !elem.IsValid() {
				return zero, false
			}

			if !elem.Type().AssignableTo(elemType) {
				return zero, false
			}

			result.Index(i).Set(elem)
		}

		return result.Interface().(Type), true
	}

	return zero, false
}
