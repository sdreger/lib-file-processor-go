package main

import (
	"database/sql"
	"embed"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/sdreger/lib-file-processor-go/app"
	"github.com/sdreger/lib-file-processor-go/config"
	"github.com/sdreger/lib-file-processor-go/filestore"
	"io"
	"log"
	"os"
)

//go:embed db/migrations/*.sql
var embedMigrations embed.FS

var logger *log.Logger

func main() {
	// Get app config
	appConfig := config.GetAppConfig()

	// Init logger
	var logDestination io.Writer
	if appConfig.LogFilePath != "" {
		logFile, err := os.OpenFile(appConfig.LogFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Can not initiate a file logger, falling back to STDOUT")
			logDestination = os.Stdout
		} else {
			logDestination = logFile
			defer logFile.Close()
		}
	} else {
		logDestination = os.Stdout
	}
	logger = log.New(logDestination, "", log.Ldate|log.Ltime|log.Lshortfile)
	defer logger.Printf("[INFO] ---------- Shutting down ----------")

	// Try to connect to database
	db, err := connectToDB(appConfig.DBConnectionString)
	if err != nil {
		logger.Printf("[WARN] - Can not connect to DB: %v", err)
	} else {
		appConfig.DBAvailable = true
	}
	defer closeDBConnection(db)

	// Try to connect to BLOB store
	minioStore, minioStoreErr := filestore.NewMinioStore(appConfig.MinioEndpoint, appConfig.MinioAccessKeyID,
		appConfig.MinioSecretAccessKey, appConfig.MinioUseSSL, logger)
	if minioStoreErr != nil {
		logger.Printf("[WARN] - Can not connect to Minio store: %v", minioStoreErr)
	} else {
		appConfig.BlobStoreAvailable = true
	}

	// Try to init a file watcher
	bookExtractor := filestore.NewCompressionService(logger)
	watcher, err :=
		filestore.NewFileSystemWatcher(bookExtractor, appConfig.ZipInputFolder, appConfig.BookInputFolder, logger)
	if err != nil {
		logger.Printf("[WARN] - Can not initialize filesystem watcher: %v", err)
	}
	defer watcher.Close()
	if err := watcher.Watch(); err != nil {
		logger.Printf("[WARN] - Can not startt watching filesystem changes: %v", err)
	}

	// Init and run the TUI application
	tuiApp, tuiErr := app.NewTuiApp(appConfig, db, minioStore, logger, watcher.BookIDChan)
	if tuiErr != nil {
		logger.Fatalf("[FATAL] - TUI app init error: %v", tuiErr)
	}
	tuiErr = tuiApp.Run()
	if tuiErr != nil {
		logger.Fatalf("[FATAL] - TUI app error: %v", tuiErr)
	}
}

func connectToDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err = runDBMigration(db); err != nil {
		return db, err
	}

	return db, nil
}

func closeDBConnection(db *sql.DB) {
	if db != nil {
		err := db.Close()
		if err != nil {
			logger.Printf("[ERROR] - Can not close DB connection: %v", err)
		}
	}
}

func runDBMigration(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	goose.SetTableName("ebook.goose_db_version")
	if err := goose.Up(db, "db/migrations"); err != nil {
		return err
	}

	return nil
}
