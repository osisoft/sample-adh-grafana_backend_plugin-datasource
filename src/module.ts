import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { SdsQuery, SdsDataSourceOptions, SdsDataSourceSecureOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, SdsQuery, SdsDataSourceOptions, SdsDataSourceSecureOptions>(
  DataSource
)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
