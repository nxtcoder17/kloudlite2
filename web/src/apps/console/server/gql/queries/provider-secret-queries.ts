import gql from 'graphql-tag';
import { IExecutor } from '~/root/lib/server/helpers/execute-query-with-context';
import { NN } from '~/root/lib/types/common';
import {
  ConsoleCreateProviderSecretMutation,
  ConsoleCreateProviderSecretMutationVariables,
  ConsoleDeleteProviderSecretMutation,
  ConsoleDeleteProviderSecretMutationVariables,
  ConsoleGetProviderSecretQuery,
  ConsoleGetProviderSecretQueryVariables,
  ConsoleListProviderSecretsQuery,
  ConsoleListProviderSecretsQueryVariables,
  ConsoleUpdateProviderSecretMutation,
  ConsoleUpdateProviderSecretMutationVariables,
  ConsoleCheckAwsAccessQueryVariables,
  ConsoleCheckAwsAccessQuery,
} from '~/root/src/generated/gql/server';

export type IProviderSecrets = NN<
  ConsoleListProviderSecretsQuery['infra_listProviderSecrets']
>;

export type IProviderSecret = NN<
  ConsoleGetProviderSecretQuery['infra_getProviderSecret']
>;

export const providerSecretQueries = (executor: IExecutor) => ({
  checkAwsAccess: executor(
    gql`
      query Infra_checkAwsAccess($cloudproviderName: String!) {
        infra_checkAwsAccess(cloudproviderName: $cloudproviderName) {
          result
          installationUrl
        }
      }
    `,
    {
      transformer(data: ConsoleCheckAwsAccessQuery) {
        return data.infra_checkAwsAccess;
      },
      vars(_: ConsoleCheckAwsAccessQueryVariables) {},
    }
  ),
  listProviderSecrets: executor(
    gql`
      query Infra_listProviderSecrets(
        $search: SearchProviderSecret
        $pagination: CursorPaginationIn
      ) {
        infra_listProviderSecrets(search: $search, pagination: $pagination) {
          pageInfo {
            endCursor
            hasNextPage
            hasPrevPage
            startCursor
          }
          totalCount
          edges {
            cursor
            node {
              cloudProviderName
              createdBy {
                userEmail
                userId
                userName
              }
              aws {
                authMechanism
              }
              creationTime
              gcp {
                serviceAccountJSON
              }
              displayName
              lastUpdatedBy {
                userEmail
                userId
                userName
              }
              metadata {
                namespace
                name
                labels
                annotations
              }
              updateTime
            }
          }
        }
      }
    `,
    {
      transformer: (data: ConsoleListProviderSecretsQuery) => {
        return data.infra_listProviderSecrets;
      },
      vars: (_: ConsoleListProviderSecretsQueryVariables) => {},
    }
  ),
  createProviderSecret: executor(
    gql`
      mutation createProviderSecret($secret: CloudProviderSecretIn!) {
        infra_createProviderSecret(secret: $secret) {
          metadata {
            name
          }
        }
      }
    `,
    {
      transformer: (data: ConsoleCreateProviderSecretMutation) =>
        data.infra_createProviderSecret,
      vars(_: ConsoleCreateProviderSecretMutationVariables) {},
    }
  ),
  updateProviderSecret: executor(
    gql`
      mutation Infra_updateProviderSecret($secret: CloudProviderSecretIn!) {
        infra_updateProviderSecret(secret: $secret) {
          id
        }
      }
    `,
    {
      transformer: (data: ConsoleUpdateProviderSecretMutation) =>
        data.infra_updateProviderSecret,
      vars(_: ConsoleUpdateProviderSecretMutationVariables) {},
    }
  ),
  deleteProviderSecret: executor(
    gql`
      mutation Infra_deleteProviderSecret($secretName: String!) {
        infra_deleteProviderSecret(secretName: $secretName)
      }
    `,
    {
      transformer: (data: ConsoleDeleteProviderSecretMutation) =>
        data.infra_deleteProviderSecret,
      vars(_: ConsoleDeleteProviderSecretMutationVariables) {},
    }
  ),
  getProviderSecret: executor(
    gql`
      query Metadata($name: String!) {
        infra_getProviderSecret(name: $name) {
          cloudProviderName
          createdBy {
            userEmail
            userId
            userName
          }
          creationTime
          displayName
          lastUpdatedBy {
            userEmail
            userId
            userName
          }
          metadata {
            namespace
            name
            labels
          }
          updateTime
        }
      }
    `,

    {
      transformer: (data: ConsoleGetProviderSecretQuery) =>
        data.infra_getProviderSecret,
      vars(_: ConsoleGetProviderSecretQueryVariables) {},
    }
  ),
});
