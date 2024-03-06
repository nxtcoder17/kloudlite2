import { Trash, PencilLine } from '@jengaicons/react';
import { useState } from 'react';
import { generateKey, titleCase } from '~/components/utils';
import {
  ListItem,
  ListTitle,
  listClass,
  listFlex,
} from '~/console/components/console-list-components';
import DeleteDialog from '~/console/components/delete-dialog';
import Grid from '~/console/components/grid';
import List from '~/console/components/list';
import ListGridView from '~/console/components/list-grid-view';
import ResourceExtraAction from '~/console/components/resource-extra-action';
import {
  ExtractNodeType,
  parseName,
  parseUpdateOrCreatedBy,
  parseUpdateOrCreatedOn,
} from '~/console/server/r-utils/common';
import { handleError } from '~/root/lib/utils/common';
import { IRouters } from '~/console/server/gql/queries/router-queries';
import { Link, useParams } from '@remix-run/react';
import { listStatus } from '~/console/components/sync-status';
import { useConsoleApi } from '~/console/server/gql/api-provider';
import { useReload } from '~/root/lib/client/helpers/reloader';
import { toast } from '~/components/molecule/toast';
import { Button } from '~/components/atoms/button';
import Tooltip from '~/components/atoms/tooltip';
import HandleRouter from './handle-router';

const RESOURCE_NAME = 'domain';
type BaseType = ExtractNodeType<IRouters>;

const parseItem = (item: BaseType) => {
  return {
    name: item.displayName,
    id: parseName(item),
    updateInfo: {
      author: `Updated by ${parseUpdateOrCreatedBy(item)}`,
      time: parseUpdateOrCreatedOn(item),
    },
  };
};

type OnAction = ({
  action,
  item,
}: {
  action: 'edit' | 'delete' | 'detail';
  item: BaseType;
}) => void;

type IExtraButton = {
  onAction: OnAction;
  item: BaseType;
};

const formatDomain = (domain: string) => {
  const d = domain.startsWith('https://') ? domain : `https://${domain}`;
  return { full: d, short: d.replace('https://', '') };
};

const ExtraButton = ({ onAction, item }: IExtraButton) => {
  return (
    <ResourceExtraAction
      options={[
        {
          label: 'Edit',
          icon: <PencilLine size={16} />,
          type: 'item',
          onClick: () => onAction({ action: 'edit', item }),
          key: 'edit',
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
}

const GridView = ({ items, onAction }: IResource) => {
  const { account, project, environment } = useParams();
  return (
    <Grid.Root className="!grid-cols-1 md:!grid-cols-3" linkComponent={Link}>
      {items.map((item, index) => {
        const { name, id, updateInfo } = parseItem(item);
        const keyPrefix = `${RESOURCE_NAME}-${id}-${index}`;
        const firstDomain = item.spec.domains?.[0];
        return (
          <Grid.Column
            key={id}
            to={`/${account}/${project}/${environment}/router/${id}/routes`}
            rows={[
              {
                key: generateKey(keyPrefix, name + id),
                render: () => (
                  <ListTitle
                    title={name}
                    action={<ExtraButton onAction={onAction} item={item} />}
                  />
                ),
              },
              {
                key: generateKey(keyPrefix, 'extra_domain'),
                render: () => (
                  <ListItem
                    data={
                      <div className="flex flex-row items-center gap-md">
                        <Button
                          LinkComponent={Link}
                          target="_blank"
                          size="sm"
                          content={formatDomain(firstDomain).short}
                          variant="primary-plain"
                          to={formatDomain(firstDomain).full}
                        />

                        {item.spec.domains.length > 1 && (
                          <Tooltip.Root
                            content={
                              <div className="flex flex-col gap-md">
                                {item.spec.domains
                                  .filter((d) => d !== firstDomain)
                                  .map((d) => (
                                    <Button
                                      key={d}
                                      LinkComponent={Link}
                                      target="_blank"
                                      size="sm"
                                      content={formatDomain(d).short}
                                      variant="primary-plain"
                                      to={formatDomain(d).full}
                                    />
                                  ))}
                              </div>
                            }
                          >
                            <Button
                              content={`+${item.spec.domains.length - 1} more`}
                              variant="plain"
                              size="sm"
                            />
                          </Tooltip.Root>
                        )}
                      </div>
                    }
                  />
                ),
              },
              {
                key: generateKey(keyPrefix, updateInfo.author),
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

const ListView = ({ items, onAction }: IResource) => {
  const { account, project, environment } = useParams();
  return (
    <List.Root linkComponent={Link}>
      {items.map((item, index) => {
        const { name, id, updateInfo } = parseItem(item);
        const keyPrefix = `${RESOURCE_NAME}-${id}-${index}`;
        const status = listStatus({ key: `${keyPrefix}status`, item });
        const firstDomain = item.spec.domains?.[0];
        return (
          <List.Row
            key={id}
            className="!p-3xl"
            to={`/${account}/${project}/${environment}/router/${id}/routes`}
            columns={[
              {
                key: generateKey(keyPrefix, name + id),
                className: listClass.title,
                render: () => <ListTitle title={name} />,
              },
              status,
              {
                key: generateKey(keyPrefix, 'extra_domain'),
                render: () => (
                  <ListItem
                    data={
                      <div className="flex flex-row items-center gap-md">
                        <Button
                          size="sm"
                          content={formatDomain(firstDomain).short}
                          variant="primary-plain"
                          onClick={(e) => {
                            e.preventDefault();
                            window.open(
                              formatDomain(firstDomain).full,
                              '_blank',
                              'noopener,noreferrer'
                            );
                          }}
                        />

                        {item.spec.domains.length > 1 && (
                          <Tooltip.Root
                            content={
                              <div className="flex flex-col gap-md">
                                {item.spec.domains
                                  .filter((d) => d !== firstDomain)
                                  .map((d) => (
                                    <Button
                                      key={d}
                                      size="sm"
                                      content={formatDomain(d).short}
                                      variant="primary-plain"
                                      onClick={(e) => {
                                        e.preventDefault();
                                        window.open(
                                          formatDomain(d).full,
                                          '_blank',
                                          'noopener,noreferrer'
                                        );
                                      }}
                                    />
                                  ))}
                              </div>
                            }
                          >
                            <Button
                              content={`+${item.spec.domains.length - 1} more`}
                              variant="plain"
                              size="sm"
                            />
                          </Tooltip.Root>
                        )}
                      </div>
                    }
                  />
                ),
              },
              listFlex({ key: 'flex-1' }),
              {
                key: generateKey(keyPrefix, updateInfo.author),
                className: listClass.author,
                render: () => (
                  <ListItem
                    data={updateInfo.author}
                    subtitle={updateInfo.time}
                  />
                ),
              },
              {
                key: generateKey(keyPrefix, 'action'),
                render: () => <ExtraButton onAction={onAction} item={item} />,
              },
            ]}
          />
        );
      })}
    </List.Root>
  );
};

const RouterResources = ({ items = [] }: { items: BaseType[] }) => {
  const [showDeleteDialog, setShowDeleteDialog] = useState<BaseType | null>(
    null
  );
  const [visible, setVisible] = useState<BaseType | null>(null);
  const api = useConsoleApi();
  const reloadPage = useReload();
  const { environment, project } = useParams();

  const props: IResource = {
    items,
    onAction: ({ action, item }) => {
      switch (action) {
        case 'edit':
          setVisible(item);
          break;
        case 'delete':
          setShowDeleteDialog(item);
          break;
        default:
      }
    },
  };
  return (
    <>
      <ListGridView
        listView={<ListView {...props} />}
        gridView={<GridView {...props} />}
      />
      <DeleteDialog
        resourceName={showDeleteDialog?.displayName}
        resourceType={RESOURCE_NAME}
        show={showDeleteDialog}
        setShow={setShowDeleteDialog}
        onSubmit={async () => {
          if (!environment || !project) {
            throw new Error('Project and Environment is required!.');
          }
          try {
            const { errors } = await api.deleteRouter({
              envName: environment,
              projectName: project,
              routerName: parseName(showDeleteDialog),
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
      <HandleRouter
        {...{
          isUpdate: true,
          data: visible!,
          visible: !!visible,
          setVisible: () => setVisible(null),
        }}
      />
    </>
  );
};

export default RouterResources;
