import { useNavigate, useOutletContext, useParams } from '@remix-run/react';
import {
  FormEventHandler,
  ReactNode,
  useCallback,
  useEffect,
  useRef,
  useState,
} from 'react';
import { toast } from 'react-toastify';
import { NumberInput, TextInput } from '@kloudlite/design-system/atoms/input';
import Select from '@kloudlite/design-system/atoms/select';
import { Switch } from '@kloudlite/design-system/atoms/switch';
import { titleCase } from '@kloudlite/design-system/utils';
import {
  BottomNavigation,
  ReviewComponent,
} from '~/console/components/commons';
import MultiStepProgress, {
  useMultiStepProgress,
} from '~/console/components/multi-step-progress';
import MultiStepProgressWrapper from '~/console/components/multi-step-progress-wrapper';
import { NameIdView } from '~/console/components/name-id-view';
import { findClusterStatus } from '~/console/hooks/use-cluster-status';
import { ClusterSelectItem } from '~/console/page-components/handle-environment';
import { useConsoleApi } from '~/console/server/gql/api-provider';
import {
  IMSvTemplate,
  IMSvTemplates,
} from '~/console/server/gql/queries/managed-templates-queries';
import { parseName, parseNodes } from '~/console/server/r-utils/common';
import { keyconstants } from '~/console/server/r-utils/key-constants';
import { ensureAccountClientSide } from '~/console/server/utils/auth-utils';
import { flatM, flatMapValidations } from '~/console/utils/commons';
import useForm, { dummyEvent } from '~/root/lib/client/hooks/use-form';
import Yup from '~/root/lib/server/helpers/yup';
import { handleError } from '~/root/lib/utils/common';
import { IAccountContext } from '../_layout';

const valueRender = ({ label, icon }: { label: string; icon: string }) => {
  return (
    <div className="flex flex-row gap-lg items-center">
      <span>
        <img alt={label} src={icon} className="w-2xl h-w-2xl" />
      </span>
      <div>{label}</div>
    </div>
  );
};

const RenderField = ({
  field,
  value,
  onChange,
  errors,
  fieldKey,
}: {
  field: IMSvTemplate['fields'][number];
  onChange: (e: string) => (e: { target: { value: any } }) => void;
  value: any;
  errors: {
    [key: string]: string;
  };
  fieldKey: string;
}) => {
  const [qos, setQos] = useState(false);
  if (field.inputType === 'Number') {
    return (
      <NumberInput
        error={!!errors[fieldKey]}
        message={errors[fieldKey]}
        label={`${field.label}${field.required ? ' *' : ''}`}
        placeholder={field.label}
        value={parseFloat(value) / (field.multiplier || 1) || ''}
        size="lg"
        onChange={({ target }) => {
          onChange(`res.${field.name}`)(
            dummyEvent(
              `${parseFloat(target.value) * (field.multiplier || 1)}${
                field.unit
              }`,
            ),
          );
        }}
        suffix={field.displayUnit}
      />
    );
  }

  if (field.inputType === 'String') {
    return (
      <TextInput
        label={field.label}
        value={value || ''}
        onChange={onChange(`res.${field.name}`)}
        suffix={field.displayUnit}
        error={!!errors[fieldKey]}
        message={errors[fieldKey]}
        size="lg"
      />
    );
  }
  if (field.inputType === 'Resource') {
    return (
      <div className="flex flex-col gap-md">
        <div className="bodyMd-medium text-text-default">{`${field.label}${
          field.required ? ' *' : ''
        }`}</div>
        <div className="flex flex-row gap-xl items-center">
          <div className="flex flex-row gap-xl items-end flex-1 ">
            <div className="flex-1">
              <NumberInput
                size="lg"
                error={!!errors[`${fieldKey}.min`]}
                message={errors[`${fieldKey}.min`]}
                placeholder={qos ? field.label : `${field.label} min`}
                value={parseFloat(value.min) / (field.multiplier || 1)}
                onChange={({ target }) => {
                  onChange(`res.${field.name}.min`)(
                    dummyEvent(
                      `${parseFloat(target.value) * (field.multiplier || 1)}${
                        field.unit
                      }`,
                    ),
                  );
                  if (qos) {
                    onChange(`res.${field.name}.max`)(
                      dummyEvent(
                        `${parseFloat(target.value) * (field.multiplier || 1)}${
                          field.unit
                        }`,
                      ),
                    );
                  }
                }}
                suffix={field.displayUnit}
              />
            </div>
            {!qos && (
              <div className="flex-1">
                <NumberInput
                  size="lg"
                  error={!!errors[`${fieldKey}.max`]}
                  message={errors[`${fieldKey}.max`]}
                  placeholder={`${field.label} max`}
                  value={parseFloat(value.max) / (field.multiplier || 1)}
                  onChange={({ target }) => {
                    onChange(`res.${field.name}.max`)(
                      dummyEvent(
                        `${parseFloat(target.value) * (field.multiplier || 1)}${
                          field.unit
                        }`,
                      ),
                    );
                  }}
                  suffix={field.displayUnit}
                />
              </div>
            )}
          </div>
          <div className="flex flex-col gap-md min-w-[115px]">
            <Switch
              label="Guaranteed"
              checked={qos}
              onChange={(_value) => {
                setQos(_value);
                if (_value) {
                  onChange(`res.${field.name}.max`)(dummyEvent(`${value.min}`));
                }
              }}
            />
          </div>
        </div>
      </div>
    );
  }
  return <div>unknown input type {field.inputType}</div>;
};

type ISelectedTemplate = {
  category: string;
  categoryDisplayName: string;
  template: IMSvTemplate;
};

const TemplateView = ({
  handleSubmit,
  values,
  handleChange,
  errors,
  templates,
  isLoading,
}: {
  handleSubmit: FormEventHandler<HTMLFormElement>;
  values: Record<string, any>;
  errors: Record<string, any>;
  templates: IMSvTemplates;
  isLoading: boolean;
  handleChange: (key: string) => (e: { target: { value: any } }) => void;
}) => {
  return (
    <form className="flex flex-col gap-3xl" onSubmit={handleSubmit}>
      <div className="bodyMd text-text-soft">Create your managed services.</div>
      <Select
        label="Managed service templates"
        size="lg"
        placeholder="Select templates"
        value={values.selectedTemplate?.template.name}
        valueRender={valueRender}
        searchable
        error={!!errors.selectedTemplate}
        message={errors.selectedTemplate}
        onChange={({ item }) => {
          handleChange('selectedTemplate')(dummyEvent(item));
        }}
        options={async () =>
          templates.map((mt) => ({
            label: mt.displayName,
            options: mt.items.map((mti) => ({
              label: mti.displayName,
              value: mti.name,
              icon: mti.logoUrl,
              item: {
                categoryDisplayName: mt.displayName,
                category: mt.category,
                template: mti,
              },
              render: () => (
                <div className="flex flex-row items-center gap-xl">
                  <span>
                    <img
                      alt={mti.displayName}
                      src={mti.logoUrl}
                      className="w-2xl h-w-2xl"
                    />
                  </span>
                  <div>{mti.displayName}</div>
                </div>
              ),
            })),
          }))
        }
      />
      <BottomNavigation
        primaryButton={{
          type: 'submit',
          loading: isLoading,
          content: 'Next',
        }}
      />
    </form>
  );
};

const FieldView = ({
  selectedTemplate,
  // nodepools,
  // nodepoolIsLoading,
  clusters,
  values,
  handleSubmit,
  handleChange,
  errors,
}: {
  handleChange: (key: string) => (e: { target: { value: any } }) => void;
  handleSubmit: FormEventHandler<HTMLFormElement>;
  values: Record<string, any>;
  errors: Record<string, any>;
  selectedTemplate: ISelectedTemplate | null;
  // nodepools: { label: string; value: string }[];
  // nodepoolIsLoading: boolean;
  clusters: {
    label: string;
    value: string;
    ready?: boolean;
    render: () => ReactNode;
  }[];
}) => {
  const nameRef = useRef<HTMLInputElement>(null);
  useEffect(() => {
    nameRef.current?.focus();
  }, [nameRef.current]);
  return (
    <form
      className="flex flex-col gap-3xl"
      onSubmit={(e) => {
        if (!values.isNameError) {
          handleSubmit(e);
        } else {
          e.preventDefault();
        }
      }}
    >
      <NameIdView
        ref={nameRef}
        placeholder="Enter managed service name"
        label="Name"
        resType="cluster_managed_service"
        name={values.name}
        displayName={values.displayName}
        errors={errors.name}
        handleChange={handleChange}
        nameErrorLabel="isNameError"
      />

      <Select
        label="Select Cluster"
        size="lg"
        value={values.clusterName}
        placeholder="Select a Cluster"
        options={async () => clusters}
        // options={async () => [
        //   ...((clusters &&
        //     clusters.filter((c) => {
        //       return c.ready;
        //     })) ||
        //     []),
        // ]}
        onChange={({ value }) => {
          handleChange('clusterName')(dummyEvent(value));
          handleChange('nodepoolName')(dummyEvent(''));
        }}
        showclear
        error={!!errors.clusterName}
        message={errors.clusterName}
        // loading={cIsLoading || byokCIsLoading}
      />

      {/* <Select
        label="Nodepool Name"
        size="lg"
        placeholder="Select Nodepool"
        value={values.nodepoolName}
        creatable
        onChange={(val) => {
          handleChange('nodepoolName')(dummyEvent(val.value));
        }}
        options={async () => [...nodepools]}
        error={!!errors.nodepoolName}
        message={errors.nodepoolName}
        showclear
        loading={nodepoolIsLoading}
        noOptionMessage={
          <div className="p-2xl bodyMd text-center">
            No stateful nodepools available
          </div>
        }
      /> */}

      {selectedTemplate?.template.fields?.map((field) => {
        const k = field.name;
        const x = k.split('.').reduce((acc, curr) => {
          if (!acc) {
            return values.res?.[curr];
          }
          return acc[curr];
        }, null);
        return (
          <RenderField
            field={field}
            key={field.name}
            onChange={handleChange}
            value={x}
            errors={errors}
            fieldKey={k}
          />
        );
      })}
      <BottomNavigation
        primaryButton={{
          type: 'submit',
          content: 'Next',
        }}
      />
    </form>
  );
};

const ReviewView = ({
  handleSubmit,
  values,
  isLoading,
  onEdit,
}: {
  values: Record<string, any>;
  onEdit: (step: number) => void;
  handleSubmit: FormEventHandler<HTMLFormElement>;
  isLoading?: boolean;
}) => {
  const renderFieldView = () => {
    const fields = Object.entries(values.res).filter(
      ([k, _v]) => !['resources'].includes(k),
    );
    if (fields.length > 0) {
      return (
        <ReviewComponent
          title="Fields"
          onEdit={() => {
            onEdit(2);
          }}
        >
          <div className="flex flex-col p-xl  gap-lg rounded border border-border-default flex-1 overflow-hidden">
            {fields?.map(([key, value]) => {
              const k = key as string;
              const v = value as string;
              return (
                <div
                  key={k}
                  className="flex flex-col gap-md  [&:not(:last-child)]:pb-lg   [&:not(:last-child)]:border-b border-border-default"
                >
                  <div className="bodyMd-medium text-text-default">
                    {titleCase(k)}
                  </div>
                  <div className="bodySm text-text-soft">{v}</div>
                </div>
              );
            })}
          </div>
        </ReviewComponent>
      );
    }
    return null;
  };
  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-3xl">
      <div className="flex flex-col gap-xl">
        <ReviewComponent
          title="Basic detail"
          onEdit={() => {
            onEdit(2);
          }}
        >
          <div className="flex flex-col p-xl gap-lg rounded border border-border-default">
            <div className="flex flex-col gap-md">
              <div className="bodyMd-semibold text-text-default">
                {values.displayName}
              </div>
              <div className="bodySm text-text-soft">{values.name}</div>
            </div>
          </div>
        </ReviewComponent>

        <ReviewComponent
          title="Service details"
          onEdit={() => {
            onEdit(1);
          }}
        >
          <div className="flex flex-col gap-xl p-xl rounded border border-border-default">
            <div className="flex flex-col gap-lg pb-xl border-b border-border-default">
              <div className="flex-1 bodyMd-medium text-text-default">
                {values?.selectedTemplate?.categoryDisplayName}
              </div>
              <div className="text-text-soft bodyMd">
                {values?.selectedTemplate?.template?.displayName}
              </div>
            </div>
            <div className="flex flex-col gap-lg ">
              <div className="flex-1 bodyMd-medium text-text-default">
                Cluster Name
              </div>
              <div className="text-text-soft bodyMd">{values.clusterName}</div>
            </div>
          </div>
        </ReviewComponent>

        {renderFieldView()}
        {values?.res?.resources && (
          <ReviewComponent
            title="Fields"
            onEdit={() => {
              onEdit(2);
            }}
          >
            <div className="flex flex-col p-xl  gap-lg rounded border border-border-default flex-1 overflow-hidden">
              {Object.entries(values?.res?.resources).map(([key, value]) => {
                return (
                  <div
                    key={key}
                    className="flex flex-col gap-md  [&:not(:last-child)]:pb-lg   [&:not(:last-child)]:border-b border-border-default"
                  >
                    <div className="bodyMd-medium text-text-default">
                      {titleCase(key)}
                    </div>
                    <div className="bodySm text-text-soft">
                      {Object.entries(value || {}).map(([pKey, pValue]) => (
                        <div key={pKey}>
                          {pKey}
                          {' : '}
                          {pValue}
                        </div>
                      ))}
                    </div>
                  </div>
                );
              })}
            </div>
          </ReviewComponent>
        )}
      </div>
      <BottomNavigation
        primaryButton={{
          type: 'submit',
          loading: isLoading,
          content: 'Create',
        }}
      />
    </form>
  );
};

// const ClusterSelectItem = ({
//   label,
//   value,
// }: {
//   label: string;
//   value: string;
// }) => {
//   return (
//     <div>
//       <div className="flex flex-col">
//         <div>{label}</div>
//         <div className="bodySm text-text-soft">{value}</div>
//       </div>
//     </div>
//   );
// };

const ManagedServiceLayout = () => {
  // const { msvtemplates, cluster, account } =
  //   useOutletContext<IClusterContext>();
  const { msvtemplates, account } = useOutletContext<IAccountContext>();
  const navigate = useNavigate();
  const api = useConsoleApi();

  // const rootUrl = `/${parseName(account)}/infra/${parseName(
  //   account
  // )}/managed-services`;

  const rootUrl = `/${parseName(account)}/managed-services`;

  const { currentStep, jumpStep, nextStep } = useMultiStepProgress({
    defaultStep: 1,
    totalSteps: 3,
  });

  const [clusterList, setClusterList] = useState<any[]>([]);
  const params = useParams();

  const getClusters = useCallback(async () => {
    ensureAccountClientSide(params);
    try {
      const byokClusters = await api.listByokClusters({});
      const data = parseNodes(byokClusters.data).map((c) => ({
        label: c.displayName,
        value: parseName(c),
        ready: findClusterStatus(c),
        disabled: () => !findClusterStatus(c),
        render: ({ disabled }: { disabled: boolean }) => (
          <ClusterSelectItem
            label={c.displayName}
            value={parseName(c)}
            disabled={disabled}
          />
        ),
      }));
      setClusterList(data);
    } catch (err) {
      handleError(err);
    }
  }, []);

  useEffect(() => {
    getClusters();
  }, []);

  const { values, errors, handleSubmit, handleChange, isLoading, setValues } =
    useForm({
      initialValues: {
        name: '',
        displayName: '',
        res: {},
        selectedTemplate: null,
        isNameError: false,
        clusterName: '',
        nodepoolName: '',
      },
      validationSchema: Yup.object().shape({
        name: Yup.string().test('required', 'Name is required', (v) => {
          return !(currentStep === 2 && !v);
        }),
        displayName: Yup.string().test('required', 'Name is required', (v) => {
          return !(currentStep === 2 && !v);
        }),
        clusterName: Yup.string().test(
          'required',
          'Cluster name is required',
          (v) => {
            return !(currentStep === 2 && !v);
          },
        ),
        selectedTemplate: Yup.object({}).required('Template is required.'),
        // @ts-ignore
        res: Yup.object({}).test({
          name: 'res',
          skipAbsent: true,
          test(value, ctx) {
            // eslint-disable-next-line react/no-this-in-sfc
            const selfValue = this.parent;

            let vs = Yup.object({});

            if (selfValue.selectedTemplate && currentStep === 2) {
              vs = Yup.object(
                flatMapValidations(
                  selfValue.selectedTemplate?.template?.fields.reduce(
                    (acc: any, curr: any) => {
                      return { ...acc, [curr.name]: curr };
                    },
                    {},
                  ),
                ),
              );
            }

            const res = vs.validateSync(value, {
              abortEarly: false,
              context: ctx,
            });

            return res;
          },
        }),
      }),
      onSubmit: async (val) => {
        const selectedTemplate =
          val.selectedTemplate as unknown as ISelectedTemplate;
        const submit = async () => {
          try {
            if (
              !selectedTemplate?.template?.apiVersion ||
              !selectedTemplate?.template?.kind
            ) {
              throw new Error('Service apiversion or kind error.');
            }
            const { errors: e } = await api.createClusterMSv({
              service: {
                displayName: val.displayName,
                metadata: {
                  name: val.name,
                },
                clusterName: val.clusterName,
                spec: {
                  msvcSpec: {
                    nodeSelector: {
                      [keyconstants.nodepoolName]: val.nodepoolName,
                    },
                    serviceTemplate: {
                      apiVersion: selectedTemplate.template.apiVersion,
                      kind: selectedTemplate.template.kind,
                      spec: {
                        ...val.res,
                      },
                    },
                  },
                },
              },
            });
            if (e) {
              throw e[0];
            }
            toast.success('Managed service created successfully');
            navigate(rootUrl);
          } catch (err) {
            handleError(err);
          }
        };

        switch (currentStep) {
          case 1:
            nextStep();
            break;
          case 2:
            nextStep();
            break;
          case 3:
            await submit();
            break;
          default:
            break;
        }
      },
    });

  useEffect(() => {
    const selectedTemplate =
      values.selectedTemplate as unknown as ISelectedTemplate;
    if (selectedTemplate?.template?.fields) {
      setValues((v) => ({
        ...v,
        res: {
          ...flatM(
            selectedTemplate?.template?.fields.reduce((acc, curr) => {
              return { ...acc, [curr.name]: curr };
            }, {}),
          ),
        },
      }));
    }
  }, [values.selectedTemplate]);

  useEffect(() => {
    setValues((v) => ({
      ...v,
      clusterName:
        clusterList.length > 0
          ? clusterList.find((c) => c.ready)?.value || ''
          : '',
    }));
  }, [clusterList]);

  return (
    <MultiStepProgressWrapper
      title="Let’s create new managed service."
      subTitle="Simplify Collaboration and Enhance Productivity with Kloudlite teams"
      backButton={{
        content: 'Back to Managed services',
        to: rootUrl,
      }}
    >
      <MultiStepProgress.Root currentStep={currentStep} jumpStep={jumpStep}>
        <MultiStepProgress.Step label="Select Managed Service" step={1}>
          <TemplateView
            isLoading={isLoading}
            templates={msvtemplates}
            handleChange={handleChange}
            handleSubmit={handleSubmit}
            errors={errors}
            values={values}
          />
        </MultiStepProgress.Step>
        <MultiStepProgress.Step label="Configure managed service" step={2}>
          <FieldView
            selectedTemplate={values.selectedTemplate}
            values={values}
            errors={errors}
            handleChange={handleChange}
            handleSubmit={handleSubmit}
            // nodepools={statefulNodepools}
            // nodepoolIsLoading={nodepoolIsLoading}
            clusters={clusterList}
          />
        </MultiStepProgress.Step>
        <MultiStepProgress.Step label="Review" step={3}>
          <ReviewView
            onEdit={jumpStep}
            values={values}
            handleSubmit={handleSubmit}
            isLoading={isLoading}
          />
        </MultiStepProgress.Step>
      </MultiStepProgress.Root>
    </MultiStepProgressWrapper>
  );
};

const NewManagedService = () => {
  return <ManagedServiceLayout />;
};

export const handle = {
  noMainLayout: true,
};

export default NewManagedService;
