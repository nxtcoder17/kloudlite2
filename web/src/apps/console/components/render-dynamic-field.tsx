import { NumberInput, TextInput } from '@kloudlite/design-system/atoms/input';
import { IMSvTemplate } from '../server/gql/queries/managed-templates-queries';

const RenderDynamicField = ({
  field,
  value,
  onChange,
  error,
  message,
}: {
  field: IMSvTemplate['fields'][number];
  onChange: (e: { target: { value: any } }) => void;
  value: any;
  error?: boolean;
  message?: string;
}) => {
  switch (field.inputType) {
    case 'String':
      return (
        <TextInput
          value={value || ''}
          onChange={onChange}
          error={error}
          message={message}
          label={`${field.label}${field.required ? ' *' : ''}`}
          suffix={field.unit}
        />
      );
    case 'Number':
      return (
        <NumberInput
          error={error}
          message={message}
          label={`${field.label}${field.required ? ' *' : ''}`}
          min={field.min}
          max={field.max}
          placeholder={field.label}
          value={value || ''}
          onChange={onChange}
          suffix={field.unit}
        />
      );
    default:
      return <div>unknown field type: {field.inputType}</div>;
  }
};

export default RenderDynamicField;
