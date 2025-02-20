/* eslint-disable guard-for-in */
/* eslint-disable react/destructuring-assignment */
import { useEffect, useRef, useState } from 'react';
import {
  ITextInputBase,
  NumberInput,
  TextInput,
} from '@kloudlite/design-system/atoms/input';
import { Switch } from '@kloudlite/design-system/atoms/switch';
import Popup from '@kloudlite/design-system/molecule/popup';
import { NameIdView } from '~/console/components/name-id-view';
import { IDialogBase } from '~/console/components/types.d';
import { useConsoleApi } from '~/console/server/gql/api-provider';
import { IClusterMSvs } from '~/console/server/gql/queries/cluster-managed-services-queries';
import { IMSvTemplates } from '~/console/server/gql/queries/managed-templates-queries';
import { ExtractNodeType, parseName } from '~/console/server/r-utils/common';
import { getManagedTemplate } from '~/console/utils/commons';
import { useReload } from '~/root/lib/client/helpers/reloader';
import useForm, { dummyEvent } from '~/root/lib/client/hooks/use-form';
import Yup from '~/root/lib/server/helpers/yup';
import { NN } from '~/root/lib/types/common';
import { handleError } from '~/root/lib/utils/common';

type IDialog = IDialogBase<ExtractNodeType<IClusterMSvs>> & {
  templates: IMSvTemplates;
};

type ISelectedService = {
  category: {
    name: string;
    displayName: string;
  };

  service?: NN<IMSvTemplates>[number]['items'][number];
} | null;

const RenderField = ({
  field,
  value,
  onChange,
  errors,
  fieldKey,
}: {
  field: NN<NN<ISelectedService>['service']>['fields'][number];
  onChange: (e: string) => (e: { target: { value: any } }) => void;
  value: any;
  errors: {
    [key: string]: string | undefined;
  };
  fieldKey: string;
}) => {
  const [qos, setQos] = useState(false);

  useEffect(() => {
    if (field.inputType === 'Resource' && value.max === value.min) {
      setQos(true);
    }
  }, []);

  if (field.inputType === 'Number') {
    return (
      <NumberInput
        error={!!errors[fieldKey]}
        message={errors[fieldKey]}
        label={`${field.label}${field.required ? ' *' : ''}`}
        placeholder={field.label}
        value={parseFloat(value) / (field.multiplier || 1) || ''}
        onChange={({ target }) => {
          onChange(`res.${field.name}`)(
            dummyEvent(
              `${parseFloat(target.value) * (field.multiplier || 1)}${
                field.unit
              }`
            )
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
                error={!!errors[`${fieldKey}.min`]}
                message={errors[`${fieldKey}.min`]}
                placeholder={qos ? field.label : `${field.label} min`}
                value={parseFloat(value.min) / (field.multiplier || 1)}
                onChange={({ target }) => {
                  onChange(`res.${field.name}.min`)(
                    dummyEvent(
                      `${parseFloat(target.value) * (field.multiplier || 1)}${
                        field.unit
                      }`
                    )
                  );
                  if (qos) {
                    onChange(`res.${field.name}.max`)(
                      dummyEvent(
                        `${parseFloat(target.value) * (field.multiplier || 1)}${
                          field.unit
                        }`
                      )
                    );
                  }
                }}
                suffix={field.displayUnit}
              />
            </div>
            {!qos && (
              <div className="flex-1">
                <NumberInput
                  error={!!errors[`${fieldKey}.max`]}
                  message={errors[`${fieldKey}.max`]}
                  placeholder={`${field.label} max`}
                  value={parseFloat(value.max) / (field.multiplier || 1)}
                  onChange={({ target }) => {
                    onChange(`res.${field.name}.max`)(
                      dummyEvent(
                        `${parseFloat(target.value) * (field.multiplier || 1)}${
                          field.unit
                        }`
                      )
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

export const Fill = ({
  selectedService,
  values,
  handleChange,
  errors,
  size = 'lg',
}: {
  selectedService: ISelectedService;
  values: { [key: string]: any };
  handleChange: (key: string) => (e: { target: { value: any } }) => void;
  errors: {
    [key: string]: string | undefined;
  };
  size?: ITextInputBase['size'];
}) => {
  const nameRef = useRef<HTMLInputElement>(null);
  useEffect(() => {
    nameRef.current?.focus();
  }, [nameRef.current]);
  return (
    <div className="flex flex-col gap-3xl min-h-[30vh]">
      <NameIdView
        isUpdate
        size={size}
        ref={nameRef}
        placeholder="Enter managed service name"
        label="Name"
        resType="managed_service"
        name={values.name}
        displayName={values.displayName}
        errors={errors.name}
        handleChange={handleChange}
        nameErrorLabel="isNameError"
      />

      {selectedService?.service?.fields.map((field) => {
        const k = field.name;
        const x = k.split('.').reduce((acc, curr) => {
          if (!acc) {
            return values.res[curr];
          }

          return acc[curr];
        }, null);

        return (
          <RenderField
            errors={errors}
            value={x}
            onChange={handleChange}
            key={k}
            field={field}
            fieldKey={k}
          />
        );
      })}
    </div>
  );
};

const Root = (props: IDialog) => {
  const { isUpdate, setVisible, templates } = props;

  const api = useConsoleApi();
  const reload = useReload();

  // const { cluster } = useOutletContext<IClusterContext>();

  const { values, errors, handleChange, handleSubmit, isLoading } = useForm({
    initialValues: isUpdate
      ? {
          name: parseName(props.data),
          displayName: props.data.displayName,
          clusterName: props.data.clusterName,
          isNameError: false,
          res: {
            ...props.data.spec?.msvcSpec.serviceTemplate.spec,
          },
        }
      : {
          name: '',
          displayName: '',
          clusterName: '',
          res: {},
          isNameError: false,
        },
    validationSchema: Yup.object({}),
    onSubmit: async (val) => {
      if (isUpdate) {
        try {
          const { errors: e } = await api.updateClusterMSv({
            service: {
              displayName: val.displayName,
              metadata: {
                name: val.name,
              },
              clusterName: val.clusterName,
              spec: {
                msvcSpec: {
                  serviceTemplate: {
                    apiVersion:
                      props.data.spec?.msvcSpec.serviceTemplate.apiVersion ||
                      '',
                    kind: props.data.spec?.msvcSpec.serviceTemplate.kind || '',
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
          setVisible(false);
          reload();
        } catch (err) {
          handleError(err);
        }
      }
    },
  });

  const getService = () => {
    if (isUpdate)
      return getManagedTemplate({
        templates,
        apiVersion: props.data.spec?.msvcSpec.serviceTemplate.apiVersion || '',
        kind: props.data.spec?.msvcSpec.serviceTemplate.kind || '',
      });
    return undefined;
  };

  if (!isUpdate) {
    return null;
  }
  return (
    <Popup.Form
      onSubmit={(e) => {
        handleSubmit(e);
      }}
    >
      <Popup.Content className="!min-h-[500px] !max-h-[500px]">
        <Fill
          {...{
            templates,
            selectedService: {
              category: { displayName: '', name: '' },
              service: getService(),
            },
            values,
            errors,
            handleChange,
          }}
        />
      </Popup.Content>
      <Popup.Footer>
        <Popup.Button type="button" variant="basic" content="Cancel" closable />
        <Popup.Button
          loading={isLoading}
          type="submit"
          content="Update"
          variant="primary"
        />
      </Popup.Footer>
    </Popup.Form>
  );
};

const HandleBackendService = (props: IDialog) => {
  const { isUpdate, setVisible, visible } = props;
  return (
    <Popup.Root show={visible} onOpenChange={(v) => setVisible(v)}>
      <Popup.Header>
        {isUpdate ? 'Edit managed service' : 'Add managed service'}
      </Popup.Header>
      {(!isUpdate || (isUpdate && props.data)) && <Root {...props} />}
    </Popup.Root>
  );
};

export default HandleBackendService;
