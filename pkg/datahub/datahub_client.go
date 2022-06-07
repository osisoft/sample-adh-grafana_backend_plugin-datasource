package datahub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/osisoft/sample-adh-grafana_backend_plugin-datasource/pkg/datahub/community"
	"github.com/osisoft/sample-adh-grafana_backend_plugin-datasource/pkg/datahub/sds"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type DataHubClient struct {
	resource        string
	apiVersion      string
	tenantId        string
	clientId        string
	clientSecret    string
	token           string
	tokenExpiration int64
	client          *http.Client
}

func NewDataHubClient(resource string, apiVersion string, tenantId string, clientId string, clientSecret string) DataHubClient {
	return DataHubClient{
		resource:     resource,
		apiVersion:   apiVersion,
		tenantId:     tenantId,
		clientId:     clientId,
		clientSecret: clientSecret,
		client:       &http.Client{},
	}
}

func GetClientToken(d *DataHubClient) string {
	if (d.tokenExpiration - time.Now().Unix()) > (5 * 60) {
		return ("Bearer " + d.token)
	}

	wellKnownEndpoint := d.resource + "/identity/.well-known/openid-configuration"
	req, err := http.NewRequest("GET", wellKnownEndpoint, nil)
	if err != nil {
		log.DefaultLogger.Warn(err.Error())
	}

	resp, err := d.client.Do(req)
	if err != nil {
		log.DefaultLogger.Warn(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.DefaultLogger.Warn(err.Error())
	}

	var openIdConfig map[string]interface{}

	err = json.Unmarshal(body, &openIdConfig)
	if err != nil {
		log.DefaultLogger.Warn("Error parsing json", err.Error())
	}

	tokenEndpoint := openIdConfig["token_endpoint"].(string)

	resp, err = d.client.PostForm(tokenEndpoint,
		url.Values{
			"client_id":     {d.clientId},
			"client_secret": {d.clientSecret},
			"grant_type":    {"client_credentials"}})

	if err != nil {
		log.DefaultLogger.Warn(err.Error())
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.DefaultLogger.Warn(err.Error())
	}
	defer resp.Body.Close()

	var tokenInformation map[string]interface{}

	err = json.Unmarshal(body, &tokenInformation)
	if err != nil {
		log.DefaultLogger.Warn("Error parsing json", err.Error())
	}

	d.token = tokenInformation["access_token"].(string)
	d.tokenExpiration = int64(tokenInformation["expires_in"].(float64)) + time.Now().Unix()

	return ("Bearer " + d.token)
}

func SdsRequest(d *DataHubClient, token string, path string, headers map[string]string) ([]byte, error) {
	// request data or collection items
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.DefaultLogger.Warn(err.Error())
		return nil, err
	}

	req.Header.Add("Authorization", token)

	// add optional headers
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	resp, err := d.client.Do(req)
	if err != nil {
		log.DefaultLogger.Warn(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.DefaultLogger.Warn(err.Error())
		return nil, err
	}

	return body, nil
}

func StreamsQuery(d *DataHubClient, namespaceId string, token string, query string) (*data.Frame, error) {
	basePath := d.resource + "/api/" + d.apiVersion + "/tenants/" + d.tenantId + "/namespaces/" + namespaceId
	path := (basePath + "/streams?query=" + query)

	body, err := SdsRequest(d, token, path, nil)
	if err != nil {
		return nil, err
	}

	var streams []sds.SdsStream

	err = json.Unmarshal(body, &streams)
	if err != nil {
		log.DefaultLogger.Warn("Error parsing json", err.Error())
		return nil, err
	}

	// create a dataframe
	frame := data.NewFrame("response")

	// create property lists from streams list
	ids := make([]string, len(streams))
	names := make([]string, len(streams))
	for i := 0; i < len(streams); i++ {
		ids[i] = streams[i].Id
		names[i] = streams[i].Name
	}

	// add fields
	frame.Fields = append(frame.Fields,
		data.NewField("Id", nil, ids),
		data.NewField("Name", nil, names),
	)

	return frame, nil
}

func CommunityStreamsQuery(d *DataHubClient, communityId string, token string, query string) (*data.Frame, error) {
	basePath := d.resource + "/api/" + d.apiVersion + "/search/communities/" + communityId

	path := (basePath + "/streams?query=" + query)

	body, err := SdsRequest(d, token, path, nil)
	if err != nil {
		return nil, err
	}

	var streams []community.StreamSearchResult

	err = json.Unmarshal(body, &streams)
	if err != nil {
		log.DefaultLogger.Warn("Error parsing json", err.Error())
		return nil, err
	}

	// create a dataframe
	frame := data.NewFrame("response")

	// create property lists from streams list
	ids := make([]string, len(streams))
	names := make([]string, len(streams))
	for i := 0; i < len(streams); i++ {
		ids[i] = streams[i].Self
		names[i] = streams[i].Name
	}

	// add fields
	frame.Fields = append(frame.Fields,
		data.NewField("Id", nil, ids),
		data.NewField("Name", nil, names),
	)

	return frame, nil
}

func StreamsDataQuery(d *DataHubClient, namespaceId string, token string, id string, startIndex string, endIndex string) (*data.Frame, error) {
	basePath := d.resource + "/api/" + d.apiVersion + "/tenants/" + d.tenantId + "/namespaces/" + namespaceId

	// get type Id
	path := (basePath + "/streams/" + id)
	body, err := SdsRequest(d, token, path, nil)
	if err != nil {
		return nil, err
	}

	var stream sds.SdsStream
	err = json.Unmarshal(body, &stream)
	if err != nil {
		log.DefaultLogger.Warn("Error parsing json", err.Error())
		return nil, err
	}

	// get type info
	path = (basePath + "/types/" + stream.TypeId)
	body, err = SdsRequest(d, token, path, nil)
	if err != nil {
		return nil, err
	}

	var sdsType sds.SdsType
	err = json.Unmarshal(body, &sdsType)
	if err != nil {
		log.DefaultLogger.Warn("Error parsing json", err.Error())
		return nil, err
	}

	log.DefaultLogger.Info(fmt.Sprint(sdsType))

	// get data
	path = (basePath + "/streams/" + id + "/Data?startIndex=" + startIndex + "&endIndex=" + endIndex)
	body, err = SdsRequest(d, token, path, nil)
	if err != nil {
		return nil, err
	}

	var sdsData []map[string]interface{}
	err = json.Unmarshal(body, &sdsData)
	if err != nil {
		log.DefaultLogger.Warn("Error parsing json", err.Error())
		return nil, err
	}

	// create a dataframe
	frame := data.NewFrame(stream.Name)

	// create columns in dataframe
	for i := 0; i < len(sdsType.Properties); i++ {
		typeCodeString := sdsType.Properties[i].SdsType.SdsTypeCode
		frame.Fields = append(frame.Fields,
			data.NewField(sdsType.Properties[i].Id, nil, createSdsValueList(typeCodeString)))
	}

	// add data to rows
	for i := 0; i < len(sdsData); i++ {
		row := make([]interface{}, len(sdsType.Properties))
		for j := 0; j < len(sdsType.Properties); j++ {
			row[j] = convertSdsValue(sdsType.Properties[j].SdsType.SdsTypeCode, sdsData[i][string(sdsType.Properties[j].Id)])
		}
		frame.AppendRow(row...)
	}

	err = json.Unmarshal(body, &sdsData)
	if err != nil {
		log.DefaultLogger.Warn("Error parsing json", err.Error())
		return nil, err
	}

	return frame, nil
}

func CommunityStreamsDataQuery(d *DataHubClient, communityId string, token string, self string, startIndex string, endIndex string) (*data.Frame, error) {

	// make a community header
	communityHeader := map[string]string{
		"Community-Id": communityId,
	}

	// get stream
	path := self
	body, err := SdsRequest(d, token, path, communityHeader)
	if err != nil {
		return nil, err
	}

	var stream sds.SdsStream
	err = json.Unmarshal(body, &stream)
	if err != nil {
		log.DefaultLogger.Warn("Error parsing json", err.Error())
		return nil, err
	}

	// get resolved type info
	path = (self + "/resolved")
	body, err = SdsRequest(d, token, path, communityHeader)
	if err != nil {
		return nil, err
	}

	var sdsResolvedStream sds.SdsResolvedStream
	err = json.Unmarshal(body, &sdsResolvedStream)
	if err != nil {
		log.DefaultLogger.Warn("Error parsing json", err.Error())
		return nil, err
	}

	sdsType := sdsResolvedStream.SdsType

	// get data
	path = (self + "/Data?startIndex=" + startIndex + "&endIndex=" + endIndex)
	body, err = SdsRequest(d, token, path, communityHeader)
	if err != nil {
		return nil, err
	}

	var sdsData []map[string]interface{}
	err = json.Unmarshal(body, &sdsData)
	if err != nil {
		log.DefaultLogger.Warn("Error parsing json", err.Error())
		return nil, err
	}

	// create a dataframe
	frame := data.NewFrame(stream.Name)

	// create columns in dataframe
	for i := 0; i < len(sdsType.Properties); i++ {
		typeCodeString := sdsType.Properties[i].SdsType.SdsTypeCode
		frame.Fields = append(frame.Fields,
			data.NewField(sdsType.Properties[i].Id, nil, createSdsValueList(typeCodeString)))
	}

	// add data to rows
	for i := 0; i < len(sdsData); i++ {
		row := make([]interface{}, len(sdsType.Properties))
		for j := 0; j < len(sdsType.Properties); j++ {
			row[j] = convertSdsValue(sdsType.Properties[j].SdsType.SdsTypeCode, sdsData[i][string(sdsType.Properties[j].Id)])
		}
		frame.AppendRow(row...)
	}

	err = json.Unmarshal(body, &sdsData)
	if err != nil {
		log.DefaultLogger.Warn("Error parsing json", err.Error())
		return nil, err
	}

	return frame, nil
}

func createSdsValueList(sdsTypeCode sds.SdsTypeCode) interface{} {
	switch t := sdsTypeCode; t {
	case "DateTime":
		return []time.Time{}
	case "NullableDateTime":
		return []*time.Time{}
	case "Boolean":
		return []bool{}
	case "NullableBoolean":
		return []*bool{}
	case "Int16":
		return []int16{}
	case "NullableInt16":
		return []*int16{}
	case "UInt16":
		return []uint16{}
	case "NullableUInt16":
		return []*uint16{}
	case "Int32":
		return []int32{}
	case "NullableInt32":
		return []*int32{}
	case "UInt32":
		return []uint32{}
	case "NullableUInt32":
		return []*uint32{}
	case "Int64":
		return []int64{}
	case "NullableInt64":
		return []*int64{}
	case "UInt64":
		return []uint64{}
	case "NullableUInt64":
		return []*uint32{}
	case "Single":
		return []float32{}
	case "NullableSingle":
		return []*float32{}
	case "Double":
		return []float64{}
	case "NullableDouble":
		return []*float64{}
	default:
		return []*string{}
	}
}

func convertSdsValue(sdsTypeCode sds.SdsTypeCode, value interface{}) interface{} {

	switch t := sdsTypeCode; t {
	case "DateTime":
		if value == nil {
			return value
		}
		timestamp, _ := time.Parse(time.RFC3339, value.(string))
		return timestamp
	case "NullableDateTime":
		if value == nil {
			return value
		}
		timestamp, _ := time.Parse(time.RFC3339, value.(string))
		return &timestamp
	case "Boolean":
		if value == nil {
			return false
		}
		return true
	case "NullableBoolean":
		if value == nil {
			return value
		}
		valuePointer := true
		return &valuePointer
	case "Int16":
		if value == nil {
			return int16(0)
		}
		return int16(value.(float64))
	case "NullableInt16":
		if value == nil {
			return value
		}
		valuePointer := int16(value.(float64))
		return &valuePointer
	case "UInt16":
		if value == nil {
			return uint16(0)
		}
		return uint16(value.(float64))
	case "NullableUInt16":
		if value == nil {
			return value
		}
		valuePointer := uint16(value.(float64))
		return &valuePointer
	case "Int32":
		if value == nil {
			return int32(0)
		}
		return int32(value.(float64))
	case "NullableInt32":
		if value == nil {
			return value
		}
		valuePointer := int32(value.(float64))
		return &valuePointer
	case "UInt32":
		if value == nil {
			return uint32(0)
		}
		return uint32(value.(float64))
	case "NullableUInt32":
		if value == nil {
			return value
		}
		valuePointer := uint32(value.(float64))
		return &valuePointer
	case "Int64":
		if value == nil {
			return int64(0)
		}
		return int64(value.(float64))
	case "NullableInt64":
		if value == nil {
			return value
		}
		valuePointer := int64(value.(float64))
		return &valuePointer
	case "UInt64":
		if value == nil {
			return uint64(0)
		}
		return uint64(value.(float64))
	case "NullableUInt64":
		if value == nil {
			return value
		}
		valuePointer := uint64(value.(float64))
		return &valuePointer
	case "Single":
		if value == nil {
			return float32(0)
		}
		return float32(value.(float64))
	case "NullableSingle":
		if value == nil {
			return value
		}
		valuePointer := float32(value.(float64))
		return &valuePointer
	case "Double":
		if value == nil {
			return float64(0)
		}
		return value.(float64)
	case "NullableDouble":
		if value == nil {
			return value
		}
		valuePointer := value.(float64)
		return &valuePointer
	default:
		log.DefaultLogger.Info("Default")
		if value == nil {
			return value
		}
		valuePointer := value.(string)
		return &valuePointer
	}
}
