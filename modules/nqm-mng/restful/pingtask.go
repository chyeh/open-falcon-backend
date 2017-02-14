package restful

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"

	commonNqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	commonGin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/chyeh/cast"
)

func listPingtasks(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}

func getPingtasksById(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}

func addNewPingtask(c *gin.Context) {
	c.JSON(http.StatusCreated, "")
}

func modifyPingtask(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}

func addPingtaskToAgentForAgent(c *gin.Context) {
	/**
	 * Builds data from body of request
	 */
	var pingtaskIDStr string
	var pingtaskID int32

	if v, ok := c.GetQuery("pingtask_id"); ok {
		pingtaskIDStr = v
	}
	if v, err := cast.ToInt32E(pingtaskIDStr); err == nil {
		pingtaskID = v
	}

	var agentIDStr string
	var agentID int32
	if v := c.Param("agent_id"); v != "" {
		agentIDStr = v
	}
	if v, err := cast.ToInt32E(agentIDStr); err == nil {
		agentID = v
	}

	/* 待做: 驗證 input 範圍*/

	agentWithNewPingtask, err := commonNqmDb.AssignPingtaskToAgentForAgent(agentID, pingtaskID)
	if err != nil {
		switch err.(type) {
		case commonNqmDb.ErrDuplicatedNqmAgent:
			commonGin.JsonConflictHandler(
				c,
				commonGin.DataConflictError{
					ErrorCode:    1,
					ErrorMessage: err.Error(),
				},
			)
		default:
			panic(err)
		}

		return
	}

	c.JSON(http.StatusCreated, agentWithNewPingtask)
}

func removePingtaskFromAgentForAgent(c *gin.Context) {
	var agentIDStr string
	var agentID int32
	if v := c.Param("agent_id"); v != "" {
		agentIDStr = v
	}
	if v, err := cast.ToInt32E(agentIDStr); err == nil {
		agentID = v
	}

	var pingtaskIDStr string
	var pingtaskID int32
	if v := c.Param("pingtask_id"); v != "" {
		pingtaskIDStr = v
	}
	if v, err := cast.ToInt32E(pingtaskIDStr); err == nil {
		pingtaskID = v
	}

	/* 待做: 驗證 input 範圍*/

	agentWithRemovedPingtask, err := commonNqmDb.RemovePingtaskFromAgentForAgent(agentID, pingtaskID)
	if err != nil {
		switch err.(type) {
		case commonNqmDb.ErrDuplicatedNqmAgent:
			commonGin.JsonConflictHandler(
				c,
				commonGin.DataConflictError{
					ErrorCode:    1,
					ErrorMessage: err.Error(),
				},
			)
		default:
			panic(err)
		}

		return
	}
	c.JSON(http.StatusOK, agentWithRemovedPingtask)
}

func listTargetsOfAgent(c *gin.Context) {
	c.JSON(http.StatusCreated, "fuck you")
}

func addPingtaskToAgentForPingtask(c *gin.Context) {
	/**
	 * Builds data from body of request
	 */
	var pingtaskIDStr string
	var pingtaskID int32

	if v := c.Param("pingtask_id"); v != "" {
		pingtaskIDStr = v
	}
	if v, err := cast.ToInt32E(pingtaskIDStr); err == nil {
		pingtaskID = v
	}

	var agentIDStr string
	var agentID int32
	if v, ok := c.GetQuery("agent_id"); ok {
		agentIDStr = v
	}
	if v, err := cast.ToInt32E(agentIDStr); err == nil {
		agentID = v
	}

	/* 待做: 驗證 input 範圍*/

	agentWithNewPingtask, err := commonNqmDb.AssignPingtaskToAgentForPingtask(agentID, pingtaskID)
	if err != nil {
		switch err.(type) {
		case commonNqmDb.ErrDuplicatedNqmAgent:
			commonGin.JsonConflictHandler(
				c,
				commonGin.DataConflictError{
					ErrorCode:    1,
					ErrorMessage: err.Error(),
				},
			)
		default:
			panic(err)
		}

		return
	}

	c.JSON(http.StatusCreated, agentWithNewPingtask)
}

func removePingtaskFromAgentForPingtask(c *gin.Context) {
	var agentIDStr string
	var agentID int32
	if v := c.Param("agent_id"); v != "" {
		agentIDStr = v
	}
	if v, err := cast.ToInt32E(agentIDStr); err == nil {
		agentID = v
	}

	var pingtaskIDStr string
	var pingtaskID int32
	if v := c.Param("pingtask_id"); v != "" {
		pingtaskIDStr = v
	}
	if v, err := cast.ToInt32E(pingtaskIDStr); err == nil {
		pingtaskID = v
	}

	/* 待做: 驗證 input 範圍*/

	agentWithRemovedPingtask, err := commonNqmDb.RemovePingtaskFromAgentForPingtask(agentID, pingtaskID)
	if err != nil {
		switch err.(type) {
		case commonNqmDb.ErrDuplicatedNqmAgent:
			commonGin.JsonConflictHandler(
				c,
				commonGin.DataConflictError{
					ErrorCode:    1,
					ErrorMessage: err.Error(),
				},
			)
		default:
			panic(err)
		}

		return
	}
	c.JSON(http.StatusOK, agentWithRemovedPingtask)
}
