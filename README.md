# gopodder
Command line podcast client written in Go

### Development

go run main.go podcasts fetch https://feeds.npr.org/510366/podcast.xml
go run main.go podcasts add https://feeds.npr.org/510366/podcast.xml
go run main.go podcasts remove https://feeds.npr.org/510366/podcast.xml
go run main.go podcasts remove "NPR News Now"
go run main.go podcasts remove NPR*
go run main.go podcasts remove --all
go run main.go version