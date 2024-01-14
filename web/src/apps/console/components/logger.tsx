/* eslint-disable react/no-unused-prop-types */
/* eslint-disable no-nested-ternary */
import { ArrowsIn, ArrowsOut, List } from '@jengaicons/react';
import Anser from 'anser';
import classNames from 'classnames';
import Fuse from 'fuse.js';
import hljs from 'highlight.js';
import React, { ReactNode, useEffect, useRef, useState } from 'react';
import { ViewportList } from 'react-viewport-list';
import * as sock from 'websocket';
import { dayjs } from '~/components/molecule/dayjs';
import {
  ISearchInfProps,
  useSearch,
} from '~/root/lib/client/helpers/search-filter';
import useClass from '~/root/lib/client/hooks/use-class';
import logger from '~/root/lib/client/helpers/log';
import { socketUrl } from '~/root/lib/configs/base-url.cjs';
import generateColor from './color-generator';
import Pulsable from './pulsable';
import { logsMockData } from '../dummy/data';

const hoverClass = `hover:bg-[#ddd]`;
const hoverClassDark = `hover:bg-[#333]`;

type ILog = {
  pod_name: string;
  message: string;
  timestamp: string;
};
type ILogWithLineNumber = ILog & { lineNumber: number };

type ISocketMessage = ILog;

const padLeadingZeros = (num: number, size: number) => {
  let s = `${num}`;
  while (s.length < size) s = `0${s}`;
  return s;
};

const threshold = 0.4;

interface IHighlightIt {
  language: string;
  inlineData: string;
  className?: string;
  enableHL?: boolean;
}

const getHashId = (str: string) => {
  let hash = 0;
  let i;
  let chr;
  if (str.length === 0) return hash;
  for (i = 0; i < str.length; i += 1) {
    chr = str.charCodeAt(i);
    // eslint-disable-next-line no-bitwise
    hash = (hash << 5) - hash + chr;
    // eslint-disable-next-line no-bitwise
    hash |= 0; // Convert to 32bit integer
  }
  return hash;
};

const HighlightIt = ({
  language,
  inlineData = '',
  className = '',
  enableHL = false,
}: IHighlightIt) => {
  const ref = useRef(null);
  const data = Anser.ansiToText(inlineData);

  useEffect(() => {
    (async () => {
      if (ref.current) {
        if (enableHL) {
          // if (!isScrolledIntoView(ref.current)) return;
          // @ts-ignore
          ref.current.innerHTML = hljs.highlight(
            data,
            {
              language,
            },
            false
          ).value;
        } else {
          // @ts-ignore
          ref.current.innerHTML = Anser.ansiToHtml(inlineData);
        }
      }
    })();
  }, [inlineData, language]);

  return (
    <div ref={ref} className={classNames(className, 'inline')}>
      {data}
    </div>
  );
};

interface ILineNumber {
  lineNumber: number;
  fontSize: number;
  lines: number;
}
const LineNumber = ({ lineNumber, fontSize, lines }: ILineNumber) => {
  const ref = useRef(null);
  const [data, setData] = useState(() => padLeadingZeros(1, `${lines}`.length));

  useEffect(() => {
    setData(padLeadingZeros(lineNumber, `${lines}`.length));
  }, [lines, lineNumber]);
  return (
    <code
      key={`ind+${lineNumber}`}
      className="inline-flex gap-xl items-center whitespace-pre select-none pulsable"
      ref={ref}
    >
      <span className="flex sticky left-0" style={{ fontSize }}>
        <HighlightIt
          enableHL
          inlineData={data}
          language="accesslog"
          className={classNames('border-b border-border-tertiary px-md')}
        />
        <div className="hljs" />
      </span>
    </code>
  );
};

interface IFilterdHighlightIt {
  searchInf?: ISearchInfProps['searchInf'];
  inlineData: string;
  className?: string;
  language: string;
  searchText: string;
  showAll: boolean;
}

interface HighlightProps {
  value: string;
  indices: Array<[number, number]>;
}

const Highlighter: React.FC<HighlightProps> = ({ value, indices }) => {
  let lastIndex = 0;
  const parts = [];

  indices.forEach(([start, end]) => {
    if (lastIndex !== start) {
      parts.push(
        <span style={{ opacity: 0.7 }} key={lastIndex}>
          <HighlightIt
            language="accesslog"
            inlineData={value.substring(lastIndex, start)}
            enableHL
          />
        </span>
      );
    }
    parts.push(
      <span className="font-bold" key={start}>
        <HighlightIt
          language="accesslog"
          inlineData={value.substring(start, end + 1)}
          enableHL
        />
      </span>
    );
    lastIndex = end + 1;
  });

  if (lastIndex !== value.length) {
    parts.push(<span key={lastIndex}>{value.substring(lastIndex)}</span>);
  }

  return parts;
};

const InlineSearch = ({
  inlineData = '',
  className = '',
  language,
  searchText,
}: IFilterdHighlightIt) => {
  const res = useSearch(
    {
      data: [{ message: inlineData }],
      keys: ['message'],
      searchText,
      threshold,
      remainOrder: true,
    },
    [inlineData, searchText]
  );

  if (res.length && res[0].searchInf.matches?.length) {
    const def: Fuse.RangeTuple[] = [];
    return (
      <Highlighter
        {...{
          value: inlineData,
          indices:
            res[0].searchInf.matches?.reduce((acc, curr) => {
              return [...acc, ...curr.indices.filter((i) => i[1] - i[0] > 1)];
            }, def) || def,
        }}
      />
    );
  }
  return (
    <HighlightIt
      {...{
        inlineData,
        language,
        className: classNames(className, {
          'opacity-40': !!searchText,
        }),
        enableHL: true,
      }}
    />
  );
};

const FilterdHighlightIt = ({
  searchInf,
  inlineData = '',
  className = '',
  language,
  searchText,
  showAll,
}: IFilterdHighlightIt) => {
  const def: Fuse.RangeTuple[] = [];

  if (showAll) {
    return (
      <div className={classNames('whitespace-pre', className)}>
        <InlineSearch
          {...{
            language,
            inlineData,
            searchText,
            className,
            showAll,
          }}
        />
      </div>
    );
  }

  return (
    <div className={classNames('whitespace-pre', className)}>
      {searchInf?.matches?.length ? (
        <Highlighter
          key={inlineData}
          {...{
            value: inlineData,
            indices: searchInf.matches.reduce((acc, curr) => {
              // const validIndices = curr.indices.filter((i) => {
              //   return i[1] - i[0] >= searchText.length - 1;
              // });
              // console.log(curr.indices, validIndices);
              return [...acc, ...curr.indices];
            }, def),
          }}
        />
      ) : (
        <HighlightIt
          {...{
            inlineData,
            language,
            enableHL: true,
          }}
        />
      )}
    </div>
  );
};

interface ILogLine {
  fontSize: number;
  selectableLines: boolean;
  showAll: boolean;
  searchText: string;
  language: string;
  lines: number;
  hideLines?: boolean;
  log: ILogWithLineNumber & {
    searchInf?: ISearchInfProps['searchInf'];
  };
  dark: boolean;
}

const LogLine = ({
  log,
  fontSize,
  selectableLines,
  showAll,
  searchText,
  language,
  lines,
  hideLines,
  dark,
}: ILogLine) => {
  return (
    <code
      title={`pod: ${log.pod_name}`}
      className={classNames(
        'flex py-xs items-center whitespace-pre border-b border-transparent transition-all',
        {
          'cursor-pointer': selectableLines,
          [hoverClass]: selectableLines && !dark,
          [hoverClassDark]: selectableLines && dark,
        }
      )}
      style={{
        fontSize,
        paddingLeft: fontSize / 4,
        paddingRight: fontSize / 2,
      }}
    >
      {!hideLines && (
        <LineNumber
          lineNumber={log.lineNumber}
          lines={lines}
          fontSize={fontSize}
        />
      )}

      <div
        className="w-[3px] mr-xl ml-sm h-full pulsable pulsable-hidden"
        style={{ backgroundImage: generateColor(log.pod_name) }}
      />
      <div className="inline-flex gap-xl pulsable">
        <HighlightIt
          {...{
            className: 'select-none',
            inlineData: `${dayjs(log.timestamp).format('lll')} |`,
            language: 'apache',
            enableHL: true,
          }}
        />

        <FilterdHighlightIt
          {...{
            searchText,
            inlineData: log.message,
            searchInf: log.searchInf,
            language,
            showAll,
          }}
        />
      </div>
    </code>
  );
};

interface ILogBlock {
  data: ISocketMessage[];
  maxLines?: number;
  follow: boolean;
  enableSearch: boolean;
  selectableLines: boolean;
  title: ReactNode;
  noScrollBar: boolean;
  fontSize: number;
  actionComponent: ReactNode;
  hideLines: boolean;
  language: string;
  solid: boolean;
  dark: boolean;
}

const LogBlock = ({
  data = [],
  follow,
  enableSearch,
  selectableLines,
  title,
  noScrollBar,
  maxLines,
  fontSize,
  actionComponent,
  hideLines,
  language,
  solid,
  dark,
}: ILogBlock) => {
  const [searchText, setSearchText] = useState('');

  const searchResult = useSearch(
    {
      data,
      keys: ['message'],
      searchText,
      threshold,
      remainOrder: true,
    },
    [data, searchText]
  );

  const [showAll, setShowAll] = useState(true);
  const ref = useRef(null);

  useEffect(() => {
    (async () => {
      if (follow && ref.current) {
        // @ts-ignore
        ref.current.scrollTo(0, ref.current.scrollHeight);
      }
    })();
  }, [data, maxLines]);

  return (
    <div
      className={classNames('hljs p-xs flex flex-col gap-sm h-full', {
        'rounded-md': !solid,
      })}
    >
      <div className="flex justify-between items-center border-b border-border-tertiary p-lg">
        <div className="">{data ? title : 'No logs found'}</div>

        <div className="flex items-center gap-xl">
          {actionComponent}
          {enableSearch && (
            <form
              className="flex gap-xl items-center text-sm"
              onSubmit={(e) => {
                e.preventDefault();
                setShowAll((s) => !s);
              }}
            >
              <input
                className="bg-transparent border border-surface-tertiary-default rounded-md px-lg py-xs w-[10rem]"
                placeholder="Search"
                value={searchText}
                onChange={(e) => setSearchText(e.target.value)}
              />
              <div
                onClick={() => {
                  setShowAll((s) => !s);
                }}
                className="cursor-pointer active:translate-y-[1px] transition-all"
              >
                <span
                  className={classNames('font-medium', {
                    'opacity-50': showAll,
                    'text-text-secondary': !showAll,
                  })}
                >
                  <List color="currentColor" size={16} />
                </span>
              </div>
              <code className={classNames('text-xs font-bold', {})}>
                {padLeadingZeros(searchResult.length, 2)} matches
              </code>
            </form>
          )}
        </div>
      </div>

      <div
        className={classNames('flex flex-1 overflow-auto', {
          'no-scroll-bar': noScrollBar,
          'hljs-log-scrollbar': !noScrollBar,
        })}
      >
        <div className="flex flex-1 h-full">
          <div
            className="flex-1 flex flex-col pb-8"
            style={{ lineHeight: `${fontSize * 1.5}px` }}
            ref={ref}
          >
            <ViewportList items={showAll ? data : searchResult}>
              {(log, index) => {
                return (
                  <LogLine
                    dark={dark}
                    log={{
                      ...log,
                      lineNumber: index + 1,
                    }}
                    language={language}
                    searchText={searchText}
                    fontSize={fontSize}
                    lines={data.length}
                    showAll={showAll}
                    key={getHashId(
                      `${log.message}${log.timestamp}${log.pod_name}${index}`
                    )}
                    hideLines={hideLines}
                    selectableLines={selectableLines}
                  />
                );
              }}
            </ViewportList>
          </div>
        </div>
      </div>
    </div>
  );
};

interface IuseLog {
  url?: string;
  account: string;
  cluster: string;
  trackingId: string;
}

export interface IHighlightJsLog {
  websocket: IuseLog;
  follow?: boolean;
  url?: string;
  text?: string;
  enableSearch?: boolean;
  selectableLines?: boolean;
  title?: ReactNode;
  height?: string;
  width?: string;
  noScrollBar?: boolean;
  maxLines?: number;
  fontSize?: number;
  loadingComponent?: ReactNode;
  actionComponent?: ReactNode;
  hideLines?: boolean;
  language?: string;
  solid?: boolean;
  className?: string;
  dark?: boolean;
}

const useSocketLogs = ({ url, account, cluster, trackingId }: IuseLog) => {
  const [logs, setLogs] = useState<ISocketMessage[]>([]);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(true);

  let wsclient: Promise<sock.w3cwebsocket>;

  const [socState, setSocState] = useState<sock.w3cwebsocket | null>(null);

  if (typeof window !== 'undefined') {
    try {
      wsclient = new Promise<sock.w3cwebsocket>((res, rej) => {
        try {
          // eslint-disable-next-line new-cap
          const w = new sock.w3cwebsocket(
            url || `${socketUrl}/logs`,
            '',
            '',
            {}
          );

          w.onmessage = (msg) => {
            try {
              const m: {
                timestamp: string;
                message: string;
                type: 'update' | 'error' | 'info';
              } = JSON.parse(msg.data as string);

              if (m.type === 'error') {
                setLogs([]);
                console.error(m.message);
                return;
              }

              if (m.type === 'info') {
                console.log(m.message);
                return;
              }

              if (m.type === 'update') {
                console.log(m.message);
                return;
              }

              if (m.type === 'log') {
                setIsLoading(false);
                setLogs((s) => [
                  ...s,
                  {
                    pod_name: 'main',
                    message: m.message,
                    timestamp: m.timestamp,
                  },
                ]);
                return;
              }

              console.log(m);
            } catch (err) {
              console.error(err);
            }
          };

          w.onopen = () => {
            res(w);
          };

          w.onerror = (e) => {
            rej(e);
          };

          w.onclose = () => {
            // wsclient.send(newMessage({ event: 'unsubscribe', data: 'test' }));
            logger.log('socket disconnected');
          };
        } catch (e) {
          rej(e);
        }
      });
    } catch (e) {
      console.log(e);
    }
  }

  useEffect(() => {
    (async () => {
      try {
        // const client = await wsclient.;
        // 'wss://auth-vision.devc.kloudlite.io/logs'

        if (account === '' || cluster === '' || trackingId === '') {
          return () => {};
        }

        if (logs.length) {
          setLogs([]);
        }
        setIsLoading(true);

        const client = await wsclient;

        setSocState(client);

        client.send(
          JSON.stringify({
            event: 'subscribe',
            data: {
              account,
              cluster,
              trackingId,
            },
          })
        );
      } catch (e) {
        console.error(e);
        setLogs([]);
        setError((e as Error).message);
      }

      return () => {
        (async () => {
          if (!socState) return;

          socState.send(
            JSON.stringify({
              event: 'unsubscribe',
              data: {
                account,
                cluster,
                trackingId,
              },
            })
          );

          // client.close();

          setLogs([]);
        })();
      };
    })();
  }, [account, cluster, trackingId, url]);

  return {
    logs,
    error,
    isLoading,
  };
};

const LogComp = ({
  websocket,
  follow = true,
  enableSearch = true,
  selectableLines = true,
  title = '',
  height = '400px',
  width = '600px',
  noScrollBar = false,
  maxLines,
  fontSize = 14,
  actionComponent = null,
  hideLines = false,
  language = 'accesslog',
  solid = false,
  className = '',
  dark = false,
}: IHighlightJsLog) => {
  const [fullScreen, setFullScreen] = useState(false);

  const { setClassName, removeClassName } = useClass({
    elementClass: 'loading-container',
  });

  useEffect(() => {
    const keyDownListener = (e: any) => {
      if (e.code === 'Escape') {
        e.stopPropagation();
        setFullScreen(false);
      }
    };

    if (fullScreen && window?.document?.children[0]) {
      // @ts-ignore
      window.document.children[0].style = `overflow-y:hidden`;

      document.addEventListener('keydown', keyDownListener);
    } else if (window?.document?.children[0]) {
      // @ts-ignore
      window.document.children[0].style = `overflow-y:auto`;

      document.removeEventListener('keydown', keyDownListener);
    }
  }, [fullScreen]);

  const { logs, error, isLoading } = useSocketLogs(websocket);

  const [isClientSide, setIsClientSide] = useState(false);

  useEffect(() => {
    if (!isClientSide) {
      setIsClientSide(true);
    }
  }, []);

  return isClientSide ? (
    <Pulsable isLoading={isLoading}>
      <div
        className={classNames(className, {
          'fixed w-full h-full left-0 top-0 z-[999] bg-black': fullScreen,
        })}
        style={{
          width: fullScreen ? '100%' : width,
          height: fullScreen ? '100vh' : height,
        }}
      >
        {error ? (
          <pre>{error}</pre>
        ) : (
          <LogBlock
            {...{
              data: isLoading ? logsMockData : logs,
              follow,
              dark,
              enableSearch,
              selectableLines,
              title,
              noScrollBar,
              solid,
              maxLines,
              fontSize,
              actionComponent: (
                <div className="flex gap-xl">
                  <div
                    onClick={() => {
                      if (!fullScreen) {
                        setClassName('z-50');
                      } else {
                        removeClassName('z-50');
                      }
                      setFullScreen((s) => !s);
                    }}
                    className="flex items-center justify-center font-bold text-xl cursor-pointer select-none active:translate-y-[1px] transition-all"
                  >
                    {fullScreen ? (
                      <ArrowsIn size={16} />
                    ) : (
                      <ArrowsOut size={16} />
                    )}
                  </div>
                  {actionComponent}
                </div>
              ),
              width: fullScreen ? '100vw' : width,
              height: fullScreen ? '100vh' : height,
              hideLines,
              language,
            }}
          />
        )}
      </div>
    </Pulsable>
  ) : (
    <div
      className={classNames(className, {
        'fixed w-full h-full left-0 top-0 z-[999] bg-black': fullScreen,
      })}
      style={{
        width: fullScreen ? '100%' : width,
        height: fullScreen ? '100vh' : height,
      }}
    />
  );
};

export default LogComp;
