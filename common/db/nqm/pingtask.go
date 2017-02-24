package nqm

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"

	commonDb "github.com/Cepave/open-falcon-backend/common/db"
	gormExt "github.com/Cepave/open-falcon-backend/common/gorm"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
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

var orderByDialectForPingtasks = commonModel.NewSqlOrderByDialect(
	map[string]string{
		"period":                "pt_period",
		"name":                  "pt_name",
		"enable":                "pt_enable",
		"comment":               "pt_comment",
		"num_of_enabled_agents": "pt_num_of_enabled_agents",
	},
)

func buildSortingClauseOfPingtasks(paging *commonModel.Paging) string {
	//if len(paging.OrderBy) == 0 {
	//	paging.OrderBy = append(paging.OrderBy, &commonModel.OrderByEntity{"num_of_enabled_agents", commonModel.Descending})
	//}

	querySyntax, err := orderByDialectForPingtasks.ToQuerySyntax(paging.OrderBy)
	gormExt.DefaultGormErrorConverter.PanicIfError(err)

	return querySyntax
}

// Lists the pingtasks according to the query parameters
func ListPingtasks(query *nqmModel.PingtaskQuery, paging commonModel.Paging) ([]*nqmModel.PingtaskView, *commonModel.Paging) {
	var result []*nqmModel.PingtaskView

	var funcTxLoader gormExt.TxCallbackFunc = func(txGormDb *gorm.DB) commonDb.TxFinale {
		/**
		 * Retrieves the page of data
		 */
		var selectPingtask = txGormDb.Model(&nqmModel.PingtaskView{}).
			Select(`SQL_CALC_FOUND_ROWS
				pt_id, pt_period, pt_name, pt_enable, pt_comment,
        COUNT(DISTINCT ag.ag_id) AS pt_num_of_enabled_agents,
        GROUP_CONCAT(DISTINCT isp.isp_id ORDER BY isp_id ASC SEPARATOR ',') AS pt_isp_filter_ids,
        GROUP_CONCAT(DISTINCT isp.isp_name ORDER BY isp_id ASC SEPARATOR '\0') AS pt_isp_filter_names,
        GROUP_CONCAT(DISTINCT pv.pv_id ORDER BY pv_id ASC SEPARATOR ',') AS pt_province_filter_ids,
        GROUP_CONCAT(DISTINCT pv.pv_name ORDER BY pv_id ASC SEPARATOR '\0') AS pt_province_filter_names,
        GROUP_CONCAT(DISTINCT ct.ct_id ORDER BY ct_id ASC SEPARATOR ',') AS pt_city_filter_ids,
        GROUP_CONCAT(DISTINCT ct.ct_pv_id ORDER BY ct_id ASC SEPARATOR ',') AS pt_city_filter_pv_ids,
        GROUP_CONCAT(DISTINCT ct.ct_name ORDER BY ct_id ASC SEPARATOR '\0') AS pt_city_filter_names,
        GROUP_CONCAT(DISTINCT nt.nt_id ORDER BY nt_id ASC SEPARATOR ',') AS pt_name_tag_filter_ids,
        GROUP_CONCAT(DISTINCT nt.nt_value ORDER BY nt_id ASC SEPARATOR '\0') AS pt_name_tag_filter_values,
        GROUP_CONCAT(DISTINCT gt.gt_id ORDER BY gt_id ASC SEPARATOR ',') AS pt_group_tag_filter_ids,
        GROUP_CONCAT(DISTINCT gt.gt_name ORDER BY gt_id ASC SEPARATOR '\0') AS pt_group_tag_filter_names
			`).
			Joins(`
				LEFT JOIN
				nqm_agent_ping_task AS apt
				ON pt_id = apt.apt_pt_id
				LEFT JOIN
				nqm_agent AS ag
				ON apt.apt_ag_id = ag.ag_id AND ag.ag_status=true
				LEFT JOIN
				nqm_pt_target_filter_isp AS tfisp
				ON pt_id = tfisp.tfisp_pt_id
				LEFT JOIN
				owl_isp AS isp
				ON tfisp.tfisp_isp_id = isp.isp_id
				LEFT JOIN
				nqm_pt_target_filter_province AS tfpv
				ON pt_id = tfpv.tfpv_pt_id
				LEFT JOIN
				owl_province AS pv
				ON tfpv.tfpv_pv_id = pv.pv_id
				LEFT JOIN
				nqm_pt_target_filter_city AS tfct
				ON pt_id = tfct.tfct_pt_id
				LEFT JOIN
				owl_city AS ct
				ON tfct.tfct_ct_id = ct.ct_id
				LEFT JOIN
				nqm_pt_target_filter_name_tag AS tfnt
				ON pt_id = tfnt.tfnt_pt_id
				LEFT JOIN
				owl_name_tag AS nt
				ON tfnt.tfnt_nt_id = nt.nt_id
				LEFT JOIN
				nqm_pt_target_filter_group_tag AS tfgt
				ON pt_id = tfgt.tfgt_pt_id
				LEFT JOIN
				owl_group_tag AS gt
				ON tfgt.tfgt_gt_id = gt.gt_id
			`).
			Limit(paging.Size).
			Group(`
				pt_id, pt_period, pt_name, pt_enable, pt_comment
			`).
			Order(buildSortingClauseOfPingtasks(&paging)).
			Offset(paging.GetOffset())

		if query.Period != "" {
			selectPingtask = selectPingtask.Where("pt_period = ?", query.Period)
		}
		if query.Name != "" {
			selectPingtask = selectPingtask.Where("pt_name LIKE ?", query.Name+"%")
		}
		if query.Enable != "" {
			selectPingtask = selectPingtask.Where("pt_enable = ?", query.Enable)
		}
		if query.Comment != "" {
			selectPingtask = selectPingtask.Where("pt_comment LIKE ?", query.Comment+"%")
		}
		if query.NumOfEnabledAgents != "" {
			selectPingtask = selectPingtask.Where("pt_num_of_enabled_agents = ?", query.NumOfEnabledAgents)
		}
		// :~)

		gormExt.ToDefaultGormDbExt(selectPingtask.Find(&result)).PanicIfError()

		return commonDb.TxCommit
	}

	gormExt.ToDefaultGormDbExt(DbFacade.GormDb).SelectWithFoundRows(
		funcTxLoader, &paging,
	)

	/**
	 * Loads group tags
	 */
	for _, pingtask := range result {
		pingtask.AfterLoad()
	}
	// :~)

	return result, &paging
}

func GetPingtaskById(id int16) *nqmModel.PingtaskView {
	var selectPingtask = DbFacade.GormDb.Model(&nqmModel.PingtaskView{}).
		Select(`
			pt_id, pt_period, pt_name, pt_enable, pt_comment,
			COUNT(DISTINCT ag.ag_id) AS pt_num_of_enabled_agents,
			GROUP_CONCAT(DISTINCT isp.isp_id ORDER BY isp_id ASC SEPARATOR ',') AS pt_isp_filter_ids,
			GROUP_CONCAT(DISTINCT isp.isp_name ORDER BY isp_id ASC SEPARATOR '\0') AS pt_isp_filter_names,
			GROUP_CONCAT(DISTINCT pv.pv_id ORDER BY pv_id ASC SEPARATOR ',') AS pt_province_filter_ids,
			GROUP_CONCAT(DISTINCT pv.pv_name ORDER BY pv_id ASC SEPARATOR '\0') AS pt_province_filter_names,
			GROUP_CONCAT(DISTINCT ct.ct_id ORDER BY ct_id ASC SEPARATOR ',') AS pt_city_filter_ids,
			GROUP_CONCAT(DISTINCT ct.ct_pv_id ORDER BY ct_id ASC SEPARATOR ',') AS pt_city_filter_pv_ids,
			GROUP_CONCAT(DISTINCT ct.ct_name ORDER BY ct_id ASC SEPARATOR '\0') AS pt_city_filter_names,
			GROUP_CONCAT(DISTINCT nt.nt_id ORDER BY nt_id ASC SEPARATOR ',') AS pt_name_tag_filter_ids,
			GROUP_CONCAT(DISTINCT nt.nt_value ORDER BY nt_id ASC SEPARATOR '\0') AS pt_name_tag_filter_values,
			GROUP_CONCAT(DISTINCT gt.gt_id ORDER BY gt_id ASC SEPARATOR ',') AS pt_group_tag_filter_ids,
			GROUP_CONCAT(DISTINCT gt.gt_name ORDER BY gt_id ASC SEPARATOR '\0') AS pt_group_tag_filter_names
		`).
		Joins(`
			LEFT JOIN
			nqm_agent_ping_task AS apt
			ON pt_id = apt.apt_pt_id
			LEFT JOIN
			nqm_agent AS ag
			ON apt.apt_ag_id = ag.ag_id AND ag.ag_status=true
			LEFT JOIN
			nqm_pt_target_filter_isp AS tfisp
			ON pt_id = tfisp.tfisp_pt_id
			LEFT JOIN
			owl_isp AS isp
			ON tfisp.tfisp_isp_id = isp.isp_id
			LEFT JOIN
			nqm_pt_target_filter_province AS tfpv
			ON pt_id = tfpv.tfpv_pt_id
			LEFT JOIN
			owl_province AS pv
			ON tfpv.tfpv_pv_id = pv.pv_id
			LEFT JOIN
			nqm_pt_target_filter_city AS tfct
			ON pt_id = tfct.tfct_pt_id
			LEFT JOIN
			owl_city AS ct
			ON tfct.tfct_ct_id = ct.ct_id
			LEFT JOIN
			nqm_pt_target_filter_name_tag AS tfnt
			ON pt_id = tfnt.tfnt_pt_id
			LEFT JOIN
			owl_name_tag AS nt
			ON tfnt.tfnt_nt_id = nt.nt_id
			LEFT JOIN
			nqm_pt_target_filter_group_tag AS tfgt
			ON pt_id = tfgt.tfgt_pt_id
			LEFT JOIN
			owl_group_tag AS gt
			ON tfgt.tfgt_gt_id = gt.gt_id
		`).
		Where("pt_id = ?", id).
		Group(`
			pt_id, pt_period, pt_name, pt_enable, pt_comment
		`)

	var loadedPingtask = &nqmModel.PingtaskView{}
	selectPingtask = selectPingtask.Find(loadedPingtask)

	if selectPingtask.Error == gorm.ErrRecordNotFound {
		return nil
	}
	gormExt.ToDefaultGormDbExt(selectPingtask).PanicIfError()

	loadedPingtask.AfterLoad()
	return loadedPingtask
}

type addPingtaskTx struct {
	pingtask *nqmModel.PingtaskModify
	res      sql.Result
	err      error
}

func (p *addPingtaskTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	r := tx.MustExec(
		`
		INSERT INTO nqm_ping_task(pt_period,pt_name,pt_enable,pt_comment)
		VALUES
		(?,?,?,?)
		`,
		p.pingtask.Period,
		p.pingtask.Name,
		p.pingtask.Enable,
		p.pingtask.Comment,
	)
	if p.err != nil {
		return commonDb.TxRollback
	}
	p.res = r
	return commonDb.TxCommit
}

func AddAndGetPingtask(pm *nqmModel.PingtaskModify) *nqmModel.PingtaskView {
	//if valueOfNameTag == "" {
	//	return -1
	//}

	//tx.MustExec(
	//	`
	//	INSERT INTO owl_name_tag(nt_value)
	//	SELECT ?
	//	FROM DUAL
	//	WHERE NOT EXISTS (
	//		SELECT *
	//		FROM owl_name_tag
	//		WHERE nt_value = ?
	//	)
	//	`,
	//	valueOfNameTag, valueOfNameTag,
	//)

	//var nameTagId int16
	//sqlxExt.ToTxExt(tx).Get(
	//	&nameTagId,
	//	`
	//	SELECT nt_id FROM owl_name_tag
	//	WHERE nt_value = ?
	//	`,
	//	valueOfNameTag,
	//)
	txProcessor := &addPingtaskTx{
		pingtask: pm,
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	// :~)

	if txProcessor.err != nil {
		commonDb.PanicIfError(txProcessor.err)
	}
	ptId, _ := txProcessor.res.LastInsertId()
	return GetPingtaskById(int16(ptId))
}

type updatePingtaskTx struct {
	pingtask   *nqmModel.PingtaskModify
	pingtaskID int16
	err        error
}

func (u *updatePingtaskTx) InTx(tx *sqlx.Tx) commonDb.TxFinale {
	tx.MustExec(
		`
		UPDATE nqm_ping_task SET
		pt_period = ?,
		pt_name = ?,
		pt_enable = ?,
		pt_comment = ?
		WHERE pt_id = ?
		`,
		u.pingtask.Period,
		u.pingtask.Name,
		u.pingtask.Enable,
		u.pingtask.Comment,
		u.pingtaskID,
	)
	if u.err != nil {
		return commonDb.TxRollback
	}
	return commonDb.TxCommit
}

func UpdateAndGetPingtask(id int16, pm *nqmModel.PingtaskModify) *nqmModel.PingtaskView {
	txProcessor := &updatePingtaskTx{
		pingtaskID: id,
		pingtask:   pm,
	}

	DbFacade.NewSqlxDbCtrl().InTx(txProcessor)
	// :~)

	if txProcessor.err != nil {
		commonDb.PanicIfError(txProcessor.err)
	}
	return GetPingtaskById(txProcessor.pingtaskID)
}
