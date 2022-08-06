package category

import "database/sql"

type storedData struct {
	ID       int64
	name     string
	parentID sql.NullInt64
}
