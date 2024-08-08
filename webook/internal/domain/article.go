package domain

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
}

type Author struct {
	Id   int64
	Name string
}

const (
	ArticleStatusUnknown ArticleStatus = iota
	ArticleStatusUnpublish
	ArticleStatusPublish
	ArticleStatusPrivate
)

type ArticleStatus uint8

func (status ArticleStatus) ToUint8() uint8 {
	return uint8(status)
}

func (status ArticleStatus) NonPublish() bool {
	return status != ArticleStatusPublish
}

func (status ArticleStatus) String() string {
	switch status {
	case ArticleStatusUnpublish:
		return "Unpublish"
	case ArticleStatusPublish:
		return "Publish"
	case ArticleStatusPrivate:
		return "Private"
	default:
		return "Unknown"
	}
}
