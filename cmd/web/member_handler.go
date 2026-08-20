package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/IridiumNan/zzyz-web-backend-golang/internal/database"
	"github.com/IridiumNan/zzyz-web-backend-golang/internal/models"
	"github.com/IridiumNan/zzyz-web-backend-golang/internal/utils"
	"github.com/gin-gonic/gin"
)

var memberDB = database.NewMemberDB()

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

// memberQueryHandler Handle the request for endpoint -> /member/query/:name
// query with conditions and return member list with json format
func memberQueryHandler(c *gin.Context) {
	attribute := c.Param("name")

	attrCandidates := []string{
		"id",
		"nick",
		"email",
		"power",
		"is_delete",
	}
	// NOTE: name can be (id, nick, email, power) attribute
	// Check it first

	// If attribute not in list
	// return badrequest with hint
	if !slices.Contains(attrCandidates, attribute) {
		type hintResp struct {
			Hint          string   `json:"hint"`
			AvailableAttr []string `json:"available_attr"`
		}
		hint := hintResp{
			Hint:          "the attribute must in available_attr list",
			AvailableAttr: attrCandidates,
		}
		c.JSON(http.StatusBadRequest, models.NewBadResponse(hint, fmt.Errorf("attribute not found")))

		return
	}

	// Get the value of this attribute
	// Then query the database for this
	val := c.Query("value")

	// WARN: Convert to bool if attribute is is_delete
	if attribute == "is_delete" && val != "0" && val != "1" {

		val = strings.ToLower(val)

		switch val {
		case "false":
			val = "0"
		case "true":
			val = "1"
		default:
			val = "0"
		}
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

func memberUpdateHander(c *gin.Context) {
	// TODO: update member info
}

// memberDeleteHandler Handle the request for endpoint -> /member/delete/:name
func memberDeleteHandler(c *gin.Context) {
	idStr := c.Param("name")
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
