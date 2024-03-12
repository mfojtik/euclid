#!/bin/bash

mkdir -p ~/go/bin

GOFLAGS="-mod=vendor" go build -o scrape ./scraper && mv scrape ~/go/bin/euclid-scraper
GOFLAGS="-mod=vendor" go build -o render ./renderer && mv render ~/go/bin/euclid-render