package utils

func GetType(value interface{}) string {
	switch value.(type) {
	case bool:
		return "Bool"
	case int, int8, int16, int32, int64:
		return "Int"
	case uint, uint8, uint16, uint32, uint64:
		return "Uint"
	case float32, float64:
		return "Float"
	case string:
		return "String"
	case []interface{}:
		return "Slice"
	case map[string]interface{}:
		return "Map"
	default:
		return "Unknown"
	}
}
