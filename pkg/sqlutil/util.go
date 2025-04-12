package sqlutil

import "fmt"

func TableColumn(table, column string) string {
	return fmt.Sprintf("%s.%s", table, column)
}
