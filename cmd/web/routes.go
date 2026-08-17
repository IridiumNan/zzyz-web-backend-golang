package main

import "github.com/gin-gonic/gin"

func getMainRouter() (router *gin.Engine) {
	router = gin.Default()

	// TODO: write routes

	return
}

func getInternalRouter() (router *gin.Engine) {
	router = gin.Default()

	router.GET("/member/query", memberQueryHandler)
	router.POST("/member/create", memberCreateHander)
	router.PATCH("/member/update", memberUpdateHander)
	router.DELETE("/member/delete", memberDeleteHandler)

	return
}
