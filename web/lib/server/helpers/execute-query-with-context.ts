import { ASTNode, print } from 'graphql';
import ServerCookie from 'cookie';
import axios, { AxiosError } from 'axios';
import { uuid } from '@kloudlite/design-system/utils';
import http from 'http';
import https from 'https';
import { gatewayUrl } from '../../configs/base-url.cjs';
import {
  ICookies,
  MapType,
  IRemixHeader,
  IGqlReturn,
  NN,
} from '../../types/common';
import logger from '../../client/helpers/log';

const ignoreLogsFor = ['accounts_listAccounts', 'auth_me'];

const parseCookie = (cookieString: string) => {
  const [cookie] = cookieString.split(';');
  const [name, value] = cookie.split('=');
  return { name, value };
};

export type IExecutorResp<B = any, C = MapType<any>> = (
  variables?: C
) => Promise<IGqlReturn<NN<B>>>;

type formatter<A, B, C> = {
  transformer: (data: A) => B;
  vars?: (_: C) => void;
  k?: string;
};

export type IExecutor = <A, B, C = MapType<any>>(
  q: ASTNode,
  formatter: formatter<A, B, C>,
  def?: any
) => IExecutorResp<B, C>;

export const ExecuteQueryWithContext = (
  headers: IRemixHeader,
  cookies: ICookies = []
) => {
  return function executor<A, B, C = MapType<any>>(
    q: ASTNode,
    formatter: formatter<A, B, C>,
    def?: any
  ): IExecutorResp<B, C> {
    const logId = uuid();

    const gqlName =
      (q as any)?.definitions[0]?.selectionSet?.selections[0]?.name?.value ||
      '';

    const apiName = `[#${logId.substring(0, 5)}] ${gqlName}`;

    const res: IExecutorResp<B, C> = async (variables) => {
      const { transformer } = formatter;

      try {
        console.time(apiName);
        const defCookie =
          headers.get('klsession') || headers.get('cookie') || null;

        const cookie = ServerCookie.parse(defCookie || '');

        if (cookies.length > 0) {
          for (let i = 0; i < cookies.length; i += 1) {
            const { name, value } = parseCookie(cookies[i]);
            cookie[name] = value;
          }
        }

        axios.defaults.httpAgent = new http.Agent({ keepAlive: true });
        axios.defaults.httpsAgent = new https.Agent({ keepAlive: true });

        const resp = await axios({
          url: gatewayUrl,
          method: 'POST',
          headers: {
            'Content-Type': 'application/json; charset=utf-8',
            ...{
              cookie: Object.entries(cookie)
                .map(([key, value]) => `${key}=${value}`)
                .join('; '),
            },
          },
          data: {
            query: print(q),
            variables: variables || {},
          },
          timeout: 15000,
        });

        let { data } = resp.data;
        const { errors } = resp.data;

        if (errors) {
          const e = errors as Error[];
          if (e.length === 1) {
            throw errors[0];
          }

          throw new Error(
            e.reduce((acc, curr) => {
              return `${acc}\n\n1. ${curr.name ? `${curr.name}:` : ''}:${
                curr.message
              }${curr.stack ? `\n${curr.stack}` : ''}`;
            }, 'Errors:')
          );
        }

        if (data) {
          data = transformer(data);
        } else if (def) {
          data = def;
        }

        if (resp.headers && resp.headers['set-cookie']) {
          return { ...resp.data, data, cookie: resp.headers['set-cookie'] };
        }
        return { ...resp.data, data };
      } catch (err) {
        if (!ignoreLogsFor.includes(gqlName)) {
          if ((err as AxiosError).response) {
            logger.log('\nErrorIn:', apiName, (err as Error).name, '\n');

            return (err as AxiosError).response?.data;
          }

          logger.log('\nErrorIn:', apiName, (err as Error).message, '\n');
        }

        return {
          data: null,
          errors: [
            {
              message: (err as Error).message,
              name: (err as Error).name,
              stack: (err as Error).stack,
            },
          ],
        };
      } finally {
        console.timeEnd(apiName);
      }
    };

    // @ts-ignore
    res.astNode = q;
    return res;
  };
};
