package filestore

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"runtime"
	"strings"
	"time"
)

type FileSystemWatcher struct {
	Watcher       *fsnotify.Watcher
	Extractor     BookExtractor
	BookIDChan    chan string
	PathToWatch   string
	PathToExtract string
	timer         *time.Timer
	logger        *log.Logger
}

func NewFileSystemWatcher(bookExtractor BookExtractor, pathToWatch, pathToExtract string,
	logger *log.Logger) (*FileSystemWatcher, error) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &FileSystemWatcher{
		Watcher:       watcher,
		Extractor:     bookExtractor,
		BookIDChan:    make(chan string),
		PathToWatch:   pathToWatch,
		PathToExtract: pathToExtract,
		logger:        logger,
	}, nil
}

func (w *FileSystemWatcher) Watch() error {
	err := w.Watcher.Add(w.PathToWatch)
	if err != nil {
		return err
	}
	go w.watchFileSystem()

	return nil
}

func (w *FileSystemWatcher) watchFileSystem() {
	errorChan := make(chan error)
	for {
		select {
		case event, ok := <-w.Watcher.Events:
			if !ok {
				return
			}
			var handlerError error
			switch runtime.GOOS {
			case "darwin":
				handlerError = w.handleMacosChanges(event)
			case "windows":
				w.handleWindowsChanges(event, errorChan)
			}
			if handlerError != nil {
				w.logger.Printf("[ERROR] - Error sync file handling: %v", handlerError)
			}
		case err, ok := <-w.Watcher.Errors:
			if !ok {
				return
			}
			w.logger.Println("[ERROR] - ", err)
		case err := <-errorChan:
			w.logger.Printf("[ERROR] - Error async file handling: %v", err)
		}
	}
}

func (w *FileSystemWatcher) handleWindowsChanges(event fsnotify.Event, errChan chan error) {
	if event.Op&fsnotify.Create == fsnotify.Create {
		w.timer = time.NewTimer(5 * time.Second)
		go func() {
			select {
			case <-w.timer.C:
				err := w.handleNewFile(event.Name)
				if err != nil {
					errChan <- err
				}
			}
		}()
	}
	if event.Op&fsnotify.Write == fsnotify.Write {
		w.timer.Reset(2 * time.Second)
	}
}

func (w *FileSystemWatcher) handleMacosChanges(event fsnotify.Event) error {
	if event.Op&fsnotify.Chmod == fsnotify.Chmod {
		return w.handleNewFile(event.Name)
	}

	return nil
}

func (w *FileSystemWatcher) handleNewFile(fileName string) error {
	err := w.Extractor.ExtractZipFile(fileName, w.PathToExtract)
	if err != nil {
		return fmt.Errorf("can not extract %q file: %w", fileName, err)
	}
	// Cut suffix like '.May.2022.zip'
	name := fileName[:len(fileName)-13]
	// Get ISBN/ASIN after the last dot
	bookID := name[strings.LastIndex(name, ".")+1:]
	w.BookIDChan <- bookID

	return nil
}

func (w *FileSystemWatcher) Close() error {
	close(w.BookIDChan)
	return w.Watcher.Close()
}
