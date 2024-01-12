import { DataSourceInstanceSettings, FieldType, MutableDataFrame } from '@grafana/data';
import { BackendSrvRequest, FetchResponse } from '@grafana/runtime';
import { Observable } from 'rxjs';
import { DataSource } from './datasource';
import { AdhDataSourceOptions } from './types';


const backendSrv = {
  fetch(options: BackendSrvRequest): Observable<FetchResponse<unknown>> {
    const edsResponse = {
      data: [
        {
          TimeStamp: '2020-01-01',
          Boolean: true,
          Number: 1,
          String: 'A',
        },
      ],
    } as FetchResponse;

    return new Observable((subscriber) => {
      subscriber.next(edsResponse);
    });
  },
};

jest.mock('@grafana/runtime', () => {
  const original = jest.requireActual('@grafana/runtime');
  return {
    ...original,
    getTemplateSrv: () => ({
      getVariables: () => [],
      replace: (s: string) => s,
    }),
    getBackendSrv: () => (backendSrv)
  };
});

describe('DataSource', () => {
  const url = 'URL';
  const resource = 'URL';
  const apiVersion = 'VERSION';
  const tenantId = 'TENANT';
  const clientId = 'CLIENT';
  const oauthPassThru = false;
  const namespaceId = 'NAMESPACE';
  const useCommunity = false;
  const communityId = 'COMMUNITY';
  const adhSettings: DataSourceInstanceSettings<AdhDataSourceOptions> = {
    id: 0,
    uid: '',
    name: '',
    access: 'direct',
    type: '',
    url,
    meta: null as any,
    jsonData: {
      resource: resource,
      apiVersion: apiVersion,
      tenantId: tenantId,
      clientId: clientId,
      oauthPassThru: oauthPassThru,
      namespaceId: namespaceId,
      useCommunity: useCommunity,
      communityId: communityId,
    },
    readOnly: false,
  };

  describe('getStreams', () => {
    it('should query for streams', (done) => {
      const datasource = new DataSource(adhSettings);

      datasource.query = jest.fn(() => {
        return new Observable((subscriber) => {
          subscriber.next({
            data: [
              new MutableDataFrame({
                refId: 'REFID',
                name: 'STREAM',
                fields: [
                  {
                    name: 'Id',
                    type: FieldType.string,
                    values: ['Id1', 'Id2'],
                  },
                  {
                    name: 'Name',
                    type: FieldType.string,
                    values: ['Name1', 'Name2'],
                  },
                ],
              }),
            ],
          });

          subscriber.complete();
        });
      });

      const results = datasource.getStreams('QUERY', () => { });

      results.then((r) => {
        expect(r).toEqual([
          { value: 'Id1', label: 'Name1' },
          { value: 'Id2', label: 'Name2' },
        ]);
        done();
      });
    });
  });
});
