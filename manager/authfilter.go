package manager

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

//URL for Rancher server
var URL = "http://54.255.182.226:8080/"

//Port for rancher auth
var Port = "8080"

//CacheProjectID the project id
var CacheProjectID = cache.New(60*time.Hour, 1*time.Hour)
