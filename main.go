package main

import (
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/sdreger/lib-file-processor-go/app"
	"github.com/sdreger/lib-file-processor-go/config"
	"github.com/sdreger/lib-file-processor-go/filestore"
	"log"
)

//go:embed db/migrations/*.sql
var embedMigrations embed.FS

func main() {
	appConfig := config.GetAppConfig()

	db, err := connectToDB(appConfig.DBConnectionString)
	if err != nil {
		log.Printf("[WARN] - Can not connect to DB: %v", err)
	} else {
		appConfig.DBAvailable = true
	}
	defer closeDBConnection(db)

	minioStore, minioStoreErr := filestore.NewMinioStore(appConfig.MinioEndpoint, appConfig.MinioAccessKeyID,
		appConfig.MinioSecretAccessKey, appConfig.MinioUseSSL)
	if minioStoreErr != nil {
		log.Printf("[WARN] - Can not get Minio store: %v", err)
	} else {
		appConfig.BlobStoreAvailable = true
	}

	bookExtractor := filestore.NewCompressionService()
	watcher, err := filestore.NewFileSystemWatcher(bookExtractor, appConfig.ZipInputFolder, appConfig.BookInputFolder)
	if err != nil {
		log.Printf("[WARN] - Can not initialize filesystem watcher: %v", err)
	}
	defer watcher.Close()
	if err := watcher.Watch(); err != nil {
		log.Printf("[WARN] - Can not startt watching filesystem changes: %v", err)
	}

	// Uncomment the section below to use the app in CLI mode instead of TUI
	/*
		cliApp, cliErr := app.NewCliApp(appConfig, db, minioStore)
		if cliErr != nil {
			log.Fatalf("CLI app init error: %v", cliErr)
		}
		cliApp.Run()
	*/

	tuiApp, tuiErr := app.NewTuiApp(appConfig, db, minioStore, watcher.BookIDChan)
	if tuiErr != nil {
		log.Fatalf("TUI app init error: %v", tuiErr)
	}
	tuiErr = tuiApp.Run()
	if tuiErr != nil {
		log.Fatalf("TUI app error: %v", tuiErr)
	}
}

func connectToDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	fmt.Println(err)

	if err = runDBMigration(db); err != nil {
		return db, err
	}

	return db, nil
}

func closeDBConnection(db *sql.DB) {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Printf("[ERROR] - Can not close DB connection: %v", err)
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
