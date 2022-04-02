package migration

import "fmt"

const tableMigration string = "migration"

var QueryEnsureMigrationTable = fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		parameter VARCHAR(32) NOT NULL PRIMARY KEY,
		value INT NOT NULL
	);
	
	INSERT INTO %s(parameter, value) VALUES ('highest_migration', 0)
	ON CONFLICT DO NOTHING;
`, tableMigration, tableMigration)

var QueryGetHighestMigration = fmt.Sprintf(`
	SELECT MAX(value) FROM %s WHERE parameter = 'highest_migration'; 
`, tableMigration)

var QueryUpdateHighestMigration = fmt.Sprintf(`
	UPDATE %s SET value = $1 WHERE parameter = 'highest_migration';
`, tableMigration)
