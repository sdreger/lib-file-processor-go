## Ebook Processor

### Description

The application is designed for getting and storing a book data. It can perform the following actions:

- scrape a book information by its ISBN10/ASIN identifier:
    - title;
    - subtitle;
    - description;
    - ISBN10;
    - ISBN13;
    - ASIN;
    - page count;
    - language;
    - publisher;
    - book URL;
    - edition;
    - publish date;
    - authors;
    - categories list;
    - book cover URL.
- process book files (compress / get metadata):
    - book file formats;
    - book archive name;
    - book archive size;
    - book cover file name.
- store the book information in the Postgres database:
    - book information;
    - book archive information;
    - book cover information;
    - book file formats information;
    - authors information;
    - categories information;
    - tags information;
- store the book files and book cover files into Minio BLOB store:
    - book file archive;
    - book cover file.

### Requirements

The application can run on macOS/Windows/Linux operating systems.
To run the TUI(default) version of the application, it's recommended to have at least 140x37 terminal.

For macOS/Linux OS the default terminal application is good enough, but for Windows it's recommended to
use [Windows Terminal](https://docs.microsoft.com/en-us/windows/terminal/install).

### Quick Start

- clone the repository: `git clone git@github.com:sdreger/lib-file-processor-go.git`;
- CD into the cloned folder: `cd lib-file-processor-go`;
- create folders inside the repository folder: `mkdir in_book in_zip in_temp out_book out_cover postgres14 minio`;
- start Docker containers using docker-compose: `docker-compose up -d`;
- wait for 5 seconds, after the previous command is finished (for the Postgres startup);
- run the application: `go run .`;
- you should see the application window with the following information in header:
    - DB: Available;
    - BLOB Store: Available;
- add two dummy files into the `in_book` folder: `touch in_book/dummy.pdf in_book/dummy.epub`. Or put the real book
  file(s) in the folder.
- enter a valid ISBN10/ASIN (for example: 1718502648) into the input field, and press `Enter`;
- after a couple of seconds you'll see the scrapped / parsed book information on the left side of the window;
- you can navigate between book info fields with `Tab`/`Shift-Tab` and make the necessary changes;
- then navigate back to the `Add` button and press `Enter` (the `Add` button is already in focus, if you've skipped the
  previous step);
- you should see the following message in the status bar: `The book is added successfully`;
- connect to the Postgres instance using credentials from the `docker-compose.yaml` file, and investigate
  the `sandbox.ebook` DB schema;
- navigate to the [Minio console](http://localhost:9001), login using credentials from the `docker-compose.yaml` file,
  and look at the newly created buckets (`ebook-covers` and `ebooks`) with corresponding content inside;
- the book archive and the book cover also stored at: `./out_book` and `./out_cover` folders, in case you don't want to use the BLOB store;
- some useful information could be found in the log file: `./lib_file_processor.log`.

### Application Configuration

| Environment <br/>Variable | Description                                 | Defaul Value                             |
|---------------------------|---------------------------------------------|------------------------------------------|
| DB_HOST                   | Postgres DB Hostname                        | 127.0.0.1:5432                           |
| DB_USER                   | Postgres DB User                            | postgres                                 |
| DB_PASSWORD               | Postgres DB Password                        | postgres                                 |
| DB_NAME                   | Postgres DB Name                            | sandbox                                  |
| DB_SCHEMA                 | Postgres DB Schema                          | ebook                                    |
| MINIO_ENDPOINT            | Minio Instance Endpoint                     | 127.0.0.1:9000                           |
| MINIO_ACCESS_KEY_ID       | Minio Access Key ID                         | AKIAIOSFODNN7EXAMPLE                     |
| MINIO_SECRET_ACCESS_KEY   | Minio Access Key Secret                     | wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY |
| MINIO_USE_SSL             | Use SSL for Minio connection                | false                                    |
| DIR_INPUT_TEMP            | Folder to store application temporary files | ./in_temp                                |
| DIR_INPUT_ZIP             | Folder to monitor for book zip archives     | ./in_zip                                 |
| DIR_INPUT_BOOK            | Book Files Input folder                     | ./in_book                                |
| DIR_OUTPUT_ARCHIVE        | Book Archive output folder                  | ./out_book                               |
| DIR_OUTPUT_COVER          | Book Cover output folder                    | ./out_cover                              |
| LOG_FILE_PATH             | Application log file path                   | ./lib_file_processor.log                 |

### Database Management

The application DB state is managed by [Goose](https://github.com/pressly/goose) DB migration tool. The migration files
could be found here: `./db/migrations/`.

### Work Modes
The application can work in several modes:
- Stateful mode, with book files compression, and storing book information to the DB and BLOB store (at least one file should be in the DIR_INPUT_BOOK folder).
If the DB is not available, the DB related operations will be omitted, if the BLOB is not available, the BLOB store operations will be omitted.
Regardless of the DB and BLOB store availability, the book archive will be created, and the book cover will be downloaded.
- Stateless mode, in this mode the application just gets the book information, shows it, and copy the book filename to the clipboard. 
The mode activates automatically if there are no files in the DIR_INPUT_BOOK folder. No file, DB or BLOB store operations will be performed in this mode.
- There is one more special mode. If you copy a _properly formatted_ filename to the DIR_INPUT_ZIP folder, 
the application extracts zip file content to the DIR_INPUT_BOOK folder, and puts the book identifier in the input field.
After pressing `Enter` application continue to work in the stateful mode.
The filename should end with book identifier and publication date, separated with dots. 
For example: `NSP.The.Book.of.Kubernetes.1718502648.Sep.2022.zip`.
