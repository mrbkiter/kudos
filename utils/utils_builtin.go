package utils

import "encoding/json"

//Bool return Bool Pointer
func Bool(b bool) *bool { return &b }

//Int return Int
func Int(n int) *int { return &n }

//Int64 return Int64 pointer
func Int64(n int64) *int64 { return &n }

//String return string pointer
func String(s string) *string { return &s }

func Float64(f float64) *float64 { return &f }
func MapToJson(objmap map[string]string) string {
	b, _ := json.Marshal(objmap)
	return string(b)
}
