package config

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
)

const (
	defaultTempFolder        = "temp"
	defaultBookInputFolder   = "in_book"
	defaultBookZipFolder     = "in_zip"
	defaultBookOutputFolder  = "out_book"
	defaultCoverOutputFolder = "out_cover"

	defaultDBHost     = "127.0.0.1"
	defaultDBUser     = "postgres"
	defaultDBPassword = "postgres"
	defaultDBName     = "sandbox"
	defaultDBSchema   = "ebook"

	defaultMinioEndpoint        = "localhost:9000"
	defaultMinioAccessKeyID     = "AKIAIOSFODNN7EXAMPLE"
	defaultMinioSecretAccessKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	defaultMinioUseSSL          = false

	EnvVarKeyDBHost     = "DB_HOST"
	EnvVarKeyDBUser     = "DB_USER"
	EnvVarKeyDBPassword = "DB_PASSWORD"
	EnvVarKeyDBName     = "DB_NAME"
	EnvVarKeyDBSchema   = "DB_SCHEMA"

	EnvVarKeyMinioEndpoint        = "MINIO_ENDPOINT"
	EnvVarKeyMinioAccessKeyID     = "MINIO_ACCESS_KEY_ID"
	EnvVarKeyMinioSecretAccessKey = "MINIO_SECRET_ACCESS_KEY"
	EnvVarKeyMinioUseSSL          = "MINIO_USE_SSL"
)

func GetAppConfig() AppConfig {
	var bookZipFolder, bookInputFolder, bookOutputFolder, coverOutputFolder, tempFolder string
	flag.StringVar(&tempFolder, "temp", defaultTempFolder, "temp folder for intermediate files")
	flag.StringVar(&bookZipFolder, "input-zip", defaultBookZipFolder, "input folder with a book zip file")
	flag.StringVar(&bookInputFolder, "input-book", defaultBookInputFolder, "input folder with a book files")
	flag.StringVar(&bookOutputFolder, "output-archive", defaultBookOutputFolder, "output folder for a book archive")
	flag.StringVar(&coverOutputFolder, "output-cover", defaultCoverOutputFolder, "output folder for a book cover")
	flag.Parse()

	DBHost := defaultDBHost
	DBUser := defaultDBUser
	DBPassword := defaultDBPassword
	DBName := defaultDBName
	DBSchema := defaultDBSchema
	if DBHostVal, DBHostValSet := os.LookupEnv(EnvVarKeyDBHost); DBHostValSet {
		DBHost = DBHostVal
	}
	if DBUserVal, DBUserValSet := os.LookupEnv(EnvVarKeyDBUser); DBUserValSet {
		DBUser = DBUserVal
	}
	if DBPasswordVal, DBPasswordValSet := os.LookupEnv(EnvVarKeyDBPassword); DBPasswordValSet {
		DBPassword = DBPasswordVal
	}
	if DBNameVal, DBNameValSet := os.LookupEnv(EnvVarKeyDBName); DBNameValSet {
		DBName = DBNameVal
	}
	if DBSchemaVal, DBSchemaValSet := os.LookupEnv(EnvVarKeyDBSchema); DBSchemaValSet {
		DBSchema = DBSchemaVal
	}

	MinioEndpoint := defaultMinioEndpoint
	MinioAccessKeyID := defaultMinioAccessKeyID
	MinioSecretAccessKey := defaultMinioSecretAccessKey
	MinioUseSSL := defaultMinioUseSSL
	if MinioEndpointVal, MinioEndpointValSet := os.LookupEnv(EnvVarKeyMinioEndpoint); MinioEndpointValSet {
		MinioEndpoint = MinioEndpointVal
	}
	if MinioAccessKeyIDVal, MinioAccessKeyIDValSet := os.LookupEnv(EnvVarKeyMinioAccessKeyID); MinioAccessKeyIDValSet {
		MinioAccessKeyID = MinioAccessKeyIDVal
	}
	if MinioSecretAccessKeyVal, MinioSecretAccessKeyValSet :=
		os.LookupEnv(EnvVarKeyMinioSecretAccessKey); MinioSecretAccessKeyValSet {
		MinioSecretAccessKey = MinioSecretAccessKeyVal
	}
	if MinioUseSSLVal, MinioUseSSLValSet := os.LookupEnv(EnvVarKeyMinioUseSSL); MinioUseSSLValSet {
		if useSSL, err := strconv.ParseBool(MinioUseSSLVal); err != nil {
			MinioUseSSL = useSSL
		}
	}

	return AppConfig{
		ZipInputFolder:       bookZipFolder,
		BookInputFolder:      bookInputFolder,
		BookOutputFolder:     bookOutputFolder,
		CoverOutputFolder:    coverOutputFolder,
		TempFolder:           tempFolder,
		NewLineDelimiter:     getNewLineDelimiter(),
		DBConnectionString:   getDBConnectionString(DBHost, DBUser, DBPassword, DBName, DBSchema),
		MinioEndpoint:        MinioEndpoint,
		MinioAccessKeyID:     MinioAccessKeyID,
		MinioSecretAccessKey: MinioSecretAccessKey,
		MinioUseSSL:          MinioUseSSL,
		DBAvailable:          false,
		BlobStoreAvailable:   false,
	}
}

func getDBConnectionString(DBHost, DBUser, DBPassword, DBName, DBSchema string) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable&search_path=%s",
		DBUser, DBPassword, DBHost, DBName, DBSchema)
}

func getNewLineDelimiter() byte {
	var delimiter byte = '\n'
	if runtime.GOOS == "windows" {
		delimiter = '\r'
	}

	return delimiter
}
