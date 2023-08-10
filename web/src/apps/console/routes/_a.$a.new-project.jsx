import { ArrowLeft, ArrowRight, CircleDashed } from '@jengaicons/react';
import { Button } from '~/components/atoms/button';
import { TextInput } from '~/components/atoms/input';
import { BrandLogo } from '~/components/branding/brand-logo';
import { ProgressTracker } from '~/components/organisms/progress-tracker';
import * as AlertDialog from '~/components/molecule/alert-dialog';
import { useState } from 'react';
import {
  useLoaderData,
  useNavigate,
  useOutletContext,
  useParams,
} from '@remix-run/react';
import * as Radio from '~/components/atoms/radio';
import useForm, { dummyEvent } from '~/root/lib/client/hooks/use-form';
import Yup from '~/root/lib/server/helpers/yup';
import * as Tooltip from '~/components/atoms/tooltip';
import { toast } from '~/components/molecule/toast';
import logger from '~/root/lib/client/helpers/log';
import { dayjs } from '~/components/molecule/dayjs';
import { useAPIClient } from '~/root/lib/client/hooks/api-provider';
import {
  ensureAccountClientSide,
  ensureAccountSet,
  ensureClusterClientSide,
} from '../server/utils/auth-utils';
import {
  getMetadata,
  getPagination,
  getSearch,
  parseName,
  parseUpdationTime,
} from '../server/r-urils/common';
import { GQLServerHandler } from '../server/gql/saved-queries';
import { IdSelector } from '../components/id-selector';
import { SearchBox } from '../components/search-box';
import { getProject, getProjectSepc } from '../server/r-urils/project';
import { keyconstants } from '../server/r-urils/key-constants';

const NewProject = () => {
  const { clustersData } = useLoaderData();
  const clusters = clustersData?.edges?.map(({ node }) => node || []);

  const api = useAPIClient();
  const navigate = useNavigate();

  const [showUnsavedChanges, setShowUnsavedChanges] = useState(false);

  // @ts-ignore
  const { user } = useOutletContext();
  const { a: account } = useParams();

  const { values, handleSubmit, handleChange, isLoading } = useForm({
    initialValues: {
      name: '',
      displayName: '',
      clusterName: '',
    },
    validationSchema: Yup.object({
      name: Yup.string().required(),
      displayName: Yup.string().required(),
      clusterName: Yup.string().required(),
    }),
    onSubmit: async (val) => {
      try {
        ensureClusterClientSide({ cluster: val.clusterName });
        ensureAccountClientSide({ account });
        const { errors: e } = await api.createProject({
          project: getProject({
            metadata: getMetadata({
              name: val.name,
              annotations: {
                [keyconstants.displayName]: val.displayName,
                [keyconstants.author]: user.name,
                [keyconstants.node_type]: val.node_type,
              },
            }),
            spec: getProjectSepc({
              clusterName: val.clusterName,
              displayName: val.displayName,
              accountName: account,
              targetNamespace: val.name,
            }),
          }),
        });

        if (e) {
          throw e[0];
        }
        toast.success('project added successfully');
        navigate('/projects');
      } catch (err) {
        toast.error(err.message);
      }
    },
  });

  return (
    <Tooltip.TooltipProvider>
      <div className="h-full flex flex-row">
        <div className="h-full w-[571px] flex flex-col bg-surface-basic-subdued py-11xl px-10xl">
          <div className="flex flex-col gap-8xl">
            <div className="flex flex-col gap-4xl items-start">
              <BrandLogo detailed={false} size={48} />
              <div className="flex flex-col gap-3xl">
                <div className="text-text-default heading4xl">
                  Let’s create new project.
                </div>
                <div className="text-text-default bodyLg">
                  Create your project to production effortlessly
                </div>
              </div>
            </div>
            <ProgressTracker
              items={[
                { label: 'Configure project', active: true, id: 1 },
                { label: 'Review', active: false, id: 2 },
              ]}
            />
            <Button
              variant="outline"
              content="Back"
              prefix={ArrowLeft}
              onClick={() => setShowUnsavedChanges(true)}
            />
          </div>
        </div>
        <form className="py-11xl px-10xl flex-1" onSubmit={handleSubmit}>
          <div className="gap-6xl flex flex-col p-3xl">
            <div className="flex flex-col gap-4xl">
              <div className="h-7xl" />
              <div className="flex flex-col gap-3xl">
                <TextInput
                  label="Project name"
                  name="name"
                  value={values.displayName}
                  onChange={handleChange('displayName')}
                />
                <IdSelector
                  name={values.displayName}
                  onChange={(v) => {
                    handleChange('name')(dummyEvent(v));
                  }}
                />
              </div>
            </div>
            <div className="flex flex-col border border-border-disabled bg-surface-basic-default rounded-md">
              <SearchBox InputElement={TextInput} />
              <Radio.RadioGroup
                value={values.clusterName}
                onChange={(e) => {
                  handleChange('clusterName')(dummyEvent(e));
                }}
                className="flex flex-col pr-2xl !gap-y-0"
                labelPlacement="left"
              >
                {clusters.map((cluster) => {
                  return (
                    <Radio.RadioItem
                      value={parseName(cluster)}
                      withBounceEffect={false}
                      className="justify-between w-full"
                      key={parseName(cluster)}
                    >
                      <div className="p-2xl pl-lg flex flex-row gap-lg items-center">
                        <CircleDashed size={24} />
                        <div className="flex flex-row flex-1 items-center gap-lg">
                          <span className="headingMd text-text-default">
                            {parseName(cluster)}
                          </span>
                          <span className="bodyMd text-text-default ">
                            {dayjs(parseUpdationTime(cluster)).fromNow()}
                          </span>
                        </div>
                      </div>
                    </Radio.RadioItem>
                  );
                })}
              </Radio.RadioGroup>
            </div>
            <div className="flex flex-row justify-end">
              <Button
                loading={isLoading}
                variant="primary"
                content="Create"
                suffix={ArrowRight}
                type="submit"
              />
            </div>
          </div>
        </form>

        {/* Unsaved change alert dialog */}
        <AlertDialog.DialogRoot
          show={showUnsavedChanges}
          onOpenChange={setShowUnsavedChanges}
        >
          <AlertDialog.Header>
            Leave page with unsaved changes?
          </AlertDialog.Header>
          <AlertDialog.Content>
            Leaving this page will delete all unsaved changes.
          </AlertDialog.Content>
          <AlertDialog.Footer>
            <AlertDialog.Button variant="basic" content="Cancel" />
            <AlertDialog.Button variant="critical" content="Delete" />
          </AlertDialog.Footer>
        </AlertDialog.DialogRoot>
      </div>
    </Tooltip.TooltipProvider>
  );
};

export const loader = async (ctx) => {
  ensureAccountSet(ctx);
  const { data, errors } = await GQLServerHandler(ctx.request).listClusters({
    pagination: getPagination(ctx),
    search: getSearch(ctx),
  });

  if (errors) {
    logger.error(errors);
  }

  return {
    clustersData: data,
  };
};

export default NewProject;
