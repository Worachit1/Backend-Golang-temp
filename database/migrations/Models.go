package migrations

import "app/app/model"

func Models() []any {
	return []any{
		(*model.User)(nil),
		// (*model.Emergency)(nil),
		// (*model.Officer)(nil),
	}
}

func RawBeforeQueryMigrate() []string {
	return []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
	}
}

func RawAfterQueryMigrate() []string {
	return []string{}
}
