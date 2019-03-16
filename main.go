package main

import (
	"github.com/PUMATeam/catapult/model"
	"github.com/gin-gonic/gin"
)

func main() {
	model.InitDB()
	gin := gin.Default()
	gin.Run()
}
