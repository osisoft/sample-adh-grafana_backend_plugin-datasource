package datahub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var (
	_ backend.QueryDataHandler      = (*DataHubDataSource)(nil)
	_ backend.CheckHealthHandler    = (*DataHubDataSource)(nil)
	_ instancemgmt.InstanceDisposer = (*DataHubDataSource)(nil)
)

type DataHubDataSource struct {
	dataHubClient *DataHubClient
	oauthPassThru bool
}

type DataHubDataSourceOptions struct {
	Resource      string `json:"resource"`
	AccountId     string `json:"accountId"`
	ClientId      string `json:"clientId"`
	OauthPassThru bool   `json:"oauthPassThru"`
}

type ServiceId string

const (
	Sds ServiceId = "sds"
)

type ServiceRequest string

const (
	ServiceInstances ServiceRequest = "serviceInstances"
	Streams          ServiceRequest = "streams"
	StreamData       ServiceRequest = "streamData"
)

type QueryModel struct {
	ServiceId       ServiceId         `json:"serviceId"`
	ServiceInstance string            `json:"serviceInstance"`
	ServiceRequest  ServiceRequest    `json:"serviceRequest"`
	UrlParameters   map[string]string `json:"urlParameters"`
}

type CheckHealthResponseBody struct {
	Id string `json:"Id"`
}

// Creates a new datasource instance.
func NewDataHubDataSource(_ context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	// Get JSON Data to read datasource settings
	var options DataHubDataSourceOptions
	err := json.Unmarshal(settings.JSONData, &options)
	if err != nil {
		log.DefaultLogger.Warn("error marshalling", "err", err)
		return nil, err
	}

	// Read Client Secret from Secure JSON Data
	var secureData = settings.DecryptedSecureJSONData
	clientSecret := secureData["clientSecret"]

	client := NewDataHubClient(options.Resource, options.AccountId, options.ClientId, clientSecret)
	return &DataHubDataSource{
		dataHubClient: &client,
		oauthPassThru: options.OauthPassThru,
	}, nil
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using the new instance factory function.
func (d *DataHubDataSource) Dispose() {
	// Clean up datasource instance resources.
}

// Handles multiple queries and returns multiple responses.
func (d *DataHubDataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	log.DefaultLogger.Info("QueryData called", "request", req)

	// retrieve token
	var token string
	if d.oauthPassThru {
		token = req.Headers["Authorization"]
		if len(token) == 0 {
			return nil, fmt.Errorf("unable to retrieve token")
		}
	} else {
		var err error
		token, err = GetClientToken(d.dataHubClient)
		if err != nil {
			log.DefaultLogger.Warn("Unable to retrieve token", err.Error())
			return nil, err
		}
	}

	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res, err := d.query(ctx, req.PluginContext, q, token)

		if err != nil {
			return nil, err
		}

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

// Handles the individual queries from QueryData.
func (d *DataHubDataSource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery, token string) (backend.DataResponse, error) {
	log.DefaultLogger.Info("Running query", "query", query)
	response := backend.DataResponse{}

	// unmarshal the JSON into our QueryModel.
	var qm QueryModel

	response.Error = json.Unmarshal(query.JSON, &qm)
	if response.Error != nil {
		return response, nil
	}

	// determine what type of query to use
	frame := data.NewFrame("response")
	var err error
	switch qm.ServiceRequest {
	case ServiceRequest(ServiceInstances):
		log.DefaultLogger.Debug("Service instances query")
		frame, err = ServiceInstanceQuery(d.dataHubClient, token)
	case ServiceRequest(Streams):
		log.DefaultLogger.Debug("Streams query")
		frame, err = StreamsQuery(d.dataHubClient, qm.ServiceInstance, token, qm.UrlParameters["query"])
	case ServiceRequest(StreamData):
		log.DefaultLogger.Debug("Stream data query")
		frame, err = StreamsDataQuery(d.dataHubClient,
			qm.ServiceInstance,
			token,
			qm.UrlParameters["id"],
			query.TimeRange.From.Format(time.RFC3339),
			query.TimeRange.To.Format(time.RFC3339))
	}

	// add the frames to the response.
	response.Frames = append(response.Frames, frame)

	log.DefaultLogger.Info("We made it")

	return response, err
}

// Handles health checks sent from Grafana to the plugin.
func (d *DataHubDataSource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	log.DefaultLogger.Info("CheckHealth called", "request", req)

	var status = backend.HealthStatusOk
	var message = "Data source is working"

	// Retrieve token
	var token string
	if d.oauthPassThru {
		return &backend.CheckHealthResult{
			Status:  status,
			Message: message,
		}, nil
	} else {
		var err error
		token, err = GetClientToken(d.dataHubClient)
		if err != nil {
			log.DefaultLogger.Warn("Error unable to get token health check", err.Error())
			return &backend.CheckHealthResult{
				Status:  backend.HealthStatusError,
				Message: "Unable to retrieve token",
			}, nil
		}
	}

	// Make a request to test the token
	_, err := ServiceInstanceQuery(d.dataHubClient, token)
	if err != nil {
		log.DefaultLogger.Warn("Error test request health check", err.Error())
		status = backend.HealthStatusError
		message = "Invalid Configuration"
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}
