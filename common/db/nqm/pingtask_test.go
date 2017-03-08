package nqm

import (
	"encoding/json"
	"fmt"

	owlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"
	. "gopkg.in/check.v1"
)

type TestPingtaskSuite struct{}

var _ = Suite(&TestPingtaskSuite{})

var DeleteNqmAgentPingtaskSQL = `DELETE FROM nqm_agent_ping_task WHERE apt_ag_id >= 24021 AND apt_ag_id <= 24025`

var InsertPingtaskSQL = `
INSERT INTO nqm_ping_task(pt_id, pt_name, pt_period)
VALUES(10119, 'test1-pingtask_name', 40),
(10120, 'test2-pingtask_name', 3)
`
var DeletePingtaskSQL = `DELETE FROM nqm_ping_task WHERE pt_id >= 10119`

var InsertHostSQL = `
INSERT INTO host(id, hostname, agent_version, plugin_version)
VALUES(36091, 'ct-agent-1', '', ''),
	(36092, 'ct-agent-2', '', ''),
	(36093, 'ct-agent-3', '', ''),
	(36094, 'ct-agent-4', '', '')
`
var DeleteHostSQL = `DELETE FROM host WHERE id >= 36091 AND id <= 36095`

var InsertNqmAgentSQL = `
INSERT INTO nqm_agent(
	ag_id, ag_hs_id, ag_name, ag_connection_id, ag_hostname, ag_ip_address,
	ag_pv_id, ag_ct_id
)
VALUES(24021, 36091, 'ct-255-1', 'ct-255-1@201.3.116.1', 'ct-1', x'C9037401', 7, 255),
	(24022, 36092, 'ct-255-2', 'ct-255-2@201.3.116.2', 'ct-2', x'C9037402', 7, 255),
	(24023, 36093, 'ct-255-3', 'ct-255-3@201.4.23.3', 'ct-3', x'C9037403', 7, 255),
	(24024, 36094, 'ct-263-1', 'ct-63-1@201.77.23.3', 'ct-4', x'C9022403', 7, 263)
`
var DeleteNqmAgentSQL = `DELETE FROM nqm_agent WHERE ag_id >= 24021 AND ag_id <= 24025`

var InitNqmAgentAndPingtaskSQL = []string{InsertPingtaskSQL, InsertHostSQL, InsertNqmAgentSQL}
var CleanNqmAgentAndPingtaskSQL = []string{DeleteNqmAgentPingtaskSQL, DeleteNqmAgentSQL, DeleteHostSQL, DeletePingtaskSQL}

func (s *TestPingtaskSuite) SetUpTest(c *C) {
	var inTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case
		"TestPingtaskSuite.TestGetPingtaskById",
		"TestPingtaskSuite.TestListPingtasks",
		"TestPingtaskSuite.TestUpdateAndGetPingtask":
		inTx(InsertPingtaskSQL)
	case
		"TestPingtaskSuite.TestAssignPingtaskToAgentForAgent",
		"TestPingtaskSuite.TestRemovePingtaskFromAgentForAgent",
		"TestPingtaskSuite.TestAssignPingtaskToAgentForPingtask",
		"TestPingtaskSuite.TestRemovePingtaskFromAgentForPingtask":
		inTx(InitNqmAgentAndPingtaskSQL...)
	case
		"TestPingtaskSuite.TestAddAndGetPingtask":
	}
}

func (s *TestPingtaskSuite) TearDownTest(c *C) {
	var inTx = DbFacade.SqlDbCtrl.ExecQueriesInTx

	switch c.TestName() {
	case
		"TestPingtaskSuite.TestGetPingtaskById",
		"TestPingtaskSuite.TestListPingtasks",
		"TestPingtaskSuite.TestUpdateAndGetPingtask":
		inTx(DeletePingtaskSQL)
	case
		"TestPingtaskSuite.TestAssignPingtaskToAgentForAgent",
		"TestPingtaskSuite.TestRemovePingtaskFromAgentForAgent",
		"TestPingtaskSuite.TestAssignPingtaskToAgentForPingtask",
		"TestPingtaskSuite.TestRemovePingtaskFromAgentForPingtask":
		inTx(CleanNqmAgentAndPingtaskSQL...)
	case
		"TestPingtaskSuite.TestAddAndGetPingtask":
	}
}

func (s *TestPingtaskSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
	owlDb.DbFacade = DbFacade
}
func (s *TestPingtaskSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
	owlDb.DbFacade = nil
}

func (suite *TestPingtaskSuite) TestAssignPingtaskToAgentForAgent(c *C) {
	testCases := []*struct {
		inputAID      int32
		inputPID      int16
		expectedAgent *nqmModel.Agent
		expectedErr   error
	}{
		{24021, 10119, GetAgentById(24021), nil},
		{24022, 10119, GetAgentById(24022), nil},
		{24023, 10120, GetAgentById(24023), nil},
		// i > 2: cases for panic
		{24024, 10121, nil, nil},
		{24025, 10120, nil, nil},
		{24026, 10121, nil, nil},
	}

	for i, v := range testCases {
		c.Logf("case[%d]\n%+v\n", i, *v)
		if i > 2 {
			c.Assert(func() (*nqmModel.Agent, error) { return AssignPingtaskToAgentForAgent(v.inputAID, v.inputPID) }, PanicMatches, `*.FOREIGN KEY.*`)
			continue
		}
		actualAgent, actualErr := AssignPingtaskToAgentForAgent(v.inputAID, v.inputPID)
		c.Assert(actualAgent, NotNil)
		c.Assert(actualErr, IsNil)
	}
}

func (suite *TestPingtaskSuite) TestRemovePingtaskFromAgentForAgent(c *C) {
	testCases := []*struct {
		inputAID      int32
		inputPID      int16
		expectedAgent *nqmModel.Agent
		expectedErr   error
	}{
		{24021, 10119, GetAgentById(24021), nil},
		{24022, 10119, GetAgentById(24022), nil},
		{24023, 10120, GetAgentById(24023), nil},
		// i > 2: Not deleting
		{24024, 10121, GetAgentById(24024), nil},
		{24025, 10120, nil, nil},
		{24026, 10121, nil, nil},
	}

	for i, v := range testCases {
		c.Logf("case[%d]\n%+v\n", i, *v)
		if i == 3 {
			actualAgent, actualErr := RemovePingtaskFromAgentForAgent(v.inputAID, v.inputPID)
			c.Assert(actualAgent, NotNil)
			c.Assert(actualErr, IsNil)
			continue
		}
		if i > 3 {
			actualAgent, _ := RemovePingtaskFromAgentForAgent(v.inputAID, v.inputPID)
			c.Assert(actualAgent, IsNil)
			continue
		}
		AssignPingtaskToAgentForAgent(v.inputAID, v.inputPID)
		actualAgent, actualErr := RemovePingtaskFromAgentForAgent(v.inputAID, v.inputPID)
		c.Assert(actualAgent, NotNil)
		c.Assert(actualErr, IsNil)
	}
}

func (suite *TestPingtaskSuite) TestGetPingtaskById(c *C) {
	testCases := []*struct {
		input int16
	}{
		{10119}, // NotNil
		{10120}, // NotNil
		// i > 1: cases for panic
		{10121}, //IsNil
	}

	for i, v := range testCases {
		c.Logf("case[%d]\n%+v\n", i, *v)
		actual := GetPingtaskById(v.input)
		if i > 1 {
			c.Assert(actual, IsNil)
			continue
		}
		c.Assert(actual, NotNil)
	}
}

func (suite *TestPingtaskSuite) TestAssignPingtaskToAgentForPingtask(c *C) {
	testCases := []*struct {
		inputAID         int32
		inputPID         int16
		expectedPingtask *nqmModel.PingtaskView
		expectedErr      error
	}{
		{24021, 10119, GetPingtaskById(10119), nil},
		{24022, 10119, GetPingtaskById(10119), nil},
		{24023, 10120, GetPingtaskById(10120), nil},
		// i > 2: cases for panic
		{24024, 10121, nil, nil},
		{24025, 10120, GetPingtaskById(10120), nil},
		{24026, 10121, nil, nil},
	}

	for i, v := range testCases {
		c.Logf("case[%d]\n%+v\n", i, *v)
		if i > 2 {
			c.Assert(func() (*nqmModel.PingtaskView, error) {
				return AssignPingtaskToAgentForPingtask(v.inputAID, v.inputPID)
			}, PanicMatches, `*.FOREIGN KEY.*`)
			continue
		}
		actualPingtask, actualErr := AssignPingtaskToAgentForPingtask(v.inputAID, v.inputPID)
		c.Assert(actualPingtask, NotNil)
		c.Assert(actualErr, IsNil)
	}
}

func (suite *TestPingtaskSuite) TestRemovePingtaskFromAgentForPingtask(c *C) {
	testCases := []*struct {
		inputAID         int32
		inputPID         int16
		expectedPingtask *nqmModel.PingtaskView
		expectedErr      error
	}{
		{24021, 10119, GetPingtaskById(10119), nil},
		{24022, 10119, GetPingtaskById(10119), nil},
		{24023, 10120, GetPingtaskById(10120), nil},
		// i > 2: Not deleting
		{24024, 10121, GetPingtaskById(10121), nil},
		{24025, 10120, nil, nil},
		{24026, 10121, nil, nil},
	}

	for i, v := range testCases {
		c.Logf("case[%d]\n%+v\n", i, *v)
		if i == 3 || i == 5 {
			actualPingtask, _ := RemovePingtaskFromAgentForPingtask(v.inputAID, v.inputPID)
			c.Assert(actualPingtask, IsNil)
			continue
		}
		if i == 4 {
			actualPingtask, actualErr := RemovePingtaskFromAgentForPingtask(v.inputAID, v.inputPID)
			c.Assert(actualPingtask, NotNil)
			c.Assert(actualErr, IsNil)
			continue
		}
		AssignPingtaskToAgentForPingtask(v.inputAID, v.inputPID)
		actualPingtask, actualErr := RemovePingtaskFromAgentForPingtask(v.inputAID, v.inputPID)
		c.Assert(actualPingtask, NotNil)
		c.Assert(actualErr, IsNil)
	}
}

func (suite *TestPingtaskSuite) TestListPingtasks(c *C) {
	testCases := []*struct {
		query                      *nqmModel.PingtaskQuery
		paging                     commonModel.Paging
		expectedCountOfCurrentPage int
		expectedCountOfAll         int32
	}{
		{
			&nqmModel.PingtaskQuery{},
			commonModel.Paging{Size: 2, Position: 1, OrderBy: []*commonModel.OrderByEntity{}},
			2, 2,
		},
		{
			&nqmModel.PingtaskQuery{},
			commonModel.Paging{Size: 1, Position: 1, OrderBy: []*commonModel.OrderByEntity{}},
			1, 2,
		},
		{
			&nqmModel.PingtaskQuery{},
			commonModel.Paging{Size: 10, Position: 10, OrderBy: []*commonModel.OrderByEntity{}},
			0, 2,
		},
		{
			&nqmModel.PingtaskQuery{
				Period: "3",
				Name:   "test2",
			},
			commonModel.Paging{Size: 2, Position: 1, OrderBy: []*commonModel.OrderByEntity{}},
			1, 1,
		},
		{
			&nqmModel.PingtaskQuery{
				Period: "40",
			},
			commonModel.Paging{Size: 2, Position: 1, OrderBy: []*commonModel.OrderByEntity{}},
			1, 1,
		},
		{
			&nqmModel.PingtaskQuery{
				Name: "test1",
			},
			commonModel.Paging{Size: 2, Position: 1, OrderBy: []*commonModel.OrderByEntity{}},
			1, 1,
		},
		{
			&nqmModel.PingtaskQuery{
				Enable: "true",
			},
			commonModel.Paging{Size: 2, Position: 1, OrderBy: []*commonModel.OrderByEntity{}},
			2, 2,
		},
		//	{	// num_of_enabled_agents }
	}

	for i, v := range testCases {
		fmt.Printf("%+v\n", v)
		actualResult, actualPaging := ListPingtasks(v.query, v.paging)
		c.Logf("case [%d]:", i)
		c.Logf("[List] Query condition: %v. Number of agents: %d", v.query, len(actualResult))

		for _, pingtask := range actualResult {
			c.Logf("[List] Pingtask: %v.", pingtask)
		}
		c.Assert(actualResult, HasLen, v.expectedCountOfCurrentPage)
		c.Assert(actualPaging.TotalCount, Equals, v.expectedCountOfAll)
	}

}

func (suite *TestPingtaskSuite) TestAddAndGetPingtask(c *C) {
	var pm1 *nqmModel.PingtaskModify
	if err := json.Unmarshal([]byte(`
		{
		  "period" : 15,
		  "name" : "廣東",
		  "enable" : true,
		  "comment" : "This is for some purpose",
		  "filter" : {
		    "ids_of_isp" : [ 7, 8, 9 ],
		    "ids_of_province" : [ 2, 3 ],
		    "ids_of_city" : [ 11, 12, 13 ]
		  }
		}
	`), &pm1); err != nil {
		c.Error(err)
	}
	testCases := []*struct {
		inputPm *nqmModel.PingtaskModify
	}{
		{pm1},
	}
	for _, v := range testCases {
		actual := AddAndGetPingtask(v.inputPm)
		c.Assert(actual, NotNil)
	}
}

func (suite *TestPingtaskSuite) TestUpdateAndGetPingtask(c *C) {
	var pm1 *nqmModel.PingtaskModify
	if err := json.Unmarshal([]byte(`
		{
		  "period" : 15,
		  "name" : "廣東",
		  "enable" : true,
		  "comment" : "This is for some purpose",
		  "filter" : {
		    "ids_of_isp" : [ 7, 8, 9 ],
		    "ids_of_province" : [ 2, 3 ],
		    "ids_of_city" : [ 11, 12, 13 ]
		  }
		}
	`), &pm1); err != nil {
		c.Error(err)
	}
	testCases := []*struct {
		inputPm *nqmModel.PingtaskModify
	}{
		{pm1},
	}
	for _, v := range testCases {
		actual := UpdateAndGetPingtask(10120, v.inputPm)
		c.Assert(actual, NotNil)
	}
}
