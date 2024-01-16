import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export interface AdhQuery extends DataQuery {
  serviceId: string;
  serviceInstance: string;
  serviceRequest: string;
  urlParameters: Record<string, string>;
}

export const DEFAULT_QUERY: Partial<AdhQuery> = {
  serviceId: 'sds',
  serviceInstance: '',
  serviceRequest: '',
  urlParameters: {},
};

export interface AdhDataSourceOptions extends DataSourceJsonData {
  resource: string;
  accountId: string;
  clientId: string;
  oauthPassThru: boolean;
}

export interface AdhDataSourceSecureOptions {
  clientSecret: string;
}
