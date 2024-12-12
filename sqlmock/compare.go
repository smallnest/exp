package sqlmock

import (
	"regexp"
	"sort"
	"strings"
)

// Remove extra whitespace and punctuation
func normalizeSQL(sql string) string {
	// Convert to lowercase
	sql = strings.ToLower(sql)

	// Remove extra whitespace
	sql = regexp.MustCompile(`\s+`).ReplaceAllString(sql, " ")

	// Remove whitespace inside parentheses
	sql = regexp.MustCompile(`\(\s+`).ReplaceAllString(sql, "(")
	sql = regexp.MustCompile(`\s+\)`).ReplaceAllString(sql, ")")

	// Remove whitespace around commas
	sql = regexp.MustCompile(`\s*,\s*`).ReplaceAllString(sql, ",")

	// Remove trailing semicolon
	sql = strings.TrimSuffix(sql, ";")

	return strings.TrimSpace(sql)
}

// Extract key components of a SELECT statement
func extractSelectComponents(sql string) map[string]string {
	components := make(map[string]string)

	// Use regular expressions to extract parts
	selectRegex := regexp.MustCompile(`select\s+(.*?)\s+from\s+(.*?)(?:\s+where\s+(.*?))?(?:\s+order\s+by\s+(.*?))?(?:\s+limit\s+(.*?))?$`)
	matches := selectRegex.FindStringSubmatch(sql)

	if len(matches) > 1 {
		// Extract columns
		columns := extractAndSortColumns(matches[1])
		components["columns"] = strings.Join(columns, ",")

		// Extract table
		components["from"] = normalizeSQL(matches[2])

		// Extract WHERE clause (if exists)
		if len(matches) > 3 && matches[3] != "" {
			components["where"] = normalizeWhereClause(matches[3])
		}

		// Extract ORDER BY (if exists)
		if len(matches) > 4 && matches[4] != "" {
			components["order_by"] = normalizeOrderBy(matches[4])
		}

		// Extract LIMIT (if exists)
		if len(matches) > 5 && matches[5] != "" {
			components["limit"] = matches[5]
		}
	}

	return components
}

// Extract and sort column names
func extractAndSortColumns(columnsStr string) []string {
	// Handle * case
	if strings.TrimSpace(columnsStr) == "*" {
		return []string{"*"}
	}

	// Split column names
	columns := strings.Split(columnsStr, ",")

	// Remove aliases and whitespace
	cleanColumns := make([]string, 0)
	for _, col := range columns {
		col = strings.TrimSpace(col)

		// Remove possible aliases
		if strings.Contains(col, " as ") {
			col = strings.Split(col, " as ")[0]
		} else if strings.Contains(col, " ") {
			col = strings.Split(col, " ")[0]
		}

		cleanColumns = append(cleanColumns, col)
	}

	// Sort
	sort.Strings(cleanColumns)
	return cleanColumns
}

// Normalize WHERE clause
func normalizeWhereClause(whereClause string) string {
	// Remove extra whitespace
	whereClause = normalizeSQL(whereClause)

	// Split conditions
	conditions := strings.Split(whereClause, " and ")

	// Sort conditions
	sort.Strings(conditions)

	return strings.Join(conditions, " and ")
}

// Normalize ORDER BY clause
func normalizeOrderBy(orderBy string) string {
	// Remove extra whitespace
	orderBy = normalizeSQL(orderBy)

	// Split sorting conditions
	sorts := strings.Split(orderBy, ",")

	// Sort and remove ASC/DESC
	cleanSorts := make([]string, 0)
	for _, sort := range sorts {
		// Remove ASC/DESC
		sort = strings.TrimSuffix(sort, " asc")
		sort = strings.TrimSuffix(sort, " desc")
		cleanSorts = append(cleanSorts, strings.TrimSpace(sort))
	}

	// Sort
	sort.Strings(cleanSorts)
	return strings.Join(cleanSorts, ",")
}

// Compare if two SQL statements are semantically equal
func CompareSQL(sql1, sql2 string) bool {
	// Normalize SQL statements
	normSql1 := normalizeSQL(sql1)
	normSql2 := normalizeSQL(sql2)

	// If exactly the same, return true
	if normSql1 == normSql2 {
		return true
	}

	// Determine SQL type
	switch {
	case strings.HasPrefix(normSql1, "select") && strings.HasPrefix(normSql2, "select"):
		return compareSelectStatements(normSql1, normSql2)
	// Add other types of SQL statement comparisons as needed
	default:
		return false
	}
}

// Compare SELECT statements
func compareSelectStatements(sql1, sql2 string) bool {
	// Extract components
	comp1 := extractSelectComponents(sql1)
	comp2 := extractSelectComponents(sql2)

	// Compare key components
	return comp1["columns"] == comp2["columns"] &&
		comp1["from"] == comp2["from"] &&
		comp1["where"] == comp2["where"] &&
		comp1["order_by"] == comp2["order_by"] &&
		comp1["limit"] == comp2["limit"]
}
