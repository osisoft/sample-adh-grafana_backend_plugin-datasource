# AVEVA Data Hub Grafana Data Source Backend Plugin Sample

**Version:** 1.0.3

[![Build Status](https://dev.azure.com/osieng/engineering/_apis/build/status/product-readiness/ADH/osisoft.sample-adh-grafana_backend_plugin-datasource?repoName=osisoft%2Fsample-adh-grafana_backend_plugin-datasource&branchName=main)](https://dev.azure.com/osieng/engineering/_build/latest?definitionId=4858&repoName=osisoft%2Fsample-adh-grafana_backend_plugin-datasource&branchName=main)

This sample demonstrates how to build a [Grafana](https://grafana.com/) data source backend plugin that runs queries against AVEVA Data Hub (ADH). The sample performs normal "Get Values" calls against a specified stream in SDS, using the time range of the Grafana dashboard. For more information about backend plugins, refer to the documentation on [Backend plugins](https://grafana.com/docs/grafana/latest/developers/plugins/backend/).

## Requirements

- [Grafana 8.3+](https://grafana.com/grafana/download)
- Web Browser with JavaScript enabled
- [NodeJS](https://nodejs.org/en/)
- [Go](https://go.dev/)
- [Mage](https://magefile.org/)
- [Git](https://git-scm.com/download/win)
- If building the plugin, a Linux OS or WSL for Windows is required. Refer to [grafana's troubleshooting documentation](https://grafana.com/developers/plugin-tools/troubleshooting) for more information. 
- If not using OAuth passthrough, register a Client Credentials Client in AVEVA Data Hub; a client secret will need to be provided to the sample plugin configuration

## Getting Started

If you plan to run this sample plugin on an existing Grafana server, the recommended sections are:
1. [Building and Deploying the Sample Plugin on a Grafana Server](#building-and-deploying-the-sample-plugin-on-a-grafana-server)
2. If you plan to use ADH OAuth [Using ADH OAuth login to Grafana](#using-adh-oauth-login-to-grafana)
3. [Configuring the Sample Plugin in Grafana](#configuring-the-sample-plugin-in-grafana)

If you wish to deploy the sample plugin with a docker container also containing Grafana, then the recommended sections are:
1. [Deploying the Sample Plugin with Docker](#deploying-the-sample-plugin-with-docker)
1. [Configuring the Sample Plugin in Grafana](#configuring-the-sample-plugin-in-grafana)

If you are a developer that wants to customize the plugin or add functionality to it, then the recommend sections are:
1. [Developing the Sample with Docker](#deploying-the-sample-plugin-with-docker)
1. [Configuring the Sample Plugin in Grafana](#configuring-the-sample-plugin-in-grafana)

## Building and Deploying the Sample Plugin on a Grafana Server

1. Copy this folder to your Grafana server's plugins directory, like `.../grafana/data/plugins`
1. (Optional) If using other plugins, rename the folder to `aveva-adh-datasource`
1. Open a command prompt inside that folder
1. Install dependencies, using `npm install`
1. Build the plugin, using `npm run build` (or `npm run dev` for browser debugging)
1. Update Grafana plugin SDK for Go dependency to the latest minor version, using `go get -u github.com/grafana/grafana-plugin-sdk-go` and `go mod tidy`
1. Build backend plugin binaries for Linux, Windows, and Darwin, using `mage -v`
1. Open the Grafana configuration and set the parameter `allow_loading_unsigned_plugins` equal to `aveva-adh-datasource` or to the name of the folder set in step 2 (see [Grafana docs](https://grafana.com/docs/grafana/latest/administration/configuration/#allow_loading_unsigned_plugins))
1. Restart the Grafana server to load the new plugin

## Deploying the Sample Plugin with Docker

1. Open a command prompt inside this folder
1. Build the container using `docker build -t grafana-adh .`  
   _Note: The dockerfile being built contains an ENV statement that creates an [environment variable](https://grafana.com/docs/grafana/latest/administration/configuration/#configure-with-environment-variables) that overrides an option in the grafana config. In this case, the `allow_loading_unsigned_plugins` option is being overridden to allow the [unsigned plugin](https://grafana.com/docs/grafana/latest/administration/configuration/#allow_loading_unsigned_plugins) in this sample to be used._
1. Run the container using `docker run -d --name=grafana-adh -p 3000:3000 grafana-adh`
1. Navigate to localhost:3000 to configure data sources and view data

## Configuring the Sample Plugin in Grafana

1. Add a new Grafana datasource using the sample (see [Grafana docs](https://grafana.com/docs/grafana/latest/features/datasources/add-a-data-source/))
1. Enter the relevant required information; if using a client credentials client, the client secret will be encrypted in the Grafana server and HTTP requests to ADH will be made by a server-side proxy, as described in the [Grafana docs](https://grafana.com/docs/grafana/latest/developers/plugins/authentication/)
1. Open a new or existing Grafana dashboard, and choose `AVEVA Data Hub (SAMPLE)` as the data source
1. Enter your Namespace (if querying ADH) and Stream, and data will populate into the dashboard from the stream for the dashboard's time range

## Using Community Data

1. Add a new Grafana datasource using the sample (see [Grafana docs](https://grafana.com/docs/grafana/latest/features/datasources/add-a-data-source/))
1. Choose AVEVA Data Hub
1. Toggle the "Community Data" switch to 'true'
1. Enter the relevant required information. You can find the Community ID in the URL of the Community Details page.

## Using ADH OAuth login to Grafana

To use AVEVA Data Hub as an Identity provider through OAuth, add the following generic OAuth configuration to your grafana server's custom.ini. Please note, you may need to create a new custom.ini if one does not already exist and an Authorization Code Client with the appropriate redirect URLs will need to be generated. For more information please refer to Grafana's [configuration documentation](https://grafana.com/docs/grafana/latest/setup-grafana/configure-grafana/) or their [Generic OAuth documentation](https://grafana.com/docs/grafana/latest/auth/generic-oauth/).

```ini
[auth.generic_oauth]
enabled = true
name = AVEVA Data Hub
allow_sign_up = true
client_id = <PLACEHOLDER_CLIENT_ID>
scopes = openid profile email ocsapi offline_access
auth_url = https://uswe.datahub.connect.aveva.com/identity/connect/authorize
token_url = https://uswe.datahub.connect.aveva.com/identity/connect/token
api_url = https://uswe.datahub.connect.aveva.com/identity/connect/userinfo
role_attribute_path = contains(role_type[*], '2dc742ab-39ea-4fc0-a39e-2bcb71c26a5f') && 'Admin' || contains(role_type[*], 'f1439595-e5a2-487f-8a4f-0627fefe75df') && 'Editor' || 'Viewer'
use_pkce = true
```

| Setting             | Description                                                                                                                                                                                                                                                                                                                                                    |
| ------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| enabled             | Whether generic OAuth is enabled.                                                                                                                                                                                                                                                                                                                              |
| name                | The name of the Identity Provider. This is also what is shown on the button when prompted for login.                                                                                                                                                                                                                                                           |
| allow_sign_up       | This setting allows Grafana users to be automatically created upon login. With this set to false, an administrator would have to create an account within Grafana for a user before said user could access Grafana.                                                                                                                                            |
| client_id           | The Authorization Code Client Id. By default, refresh tokens are not issued and the token lifetime is 1 hour. To enable refresh tokens and allow the token to be refreshed for up to 8 hours, AllowOfflineAccess must be set to true on the client's configuration, which can be set within the API console.                                                   |
| scopes              | The OAuth scopes to be designate what access the application should have to the user’s account. OpenId, Profile, and Email are used to gather information about the user and determine what their role should be if role_attribute_path is specified. Ocsapi gives the user access to the AVEVA Data Hub API. Offline_access is used to enable refresh tokens. |
| auth_url            | The well-known authorization URL of AVEVA Data Hub (may depend on region).                                                                                                                                                                                                                                                                                     |
| token_url           | The well-known token URL of AVEVA Data Hub (may depend on region).                                                                                                                                                                                                                                                                                             |
| api_url             | The well-known user information URL of AVEVA Data Hub (may depend on region).                                                                                                                                                                                                                                                                                  |
| role_attribute_path | Defines how roles are mapped between AVEVA Data Hub and Grafana.                                                                                                                                                                                                                                                                                               |
| use_pkce            | Enables and forces Grafana to use PKCE.                                                                                                                                                                                                                                                                                                                        |

## Developing the Sample with Docker

1. Open a command prompt inside this folder
1. Install dependencies, using `npm install`
1. Build the plugin for development, using `npm run dev`
1. In a new terminal window, update Grafana plugin SDK for Go dependency to the latest minor version, using `go get -u github.com/grafana/grafana-plugin-sdk-go` and `go mod tidy`
1. Build backend plugin binaries for Linux `mage -v build:Linux`
1. Start grafana by running `docker-compose up --build`
1. Open Grafana, by default `http://localhost:3000/`

## Running the Automated Tests on Frontend Components

1. Open a command prompt inside this folder
1. Install dependencies, using `npm ci`
1. Run the tests, using `npm test`

## Running the Automated Tests on Backend Components

1. Open a command prompt inside this folder
1. Install dependencies, using `go mod tidy`
1. Run the tests, using `go test`

---

For the main ADH page [ReadMe](https://github.com/osisoft/OSI-Samples-OCS)  
For the main samples page [ReadMe](https://github.com/osisoft/OSI-Samples)
