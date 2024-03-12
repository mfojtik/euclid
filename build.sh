#!/bin/bash

mkdir -p ~/go/bin

go build -o scrape ./scraper && mv scrape ~/go/bin/euclid-scraper
go build -o render ./renderer && mv render ~/go/bin/euclid-render