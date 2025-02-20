import { Plus } from '~/iotconsole/components/icons';
import { defer } from '@remix-run/node';
import { useLoaderData } from '@remix-run/react';
import { useState } from 'react';
import { Button } from '@kloudlite/design-system/atoms/button';
import {
  LoadingComp,
  pWrapper,
} from '~/iotconsole/components/loading-component';
import Wrapper from '~/iotconsole/components/wrapper';
import { GQLServerHandler } from '~/iotconsole/server/gql/saved-queries';
import { ensureAccountSet } from '~/iotconsole/server/utils/auth-utils';
import { getPagination, getSearch } from '~/iotconsole/server/utils/common';
import logger from '~/root/lib/client/helpers/log';
import { IRemixCtx } from '~/root/lib/types/common';
import fake from '~/root/fake-data-generator/fake';
import SecondarySubHeader from '~/iotconsole/components/secondary-sub-header';
import HandleCrCred from './handle-cr-cred';
import Tools from './tools';
import CredResourcesV2 from './cred-resources-V2';

export const loader = (ctx: IRemixCtx) => {
  const promise = pWrapper(async () => {
    ensureAccountSet(ctx);
    const { data, errors } = await GQLServerHandler(ctx.request).listCred({
      pagination: getPagination(ctx),
      search: getSearch(ctx),
    });
    if (errors) {
      logger.error(errors[0]);
      throw errors[0];
    }

    return {
      credentials: data || {},
    };
  });

  return defer({ promise });
};

const ContainerRegistryAccessManagement = () => {
  const [visible, setVisible] = useState(false);
  const { promise } = useLoaderData<typeof loader>();

  return (
    <>
      <LoadingComp
        data={promise}
        skeletonData={{
          credentials: fake.ConsoleListCredQuery.cr_listCreds as any,
        }}
      >
        {({ credentials }) => {
          const creds = credentials.edges?.map(({ node }) => node);
          return (
            <div className="flex flex-col gap-6xl">
              <SecondarySubHeader
                title="Container registry"
                action={
                  creds.length > 0 && (
                    <Button
                      content="Create new credential"
                      variant="primary"
                      onClick={() => {
                        setVisible(true);
                      }}
                    />
                  )
                }
              />
              <Wrapper
                empty={{
                  is: creds?.length === 0,
                  title: 'This is where you’ll manage your credentials.',
                  content: (
                    <p>
                      You can create a new credential and manage the listed
                      credentials.
                    </p>
                  ),
                  action: {
                    content: 'Create credential',
                    prefix: <Plus />,
                    onClick: () => {
                      setVisible(true);
                    },
                  },
                }}
                tools={<Tools />}
              >
                <CredResourcesV2 items={creds} />
              </Wrapper>
            </div>
          );
        }}
      </LoadingComp>
      <HandleCrCred {...{ isUpdate: false, visible, setVisible }} />
    </>
  );
};

export default ContainerRegistryAccessManagement;
