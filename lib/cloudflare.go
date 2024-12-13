package lib

import (
	"log"
)

func CloudFlare(APIKey string) {
	type CloudFlare struct {
		APIKey string
	}
	
	log.Printf("Connected to Cloudflare")
}
