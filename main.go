package main

import (
	"log"
	"os"

	"github.com/869413421/wechatbot/bootstrap"
)

func main() {
	log.Printf("Awesome Yuer")

	//var storageJson string = "storage.json"
	storageJson := "storage.json"
	if FileIsExisted(storageJson) {
		os.Remove(storageJson)
	}

	bootstrap.Run()
}

func FileIsExisted(filename string) bool {
	existed := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		existed = false
	}
	return existed
}
