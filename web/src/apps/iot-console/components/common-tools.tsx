import { useSearchParams } from '@remix-run/react';
import { useState } from 'react';
import OptionList from '@kloudlite/design-system/atoms/option-list';
import Toolbar from '@kloudlite/design-system/atoms/toolbar';
import { cn } from '@kloudlite/design-system/utils';
import { CommonFilterOptions } from '~/iotconsole/components/common-filter';
import Filters, {
  FilterType,
  IAppliedFilters,
  useSetAppliedFilters,
} from '~/iotconsole/components/filters';
import { SearchBox } from '~/iotconsole/components/search-box';
import ViewMode from '~/iotconsole/components/view-mode';
import {
  decodeUrl,
  encodeUrl,
  useQueryParameters,
} from '~/root/lib/client/hooks/use-search';
import { NonNullableString } from '~/root/lib/types/common';
import {
  ArrowDown,
  ArrowUp,
  ArrowsDownUp,
  Plus,
} from '~/iotconsole/components/icons';

interface ISortbyOptionList {
  open?: boolean;
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

const SortbyOptionList = (_: ISortbyOptionList) => {
  const { setQueryParameters } = useQueryParameters();
  const [searchparams] = useSearchParams();
  const page = decodeUrl(searchparams.get('page')) || {};

  const { orderBy = 'updateTime', sortDirection = 'DESC' } = page || {};

  const updateOrder = ({ order, direction }: any) => {
    setQueryParameters({
      page: encodeUrl({
        ...page,
        orderBy: order,
        sortDirection: direction,
      }),
    });
  };

  return (
    <OptionList.Root>
      <OptionList.Trigger>
        <Toolbar.Button
          content="Sortby"
          variant="basic"
          prefix={<ArrowsDownUp />}
        />
      </OptionList.Trigger>
      <OptionList.Content>
        <OptionList.RadioGroup
          value={orderBy}
          onValueChange={(v) =>
            updateOrder({
              direction: sortDirection,
              order: v,
            })
          }
        >
          <OptionList.RadioGroupItem
            value="metadata.name"
            onClick={(e) => e.preventDefault()}
          >
            Name
          </OptionList.RadioGroupItem>
          <OptionList.RadioGroupItem
            value="updateTime"
            onClick={(e) => e.preventDefault()}
          >
            Updated
          </OptionList.RadioGroupItem>
        </OptionList.RadioGroup>
        <OptionList.Separator />
        <OptionList.RadioGroup
          value={sortDirection}
          onValueChange={(v) =>
            updateOrder({
              order: orderBy,
              direction: v,
            })
          }
        >
          <OptionList.RadioGroupItem
            showIndicator={false}
            value="ASC"
            onClick={(e) => e.preventDefault()}
          >
            <ArrowUp size={16} />
            {orderBy === 'updateTime' ? 'Oldest' : 'Ascending'}
          </OptionList.RadioGroupItem>
          <OptionList.RadioGroupItem
            value="DESC"
            showIndicator={false}
            onClick={(e) => e.preventDefault()}
          >
            <ArrowDown size={16} />
            {orderBy === 'updateTime' ? 'Newest' : 'Descending'}
          </OptionList.RadioGroupItem>
        </OptionList.RadioGroup>
      </OptionList.Content>
    </OptionList.Root>
  );
};

export interface IModeProps<T = 'list' | 'grid' | NonNullableString> {
  // eslint-disable-next-line react/no-unused-prop-types
  viewMode?: T;
  // eslint-disable-next-line react/no-unused-prop-types
  setViewMode?: (fn: T) => void;
}

interface ICommonTools extends IModeProps {
  options: FilterType[];
  noViewMode?: boolean;
  noSort?: boolean;
}

const CommonTools = ({
  options,
  noViewMode = false,
  noSort = false,
}: ICommonTools) => {
  const [appliedFilters, setAppliedFilters] = useState<IAppliedFilters>({});
  const [sortbyOptionListOpen, setSortybyOptionListOpen] = useState(false);

  useSetAppliedFilters({
    setAppliedFilters,
    types: options,
  });

  return (
    <div className={cn('flex flex-col bg-surface-basic-subdued pb-6xl')}>
      <div>
        {/* Toolbar for md and up */}
        <div className="hidden md:flex">
          <Toolbar.Root>
            <SearchBox />
            <CommonFilterOptions options={options} />
            {!noSort && (
              <SortbyOptionList
                open={sortbyOptionListOpen}
                setOpen={setSortybyOptionListOpen}
              />
            )}
            {!noViewMode && <ViewMode />}
          </Toolbar.Root>
        </div>

        {/* Toolbar for mobile screen */}
        <div className="flex md:hidden">
          <Toolbar.Root>
            <div className="flex-1">
              <SearchBox />
            </div>
            <Toolbar.Button
              content="Add filters"
              prefix={<Plus />}
              variant="basic"
            />
            <SortbyOptionList
              open={sortbyOptionListOpen}
              setOpen={setSortybyOptionListOpen}
            />
          </Toolbar.Root>
        </div>
      </div>

      <Filters appliedFilters={appliedFilters} />
    </div>
  );
};

export default CommonTools;
