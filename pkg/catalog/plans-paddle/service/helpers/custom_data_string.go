package helpers

func CustomDataString(data map[string]any, key string) string {
	if data == nil {
		return ""
	}
	value, ok := data[key]
	if !ok {
		return ""
	}
	stringValue, _ := value.(string)
	return stringValue
}
