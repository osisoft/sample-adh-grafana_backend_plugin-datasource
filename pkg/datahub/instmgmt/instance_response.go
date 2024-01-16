package instmgmt

type InstanceResponse struct {
	AccountId string `json:"accountId"`
	ServiceId string `json:"serviceId"`
	Id        string `json:"id"`
	Geography string `json:"geography"`
	Name      string `json:"name"`
}
