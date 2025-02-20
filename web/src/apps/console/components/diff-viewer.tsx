/* eslint-disable react/jsx-pascal-case */
import hljs from 'highlight.js';
import yaml from 'js-yaml';
import classNames from 'classnames';
import { useEffect, useRef } from 'react';
import ReactDiffViewerBase, {
  ReactDiffViewerProps,
} from 'react-diff-viewer-continued';

export const yamlDump = (data: any) => {
  return yaml.dump(data);
};

export const yamlParse = (data: string) => {
  let resp;
  try {
    resp = yaml.load(data) as any;
  } catch (err) {
    return {
      data: null,
      error: (err as Error).message,
    };
  }

  return {
    data: resp,
    error: null,
  };
};
const HighlightIt = ({
  language,
  inlineData = '',
  className = '',
}: {
  language: string;
  inlineData?: string;
  className?: string;
}) => {
  const ref = useRef(null);

  useEffect(() => {
    (async () => {
      if (ref.current) {
        // @ts-ignore
        ref.current.innerHTML = hljs.highlight(inlineData, {
          language,
        }).value;
      }
    })();
  }, [inlineData, language]);

  return (
    <div ref={ref} className={classNames(className, 'inline')}>
      {inlineData}
    </div>
  );
};

export const DiffViewer = (props: ReactDiffViewerProps): JSX.Element => {
  const highlightSyntax = (str: string) => {
    return <HighlightIt language="yaml" inlineData={str} />;
  };

  return (
    <div className="p-sm border rounded-sm border-text-primary theme-atom-one-dark">
      {/* @ts-ignore */}
      <ReactDiffViewerBase.default {...props} renderContent={highlightSyntax} />
    </div>
  );
};

const _Test = () => {
  return <DiffViewer newValue="" oldValue="" useDarkTheme />;
};
