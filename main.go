package main

import (
	"github.com/PUMATeam/catapult/model"
	"github.com/gin-gonic/gin"
)

func Migrate() {

}

func main() {
	model.Setup()
	gin := gin.Default()
	gin.Run()
}
