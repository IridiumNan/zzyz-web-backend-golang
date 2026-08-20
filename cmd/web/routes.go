package main

import "github.com/gin-gonic/gin"

func getMainRouter() (router *gin.Engine) {
	router = gin.Default()

	// TODO: write routes
	// for POST

	return
}

func getInternalRouter() (router *gin.Engine) {
	router = gin.Default()

	router.GET("/member/list", memberListHandler)

	// WARN: The /member/query/:name endpoint
	// :name is attr which is quered for, not id
	// like "id", "nick", "power" etc..
	router.GET("/member/query/:name", memberQueryHandler)

	router.POST("/member/create", memberCreateHander)
	router.PATCH("/member/update/:name/*action", memberUpdateHander)
	router.DELETE("/member/delete/:name", memberDeleteHandler)

	return
}
