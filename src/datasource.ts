import {
  CoreApp,
  DataQueryRequest,
  DataSourceInstanceSettings,
  SelectableValue,
  DataFrame,
} from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { DEFAULT_QUERY, AdhDataSourceOptions, AdhQuery } from './types';
import { lastValueFrom } from 'rxjs';
import { Dispatch, SetStateAction } from 'react';

export class DataSource extends DataSourceWithBackend<AdhQuery, AdhDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<AdhDataSourceOptions>) {
    super(instanceSettings);
  }

  async getServiceInstances(
    stateAction: Dispatch<SetStateAction<boolean | Array<SelectableValue<string>>>>
  ): Promise<Array<SelectableValue<string>>> {

  }

  async getStreams(
    query: string,
    stateAction: Dispatch<SetStateAction<boolean | Array<SelectableValue<string>>>>
  ): Promise<Array<SelectableValue<string>>> {
    const observableResponse = this.query({
      targets: [{ ...DEFAULT_QUERY, refId: 'adh-stream-autocomplete', urlParameters: { 'query': query }, serviceId: 'sds', serviceInstance: '' }],
    } as DataQueryRequest<AdhQuery>);

    const response = await lastValueFrom(observableResponse);

    if (!Array.isArray(response?.data) || response?.data.length === 0) {
      return [];
    }

    const dataFrame = response.data[0] as DataFrame;

    if (!Array.isArray(dataFrame?.fields) || dataFrame?.fields.length === 0) {
      return [];
    }

    const ids = dataFrame.fields[dataFrame.fields.findIndex((field) => field.name === 'Id')].values;
    const names = dataFrame.fields[dataFrame.fields.findIndex((field) => field.name === 'Name')].values;

    const selectables = [];
    for (let i = 0; i < ids.length; i++) {
      selectables.push({ value: ids[i], label: names[i] });
    }

    // Set state to persist selectables list
    stateAction(selectables);

    return selectables;
  }

  getDefaultQuery(_: CoreApp): Partial<AdhQuery> {
    return DEFAULT_QUERY;
  }
}
