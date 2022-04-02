package migration

import (
	"github.com/Depal/quotebot/pkg/logger"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"path"
	"sort"
	"strconv"
	"strings"
)

const schemaFolder string = "schema"

func Apply(db *sqlx.DB, logger logger.ILogger) (err error) {
	err = ensureMigrationTable(db)
	if err != nil {
		return err
	}

	fileInfos, err := ioutil.ReadDir("schema")
	if err != nil {
		return err
	}

	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].Name() < fileInfos[j].Name()
	})

	var highestMigration int
	err = db.Get(&highestMigration, QueryGetHighestMigration)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		scriptOrderPart := strings.Split(fileInfo.Name(), "_")[0]
		scriptOrder, err := strconv.Atoi(scriptOrderPart)
		if err != nil {
			return err
		}
		if scriptOrder > highestMigration {
			err = executeScript(db, fileInfo.Name(), logger)
			if err != nil {
				return err
			}

			_, err = db.Exec(QueryUpdateHighestMigration, scriptOrder)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ensureMigrationTable(db *sqlx.DB) (err error) {
	_, err = db.Exec(QueryEnsureMigrationTable)
	return err
}

func executeScript(db *sqlx.DB, scriptName string, logger logger.ILogger) (err error) {
	bytes, err := ioutil.ReadFile(path.Join(schemaFolder, scriptName))
	if err != nil {
		return err
	}
	query := string(bytes)

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	logger.Info("Applied migration: " + scriptName)

	return nil
}
