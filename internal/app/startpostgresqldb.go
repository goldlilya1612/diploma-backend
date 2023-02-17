package app

import "github.com/goldlilya1612/diploma-backend/internal/database"

func StartPostgreSQL(config *database.Config) (*database.PostgreSQL, error) {

	psql := database.NewPostgreSQL(config)
	err := psql.Connect()
	if err != nil {
		return nil, err
	}

	err = psql.InitDB()
	if err != nil {
		return nil, err
	}

	return psql, nil
}
