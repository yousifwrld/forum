package forum

import "strings"

// will take categories, and format them to display
func joinAndTrim(categories []Category) string {
	names := make([]string, len(categories))
	for i, category := range categories {
		names[i] = category.Name
	}
	return strings.TrimSuffix(strings.Join(names, ", "), ", ")
}
