import { LockSimple, Trash } from '~/console/components/icons';
import { generateKey, titleCase } from '~/components/utils';
import {
  ListItem,
  ListItemV2,
  ListTitle,
  ListTitleV2,
  listClass,
} from '~/console/components/console-list-components';
import Grid from '~/console/components/grid';
import ListGridView from '~/console/components/list-grid-view';
import {
  ExtractNodeType,
  parseName,
  parseUpdateOrCreatedBy,
  parseUpdateOrCreatedOn,
} from '~/console/server/r-utils/common';
import DeleteDialog from '~/console/components/delete-dialog';
import ResourceExtraAction from '~/console/components/resource-extra-action';
import { useConsoleApi } from '~/console/server/gql/api-provider';
import { useReload } from '~/lib/client/helpers/reloader';
import { useState } from 'react';
import { handleError } from '~/lib/utils/common';
import { toast } from '~/components/molecule/toast';
import { useOutletContext, useParams } from '@remix-run/react';
import { useWatchReload } from '~/lib/client/helpers/socket/useWatch';
import ListV2 from '~/console/components/listV2';
import { IMSvTemplates } from '~/console/server/gql/queries/managed-templates-queries';
import { getManagedTemplateLogo } from '~/console/utils/commons';
import { IImportedManagedResources } from '~/console/server/gql/queries/imported-managed-resource-queries';
import { Badge } from '~/components/atoms/badge';
import useClusterStatus from '~/console/hooks/use-cluster-status';
import { ViewSecret } from './handle-managed-resource-v2';
import { IEnvironmentContext } from '../_layout';

const RESOURCE_NAME = 'integrated resource';
type BaseType = ExtractNodeType<IImportedManagedResources>;

const parseItem = (item: BaseType, templates: IMSvTemplates) => {
  const logoUrl = getManagedTemplateLogo(
    templates,
    item.managedResource?.spec?.resourceTemplate.apiVersion || ''
  );
  return {
    name: item?.displayName,
    id: item?.name,
    updateInfo: {
      author: `Updated by ${titleCase(parseUpdateOrCreatedBy(item))}`,
      time: parseUpdateOrCreatedOn(item),
    },
    logo: logoUrl,
  };
};

type OnAction = ({
  action,
  item,
}: {
  action: 'delete' | 'edit' | 'view_secret';
  item: BaseType;
}) => void;

type IExtraButton = {
  onAction: OnAction;
  item: BaseType;
};

const ExtraButton = ({ onAction, item }: IExtraButton) => {
  return (
    <ResourceExtraAction
      options={[
        // {
        //   label: 'Edit',
        //   icon: <PencilSimple size={16} />,
        //   type: 'item',
        //   onClick: () => onAction({ action: 'edit', item }),
        //   key: 'edit',
        // },
        {
          label: 'View Secret',
          icon: <LockSimple size={16} />,
          type: 'item',
          onClick: () => onAction({ action: 'view_secret', item }),
          key: 'view_secret',
        },
        {
          label: 'Delete',
          icon: <Trash size={16} />,
          type: 'item',
          onClick: () => onAction({ action: 'delete', item }),
          key: 'delete',
          className: '!text-text-critical',
        },
      ]}
    />
  );
};

interface IResource {
  items: BaseType[];
  onAction: OnAction;
  templates: IMSvTemplates;
}

const GridView = ({ items = [], onAction, templates }: IResource) => {
  return (
    <Grid.Root className="!grid-cols-1 md:!grid-cols-3">
      {items.map((item, index) => {
        const { name, id, updateInfo } = parseItem(item, templates);
        const keyPrefix = `${RESOURCE_NAME}-${id}-${index}`;
        return (
          <Grid.Column
            key={id}
            rows={[
              {
                key: generateKey(keyPrefix, name),
                render: () => (
                  <ListTitle
                    title={name}
                    subtitle={id}
                    action={<ExtraButton onAction={onAction} item={item} />}
                  />
                ),
              },
              {
                key: generateKey(keyPrefix, 'author'),
                render: () => (
                  <ListItem
                    data={updateInfo.author}
                    subtitle={updateInfo.time}
                  />
                ),
              },
            ]}
          />
        );
      })}
    </Grid.Root>
  );
};

const ListView = ({ items = [], onAction, templates }: IResource) => {
  const { cluster } = useOutletContext<IEnvironmentContext>();
  const { findClusterStatus, clusters, loading } = useClusterStatus();

  return (
    <ListV2.Root
      data={{
        headers: [
          {
            render: () => 'Resource Name',
            name: 'name',
            className: listClass.title,
          },
          {
            render: () => 'Resource Type',
            name: 'resource',
            className: listClass.item,
          },
          {
            render: () => '',
            name: 'flex-pre',
            className: listClass.flex,
          },
          {
            render: () => 'Integrated Service',
            name: 'service',
            className: 'w-[175px]',
          },
          {
            render: () => '',
            name: 'flex-post',
            className: listClass.flex,
          },
          {
            render: () => 'Status',
            name: 'status',
            className: listClass.status,
          },
          {
            render: () => 'Updated',
            name: 'updated',
            className: listClass.updated,
          },
          {
            render: () => '',
            name: 'action',
            className: listClass.action,
          },
        ],
        rows: items.map((i) => {
          const { name, id, logo, updateInfo } = parseItem(i, templates);
          const isClusterOnline = findClusterStatus(
            clusters.length > 0
              ? clusters.find((c) => parseName(c) === parseName(cluster))
              : cluster
          );

          return {
            columns: {
              name: {
                render: () => <ListTitleV2 title={name} subtitle={id} />,
              },
              resource: {
                render: () => (
                  <ListItemV2
                    data={`${i.managedResource?.spec?.resourceTemplate?.kind}`}
                  />
                ),
              },
              service: {
                render: () => (
                  <ListItemV2
                    pre={
                      <div className="pulsable">
                        <img
                          src={logo}
                          alt={`${i.managedResource?.spec?.resourceTemplate?.msvcRef?.name}`}
                          className="w-4xl h-4xl"
                        />
                      </div>
                    }
                    data={
                      i.managedResource?.spec?.resourceTemplate?.msvcRef
                        ?.name || ''
                    }
                  />
                ),
              },
              status: {
                render: () => {
                  if (loading) {
                    return null;
                  }

                  if (!isClusterOnline) {
                    return <Badge type="warning">Cluster Offline</Badge>;
                  }

                  if (i.syncStatus?.state === 'UPDATED_AT_AGENT') {
                    return <Badge type="info">Ready</Badge>;
                  }

                  return <Badge type="warning">Waiting</Badge>;
                },
              },
              updated: {
                render: () => (
                  <ListItemV2
                    data={`${updateInfo.author}`}
                    subtitle={updateInfo.time}
                  />
                ),
              },
              action: {
                render: () => <ExtraButton item={i} onAction={onAction} />,
              },
            },
          };
        }),
      }}
    />
  );
};

const ManagedResourceResourcesV2 = ({
  items = [],
  templates = [],
}: {
  items: BaseType[];
  templates: IMSvTemplates;
}) => {
  const [showDeleteDialog, setShowDeleteDialog] = useState<BaseType | null>(
    null
  );
  const [showSecret, setShowSecret] = useState<BaseType | null>(null);
  const api = useConsoleApi();
  const reloadPage = useReload();
  const params = useParams();

  const { environment, account } = useParams();

  useWatchReload(
    items.map((i) => {
      return `account:${account}.environment:${environment}.managed_resource:${i.name}`;
    })
  );

  const props: IResource = {
    items,
    onAction: ({ action, item }) => {
      switch (action) {
        case 'delete':
          setShowDeleteDialog(item);
          break;
        case 'view_secret':
          setShowSecret(item);
          break;
        default:
          break;
      }
    },
    templates,
  };
  return (
    <>
      <ListGridView
        listView={<ListView {...props} />}
        gridView={<GridView {...props} />}
      />
      <DeleteDialog
        resourceName={showDeleteDialog?.name || ''}
        resourceType={RESOURCE_NAME}
        show={showDeleteDialog}
        setShow={setShowDeleteDialog}
        onSubmit={async () => {
          if (!params.environment) {
            throw new Error('Environment is required!.');
          }
          try {
            const { errors } = await api.deleteImportedManagedResource({
              importName: showDeleteDialog?.name || '',
              envName: environment || '',
            });

            if (errors) {
              throw errors[0];
            }
            reloadPage();
            toast.success(`${titleCase(RESOURCE_NAME)} deleted successfully`);
            setShowDeleteDialog(null);
          } catch (err) {
            handleError(err);
          }
        }}
      />

      {showSecret && (
        <ViewSecret
          show={!!showSecret}
          setShow={() => {
            setShowSecret(null);
          }}
          item={showSecret!}
        />
      )}
    </>
  );
};

export default ManagedResourceResourcesV2;
