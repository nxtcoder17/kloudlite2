/* eslint-disable react/destructuring-assignment */
import { useParams } from '@remix-run/react';
import Popup from '@kloudlite/design-system/molecule/popup';
import { toast } from '@kloudlite/design-system/molecule/toast';
import {
  ExtractNodeType,
  parseName,
  parseNodes,
} from '~/console/server/r-utils/common';
import { useReload } from '~/lib/client/helpers/reloader';
import useForm, { dummyEvent } from '~/lib/client/hooks/use-form';
import Yup from '~/lib/server/helpers/yup';
import { handleError } from '~/lib/utils/common';
import { IDialogBase } from '~/console/components/types.d';
import { IRouters } from '~/console/server/gql/queries/router-queries';
import { useConsoleApi } from '~/console/server/gql/api-provider';
import Select from '@kloudlite/design-system/atoms/select';
import useCustomSwr from '~/lib/client/hooks/use-custom-swr';
import { useMapper } from '@kloudlite/design-system/utils';
import { NN } from '~/lib/types/common';
import { TextInput } from '@kloudlite/design-system/atoms/input';
import { useEffect, useState } from 'react';
import { IApps } from '~/console/server/gql/queries/app-queries';
import { Switch } from '@kloudlite/design-system/atoms/switch';
import { ModifiedRouter } from './_index';

type IDialog = IDialogBase<
  NN<ExtractNodeType<IRouters>['spec']['routes']>[number] & { id: string }
> & { router?: ModifiedRouter };

const Root = (props: IDialog) => {
  const { isUpdate, setVisible } = props;
  const api = useConsoleApi();
  const reloadPage = useReload();

  const {environment: envName } = useParams();
  const [selectedApp, setSelectedApp] = useState<ExtractNodeType<IApps>>();

  const {
    data,
    isLoading: appLoading,
    error: appLoadingError,
  } = useCustomSwr('/apps', async () => {
    if (!envName) {
      throw new Error('Project and Environment is required!.');
    }
    return api.listApps({  envName });
  });

  const { values, errors, handleSubmit, handleChange, isLoading, resetValues } =
    useForm({
      initialValues: isUpdate
        ? {
            path:
              props.data.path?.[0] === '/'
                ? props.data.path.substring(1)
                : props.data.path,
            app: props.data.app || '',
            port: `${props.data.port}`,
            reWrite: false,
          }
        : {
            path: '',
            app: '',
            port: '',
            reWrite: false,
          },
      validationSchema: Yup.object({
        path: Yup.string().test(
          'is-valid',
          'Path should not contain spaces.',
          (value) => {
            return !value?.includes(' ');
          }
        ),
        app: Yup.string().required(),
        port: Yup.string().required(),
      }),

      onSubmit: async (val) => {
        const { router } = props;
        if (!envName || !router || !router.metadata?.name) {
          throw new Error('Project, Router and Environment is required!.');
        }
        try {
          if (!isUpdate) {
            const { errors: e } = await api.updateRouter({
              envName,
              
              router: {
                displayName: router.displayName,
                spec: {
                  ...router.spec,
                  routes: [
                    ...(router.spec.routes?.map((r) => ({
                      path: r.path,
                      app: r.app,
                      port: r.port,
                      rewrite: r.rewrite,
                    })) || []),
                    {
                      path: `/${val.path}`,
                      app: val.app,
                      port: parseInt(val.port, 10),
                      rewrite: val.reWrite,
                    },
                  ],
                },
                metadata: {
                  ...router.metadata,
                },
              },
            });
            if (e) {
              throw e[0];
            }
            toast.success('Route created successfully');
          } else {
            const { errors: e } = await api.updateRouter({
              envName,
              
              router: {
                displayName: router.displayName,
                spec: {
                  ...router.spec,
                  routes: [
                    ...(router.spec.routes

                      ?.filter(
                        (
                          rou // @ts-ignore
                        ) => rou.id !== props.data.id
                      )
                      .map((route) => ({
                        app: route.app,
                        path: route.path,
                        port: route.port,
                      })) || []),
                    {
                      path: `/${val.path}`,
                      app: val.app,
                      port: parseInt(val.port, 10),
                    },
                  ],
                },
                metadata: {
                  ...router.metadata,
                },
              },
            });
            if (e) {
              throw e[0];
            }
            toast.success('Route updated successfully');
          }
          reloadPage();
          setVisible(false);
          resetValues();
        } catch (err) {
          handleError(err);
        }
      },
    });

  const apps = useMapper(parseNodes(data), (val) => ({
    label: val.displayName,
    value: parseName(val),
    app: val,
    render: () => val.displayName,
  }));

  useEffect(() => {
    const d = parseNodes(data);
    if (d.length > 0) {
      if (isUpdate) {
        setSelectedApp(d.find((app) => parseName(app) === props.data.app));
      } else if (d.length === 1) {
        handleChange('app')(dummyEvent(parseName(d[0])));
        setSelectedApp(d[0]);
      }
    }
  }, [isUpdate, data]);

  useEffect(() => {
    if (selectedApp?.spec.services?.length === 0) {
      handleChange('port')(dummyEvent(selectedApp.spec.services[0].port));
    }
  }, [selectedApp]);

  return (
    <Popup.Form onSubmit={handleSubmit}>
      <Popup.Content className="flex flex-col gap-3xl">
        <div className="flex flex-row gap-xl items-end">
          <div className="flex flex-row gap-xl items-end flex-1">
            <div className="flex-1">
              <TextInput
                label="Path"
                size="lg"
                value={values.path}
                onChange={(e) => {
                  handleChange('path')(
                    dummyEvent(e.target.value.toLowerCase())
                  );
                }}
                error={!!errors.path}
                message={errors.path}
                prefix="/"
              />
            </div>
          </div>

          <div className="flex flex-col gap-md ">
            <div className="bodyMd-medium text-text-default">Rewrite</div>
            <div className="flex items-center h-6xl">
              <Switch
                label=""
                checked={values.reWrite}
                onChange={(val) => {
                  handleChange('reWrite')(dummyEvent(val));
                }}
              />
            </div>
          </div>
        </div>
        <Select
          size="lg"
          label="App"
          value={values.app}
          options={async () => [...apps]}
          onChange={(val) => {
            handleChange('app')(dummyEvent(val.value));
            setSelectedApp(val.app);
          }}
          error={!!errors.app || !!appLoadingError}
          message={appLoadingError ? 'Error fetching apps.' : errors.app}
          loading={appLoading}
        />
        <Select
          size="lg"
          label="Port"
          disabled={!values.app}
          value={values.port}
          options={async () => [
            ...(selectedApp?.spec.services?.map((svc) => ({
              label: `${svc.port}`,
              value: `${svc.port}`,
            })) || []),
          ]}
          onChange={(_, val) => {
            handleChange('port')(dummyEvent(val));
          }}
          error={!!errors.port}
          message={errors.port}
        />
      </Popup.Content>
      <Popup.Footer>
        <Popup.Button content="Cancel" variant="basic" closable />
        <Popup.Button
          loading={isLoading}
          type="submit"
          content={!isUpdate ? 'Add' : 'Update'}
          variant="primary"
        />
      </Popup.Footer>
    </Popup.Form>
  );
};

const HandleRoute = (props: IDialog) => {
  const { isUpdate, setVisible, visible } = props;
  return (
    <Popup.Root show={visible} onOpenChange={(v) => setVisible(v)}>
      <Popup.Header>{isUpdate ? 'Edit route' : 'Add Route'}</Popup.Header>
      {(!isUpdate || (isUpdate && props.data)) && <Root {...props} />}
    </Popup.Root>
  );
};

export default HandleRoute;
