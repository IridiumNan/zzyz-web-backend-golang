package main

import "github.com/gin-gonic/gin"

func getMainRouter() (router *gin.Engine) {
	router = gin.Default()

	// TODO: write routes

	return
}

func getInternalRouter() (router *gin.Engine) {
	router = gin.Default()

	// TODO: Update this as RESTFUL API
	// /member/{id}/*action
	router.GET("/member/list", memberListHandler)
	router.GET("/member/query/:name", memberQueryHandler)
	router.POST("/member/create", memberCreateHander)
	router.PATCH("/member/update/:name/*action", memberUpdateHander)
	router.DELETE("/member/delete/:name", memberDeleteHandler)

	return
}
