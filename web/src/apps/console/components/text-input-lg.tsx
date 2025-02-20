import { ITextInput, TextInput } from '@kloudlite/design-system/atoms/input';
import { IconButton } from '@kloudlite/design-system/atoms/button';
import { ArrowRight } from '@kloudlite/design-system/icons';

const TextInputLg = ({
  value,
  onChange,
  onEnter,
  error,
  ...props
}: ITextInput & {
  onEnter?: () => void;
}) => {
  return (
    <div
      style={{ background: 'linear-gradient(#93C5FD, #3B82F6)' }}
      className="p-[2px] rounded-md"
      id="join-waitlist"
    >
      <TextInput
        value={value}
        onChange={onChange}
        placeholder="Enter Code"
        size="xl"
        className="!border-none"
        onKeyDown={(e) => {
          if (e.key === 'Enter') {
            onEnter?.();
            // e.stopPropagation();
            // e.preventDefault();
          }
        }}
        suffix={
          <IconButton
            variant="outline"
            icon={<ArrowRight />}
            onClick={onEnter}
          />
        }
        focusRing={false}
        textFieldClassName="!bodyLg"
        {...props}
      />
    </div>
  );
};

export default TextInputLg;
