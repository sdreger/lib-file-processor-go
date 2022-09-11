package config

type AppConfig struct {
	ZipInputFolder    string
	BookInputFolder   string
	BookOutputFolder  string
	CoverOutputFolder string
	TempInputFolder   string
	NewLineDelimiter  byte

	DBConnectionString string

	MinioEndpoint        string
	MinioAccessKeyID     string
	MinioSecretAccessKey string
	MinioUseSSL          bool

	DBAvailable        bool
	BlobStoreAvailable bool

	LogFilePath string
}

func (a AppConfig) IsStatelessMode() bool {
	return !a.DBAvailable || !a.BlobStoreAvailable
}
