import { useOutletContext } from '@remix-run/react';
import Chart from '~/console/components/charts/charts-client';
import { useState } from 'react';
import { dayjs } from '@kloudlite/design-system/molecule/dayjs';
import { ApexOptions } from 'apexcharts';
import { parseName } from '~/console/server/r-utils/common';
import { useDataState } from '~/console/page-components/common-state';
import LogComp from '~/root/lib/client/components/logger';
import LogAction from '~/console/page-components/log-action';
import { IClusterContext } from '../../_layout';

const LogsAndMetrics = () => {
  const { cluster, account } = useOutletContext<IClusterContext>();
  const [cpuData, _setCpuData] = useState<number[]>([]);
  const [memoryData, _setMemoryData] = useState<number[]>([]);

  const xAxisFormatter = (_: string, __?: number) => {
    // return dayjs((val || 0) * 1000).format('hh:mm A');
    return '';
  };

  const tooltipXAixsFormatter = (val: number) =>
    dayjs(val * 1000).format('DD/MM/YY hh:mm A');

  // useDebounce(
  //   () => {
  //     (async () => {
  //       try {
  //         const resp = await axios({
  //           url: `https://observe.dev.kloudlite.io/observability/metrics/cpu?cluster_name=${parseName(
  //             cluster
  //           )}&tracking_id=${cluster.id}`,
  //           method: 'GET',
  //           withCredentials: true,
  //         });
  //
  //         setCpuData(resp?.data?.data?.result[0]?.values || []);
  //       } catch (err) {
  //         console.error(err);
  //       }
  //     })();
  //     (async () => {
  //       try {
  //         const resp = await axios({
  //           url: `https://observe.dev.kloudlite.io/observability/metrics/memory?cluster_name=${parseName(
  //             cluster
  //           )}&tracking_id=${cluster.id}`,
  //           method: 'GET',
  //           withCredentials: true,
  //         });
  //
  //         setMemoryData(resp?.data?.data?.result[0]?.values || []);
  //       } catch (err) {
  //         console.error(err);
  //       }
  //     })();
  //   },
  //   1000,
  //   []
  // );

  const chartOptions: ApexOptions = {
    chart: {
      type: 'area',
      zoom: {
        enabled: false,
      },
      toolbar: {
        show: false,
      },
      redrawOnWindowResize: true,
    },
    dataLabels: {
      enabled: false,
    },
    stroke: {
      curve: 'smooth',
    },

    xaxis: {
      type: 'datetime',
      labels: {
        show: false,
        formatter: xAxisFormatter,
      },
    },
  };

  const { state } = useDataState<{
    linesVisible: boolean;
    timestampVisible: boolean;
  }>('logs');

  return (
    <div className="flex flex-col gap-6xl pt-6xl">
      <div className="gap-6xl items-center flex-col grid sm:grid-cols-2 lg:grid-cols-4">
        <Chart
          disabled
          title="CPU Usage"
          options={{
            ...chartOptions,
            series: [
              {
                color: '#1D4ED8',
                name: 'CPU',
                data: cpuData,
              },
            ],
            tooltip: {
              x: {
                formatter: tooltipXAixsFormatter,
              },
              y: {
                formatter(val) {
                  return `${val.toFixed(2)} m`;
                },
              },
            },
            yaxis: {
              min: 0,

              floating: false,
              labels: {
                formatter: (val) => `${val} m`,
              },
            },
          }}
        />

        <Chart
          disabled
          title="Memory Usage"
          options={{
            ...chartOptions,
            series: [
              {
                color: '#1D4ED8',
                name: 'Memory',
                data: memoryData,
              },
            ],
            yaxis: {
              min: 0,

              floating: false,
              labels: {
                formatter: (val) => `${val} MB`,
              },
            },
            tooltip: {
              x: {
                formatter: tooltipXAixsFormatter,
              },
              y: {
                formatter(val) {
                  return `${val.toFixed(2)} MB`;
                },
              },
            },
          }}
        />

        <Chart
          disabled
          title="Network IO"
          options={{
            ...chartOptions,
            series: [
              {
                color: '#1D4ED8',
                name: 'DiskIO',
                data: memoryData,
              },
            ],
            yaxis: {
              min: 0,

              floating: false,
              labels: {
                formatter: (val) => `${val} MB`,
              },
            },
            tooltip: {
              x: {
                formatter: tooltipXAixsFormatter,
              },
              y: {
                formatter(val) {
                  return `${val.toFixed(2)} MB`;
                },
              },
            },
          }}
        />

        <Chart
          disabled
          title="Disk IO"
          options={{
            ...chartOptions,
            series: [
              {
                color: '#1D4ED8',
                name: 'DiskIO',
                data: memoryData,
              },
            ],
            yaxis: {
              min: 0,

              floating: false,
              labels: {
                formatter: (val) => `${val} MB`,
              },
            },
            tooltip: {
              x: {
                formatter: tooltipXAixsFormatter,
              },
              y: {
                formatter(val) {
                  return `${val.toFixed(2)} MB`;
                },
              },
            },
          }}
        />
      </div>

      <div className="flex-1 max-w-full">
        <LogComp
          {...{
            hideLineNumber: !state.linesVisible,
            hideTimestamp: !state.timestampVisible,
            dark: true,
            width: '100%',
            height: '70vh',
            title: 'Logs',
            actionComponent: <LogAction />,
            websocket: {
              account: parseName(account),
              cluster: parseName(cluster),
              trackingId: cluster.id,
            },
          }}
        />
      </div>
    </div>
  );
};

export default LogsAndMetrics;
