/* eslint-disable react/destructuring-assignment */
import { useMemo } from 'react';
import { toast } from 'react-toastify';
import { NumberInput, TextInput } from '@kloudlite/design-system/atoms/input';
import Select from '@kloudlite/design-system/atoms/select';
import Popup from '@kloudlite/design-system/molecule/popup';
import { useConsoleApi } from '~/console/server/gql/api-provider';
import { ExtractNodeType, parseName } from '~/console/server/r-utils/common';
import { useReload } from '~/root/lib/client/helpers/reloader';
import useForm, { dummyEvent } from '~/root/lib/client/hooks/use-form';
import Yup from '~/root/lib/server/helpers/yup';
import { handleError } from '~/root/lib/utils/common';
import {
  Github__Com___Kloudlite___Operator___Apis___Clusters___V1__AwsPoolType as awsPoolType,
  Github__Com___Kloudlite___Operator___Apis___Clusters___V1__GcpPoolType as gcpPoolType,
} from '~/root/src/generated/gql/server';
import { Link, useOutletContext } from '@remix-run/react';
import { INodepools } from '~/console/server/gql/queries/nodepool-queries';
import { awsRegions } from '~/console/dummy/consts';
import { mapper } from '@kloudlite/design-system/utils';
import { IDialogBase } from '~/console/components/types.d';
import { Switch } from '@kloudlite/design-system/atoms/switch';
import { NameIdView } from '~/console/components/name-id-view';
import { keyconstants } from '~/console/server/r-utils/key-constants';
import KeyValuePair from '~/console/components/key-value-pair';
import { InfoLabel } from '~/console/components/commons';
import { Button } from '@kloudlite/design-system/atoms/button';
import { IClusterContext } from '../_layout';
import {
  findNodePlan,
  nodePlans,
  provisionTypes,
  gcpPoolTypes,
} from './nodepool-utils';

type IDialog = IDialogBase<ExtractNodeType<INodepools>>;

const Root = (props: IDialog) => {
  const { setVisible, isUpdate } = props;

  const api = useConsoleApi();
  const reloadPage = useReload();
  const { cluster } = useOutletContext<IClusterContext>();
  const clusterRegion = cluster.spec?.aws?.region;
  const cloudProvider = cluster.spec?.cloudProvider;

  const filterLabels = (labels: Array<string>) => {
    if (isUpdate) {
      const org = { ...props.data.spec.nodeLabels };
      labels.forEach((label) => {
        delete org[label];
      });
      return org;
    }

    return {};
  };

  const { values, errors, handleChange, handleSubmit, resetValues, isLoading } =
    useForm({
      initialValues: isUpdate
        ? {
            nvidiaGpuEnabled: props.data.spec.aws?.nvidiaGpuEnabled || false,
            displayName: props.data.displayName,
            name: parseName(props.data),
            maximum: `${props.data.spec.maxCount}`,
            minimum: `${props.data.spec.minCount}`,
            poolType: props.data.spec.aws?.poolType || 'ec2',
            gcpMachineType: props.data.spec.gcp?.machineType || '',
            gcpAvailablityZone: props.data.spec.gcp?.availabilityZone || '',
            gcpPoolType: props.data.spec.gcp?.poolType || 'STANDARD',
            awsAvailabilityZone:
              props.data.spec.aws?.availabilityZone ||
              awsRegions.find((v) => v.Name === clusterRegion)?.Zones[0] ||
              '',
            instanceType:
              props.data.spec.aws?.ec2Pool?.instanceType || 'c6a.large',

            labels: filterLabels([keyconstants.nodepoolStateType]),
            labelsTemp: Object.entries(
              filterLabels([keyconstants.nodepoolStateType]) || {}
            ).map(([key, value]) => ({
              key,
              value,
            })),
            autoScale: props.data.spec.minCount !== props.data.spec.maxCount,
            isNameError: false,
            stateful:
              props.data.spec.nodeLabels[keyconstants.nodepoolStateType] ===
              'stateful',
          }
        : {
            nvidiaGpuEnabled: false,
            autoScale: false,
            name: '',
            displayName: '',
            minimum: '1',
            maximum: '1',
            gcpMachineType: '',
            gcpAvailablityZone: '',
            gcpPoolType: 'STANDARD',
            poolType: 'ec2',
            awsAvailabilityZone:
              awsRegions.find((v) => v.Name === clusterRegion)?.Zones[0] || '',

            // onDemand specs
            instanceType: 'c6a.large',

            labels: {},
            labelsTemp: [],
            isNameError: false,
            stateful: false,
          },
      validationSchema: Yup.object({
        name: Yup.string().required('id is required'),
        displayName: Yup.string().required('name is required'),
        minimum: Yup.number()
          .max(10, "you can't use more than 10 nodes for now")
          .min(0, 'minimum node count should be 0'),
        maximum: Yup.number()
          .max(10, "you can't use more than 10 nodes for now")
          .min(0, 'minimum node count should be 0'),
        poolType: Yup.string().required().oneOf(['ec2', 'spot']),
      }),
      onSubmit: async (val) => {
        const getNodeConf = () => {
          const getAwsNodeSpecs = () => {
            switch (val.poolType) {
              case 'ec2':
                return {
                  ec2Pool: {
                    instanceType: val.instanceType,
                    nodes: {},
                  },
                };
              case 'spot':
                const plan = findNodePlan(val.instanceType);
                return val.nvidiaGpuEnabled
                  ? {
                      gpuNode: {
                        instanceTypes: [plan?.value],
                      },
                    }
                  : {
                      spotPool: {
                        cpuNode: {
                          vcpu: {
                            max: `${plan?.spotSpec.cpuMax}`,
                            min: `${plan?.spotSpec.cpuMin}`,
                          },
                          memoryPerVcpu: {
                            max: `${plan?.spotSpec.memMax}`,
                            min: `${plan?.spotSpec.memMin}`,
                          },
                        },
                        nodes: {},
                      },
                    };
              default:
                return {};
            }
          };

          switch (cloudProvider) {
            case 'aws':
              return {
                aws: {
                  availabilityZone: val.awsAvailabilityZone,
                  nvidiaGpuEnabled: val.nvidiaGpuEnabled,
                  poolType: (val.poolType === 'ec2'
                    ? 'ec2'
                    : 'spot') as awsPoolType,
                  ...getAwsNodeSpecs(),
                },
              };
            case 'gcp':
              return {
                gcp: {
                  availabilityZone: val.gcpAvailablityZone,
                  machineType: val.gcpMachineType,
                  poolType: val.gcpPoolType as gcpPoolType,
                },
              };
            default:
              return {};
          }
        };

        try {
          if (!isUpdate) {
            const { errors: e } = await api.createNodePool({
              clusterName: parseName(cluster),
              pool: {
                displayName: val.displayName,
                metadata: {
                  name: val.name,
                },
                spec: {
                  maxCount: Number.parseInt(val.maximum, 10),
                  minCount: Number.parseInt(val.minimum, 10),
                  cloudProvider,
                  nodeLabels: {
                    ...val.labels,
                    [keyconstants.nodepoolStateType]: val.stateful
                      ? 'stateful'
                      : 'stateless',
                    // ...val.labels.map((l: any) => ({ [l.key]: l.value })),
                  },
                  ...getNodeConf(),
                },
              },
            });
            if (e) {
              throw e[0];
            }
          } else if (isUpdate) {
            const { errors: e } = await api.updateNodePool({
              clusterName: parseName(cluster),
              pool: {
                displayName: val.displayName,
                metadata: {
                  name: val.name,
                },
                spec: {
                  ...props.data.spec,
                  nodeLabels: {
                    ...(props.data.spec.nodeLabels || {}),
                    [keyconstants.nodepoolStateType]: val.stateful
                      ? 'stateful'
                      : 'stateless',
                  },
                  maxCount: Number.parseInt(val.maximum, 10),
                  minCount: Number.parseInt(val.minimum, 10),
                  ...getNodeConf(),
                },
              },
            });
            if (e) {
              throw e[0];
            }
          }
          reloadPage();
          resetValues();
          toast.success(
            `nodepool ${isUpdate ? 'updated' : 'created'} successfully`
          );
          setVisible(false);
        } catch (err) {
          handleError(err);
        }
      },
    });

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
      <Popup.Content>
        <div className="flex flex-col gap-2xl">
          <NameIdView
            resType="nodepool"
            displayName={values.displayName}
            name={values.name}
            label="Nodepool name"
            placeholder="Enter nodepool name"
            errors={errors.name}
            handleChange={handleChange}
            nameErrorLabel="isNameError"
            isUpdate={isUpdate}
          />

          {cloudProvider === 'aws' && (
            <>
              <Select
                label="Provision Mode"
                // eslint-disable-next-line react-hooks/rules-of-hooks
                value={values.poolType}
                options={async () => provisionTypes}
                onChange={(_, value) => {
                  handleChange('poolType')(dummyEvent(value));
                }}
              />

              <Select
                label="Availability Zone"
                value={values.awsAvailabilityZone}
                options={async () =>
                  mapper(
                    awsRegions.find((v) => v.Name === clusterRegion)?.Zones ||
                      [],
                    (v) => ({
                      value: v,
                      label: v,
                    })
                  )
                }
                onChange={(_, v) => {
                  handleChange('awsAvailabilityZone')(dummyEvent(v));
                }}
              />

              <div className="flex flex-row gap-xl items-end">
                <div className="flex flex-row gap-xl items-end flex-1">
                  <div className="flex-1">
                    <Select
                      // eslint-disable-next-line react-hooks/rules-of-hooks
                      value={useMemo(() => {
                        const plan = findNodePlan(values.instanceType);
                        return plan?.value;
                      }, [values.instanceType])}
                      label="Node plan"
                      options={async () => nodePlans}
                      onChange={(value) => {
                        handleChange('instanceType')(dummyEvent(value.value));
                        handleChange('nvidiaGpuEnabled')(
                          dummyEvent(!!value.gpuEnabled)
                        );
                      }}
                    />
                  </div>
                </div>
              </div>
            </>
          )}

          {cloudProvider === 'gcp' && (
            <>
              <div className="flex flex-row gap-2xl">
                <div className="flex-1">
                  <TextInput
                    // label="Availability zone"
                    label={
                      <InfoLabel
                        title="Availability zone"
                        info={
                          <span>
                            An availability zone is a distinct, isolated
                            location within a region, providing redundancy and
                            reliability for deploying resources{' '}
                            {/* <span className="bodySm-semibold">
                              single master
                            </span>{' '} */}
                            {/* <Button
                              // linkComponent={Link}
                              // target="_blank"
                              onClick={() => {
                                window.open(
                                  'https://cloud.google.com/compute/docs/regions-zones',
                                  '_blank'
                                );
                              }}
                              size="sm"
                              content="Click here"
                              variant="primary-plain"
                              className="!p-0"
                              // to="https://cloud.google.com/compute/docs/regions-zones"
                            /> */}
                            <a
                              className="text-text-primary hover:underline text-nowrap  cursor-pointer"
                              // href="https://cloud.google.com/compute/docs/regions-zones"
                              onClick={() => {
                                window.open(
                                  'https://cloud.google.com/compute/docs/regions-zones',
                                  '_blank'
                                );
                              }}
                            >
                              Click here{' '}
                            </a>
                            to get availability zone list.
                          </span>
                        }
                        label="Availability zone"
                      />
                    }
                    size="lg"
                    placeholder="Availability zone"
                    value={values.gcpAvailablityZone}
                    onChange={handleChange('gcpAvailablityZone')}
                    disabled={isUpdate}
                  />
                </div>

                <div className="flex-1">
                  <TextInput
                    label="Machine Type"
                    size="lg"
                    placeholder="Machine type"
                    value={values.gcpMachineType}
                    onChange={handleChange('gcpMachineType')}
                    disabled={isUpdate}
                  />
                </div>
              </div>

              <div className="flex flex-row gap-2xl">
                <div className="flex-1">
                  <Select
                    // eslint-disable-next-line react-hooks/rules-of-hooks
                    value={values.gcpPoolType}
                    label="Pool type"
                    options={async () => gcpPoolTypes}
                    onChange={(_, value) => {
                      handleChange('gcpPoolType')(dummyEvent(value));
                    }}
                    disabled={isUpdate}
                  />
                </div>

                <div className="flex flex-col gap-md ">
                  <div className="bodyMd-medium text-text-default">
                    Stateful
                  </div>
                  <div className="flex items-center h-6xl">
                    <Switch
                      label=""
                      disabled={isUpdate}
                      checked={values.stateful}
                      onChange={(val) => {
                        handleChange('stateful')(dummyEvent(val));
                      }}
                    />
                  </div>
                </div>
              </div>
            </>
          )}

          {cloudProvider === 'aws' && (
            <div className="flex flex-col gap-md ">
              <div className="bodyMd-medium text-text-default">Stateful</div>
              <div className="flex items-center h-6xl">
                <Switch
                  label=""
                  disabled={isUpdate}
                  checked={values.stateful}
                  onChange={(val) => {
                    handleChange('stateful')(dummyEvent(val));
                  }}
                />
              </div>
            </div>
          )}

          <div className="flex flex-row gap-xl items-end">
            <div className="flex flex-row gap-xl items-end flex-1 ">
              <div className="flex-1">
                <NumberInput
                  label={values.autoScale ? 'Min Node Count' : `Node Count`}
                  placeholder="Minimum"
                  value={values.minimum}
                  error={!!errors.minimum}
                  message={errors.minimum}
                  onChange={(e) => {
                    handleChange('minimum')(e);
                    if (!values.autoScale) {
                      handleChange('maximum')(e);
                    }
                  }}
                />
              </div>
              {values.autoScale && (
                <div className="flex-1">
                  <NumberInput
                    error={!!errors.maximum}
                    message={errors.maximum}
                    label="Max Node Count"
                    placeholder="Maximum"
                    value={values.maximum}
                    onChange={handleChange('maximum')}
                  />
                </div>
              )}
            </div>
            <div className="flex flex-col gap-md min-w-[115px]">
              <div className="bodyMd-medium text-text-default">AutoScale</div>
              <div className="flex items-center h-form-text-field-height">
                <Switch
                  label={values.autoScale ? 'Enabled' : 'Disabled'}
                  checked={values.autoScale}
                  onChange={(v) => handleChange('autoScale')(dummyEvent(v))}
                />
              </div>
            </div>
          </div>
          {!isUpdate && (
            <KeyValuePair
              addText="Add new"
              label="Labels"
              size="lg"
              value={values.labelsTemp}
              onChange={(val, v) => {
                handleChange('labels')(dummyEvent(v));
                handleChange('labelsTemp')(dummyEvent(val));
              }}
            />
          )}
        </div>
      </Popup.Content>
      <Popup.Footer>
        <Popup.Button closable content="Cancel" variant="basic" />
        <Popup.Button
          loading={isLoading}
          type="submit"
          content={isUpdate ? 'Update' : 'Create'}
          variant="primary"
        />
      </Popup.Footer>
    </Popup.Form>
  );
};

const HandleNodePool = (props: IDialog) => {
  const { isUpdate, setVisible, visible } = props;

  return (
    <Popup.Root show={visible} onOpenChange={(v) => setVisible(v)}>
      <Popup.Header>{isUpdate ? 'Edit nodepool' : 'Add nodepool'}</Popup.Header>
      {(!isUpdate || (isUpdate && props.data)) && <Root {...props} />}
    </Popup.Root>
  );
};

export default HandleNodePool;
