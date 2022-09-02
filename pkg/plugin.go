package main

import (
	"os"

	"github.com/osisoft/sample-adh-grafana_backend_plugin-datasource/pkg/datahub"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func main() {

	if err := datasource.Manage("data-hub-backend-datasource", datahub.NewDataHubDataSource, datasource.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
