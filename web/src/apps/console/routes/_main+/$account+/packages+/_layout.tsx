import { Outlet, useOutletContext } from '@remix-run/react';
import SidebarLayout from '~/console/components/sidebar-layout';
import { IAccountContext } from '../_layout';

export interface IPackageContext extends IAccountContext {}

const ContainerRegistry = () => {
  const rootContext = useOutletContext<IPackageContext>();

  return (
    <SidebarLayout
      navItems={[
        { label: 'Container Repos', value: 'repos' },
        { label: 'Helm Repos', value: 'helm-repos' },
        { label: 'Access management', value: 'access-management' },
      ]}
      parentPath="/packages"
    >
      <Outlet context={rootContext} />
    </SidebarLayout>
  );
};

export default ContainerRegistry;
