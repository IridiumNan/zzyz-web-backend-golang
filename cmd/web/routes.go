package main

import "github.com/gin-gonic/gin"

func getMainRouter() (router *gin.Engine) {
	router = gin.Default()

	// TODO: write routes
	// for posts

	return
}

func getInternalRouter() (router *gin.Engine) {
	router = gin.Default()

	router.GET("/member/list", memberListHandler)

	// WARN: The /member/query/:name endpoint
	// :name is attr which is quered for, not id
	// like "id", "nick", "power" etc..
	router.GET("/member/query/:attr", memberQueryHandler)

	router.POST("/member/create", memberCreateHander)
	router.PATCH("/member/update/:id/:attr", memberUpdateHander)
	router.DELETE("/member/delete/:id", memberDeleteHandler)

	return
}
