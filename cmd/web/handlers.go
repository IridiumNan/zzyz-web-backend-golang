package main

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
)

func memberQueryHandler(c *gin.Context) {
	// TODO: return msg of the member include id power nick email is_delete
	// receive the nick, email or power
	// If no condition provide, return all member's base information
	c.String(200, "hello")
}

func memberCreateHander(c *gin.Context) {
	// TODO: create new member
	// insert it by calling InsertMember func
}

func memberUpdateHander(c *gin.Context) {
	// TODO: update member info
}

func memberDeleteHandler(c *gin.Context) {
	idStr := c.Param("id")

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("error when convet id to int", "err", err, "idStr", idStr)
		return
	}

	fmt.Println("get id -> ", idInt)

	type resp struct {
		Greet string `json:"greet"`
	}

	c.JSON(200, resp{Greet: "hello"})
}
