/* eslint-disable jsx-a11y/tabindex-no-positive */
import {
  Envelope,
  GithubLogoFill,
  GitlabLogoFill,
  GoogleLogo,
} from '@jengaicons/react';
import { Link, useOutletContext, useSearchParams } from '@remix-run/react';
import { useEffect } from 'react';
import { useAuthApi } from '~/auth/server/gql/api-provider';
import { Button } from '@kloudlite/design-system/atoms/button';
import { PasswordInput, TextInput } from '@kloudlite/design-system/atoms/input';
import { ArrowLeft, ArrowRight } from '@kloudlite/design-system/icons';
import { toast } from '@kloudlite/design-system/molecule/toast';
import { cn } from '@kloudlite/design-system/utils';
import { getCookie } from '~/root/lib/app-setup/cookies';
import { useReload } from '~/root/lib/client/helpers/reloader';
import useForm from '~/root/lib/client/hooks/use-form';
import Yup from '~/root/lib/server/helpers/yup';
import { handleError } from '~/root/lib/utils/common';
import Container from '../../components/container';
import { IProviderContext } from './_layout';
import grecaptcha from '~/root/lib/client/helpers/g-recaptcha';
import { RECAPTCHA_SITE_KEY } from '~/auth/consts';

const CustomGoogleIcon = (props: any) => {
  return <GoogleLogo {...props} weight={4} />;
};

const LoginWithEmail = () => {
  const api = useAuthApi();
  const [searchParams, _setSearchParams] = useSearchParams();

  const reloadPage = useReload();
  const { values, errors, handleChange, handleSubmit, isLoading } = useForm({
    initialValues: {
      email: '',
      password: '',
    },
    validationSchema: Yup.object({
      email: Yup.string().required().email(),
      password: Yup.string().trim().required(),
    }),
    onSubmit: async (v) => {
      try {
        const token = await grecaptcha.execute(RECAPTCHA_SITE_KEY, {
          action: 'login',
        });
        const { errors: _errors } = await api.login({
          email: v.email,
          password: v.password,
          //@ts-ignore
          token,
        });
        if (_errors) {
          throw _errors[0];
        }
        toast.success('Logged in successfully');

        const callback = searchParams.get('callback');
        if (callback) {
          window.location.href = callback;
          return;
        }
        reloadPage();
      } catch (err) {
        handleError(err);
      }
    },
  });

  return (
    <form
      onSubmit={handleSubmit}
      className="flex flex-col items-stretch gap-3xl"
    >
      <TextInput
        value={values.email}
        error={!!errors.email}
        message={errors.email}
        onChange={handleChange('email')}
        label="Email"
        placeholder="ex: john@company.com"
        size="lg"
        className="h-[48px]"
      />
      <div className="flex flex-col gap-md">
        <PasswordInput
          value={values.password}
          error={!!errors.password}
          message={errors.ram}
          onChange={handleChange('password')}
          label="Password"
          placeholder="XXXXXX"
          size="lg"
          className="h-[48px]"
        />
        <Button
          content={<span className="text-text-soft">Forgot password</span>}
          size="sm"
          variant="plain"
          to="/forgot-password"
          linkComponent={Link}
        />
      </div>
      <Button
        loading={isLoading}
        size="lg"
        variant="primary"
        content={<span className="bodyLg-medium">Continue with email</span>}
        suffix={<ArrowRight />}
        block
        type="submit"
      />
    </form>
  );
};

const Login = () => {
  const { githubLoginUrl, gitlabLoginUrl, googleLoginUrl } =
    useOutletContext<IProviderContext>();
  const [searchParams, _setSearchParams] = useSearchParams();
  const callback = searchParams.get('callback');

  const loginUrl = callback
    ? `/login?mode=email&callback=${callback}`
    : `/login?mode=email`;

  useEffect(() => {
    if (callback) {
      getCookie().set('callback_url', callback, {
        expires: new Date(Date.now() + 1000 * 60),
      });
    }
  }, [callback]);

  return (
    <Container
      headerExtra={
        <Button
          variant="outline"
          content="Sign up"
          linkComponent={Link}
          to="/signup"
        />
      }
    >
      <div className="flex flex-col gap-3xl md:w-[500px] px-3xl py-5xl md:px-9xl">
        <div className="flex flex-col items-stretch">
          <div className="flex flex-col gap-lg items-center pb-6xl text-center">
            <div className={cn('text-text-strong headingXl text-center')}>
              Sign in to Kloudlite.io
            </div>
            <div className="bodyMd-medium text-text-soft">
              Start integrating local coding with remote power
            </div>
          </div>
          {searchParams.get('mode') === 'email' ? (
            <LoginWithEmail />
          ) : (
            <div className="flex flex-col items-stretch gap-3xl">
              <Button
                size="lg"
                variant="tertiary"
                content={
                  <span className="bodyLg-medium">Continue with GitHub</span>
                }
                prefix={<GithubLogoFill />}
                to={githubLoginUrl}
                disabled={!githubLoginUrl}
                block
                linkComponent={Link}
              />
              <Button
                size="lg"
                variant="purple"
                content={
                  <span className="bodyLg-medium">Continue with GitLab</span>
                }
                prefix={<GitlabLogoFill />}
                to={gitlabLoginUrl}
                disabled={!gitlabLoginUrl}
                block
                linkComponent={Link}
              />
              <Button
                size="lg"
                variant="primary"
                content={
                  <span className="bodyLg-medium">Continue with Google</span>
                }
                prefix={<CustomGoogleIcon />}
                to={googleLoginUrl}
                disabled={!googleLoginUrl}
                block
                linkComponent={Link}
              />
            </div>
          )}
        </div>
        {searchParams.get('mode') === 'email' ? (
          <Button
            size="lg"
            variant="plain"
            content={
              <span className="bodyLg-medium">Other sign in options</span>
            }
            prefix={<ArrowLeft />}
            to="/login"
            block
            linkComponent={Link}
          />
        ) : (
          <Button
            size="lg"
            variant="outline"
            content={<span className="bodyLg-medium">Continue with email</span>}
            prefix={<Envelope />}
            // to="/login/?mode=email"
            to={loginUrl}
            block
            linkComponent={Link}
          />
        )}
      </div>
    </Container>
  );
};

export default Login;
