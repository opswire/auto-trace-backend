package repo

import "fmt"

const (
	adTableName            = "ads"
	carTableName           = "cars"
	userFavoritesTableName = "user_favorites"
)

func tableColumn(table, column string) string {
	return fmt.Sprintf("%s.%s", table, column)
}
