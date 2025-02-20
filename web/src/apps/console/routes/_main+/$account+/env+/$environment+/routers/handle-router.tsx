/* eslint-disable react/destructuring-assignment */
import { useParams } from '@remix-run/react';
import { useEffect, useRef } from 'react';
import { Checkbox } from '@kloudlite/design-system/atoms/checkbox';
import Select from '@kloudlite/design-system/atoms/select';
import Banner from '@kloudlite/design-system/molecule/banner';
import Popup from '@kloudlite/design-system/molecule/popup';
import { toast } from '@kloudlite/design-system/molecule/toast';
import { useAppend, useMapper } from '@kloudlite/design-system/utils';
import CommonPopupHandle from '~/console/components/common-popup-handle';
import { NameIdView } from '~/console/components/name-id-view';
import { IDialogBase } from '~/console/components/types.d';
import { useConsoleApi } from '~/console/server/gql/api-provider';
import { IRouters } from '~/console/server/gql/queries/router-queries';
import {
  ExtractNodeType,
  parseName,
  parseNodes,
} from '~/console/server/r-utils/common';
import { useReload } from '~/lib/client/helpers/reloader';
import useCustomSwr from '~/lib/client/hooks/use-custom-swr';
import useForm, { dummyEvent } from '~/lib/client/hooks/use-form';
import Yup from '~/lib/server/helpers/yup';
import { handleError } from '~/lib/utils/common';

type IDialog = IDialogBase<ExtractNodeType<IRouters>>;

const Root = (props: IDialog) => {
  const { isUpdate, setVisible } = props;
  const api = useConsoleApi();
  const reloadPage = useReload();

  const { environment: envName } = useParams();

  // const { cluster } = useOutletContext<IAppContext>();

  const {
    data,
    isLoading: domainLoading,
    error: domainLoadingError,
  } = useCustomSwr('/domains', async () => {
    return api.listDomains({
      search: {
        clusterName: {
          matchType: 'exact',
          // exact: parseName(cluster),
          exact: 'localhost',
        },
      },
    });
  });

  const { values, errors, handleSubmit, handleChange, isLoading, resetValues } =
    useForm({
      initialValues: isUpdate
        ? {
          name: parseName(props.data),
          displayName: props.data.displayName,
          domains: [],
          isNameError: false,
          isTlsEnabled: props.data.spec.https?.enabled || false,
        }
        : {
          name: '',
          displayName: '',
          domains: [],
          isNameError: false,
          isTlsEnabled: false,
        },
      validationSchema: Yup.object({
        displayName: Yup.string().required(),
        name: Yup.string().required(),
        domains: Yup.array().test('required', 'domain is required', (val) => {
          return val && val?.length > 0;
        }),
      }),

      onSubmit: async (val) => {
        if (!envName || val.domains.length === 0) {
          throw new Error('Project, Environment and Domain is required!.');
        }
        try {
          if (!isUpdate) {
            const { errors: e } = await api.createRouter({
              envName,

              router: {
                displayName: val.displayName,
                metadata: {
                  name: val.name,
                },
                spec: {
                  domains: val.domains,
                  https: {
                    enabled: val.isTlsEnabled,
                  },
                },
              },
            });
            if (e) {
              throw e[0];
            }
            toast.success('Router created successfully');
          } else {
            const { errors: e } = await api.updateRouter({
              envName,

              router: {
                displayName: val.displayName,
                metadata: {
                  name: val.name,
                },
                spec: {
                  ...props.data.spec,
                  domains: val.domains,
                  https: {
                    enabled: val.isTlsEnabled,
                  },
                },
              },
            });
            if (e) {
              throw e[0];
            }
            toast.success('Router updated successfully');
          }
          reloadPage();
          setVisible(false);
          resetValues();
        } catch (err) {
          handleError(err);
        }
      },
    });

  const domains = useMapper(parseNodes(data), (val) => ({
    label: val.domainName,
    value: val.domainName,
  }));

  const combinedDomains = useAppend(
    domains,
    isUpdate
      ? props.data.spec.domains
        .filter((d) => !domains.find((f) => f.value === d))
        .map((d) => ({ label: d, value: d }))
      : []
  );

  useEffect(() => {
    if (isUpdate) {
      const d = combinedDomains
        .filter((d) => props.data.spec.domains.includes(d.value))
        .map((x) => x.value);
      handleChange('domains')(dummyEvent(d));
    }
  }, [data]);

  const nameIDRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    nameIDRef.current?.focus();
  }, [nameIDRef]);

  return (
    <Popup.Form
      onSubmit={(e) => {
        if (!values.isNameError) {
          handleSubmit(e);
        } else {
          e.preventDefault();
        }
      }}
    >
      <Popup.Content className="flex flex-col justify-start gap-3xl">
        <NameIdView
          ref={nameIDRef}
          resType="router"
          label="Name"
          displayName={values.displayName}
          name={values.name}
          errors={errors.name}
          handleChange={handleChange}
          nameErrorLabel="isNameError"
          isUpdate={isUpdate}
        />
        <Select
          creatable
          size="lg"
          label="Domains"
          multiple
          value={values.domains}
          options={async () => [...combinedDomains]}
          onChange={(val, v) => {
            handleChange('domains')(dummyEvent(v));
          }}
          error={!!errors.domains || !!domainLoadingError}
          message={
            errors.domains ||
            (domainLoadingError ? 'Error fetching domains.' : '')
          }
          loading={domainLoading}
          disableWhileLoading
        />
        <Checkbox
          label="Enable TLS"
          checked={values.isTlsEnabled}
          onChange={(val) => {
            handleChange('isTlsEnabled')(dummyEvent(val));
          }}
        />
        <Banner
          type="info"
          body={
            <span>
              All the domain CNames should be pointed to following Cluster DNS
              Name{' '}
              <span className="bodyMd-medium">
                {/* `{cluster.spec.publicDNSHost}` */}
                "localhost"
              </span>
            </span>
          }
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

const HandleRouter = (props: IDialog) => {
  return (
    <CommonPopupHandle
      {...props}
      createTitle="Create router"
      updateTitle="Update router"
      root={Root}
    />
  );
};
export default HandleRouter;
