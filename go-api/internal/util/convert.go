package util

func ToFloat(v interface{}) float64 {
    switch val := v.(type) {
    case float64:
        return val
    case float32:
        return float64(val)
    case int:
        return float64(val)
    default:
        return 0.0
    }
}