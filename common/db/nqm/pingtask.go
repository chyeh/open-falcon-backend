package nqm

import (
	"github.com/jmoiron/sqlx"

	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
)

type addAgentPingtaskTx struct {
	agentPingtask *nqmModel.AgentPingtask
	err           error
}

func (agentPingtaskTx *addAgentPingtaskTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	tx.MustExec(
		`
		INSERT INTO nqm_agent_ping_task(apt_ag_id,apt_pt_id)
		VALUES
		(?,?)
		ON DUPLICATE KEY UPDATE
		apt_ag_id=VALUES(apt_ag_id),
		apt_pt_id=VALUES(apt_pt_id)
		`,
		agentPingtaskTx.agentPingtask.AgentID,
		agentPingtaskTx.agentPingtask.PingtaskID,
	)
	//agentPingtaskTx.agentPingtask.AgentID = owlDb.BuildAndGetNameTagId(
	//	tx, agentPingtaskTx.agent.NameTagValue,
	//)

	//agentPingtaskTx.addAgent(tx)
	if agentPingtaskTx.err != nil {
		return commonDb.TxRollback
	}

	//agentPingtaskTx.prepareGroupTags(tx)
	return commonDb.TxCommit
}

func AssignPingtaskToAgentForAgent(aID int32, pID int32) (*nqmModel.Agent, error) {
	txProcessor := &addAgentPingtaskTx{
		agentPingtask: &nqmModel.AgentPingtask{AgentID: aID, PingtaskID: pID},
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	// :~)

	if txProcessor.err != nil {
		return nil, txProcessor.err
	}

	return GetAgentById(aID), nil
}

type deleteAgentPingtaskTx struct {
	agentPingtask *nqmModel.AgentPingtask
	err           error
}

func (agentPingtaskTx *deleteAgentPingtaskTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	tx.MustExec(
		`
		DELETE from nqm_agent_ping_task WHERE
		apt_ag_id=? AND apt_pt_id=? LIMIT 1;
		`,
		agentPingtaskTx.agentPingtask.AgentID,
		agentPingtaskTx.agentPingtask.PingtaskID,
	)

	if agentPingtaskTx.err != nil {
		return commonDb.TxRollback
	}

	return commonDb.TxCommit
}

func RemovePingtaskFromAgentForAgent(aID int32, pID int32) (*nqmModel.Agent, error) {
	agent := GetAgentById(aID)
	txProcessor := &deleteAgentPingtaskTx{
		agentPingtask: &nqmModel.AgentPingtask{AgentID: aID, PingtaskID: pID},
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	// :~)

	if txProcessor.err != nil {
		return nil, txProcessor.err
	}

	return agent, nil
}

func AssignPingtaskToAgentForPingtask(aID int32, pID int32) (*nqmModel.Agent, error) {
	txProcessor := &addAgentPingtaskTx{
		agentPingtask: &nqmModel.AgentPingtask{AgentID: aID, PingtaskID: pID},
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	// :~)

	if txProcessor.err != nil {
		return nil, txProcessor.err
	}

	return GetAgentById(aID), nil
}

func RemovePingtaskFromAgentForPingtask(aID int32, pID int32) (*nqmModel.Agent, error) {
	agent := GetAgentById(aID)
	txProcessor := &deleteAgentPingtaskTx{
		agentPingtask: &nqmModel.AgentPingtask{AgentID: aID, PingtaskID: pID},
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	// :~)

	if txProcessor.err != nil {
		return nil, txProcessor.err
	}

	return agent, nil
}
