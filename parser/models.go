package parser

import "time"

type BookPublishMeta struct {
	Publisher string
	Edition   uint8
	PubDate   time.Time
}
