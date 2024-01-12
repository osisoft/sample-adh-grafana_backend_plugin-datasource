import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { AdhQuery, AdhDataSourceOptions, AdhDataSourceSecureOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, AdhQuery, AdhDataSourceOptions, AdhDataSourceSecureOptions>(
  DataSource
)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
