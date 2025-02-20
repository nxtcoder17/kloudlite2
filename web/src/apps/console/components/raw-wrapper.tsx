import { Key, ReactNode } from 'react';
import { Button } from '@kloudlite/design-system/atoms/button';
import { BrandLogo } from '@kloudlite/design-system/branding/brand-logo';

import { cn } from '@kloudlite/design-system/utils';
import ProgressTracker from './console-progress-tracker';

interface IRawWrapper<V = any> {
  title: string;
  subtitle: string;
  badge?: {
    title?: string;
    subtitle?: string;
    image?: ReactNode;
  };
  progressItems?: any;
  onProgressClick?: (value: V) => void;
  onCancel?: () => void;
  rightChildren: ReactNode;
}
function RawWrapper<V = any>({
  title,
  subtitle,
  progressItems,
  onProgressClick = () => {},
  onCancel,
  badge,
  rightChildren,
}: IRawWrapper<V>) {
  return (
    <div className="min-h-screen flex flex-row">
      <div className="min-h-full flex flex-col bg-surface-basic-subdued px-11xl pt-11xl pb-10xl">
        <div className="flex flex-col items-start gap-6xl w-[379px]">
          <BrandLogo detailed={false} size={48} />
          <div
            className={cn('flex flex-col', {
              'gap-8xl': !!badge?.title || !!badge?.subtitle,
              'gap-4xl': !badge?.title && !badge?.subtitle,
            })}
          >
            <div className="flex flex-col gap-3xl">
              <div className="text-text-default heading4xl">{title}</div>
              <div className="text-text-default bodyLg">{subtitle}</div>
              {(!!badge?.title || !!badge?.subtitle) && (
                <div className="flex flex-row gap-lg p-lg rounded border border-border-default bg-surface-basic-active min-w-[120px] w-fit">
                  {badge.image && (
                    <div className="p-md text-icon-default flex items-center rounded bg-surface-basic-default">
                      {badge?.image}
                    </div>
                  )}
                  <div className="flex flex-col">
                    <div className="bodySm-semibold text-text-default">
                      {badge?.title}
                    </div>
                    <div className="bodySm text-text-soft">
                      {badge?.subtitle}
                    </div>
                  </div>
                </div>
              )}
            </div>
            {progressItems && (
              <ProgressTracker.Root
                items={progressItems}
                onClick={() => {
                  // onProgressClick(v);
                }}
              />
            )}
          </div>

          {!!onCancel && (
            <Button
              variant="outline"
              content="Cancel"
              size="lg"
              onClick={onCancel}
            />
          )}
        </div>
      </div>
      <div className="pt-11xl pb-12xl px-11xl flex flex-1 bg-surface-basic-default">
        <div className="w-[628px] flex items-center">
          <div className="flex flex-col gap-6xl w-full">{rightChildren}</div>
        </div>
      </div>
    </div>
  );
}

interface ITitleBox {
  title?: ReactNode;
  subtitle?: ReactNode;
}
export const TitleBox = ({ title, subtitle }: ITitleBox) => {
  return (
    <div className="flex flex-col gap-lg">
      {title && <div className="headingXl text-text-default">{title}</div>}
      {subtitle && <div className="bodyMd text-text-soft">{subtitle}</div>}
    </div>
  );
};

export default RawWrapper;
