package util

import (
	"html/template"
	"strings"
	"time"
)

// FuncMap returns the template.FuncMap used by the application.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		// Mathematical functions
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		// Dictionary/Map functions
		"dict": func() map[string]interface{} {
			return make(map[string]interface{})
		},
		"set": func(dict map[string]interface{}, key string, value interface{}) map[string]interface{} {
			dict[key] = value
			return dict
		},
		"get": func(dict map[string]interface{}, key string) interface{} {
			return dict[key]
		},

		// String functions
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
		"title": func(s string) string {
			return strings.Title(s)
		},

		// Utility functions
		"default": func(defaultVal, val interface{}) interface{} {
			if val == nil || val == "" {
				return defaultVal
			}
			return val
		},
		"seq": func(start, end int) []int {
			result := make([]int, 0, end-start+1)
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
		// Date/Time functions
		"now": func() time.Time {
			return time.Now()
		},
	}
}
