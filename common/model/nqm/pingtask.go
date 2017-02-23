package nqm

import (
	"fmt"
	"strings"

	owlJson "github.com/Cepave/open-falcon-backend/common/json"
	commonOwlModel "github.com/Cepave/open-falcon-backend/common/model/owl"
	"github.com/chyeh/cast"
)

type AgentPingtask struct {
	AgentID    int32 `json:"-"`
	PingtaskID int32 `json:"-"`
}

type pingtaskFilter struct {
	IspFilters      []*commonOwlModel.IspOfPingtaskView      `json:"isps"`
	ProvinceFilters []*commonOwlModel.ProvinceOfPingtaskView `json:"provinces"`
	CityFilters     []*commonOwlModel.CityOfPingtaskView     `json:"cities"`
	NameTagFilters  []*commonOwlModel.NameTagOfPingtaskView  `json:"name_tags"`
	GroupTagFilters []*commonOwlModel.GroupTagOfPingtaskView `json:"group_tags"`
}

type PingtaskView struct {
	ID                 int16              `gorm:"primary_key:true;column:pt_id" json:"id"`
	Period             int8               `gorm:"column:pt_period" json:"period"`
	Name               owlJson.JsonString `gorm:"column:pt_name" json:"name"`
	Enable             bool               `gorm:"column:pt_enable" json:"enable"`
	Comment            owlJson.JsonString `gorm:"column:pt_comment" json:"comment"`
	NumOfEnabledAgents int8               `gorm:"column:pt_num_of_enabled_agents" json:"num_of_enabled_agents"`

	IdsOfIspFilters  string `gorm:"column:pt_isp_filter_ids" json:"-"`
	NamesOfIspFilter string `gorm:"column:pt_isp_filter_names" json:"-"`
	//IspFilters       []*commonOwlModel.IspOfPingtaskView

	IdsOfProvinceFilters  string `gorm:"column:pt_province_filter_ids" json:"-"`
	NamesOfProvinceFilter string `gorm:"column:pt_province_filter_names" json:"-"`
	//ProvinceFilters       []*commonOwlModel.ProvinceOfPingtaskView

	IdsOfCityFilters         string `gorm:"column:pt_city_filter_ids" json:"-"`
	ProvinceIdsOfCityFilters string `gorm:"column:pt_city_filter_pv_ids" json:"-"`
	NamesOfCityFilter        string `gorm:"column:pt_city_filter_names" json:"-"`
	//CityFilters              []*commonOwlModel.CityOfPingtaskView

	IdsOfNameTagFilters  string `gorm:"column:pt_name_tag_filter_ids" json:"-"`
	NamesOfNameTagFilter string `gorm:"column:pt_name_tag_filter_values" json:"-"`
	//NameTagFilters       []*commonOwlModel.NameTagOfPingtaskView

	IdsOfGroupTagFilters  string `gorm:"column:pt_group_tag_filter_ids" json:"-"`
	NamesOfGroupTagFilter string `gorm:"column:pt_group_tag_filter_names" json:"-"`
	//GroupTagFilters       []*commonOwlModel.GroupTagOfPingtaskView

	Filter pingtaskFilter `json:"filter"`
}

func (PingtaskView) TableName() string {
	return "nqm_ping_task"
}

//func (p *PingtaskView) MarshalJSON() ([]byte, error) {
//	jsonObject := json.New()
//
//	jsonObject.Set("id", p.ID)
//	jsonObject.Set("period", p.Period)
//	jsonObject.Set("name", p.Name)
//	jsonObject.Set("enable", p.Enable)
//	jsonObject.Set("comment", p.Comment)
//	jsonObject.Set("num_of_enabled_agents", p.NumOfEnabledAgents)
//
//	jsonFilter := json.New()
//	jsonObject.Set("filter", jsonFilter)
//
//	jsonIsp := json.New()
//	jsonIsp.Set("id", p.Filter.IspFilters.I.Id)
//	jsonIsp.Set("name", p.Filter.IspFilters.Name)
//	jsonObject.Set("isp", jsonIsp)
//
//	jsonProvince := json.New()
//	jsonProvince.Set("id", p.ProvinceId)
//	jsonProvince.Set("name", p.ProvinceName)
//	jsonObject.Set("province", jsonProvince)
//
//	jsonCity := json.New()
//	jsonCity.Set("id", p.CityId)
//	jsonCity.Set("name", p.CityName)
//	jsonObject.Set("city", jsonCity)
//
//	jsonNameTag := json.New()
//	jsonNameTag.Set("id", p.NameTagId)
//	jsonNameTag.Set("value", p.NameTagValue)
//	jsonObject.Set("name_tag", jsonNameTag)
//
//	jsonGroupTags := json.New()
//	jsonNameTag.Set("id", p.NameTagId)
//	jsonNameTag.Set("value", p.NameTagValue)
//	jsonObject.Set("name_tag", jsonNameTag)
//
//	return jsonObject.MarshalJSON()
//}

func (p *PingtaskView) AfterLoad() {
	var ids []string
	var names []string
	p.Filter.IspFilters = make([]*commonOwlModel.IspOfPingtaskView, 0)
	p.Filter.ProvinceFilters = make([]*commonOwlModel.ProvinceOfPingtaskView, 0)
	p.Filter.CityFilters = make([]*commonOwlModel.CityOfPingtaskView, 0)
	p.Filter.NameTagFilters = make([]*commonOwlModel.NameTagOfPingtaskView, 0)
	p.Filter.GroupTagFilters = make([]*commonOwlModel.GroupTagOfPingtaskView, 0)
	//p.IspFilters = owlModel.SplitToArrayOfGroupTags(
	//	p.IdsOfIspFilters, ",",
	//	p.NamesOfIspFilter, "\000",
	//)
	if p.IdsOfIspFilters != "" {
		ids = strings.Split(p.IdsOfIspFilters, ",")
		names = strings.Split(p.NamesOfIspFilter, "\000")
		if len(ids) != len(names) {
			panic(fmt.Errorf("Error on parsing: Can't match ids and names"))
		}
		for i := range ids {
			p.Filter.IspFilters = append(
				p.Filter.IspFilters,
				&commonOwlModel.IspOfPingtaskView{
					Id:   cast.ToInt(ids[i]),
					Name: names[i],
				},
			)
		}
	}

	//p.ProvinceFilters = owlModel.SplitToArrayOfGroupTags(
	//	p.IdsOfProvinceFilters, ",",
	//	p.NamesOfProvinceFilter, "\000",
	//)
	if p.IdsOfProvinceFilters != "" {
		ids = strings.Split(p.IdsOfProvinceFilters, ",")
		names = strings.Split(p.NamesOfProvinceFilter, "\000")
		if len(ids) != len(names) {
			panic(fmt.Errorf("Error on parsing: Can't match ids and names"))
		}
		for i := range ids {
			p.Filter.ProvinceFilters = append(
				p.Filter.ProvinceFilters,
				&commonOwlModel.ProvinceOfPingtaskView{
					Id:   cast.ToInt(ids[i]),
					Name: names[i],
				},
			)
		}
	}

	//p.CityFilters = owlModel.SplitToArrayOfGroupTags(
	//	p.IdsOfCityFilters, ",",
	//	p.ProvinceIdsOfCityFilters, "\000",
	//)
	if p.IdsOfCityFilters != "" {
		ids = strings.Split(p.IdsOfCityFilters, ",")
		pvIds := strings.Split(p.ProvinceIdsOfCityFilters, ",")
		names = strings.Split(p.NamesOfCityFilter, "\000")
		if len(ids) != len(names) {
			panic(fmt.Errorf("Error on parsing: Can't match ids and names"))
		}
		for i := range ids {
			p.Filter.CityFilters = append(
				p.Filter.CityFilters,
				&commonOwlModel.CityOfPingtaskView{
					Id:         cast.ToInt(ids[i]),
					ProvinceId: cast.ToInt(pvIds[i]),
					Name:       names[i],
				},
			)
		}
	}

	//p.NameTagFilters = owlModel.SplitToArrayOfGroupTags(
	//	p.IdsOfNameTagFilters, ",",
	//	p.NamesOfNameTagFilter, "\000",
	//)
	if p.IdsOfNameTagFilters != "" {
		ids = strings.Split(p.IdsOfNameTagFilters, ",")
		values := strings.Split(p.NamesOfNameTagFilter, "\000")
		if len(ids) != len(values) {
			panic(fmt.Errorf("Error on parsing: Can't match ids and names"))
		}
		for i := range ids {
			p.Filter.NameTagFilters = append(
				p.Filter.NameTagFilters,
				&commonOwlModel.NameTagOfPingtaskView{
					Id:    cast.ToInt(ids[i]),
					Value: values[i],
				},
			)
		}
	}

	//p.GroupTagFilters = owlModel.SplitToArrayOfGroupTags(
	//	p.IdsOfGroupTagFilters, ",",
	//	p.NamesOfGroupTagFilter, "\000",
	//)
	if p.IdsOfGroupTagFilters != "" {
		ids = strings.Split(p.IdsOfGroupTagFilters, ",")
		names = strings.Split(p.NamesOfGroupTagFilter, "\000")
		if len(ids) != len(names) {
			panic(fmt.Errorf("Error on parsing: Can't match ids and names"))
		}
		for i := range ids {
			p.Filter.GroupTagFilters = append(
				p.Filter.GroupTagFilters,
				&commonOwlModel.GroupTagOfPingtaskView{
					Id:   cast.ToInt(ids[i]),
					Name: names[i],
				},
			)
		}
	}
}
