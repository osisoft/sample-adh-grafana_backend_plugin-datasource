import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export interface AdhQuery extends DataQuery {
  collection: string;
  queryText: string;
  id: string;
  name: string;
}

export const DEFAULT_QUERY: Partial<AdhQuery> = {
  collection: 'streams',
  queryText: '',
  id: '',
  name: '',
};

export interface AdhDataSourceOptions extends DataSourceJsonData {
  resource: string;
  apiVersion: string;
  tenantId: string;
  clientId: string;
  useCommunity: boolean;
  communityId: string;
  oauthPassThru: boolean;
  namespaceId: string;
}

export interface AdhDataSourceSecureOptions {
  clientSecret: string;
}
