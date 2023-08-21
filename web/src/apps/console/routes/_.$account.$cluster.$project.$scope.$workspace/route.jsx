import { Outlet, useOutletContext, useLoaderData } from '@remix-run/react';
import OptionList from '~/components/atoms/option-list';
import { ChevronDown, Plus, Search } from '@jengaicons/react';
import Breadcrum from '~/console/components/breadcrum';
import { useState } from 'react';
import {
  BlackProdLogo,
  BlackWorkspaceLogo,
} from '~/console/components/commons';
import { GQLServerHandler } from '~/console/server/gql/saved-queries';
import logger from '~/root/lib/client/helpers/log';
import { HandlePopup } from './handle-wrkspc-env';

const Workspace = () => {
  const rootContext = useOutletContext();
  const { workspace } = useLoaderData();

  // @ts-ignore
  return <Outlet context={{ ...rootContext, workspace }} />;
};

export default Workspace;

export const handle = ({ account, project, cluster, scope }) => {
  return {
    navbar: {
      backurl: {
        href: `/${account}/${cluster}/${project}/${
          scope === 'workspace' ? 'workspaces' : 'environments'
        }`,
        name: scope === 'workspace' ? 'Workspaces' : 'Environments',
      },
      items: [
        {
          label: 'Apps',
          to: '/apps',
          key: 'apps',
          value: '/apps',
        },
        {
          label: 'Routers',
          to: '/routers',
          key: 'routers',
          value: '/routers',
        },
        {
          label: 'Config & Secrets',
          to: '/cs',
          key: 'config-and-secrets',
          value: '/cs',
        },
        {
          label: 'Backing services',
          to: '/backing-services',
          key: 'backing-services',
          value: '/backing-services',
        },
        {
          label: 'Settings',
          to: '/settings',
          key: 'settings',
          value: '/settings',
        },
      ],
    },
    breadcrum: () => <CurrentBreadcrum />,
  };
};

const CurrentBreadcrum = () => {
  const [showPopup, setShowPopup] = useState(null);
  const [activeTab, setActiveTab] = useState('environments');
  return (
    <>
      <OptionList.Root>
        <OptionList.Trigger>
          <Breadcrum.Button
            content="Staging"
            prefix={BlackProdLogo}
            suffix={ChevronDown}
          />
        </OptionList.Trigger>
        <OptionList.Content className="!pt-0 !pb-md" align="center">
          <div className="p-[3px] pb-0">
            <OptionList.TextInput
              value=""
              prefixIcon={Search}
              placeholder="Search"
              compact
              className="border-0 rounded-none"
            />
          </div>
          <OptionList.Separator />
          <OptionList.Tabs.Root
            value={activeTab}
            size="sm"
            className="!overflow-x-visible"
            onChange={setActiveTab}
            // LinkComponent={Link}
          >
            <OptionList.Tabs.Tab
              prefix={BlackProdLogo}
              label="Environments"
              value="environments"
            />
            <OptionList.Tabs.Tab
              prefix={BlackWorkspaceLogo}
              label="Workspaces"
              value="workspaces"
            />
          </OptionList.Tabs.Root>
          <OptionList.Item>Staging</OptionList.Item>
          <OptionList.Item>Hustle</OptionList.Item>
          <OptionList.Item>Visionary</OptionList.Item>
          <OptionList.Separator />
          <OptionList.Item
            className="text-text-primary"
            onSelect={() => setShowPopup({ type: activeTab })}
          >
            <Plus size={16} />{' '}
            <span>
              {activeTab === 'workspaces' ? 'New Workspace' : 'New Environment'}
            </span>
          </OptionList.Item>
        </OptionList.Content>
      </OptionList.Root>
      <HandlePopup show={showPopup} setShow={setShowPopup} />
    </>
  );
};

export const loader = async (ctx) => {
  const { account, cluster, project, workspace, scope } = ctx.params;

  const api =
    scope === 'workspace'
      ? GQLServerHandler(ctx.request).getWorkspace
      : GQLServerHandler(ctx.request).getEnvironment;

  const { data, errors } = await api({
    project: {
      value: project,
      type: 'name',
    },
    name: workspace,
  });
  if (errors) {
    logger.error(errors);
  }

  return {
    baseurl: `/${account}/${cluster}/${project}/${scope}/${workspace}`,
    workspace: data || {},
    account,
    project,
    cluster,
    scope,
  };
};
