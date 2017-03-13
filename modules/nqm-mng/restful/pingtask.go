package restful

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"

	commonNqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	commonGin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	commonNqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	"github.com/chyeh/cast"
)

func listPingtasks(c *gin.Context) {
	/**
	 * Set-up paging
	 */
	paging := commonGin.PagingByHeader(
		c,
		&commonModel.Paging{
			Size:     50,
			Position: 1,
		},
	)
	// :~)

	query := buildQueryForListPingtasks(c)
	pingtasks, resultPaging := commonNqmDb.ListPingtasks(query, *paging)

	commonGin.HeaderWithPaging(c, resultPaging)
	c.JSON(http.StatusOK, pingtasks)
}

//func listPingtasks(
//	p *struct {
//		Period             string        `mvc:"query[period]"`
//		Name               string        `mvc:"query[name]"`
//		Enable             string        `mvc:"query[enable]"`
//		Comment            string        `mvc:"query[comment]"`
//		NumOfEnabledAgents string        `mvc:"query[num_of_enabled_agents]"`
//		Paging             *model.Paging `mvc:"pageSize[50] pageOrderBy[period#asc:name#asc:enable#desc:comment#asc:num_of_enabled_agents#desc]"`
//	},
//) (*model.Paging, mvc.OutputBody) {
//	return p.Paging,
//		mvc.JsonOutputBody(
//			commonNqmDb.ListPingtasks(p.Name, p.Paging),
//		)
//}

func buildQueryForListPingtasks(c *gin.Context) *commonNqmModel.PingtaskQuery {
	query := &commonNqmModel.PingtaskQuery{}

	if v, ok := c.GetQuery("period"); ok {
		query.Period = v
	}
	if v, ok := c.GetQuery("name"); ok {
		query.Name = v
	}
	if v, ok := c.GetQuery("enable"); ok {
		query.Enable = v
	}
	if v, ok := c.GetQuery("comment"); ok {
		query.Comment = v
	}
	if v, ok := c.GetQuery("num_of_enabled_agents"); ok {
		query.NumOfEnabledAgents = v
	}

	return query
}

//func getPingtasksById(c *gin.Context) {
//pingtaskId, err := strconv.Atoi(c.Param("pingtask_id"))
//if err != nil {
//	commonGin.OutputJsonIfNotNil(c, nil)
//}

//var pingtaskIDStr string
//var pingtaskID int32
//if v := c.Param("pingtask_id"); v != "" {
//	pingtaskIDStr = v
//}
//if v, err := cast.ToInt32E(pingtaskIDStr); err == nil {
//	pingtaskID = v
//}

//pingtask := commonNqmDb.GetPingtaskById(pingtaskID)

//commonGin.OutputJsonIfNotNil(c, pingtask)
//c.JSON(http.StatusOK, "")
//}

func getPingtasksById(
	p *struct {
		PingtaskID int32 `mvc:"param[pingtask_id]"`
	},
) mvc.OutputBody {
	return mvc.JsonOutputOrNotFound(commonNqmDb.GetPingtaskById(p.PingtaskID))
}

//func addNewPingtask(c *gin.Context) {
//	fmt.Println(c.Value("period"))
//	c.JSON(http.StatusCreated, "fuck")
//}

func addNewPingtask(
	pm *commonNqmModel.PingtaskModify,
) mvc.OutputBody {
	pingtask := commonNqmDb.AddAndGetPingtask(pm)
	return mvc.JsonOutputBody2(http.StatusCreated, pingtask)
}

func modifyPingtask(
	p *struct {
		ID int32 `mvc:"param[pingtask_id]"`
	},
	pm *commonNqmModel.PingtaskModify,
) mvc.OutputBody {
	pingtask := commonNqmDb.UpdateAndGetPingtask(p.ID, pm)
	return mvc.JsonOutputBody(pingtask)
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
