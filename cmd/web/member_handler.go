package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/IridiumNan/zzyz-web-backend-golang/internal/database"
	"github.com/IridiumNan/zzyz-web-backend-golang/internal/models"
	"github.com/IridiumNan/zzyz-web-backend-golang/internal/utils"
	"github.com/gin-gonic/gin"
)

var (
	// memberDB is not real database pointor
	// It act as the task producer and push task to chan
	memberDB = database.NewMemberDB()

	// attrChecker for checking if specific attribute is validate on endpoint
	attrChecker = models.NewAttrChecker()
)

// parseIsDelete utils function for convert the string
// false or true into 0 or 1
// type string
func parseIsDelete(value string) (val string) {
	if value == "0" || value == "1" {
		return value
	}
	val = strings.ToLower(value)

	switch val {
	case "false":
		val = "0"
	case "true":
		val = "1"
	default:
		val = "0"
	}

	return
}

// memberListHandler Handle the request for endpoint -> /member/list
// list all list all members and return json response
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

// memberQueryHandler Handle the request for endpoint -> /member/query/:attr
// query with conditions and return member list with json format
func memberQueryHandler(c *gin.Context) {
	attribute := c.Param("attr")

	// NOTE: name can be (id, nick, email, power) attribute
	// Check it first
	if ok, hint := attrChecker.MemberCheck(models.MemberQuery, attribute); !ok {
		c.JSON(http.StatusBadRequest, models.NewBadResponse(hint, fmt.Errorf("attribute not found")))
		return
	}

	// Get the value of this attribute
	// Then query the database for this
	val := c.Query("value")

	// WARN: Convert to bool if attribute is is_delete
	if attribute == "is_delete" {
		val = parseIsDelete(val)
	}

	var members []database.Member
	var err error
	switch attribute {
	case "id", "is_delete", "power":
		members, err = memberDB.QeuryMember(attribute, val, false)
	default:
		members, err = memberDB.QeuryMember(attribute, val, true)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewBadResponse(members, err))
		return
	}

	c.JSON(http.StatusOK, models.NewDataResponse(members))
}

// memberCreateHander Handle the request for endpoint -> /member/create
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

// memberUpdateHander Handle the request for endpoint -> /member/update/:id/:attr
func memberUpdateHander(c *gin.Context) {
	idStr := c.Param("id")
	attrStr := c.Param("attr")

	val := c.Query("value")

	if attrStr == "is_delete" {
		val = parseIsDelete(val)
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadResponse(nil, fmt.Errorf("error when convert id into int")))
	}

	if ok, hint := attrChecker.MemberCheck(models.MemberUpdate, attrStr); !ok {
		c.JSON(http.StatusBadRequest, models.NewBadResponse(hint, fmt.Errorf("attribute not found")))
	}

	err = memberDB.UpdateMember(idInt, attrStr, val)
	if err != nil {
		utils.TextLogger.Error("error when update member infor", "err", err)
		c.JSON(http.StatusInternalServerError, models.NewBadResponse(nil, fmt.Errorf("error when update member: %w", err)))
		return
	}

	c.JSON(http.StatusOK, models.NewDataResponse(fmt.Sprintf("update the member with id : %d, %s", idInt, attrStr)))
}

// memberDeleteHandler Handle the request for endpoint -> /member/delete/:id
func memberDeleteHandler(c *gin.Context) {
	idStr := c.Param("id")
	isSoft := c.Query("soft")

	var isSoftBool bool

	switch isSoft {
	case "True":
		isSoftBool = true
	case "False":
		isSoftBool = false
	}

	if isSoft == "True" {
		isSoftBool = true
	}

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
