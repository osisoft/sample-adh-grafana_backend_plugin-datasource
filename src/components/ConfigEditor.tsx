import {
  DataSourcePluginOptionsEditorProps,
  onUpdateDatasourceJsonDataOption,
  onUpdateDatasourceJsonDataOptionChecked,
  onUpdateDatasourceSecureJsonDataOption
} from '@grafana/data';
import { InlineField, InlineFieldRow, InlineSwitch, Input } from '@grafana/ui';
import React, { SyntheticEvent } from 'react';
import { AdhDataSourceOptions, AdhDataSourceSecureOptions } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<AdhDataSourceOptions, AdhDataSourceSecureOptions> { }

export const ConfigEditor = (props: Props) => {

  const warningStyle = {
    color: 'orange',
    alignSelf: 'center',
  };

  const onResetClientSecret = (event: SyntheticEvent) => {
    event.preventDefault();
    const { onOptionsChange, options } = props;
    const secureJsonData = {
      ...options.secureJsonData,
      clientSecret: '',
    };
    const secureJsonFields = {
      ...options.secureJsonFields,
      clientSecret: false,
    };
    onOptionsChange({ ...options, secureJsonData, secureJsonFields });
  };

  const { options } = props;
  const { jsonData, secureJsonData } = options;

  // Fill in defaults
  if (!jsonData.resource) {
    jsonData.resource = 'https://uswe.datahub.connect.aveva.com';
  }
  if (!jsonData.apiVersion) {
    jsonData.apiVersion = 'v1';
  }
  if (jsonData.oauthPassThru == null) {
    jsonData.oauthPassThru = false;
  }

  return (
    <div>
      <div className="gf-form-group">
        <h3 className="page-heading">AVEVA Data Hub</h3>
        <InlineField label="URL" tooltip="The URL for AVEVA Data Hub" labelWidth={20}>
          <Input
            required={true}
            placeholder="https://uswe.datahub.connect.aveva.com"
            width={40}
            onChange={onUpdateDatasourceJsonDataOption(props, 'resource')}
            value={jsonData.resource || ''}
          />
        </InlineField>
        <InlineField label="API Version" tooltip="The version of the ADH API to use" labelWidth={20}>
          <Input
            required={true}
            placeholder="v1"
            width={40}
            onChange={onUpdateDatasourceJsonDataOption(props, 'apiVersion')}
            value={jsonData.apiVersion || ''}
          />
        </InlineField>
        <InlineField label="Tenant ID" tooltip="The ID of your AVEVA Data Hub tenant" labelWidth={20}>
          <Input
            required={true}
            placeholder="00000000-0000-0000-0000-000000000000"
            width={40}
            onChange={onUpdateDatasourceJsonDataOption(props, 'tenantId')}
            value={jsonData.tenantId || ''}
          />
        </InlineField>
        <InlineFieldRow>
          <InlineField
            label="Community Data"
            tooltip="Switch to toggle reading from a Namespace to a Community"
            labelWidth={20}
          >
            <InlineSwitch
              onChange={onUpdateDatasourceJsonDataOptionChecked(props, 'useCommunity')}
              value={jsonData.useCommunity}
            />
          </InlineField>
        </InlineFieldRow>
        {jsonData.useCommunity && (
          <InlineField label="Community ID" tooltip="The ID of the AVEVA Data Hub Community" labelWidth={20}>
            <Input
              placeholder="00000000-0000-0000-0000-000000000000"
              width={40}
              onChange={onUpdateDatasourceJsonDataOption(props, 'communityId')}
              value={jsonData.communityId || ''}
            />
          </InlineField>
        )}
        {!jsonData.useCommunity && (
          <InlineField label="Namespace ID" tooltip="The Namespace in your for AVEVA Data Hub tenant" labelWidth={20}>
            <Input
              required={true}
              placeholder="Enter a Namespace ID..."
              width={40}
              onChange={onUpdateDatasourceJsonDataOption(props, 'namespaceId')}
              value={jsonData.namespaceId || ''}
            />
          </InlineField>
        )}
        <InlineFieldRow>
          <InlineField label="Use OAuth token" tooltip="Switch to toggle authentication modes" labelWidth={20}>
            <InlineSwitch
              onChange={onUpdateDatasourceJsonDataOptionChecked(props, 'oauthPassThru')}
              value={jsonData.oauthPassThru}
            />
          </InlineField>
          {jsonData.oauthPassThru && (
            <div style={warningStyle}>
              Warning: Requires configuring genenric OAuth with AVEVA Data Hub in your Grafana Server
            </div>
          )}
        </InlineFieldRow>
        {!jsonData.oauthPassThru && (
          <InlineField
            label="Client ID"
            tooltip="The ID of the Client Credentials client to authenticate against your ADH tenant"
            labelWidth={20}
          >
            <Input
              placeholder="00000000-0000-0000-0000-000000000000"
              width={40}
              onChange={onUpdateDatasourceJsonDataOption(props, 'clientId')}
              value={jsonData.clientId || ''}
            />
          </InlineField>
        )}
        {!jsonData.oauthPassThru && (
          <InlineField
            label="Client Secret"
            tooltip="The secret for the specified Client Credentials client"
            labelWidth={20}
          >
            <Input
              required={true}
              type="password"
              placeholder="Enter a Client secret..."
              width={40}
              onChange={onUpdateDatasourceSecureJsonDataOption(props, 'clientSecret')}
              onReset={onResetClientSecret}
              value={secureJsonData?.clientSecret || ''}
            />
          </InlineField>
        )}
      </div>
    </div>
  );
};
