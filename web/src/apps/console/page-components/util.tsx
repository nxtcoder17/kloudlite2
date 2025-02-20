import { motion } from 'framer-motion';
import { FormEventHandler } from 'react';
import { ChildrenProps } from '@kloudlite/design-system/types';
import { cn } from '@kloudlite/design-system/utils';
import { InputMaybe } from '~/root/src/generated/gql/server';

export const FadeIn = ({
  children,
  className = '',
  onSubmit,
  notForm = false,
}: ChildrenProps & {
  className?: string;
  onSubmit?: FormEventHandler<HTMLFormElement>;
  notForm?: boolean;
}) => {
  if (notForm) {
    return (
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ ease: 'linear', duration: 0.2 }}
        className={cn('flex flex-col gap-6xl w-full justify-center', className)}
      >
        {children}
      </motion.div>
    );
  }
  return (
    <motion.form
      onSubmit={onSubmit}
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ ease: 'linear', duration: 0.2 }}
      className={cn('flex flex-col gap-6xl w-full justify-center', className)}
    >
      {children}
    </motion.form>
  );
};

export function parseValue<T>(v: any, def: T): T {
  try {
    switch (typeof def) {
      case 'number':
        const res = parseFloat(v);
        if (Number.isNaN(res)) {
          return def;
        }
        return res as T;
      default:
        return def;
    }
  } catch (_) {
    return def;
  }
}

export type ExtractArrayType<T> = T extends (infer U)[] ? U : never;

export type ExtractInputMaybe<Type> = Type extends InputMaybe<infer U>
  ? U
  : never;
