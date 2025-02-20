import Wrapper from '~/iotconsole/components/wrapper';
import { useLoaderData } from '@remix-run/react';
import { IRemixCtx } from '~/lib/types/common';
import { ensureAccountSet } from '~/iotconsole/server/utils/auth-utils';
import {
  LoadingComp,
  pWrapper,
} from '~/iotconsole/components/loading-component';
import { GQLServerHandler } from '~/iotconsole/server/gql/saved-queries';
import { base64Decrypt, getPagination } from '~/iotconsole/server/utils/common';
import { defer } from '@remix-run/node';
import fake from '~/root/fake-data-generator/fake';
import Tools from './tools';
import BuildRunResources from './buildruns-resources';

export const loader = async (ctx: IRemixCtx) => {
  ensureAccountSet(ctx);

  const promise = pWrapper(async () => {
    const { data, errors } = await GQLServerHandler(ctx.request).listBuildRuns({
      pq: getPagination(ctx),
      search: {
        repoName: {
          exact: base64Decrypt(ctx.params.repo),
          matchType: 'exact',
        },
      },
    });
    if (errors) {
      throw errors[0];
    }
    return { buildRunData: data };
  });

  return defer({ promise });
};

const BuildRuns = () => {
  const { promise } = useLoaderData<typeof loader>();

  return (
    <LoadingComp
      data={promise}
      skeletonData={{
        buildRunData: fake.ConsoleListBuildRunsQuery.cr_listBuildRuns as any,
      }}
    >
      {({ buildRunData }) => {
        const buildruns = buildRunData?.edges?.map(({ node }) => node);
        if (!buildruns) {
          return null;
        }
        const { pageInfo, totalCount } = buildRunData;
        return (
          <Wrapper
            header={{
              title: 'Build Runs',
            }}
            empty={{
              is: buildruns.length === 0,
              title: 'This is where you’ll manage your buildruns',
              content: '',
            }}
            pagination={{
              pageInfo,
              totalCount,
            }}
            tools={<Tools />}
          >
            <BuildRunResources items={buildruns} />
          </Wrapper>
        );
      }}
    </LoadingComp>
  );
};

export default BuildRuns;
