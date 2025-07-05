package dto

// Save Parameter
type SaveConfig struct {
	ConfigId    int    `json:"configId"`
	ConfigName  string `json:"configName"`
	ConfigKey   string `json:"configKey"`
	ConfigValue string `json:"configValue"`
	ConfigType  string `json:"configType"`
	CreateBy    string `json:"createBy"`
	UpdateBy    string `json:"updateBy"`
	Remark      string `json:"remark"`
}

// Parameter List
type ConfigListRequest struct {
	PageRequest
	ConfigName string `query:"configName" form:"configName"`
	ConfigKey  string `query:"configKey" form:"configKey"`
	ConfigType string `query:"configType" form:"configType"`
	BeginTime  string `query:"params[beginTime]" form:"params[beginTime]"`
	EndTime    string `query:"params[endTime]" form:"params[endTime]"`
}

// Create Parameter
type CreateConfigRequest struct {
	ConfigName  string `json:"configName"`
	ConfigKey   string `json:"configKey"`
	ConfigValue string `json:"configValue"`
	ConfigType  string `json:"configType"`
	Remark      string `json:"remark"`
}

// Update Parameter
type UpdateConfigRequest struct {
	ConfigId    int    `json:"configId"`
	ConfigName  string `json:"configName"`
	ConfigKey   string `json:"configKey"`
	ConfigValue string `json:"configValue"`
	ConfigType  string `json:"configType"`
	Remark      string `json:"remark"`
}
