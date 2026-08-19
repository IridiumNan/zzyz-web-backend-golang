package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/IridiumNan/zzyz-web-backend-golang/internal/database"
	"github.com/IridiumNan/zzyz-web-backend-golang/internal/models"
	"github.com/IridiumNan/zzyz-web-backend-golang/internal/utils"
	"github.com/gin-gonic/gin"
)

var memberDB = database.NewMemberDB()

func memberListHandler(c *gin.Context) {
	mp := database.NewMemberDB()

	members, err := mp.ListMember()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewBadResponse(nil, err))
		return
	}

	if len(members) == 0 {
		c.JSON(http.StatusOK, models.NewBadResponse(nil, fmt.Errorf("there is no any member yet")))
		return
	}

	c.JSON(http.StatusOK, models.NewDataResponseWithMessage(members, "Password will no display, if you want to get passwd, update it pls"))
}

func memberQueryHandler(c *gin.Context) {
	// The name act as the attribute
	attribute := c.Param("name")

	val := c.Query("value")
	fmt.Println("value: ", val)

	type resp struct {
		Attr  string `json:"attr"`
		Value string `json:"value"`
	}

	c.JSON(200, models.NewDataResponse(resp{Attr: attribute, Value: val}))
}

// memberCreateHander : handle the requet for create new member
func memberCreateHander(c *gin.Context) {
	type createReq struct {
		Data map[string]database.Member `json:"data"`
	}
	var reqJSON createReq

	err := c.ShouldBindBodyWithJSON(&reqJSON)
	if err != nil {
		utils.TextLogger.Error("error when parse json from request", "err", err)
	}

	mb := reqJSON.Data["member"]

	err = memberDB.InsertMember(mb)
	if err != nil {
		// NOTE: This will happen when passwd hash failed
		c.JSON(http.StatusBadRequest, models.NewBadResponse(nil, err))
		return
	}

	c.JSON(http.StatusOK, models.NewDataResponse("add the task of create member"))
}

func memberUpdateHander(c *gin.Context) {
	// TODO: update member info
}

func memberDeleteHandler(c *gin.Context) {
	idStr := c.Param("name")
	isSoft := c.Query("soft")

	isSoftBool := false

	if isSoft == "True" {
		isSoftBool = true
	}

	fmt.Println("IS SOFT", isSoft)
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("error when convet id to int", "err", err, "idStr", idStr)
		return
	}

	err = memberDB.DeleteMember(idInt, isSoftBool)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewBadResponse(nil, fmt.Errorf("error when delete id : %d, err: %w", idInt, err)))
	}

	c.JSON(http.StatusOK, models.NewDataResponse(fmt.Sprintf("add the task for delete member with id: %d, soft: %v", idInt, isSoftBool)))
}
