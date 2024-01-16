package instmgmt

type InstanceQueryResponse struct {
	Items              []InstanceResponse `json:"items"`
	ContinutationToken string             `json:"continuationToken"`
}
