package sqlmock

import (
	"regexp"
	"sort"
	"strings"
)

// 删除多余的空白和标点符号
func normalizeSQL(sql string) string {
	// 转换为小写
	sql = strings.ToLower(sql)

	// 去除多余空白
	sql = regexp.MustCompile(`\s+`).ReplaceAllString(sql, " ")

	// 去除括号内的空白
	sql = regexp.MustCompile(`\(\s+`).ReplaceAllString(sql, "(")
	sql = regexp.MustCompile(`\s+\)`).ReplaceAllString(sql, ")")

	// 去除逗号周围的空白
	sql = regexp.MustCompile(`\s*,\s*`).ReplaceAllString(sql, ",")

	// 去除行尾分号
	sql = strings.TrimSuffix(sql, ";")

	return strings.TrimSpace(sql)
}

// 提取 SELECT 语句的关键组件
func extractSelectComponents(sql string) map[string]string {
	components := make(map[string]string)

	// 使用正则表达式提取各个部分
	selectRegex := regexp.MustCompile(`select\s+(.*?)\s+from\s+(.*?)(?:\s+where\s+(.*?))?(?:\s+order\s+by\s+(.*?))?(?:\s+limit\s+(.*?))?$`)
	matches := selectRegex.FindStringSubmatch(sql)

	if len(matches) > 1 {
		// 提取列
		columns := extractAndSortColumns(matches[1])
		components["columns"] = strings.Join(columns, ",")

		// 提取表
		components["from"] = normalizeSQL(matches[2])

		// 提取 WHERE 条件（如果存在）
		if len(matches) > 3 && matches[3] != "" {
			components["where"] = normalizeWhereClause(matches[3])
		}

		// 提取 ORDER BY（如果存在）
		if len(matches) > 4 && matches[4] != "" {
			components["order_by"] = normalizeOrderBy(matches[4])
		}

		// 提取 LIMIT（如果存在）
		if len(matches) > 5 && matches[5] != "" {
			components["limit"] = matches[5]
		}
	}

	return components
}

// 提取并排序列名
func extractAndSortColumns(columnsStr string) []string {
	// 处理 * 情况
	if strings.TrimSpace(columnsStr) == "*" {
		return []string{"*"}
	}

	// 分割列名
	columns := strings.Split(columnsStr, ",")

	// 去除别名和空白
	cleanColumns := make([]string, 0)
	for _, col := range columns {
		col = strings.TrimSpace(col)

		// 去除可能的别名
		if strings.Contains(col, " as ") {
			col = strings.Split(col, " as ")[0]
		} else if strings.Contains(col, " ") {
			col = strings.Split(col, " ")[0]
		}

		cleanColumns = append(cleanColumns, col)
	}

	// 排序
	sort.Strings(cleanColumns)
	return cleanColumns
}

// 规范化 WHERE 子句
func normalizeWhereClause(whereClause string) string {
	// 去除多余空白
	whereClause = normalizeSQL(whereClause)

	// 分割条件
	conditions := strings.Split(whereClause, " and ")

	// 排序条件
	sort.Strings(conditions)

	return strings.Join(conditions, " and ")
}

// 规范化 ORDER BY 子句
func normalizeOrderBy(orderBy string) string {
	// 去除多余空白
	orderBy = normalizeSQL(orderBy)

	// 分割排序条件
	sorts := strings.Split(orderBy, ",")

	// 排序并去除 ASC/DESC
	cleanSorts := make([]string, 0)
	for _, sort := range sorts {
		// 去除 ASC/DESC
		sort = strings.TrimSuffix(sort, " asc")
		sort = strings.TrimSuffix(sort, " desc")
		cleanSorts = append(cleanSorts, strings.TrimSpace(sort))
	}

	// 排序
	sort.Strings(cleanSorts)
	return strings.Join(cleanSorts, ",")
}

// 比较两个 SQL 语句是否语义相等
func CompareSQL(sql1, sql2 string) bool {
	// 标准化 SQL 语句
	normSql1 := normalizeSQL(sql1)
	normSql2 := normalizeSQL(sql2)

	// 如果完全相同，直接返回 true
	if normSql1 == normSql2 {
		return true
	}

	// 判断 SQL 类型
	switch {
	case strings.HasPrefix(normSql1, "select") && strings.HasPrefix(normSql2, "select"):
		return compareSelectStatements(normSql1, normSql2)
	// 可以根据需要添加其他类型的 SQL 语句比较
	default:
		return false
	}
}

// 比较 SELECT 语句
func compareSelectStatements(sql1, sql2 string) bool {
	// 提取组件
	comp1 := extractSelectComponents(sql1)
	comp2 := extractSelectComponents(sql2)

	// 比较关键组件
	return comp1["columns"] == comp2["columns"] &&
		comp1["from"] == comp2["from"] &&
		comp1["where"] == comp2["where"] &&
		comp1["order_by"] == comp2["order_by"] &&
		comp1["limit"] == comp2["limit"]
}
