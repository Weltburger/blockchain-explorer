package storage

type DataSources struct {
	Postgres   *PostgresDataSource
	Clickhouse *ClickhouseDataSource
}

func InitDataSources() (*DataSources, error) {
	postgres, err := InitPostgres()
	if err != nil {
		return nil, err
	}

	clickhouse, err := InitClickhouse()
	if err != nil {
		return nil, err
	}

	return &DataSources{
		Postgres:   postgres,
		Clickhouse: clickhouse,
	}, nil
}
