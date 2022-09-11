package config

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

const (
	defaultTempFolder        = "in_temp"
	defaultBookInputFolder   = "in_book"
	defaultBookZipFolder     = "in_zip"
	defaultBookOutputFolder  = "out_book"
	defaultCoverOutputFolder = "out_cover"

	defaultDBHost     = "127.0.0.1:5432"
	defaultDBUser     = "postgres"
	defaultDBPassword = "postgres"
	defaultDBName     = "sandbox"
	defaultDBSchema   = "ebook"

	defaultMinioEndpoint        = "127.0.0.1:9000"
	defaultMinioAccessKeyID     = "AKIAIOSFODNN7EXAMPLE"
	defaultMinioSecretAccessKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	defaultMinioUseSSL          = false

	defaultLogFilePath = "lib_file_processor.log"

	EnvVarKeyDBHost     = "DB_HOST"
	EnvVarKeyDBUser     = "DB_USER"
	EnvVarKeyDBPassword = "DB_PASSWORD"
	EnvVarKeyDBName     = "DB_NAME"
	EnvVarKeyDBSchema   = "DB_SCHEMA"

	EnvVarKeyMinioEndpoint        = "MINIO_ENDPOINT"
	EnvVarKeyMinioAccessKeyID     = "MINIO_ACCESS_KEY_ID"
	EnvVarKeyMinioSecretAccessKey = "MINIO_SECRET_ACCESS_KEY"
	EnvVarKeyMinioUseSSL          = "MINIO_USE_SSL"

	EnvVarDirInputTemp     = "DIR_INPUT_TEMP"
	EnvVarDirInputZip      = "DIR_INPUT_ZIP"
	EnvVarDirInputBook     = "DIR_INPUT_BOOK"
	EnvVarDirOutputArchive = "DIR_OUTPUT_ARCHIVE"
	EnvVarDirOutputCover   = "DIR_OUTPUT_COVER"

	EnvVarLogFilePath = "LOG_FILE_PATH"
)

func GetAppConfig() AppConfig {
	bookZipFolder := defaultBookZipFolder
	tempFolder := defaultTempFolder
	bookInputFolder := defaultBookInputFolder
	bookOutputFolder := defaultBookOutputFolder
	coverOutputFolder := defaultCoverOutputFolder
	if bookZipFolderVal, bookZipFolderValSet := os.LookupEnv(EnvVarDirInputZip); bookZipFolderValSet {
		bookZipFolder = bookZipFolderVal
	}
	if tempFolderVal, tempFolderValSet := os.LookupEnv(EnvVarDirInputTemp); tempFolderValSet {
		tempFolder = tempFolderVal
	}
	if bookInputFolderVal, bookInputFolderValSet := os.LookupEnv(EnvVarDirInputBook); bookInputFolderValSet {
		bookInputFolder = bookInputFolderVal
	}
	if bookOutputFolderVal, bookOutputFolderValSet := os.LookupEnv(EnvVarDirOutputArchive); bookOutputFolderValSet {
		bookOutputFolder = bookOutputFolderVal
	}
	if coverOutputFolderVal, coverOutputFolderValSet := os.LookupEnv(EnvVarDirOutputCover); coverOutputFolderValSet {
		coverOutputFolder = coverOutputFolderVal
	}

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

	logFilePath := defaultLogFilePath
	if logFilePathVal, logFilePathValSet := os.LookupEnv(EnvVarLogFilePath); logFilePathValSet {
		logFilePath = logFilePathVal
	}

	return AppConfig{
		ZipInputFolder:       bookZipFolder,
		BookInputFolder:      bookInputFolder,
		BookOutputFolder:     bookOutputFolder,
		CoverOutputFolder:    coverOutputFolder,
		TempInputFolder:      tempFolder,
		NewLineDelimiter:     getNewLineDelimiter(),
		DBConnectionString:   getDBConnectionString(DBHost, DBUser, DBPassword, DBName, DBSchema),
		MinioEndpoint:        MinioEndpoint,
		MinioAccessKeyID:     MinioAccessKeyID,
		MinioSecretAccessKey: MinioSecretAccessKey,
		MinioUseSSL:          MinioUseSSL,
		DBAvailable:          false,
		BlobStoreAvailable:   false,
		LogFilePath:          logFilePath,
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
