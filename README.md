# gopodder
Command line podcast client written in Go

### Development

go run main.go podcasts fetch
go run main.go podcasts fetch --start-time=2023-12-01
go run main.go podcasts fetch https://feeds.npr.org/500005/podcast.xml
go run main.go podcasts fetch https://feeds.npr.org/510366/podcast.xml --start-time=2023-12-08
go run main.go podcasts fetch https://feeds.npr.org/510366/podcast.xml --start-time=2023-11-01 --end-time=2023-11-30

go run main.go podcasts fetch https://rss.wbur.org/onpoint/podcast.xml --start-time=2023-12-01
go run main.go podcasts add https://rss.wbur.org/onpoint/podcast.xml
go run main.go podcasts remove https://feeds.npr.org/510366/podcast.xml
go run main.go podcasts remove "NPR News Now"
go run main.go podcasts remove NPR*
go run main.go podcasts remove --all
go run main.go version