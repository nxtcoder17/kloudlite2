/* eslint-disable jsx-a11y/control-has-associated-label */
import { useOutletContext } from '@remix-run/react';
import { Box } from '~/console/components/common-console-components';
import { parseName } from '~/console/server/r-utils/common';
import { DownloadSimple } from '@jengaicons/react';
import { downloadFile, renderCloudProvider } from '~/console/utils/commons';
import { Chip } from '~/components/atoms/chips';
import {
  CopyButton,
  DetailItem,
  InfoLabel,
} from '~/console/components/commons';
import { useConsoleApi } from '~/console/server/gql/api-provider';
import { handleError } from '~/root/lib/utils/common';
import { toast } from '~/components/molecule/toast';
import { useState } from 'react';
import Wrapper from '~/console/components/wrapper';
import { IClusterContext } from '../_.$account.$cluster';

const KubeConfigDownload = ({ cluster }: { cluster: string }) => {
  const filename = `${cluster} kubeconfig.yaml`;
  const api = useConsoleApi();

  const [loading, setLoading] = useState(false);

  const downloadConfig = async () => {
    setLoading(true);
    try {
      const { errors, data } = await api.getKubeConfig({
        name: cluster,
      });
      if (errors) {
        throw errors[0];
      }
      if (data.adminKubeconfig) {
        const { encoding, value } = data.adminKubeconfig;
        let linkSource = '';
        switch (encoding) {
          case 'base64':
            linkSource = atob(value);
            break;
          default:
            linkSource = value;
        }

        downloadFile({
          filename,
          data: linkSource,
          format: 'text/plain',
        });
      } else {
        toast.error('Kubeconfig not found.');
      }
    } catch (err) {
      handleError(err);
    } finally {
      setLoading(false);
    }
  };
  return (
    <Chip
      type="CLICKABLE"
      item={cluster}
      label="Download"
      prefix={<DownloadSimple />}
      loading={loading}
      onClick={() => {
        downloadConfig();
      }}
    />
  );
};

const ClusterInfo = () => {
  const { cluster } = useOutletContext<IClusterContext>();

  const providerInfo = () => {
    const provider = cluster.spec?.cloudProvider;
    switch (provider) {
      case 'aws':
        return (
          <DetailItem
            title="Region"
            value={
              <CopyButton
                title={cluster.spec?.aws?.region || ''}
                value={cluster.spec?.aws?.region || ''}
              />
            }
          />
        );
      default:
        return null;
    }
  };
  return (
    <Box title={cluster.displayName}>
      <div className="flex flex-col">
        <div className="flex flex-row gap-3xl flex-wrap">
          <DetailItem
            title="Cluster ID"
            value={
              <CopyButton
                title={parseName(cluster)}
                value={parseName(cluster)}
              />
            }
          />
          <DetailItem
            title="kubeconfig"
            value={<KubeConfigDownload cluster={parseName(cluster)} />}
          />
          <DetailItem
            title="Availability mode"
            value={
              <InfoLabel
                title="Development mode"
                info={
                  <div>
                    In Development mode, we will operate your Kubernetes cluster
                    with the control plane running on a{' '}
                    <span className="bodySm-semibold">single master</span> node.
                  </div>
                }
                label={
                  cluster.spec?.availabilityMode === 'dev'
                    ? 'Development'
                    : 'Highly Available' || ''
                }
              />
            }
          />
          <DetailItem
            title="Kloudlite Release"
            value={
              <CopyButton
                title={cluster.spec?.kloudliteRelease || ''}
                value={cluster.spec?.kloudliteRelease || ''}
              />
            }
          />
          <DetailItem
            title="Public DNS Host"
            value={
              <CopyButton
                title={cluster.spec?.publicDNSHost || ''}
                value={cluster.spec?.publicDNSHost || ''}
              />
            }
          />

          <DetailItem
            title="Cloud provider"
            value={renderCloudProvider({
              cloudprovider: cluster.spec?.cloudProvider || 'unknown',
            })}
          />
          {providerInfo()}
        </div>
      </div>
    </Box>
  );
};
export default ClusterInfo;
