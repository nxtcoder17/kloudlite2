import 'kl-design-system/index.css';

import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import 'react-toastify/dist/ReactToastify.css';
import ToastifyContainer from './components/toastify-container';
import './globals.css';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Kloudlite Events',
  description: 'Generated by create next app',
};

const getServerEnv = () => {
  return {
    ...(process.env.MARKETING_API_URL
      ? { MARKETING_API_URL: process.env.MARKETING_API_URL }
      : {}),
    ...(process.env.DYTE_ORG_ID
      ? { DYTE_ORG_ID: process.env.DYTE_ORG_ID }
      : {}),
    ...(process.env.DYTE_API_KEY
      ? { DYTE_API_KEY: process.env.DYTE_API_KEY }
      : {}),
    ...(process.env.DYTE_MEETING_ID
      ? { DYTE_MEETING_ID: process.env.DYTE_MEETING_ID }
      : {}),
  };
};

const getClientEnv = (env: any) => {
  const { DYTE_ORG_ID, MARKETING_API_URL, DYTE_API_KEY, DYTE_MEETING_ID } = env;
  return `
${MARKETING_API_URL ? `window.MARKETING_API_URL = ${`'${MARKETING_API_URL}'`}` : ''}
${DYTE_ORG_ID ? `window.DYTE_ORG_ID = ${`'${DYTE_ORG_ID}'`}` : ''}
${DYTE_API_KEY ? `window.DYTE_API_KEY = ${`'${DYTE_API_KEY}'`}` : ''}
${DYTE_MEETING_ID ? `window.DYTE_MEETING_ID = ${`'${DYTE_MEETING_ID}'`}` : ''}`;
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const env = getServerEnv();
  return (
    <html lang="en">
      <body className={inter.className}>
        <script
          // eslint-disable-next-line react/no-danger
          dangerouslySetInnerHTML={{
            __html: getClientEnv(env),
          }}
        />
        {children}
        <ToastifyContainer />
      </body>

      {/* <body className={inter.className}>
        <main>{children}</main>
      </body> */}
    </html>
  );
}
