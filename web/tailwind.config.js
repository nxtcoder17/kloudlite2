// eslint-disable-next-line import/no-relative-packages
import tailwindBase from './tailwind-base.js';

const app = process.env.APP;

export default {
  ...tailwindBase,
  content: [
    './src/design-system/components/**/*.{js,ts,jsx,tsx,mdx}',
    `./src/apps/${app}/**/*.{js,ts,jsx,tsx,mdx}`,
    `./lib/**/*.{js,ts,jsx,tsx,mdx}`,
  ],
};
