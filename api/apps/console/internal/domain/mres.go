package domain

import (
	"fmt"
	"time"

	"github.com/kloudlite/api/apps/console/internal/entities"
	fc "github.com/kloudlite/api/apps/console/internal/entities/field-constants"
	iamT "github.com/kloudlite/api/apps/iam/types"
	"github.com/kloudlite/api/common"
	"github.com/kloudlite/api/common/fields"
	"github.com/kloudlite/api/grpc-interfaces/infra"
	"github.com/kloudlite/api/pkg/errors"
	fn "github.com/kloudlite/api/pkg/functions"
	"github.com/kloudlite/api/pkg/repos"
	t "github.com/kloudlite/api/pkg/types"
	"github.com/kloudlite/operator/operators/resource-watcher/types"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// query

func (d *domain) ListManagedResources(ctx ConsoleContext, search map[string]repos.MatchFilter, pq repos.CursorPagination) (*repos.PaginatedRecord[*entities.ManagedResource], error) {
	if err := d.canPerformActionInAccount(ctx, iamT.ListManagedResources); err != nil {
		return nil, errors.NewE(err)
	}

	filters := d.mresRepo.MergeMatchFilters(repos.Filter{
		fields.AccountName: ctx.AccountName,
	}, search)

	return d.mresRepo.FindPaginated(ctx, filters, pq)
}

func (d *domain) listImportedMres(ctx ConsoleContext, mresName string) ([]*entities.ManagedResource, error) {
	filter := repos.Filter{
		fields.AccountName:        ctx.AccountName,
		fc.ManagedResourceMresRef: mresName,
	}

	return d.mresRepo.Find(ctx, repos.Query{Filter: filter})
}

func (d *domain) findMRes(ctx ManagedResourceContext, name string) (*entities.ManagedResource, error) {
	filters, err := ctx.MresDBFilters()
	if err != nil {
		return nil, errors.NewE(err)
	}

	mres, err := d.mresRepo.FindOne(
		ctx,
		filters.Add(fields.MetadataName, name),
	)
	if err != nil {
		return nil, errors.NewE(err)
	}
	if mres == nil {
		return nil, errors.Newf("no managed resource with name (%s) found", name)
	}
	return mres, nil
}

func (d *domain) GetManagedResource(ctx ManagedResourceContext, name string) (*entities.ManagedResource, error) {
	if err := d.canPerformActionInAccount(ctx.ConsoleContext, iamT.GetManagedResource); err != nil {
		return nil, errors.NewE(err)
	}

	return d.findMRes(ctx, name)
}

func (d *domain) GetManagedResourceByID(ctx ConsoleContext, id repos.ID) (*entities.ManagedResource, error) {
	if err := d.canPerformActionInAccount(ctx, iamT.GetManagedResource); err != nil {
		return nil, errors.NewE(err)
	}

	return d.mresRepo.FindById(ctx, id)
}

// GetManagedResourceOutputKVs implements Domain.
func (d *domain) GetManagedResourceOutputKVs(ctx ManagedResourceContext, keyrefs []ManagedResourceKeyRef) ([]*ManagedResourceKeyValueRef, error) {
	f, err := ctx.MresDBFilters()
	if err != nil {
		return nil, err
	}

	names := make([]any, 0, len(keyrefs))
	for i := range keyrefs {
		names = append(names, keyrefs[i].MresName)
	}

	filters := d.mresRepo.MergeMatchFilters(*f, map[string]repos.MatchFilter{
		fields.MetadataName: {
			MatchType: repos.MatchTypeArray,
			Array:     names,
		},
	})

	mresSecrets, err := d.mresRepo.Find(ctx, repos.Query{Filter: filters})
	if err != nil {
		return nil, errors.NewE(err)
	}

	results := make([]*ManagedResourceKeyValueRef, 0, len(mresSecrets))

	data := make(map[string]map[string]string)

	for i := range mresSecrets {
		m := make(map[string]string, len(mresSecrets[i].SyncedOutputSecretRef.Data))
		for k, v := range mresSecrets[i].SyncedOutputSecretRef.Data {
			m[k] = string(v)
		}

		for k, v := range mresSecrets[i].SyncedOutputSecretRef.StringData {
			m[k] = v
		}

		data[mresSecrets[i].Name] = m
	}

	for i := range keyrefs {
		results = append(results, &ManagedResourceKeyValueRef{
			MresName: keyrefs[i].MresName,
			Key:      keyrefs[i].Key,
			Value:    data[keyrefs[i].MresName][keyrefs[i].Key],
		})
	}

	return results, nil
}

// GetManagedResourceOutputKeys implements Domain.
func (d *domain) GetManagedResourceOutputKeys(ctx ManagedResourceContext, name string) ([]string, error) {
	mresSecret, err := d.findMRes(ctx, name)
	if err != nil {
		return nil, errors.NewE(err)
	}

	if mresSecret.SyncedOutputSecretRef == nil {
		return nil, errors.Newf("waiting for managed resource output to sync")
	}

	results := make([]string, 0, len(mresSecret.SyncedOutputSecretRef.Data))

	for k := range mresSecret.SyncedOutputSecretRef.Data {
		results = append(results, k)
	}

	return results, nil
}

// mutations
func (d *domain) CreateRootManagedResource(ctx ConsoleContext, accountNamespace string, mres *entities.ManagedResource) (*entities.ManagedResource, error) {
	if err := d.canPerformActionInAccount(ctx, iamT.CreateManagedResource); err != nil {
		return nil, errors.NewE(err)
	}

	mres.SyncStatus = t.SyncStatus{
		SyncScheduledAt: time.Now(),
		LastSyncedAt:    time.Now(),
		Action:          t.SyncActionApply,
		RecordVersion:   1,
		State:           t.SyncStateUpdatedAtAgent,
		Error:           nil,
	}

	if mres.SyncedOutputSecretRef == nil {
		return nil, errors.Newf("managed resource (%s), not ready yet, please try again", mres.Name)
	}

	secret := *mres.SyncedOutputSecretRef
	secret.Namespace = accountNamespace

	if _, err := d.secretRepo.Upsert(ctx, repos.Filter{
		fc.AccountName:       ctx.AccountName,
		fc.MetadataName:      secret.Name,
		fc.MetadataNamespace: secret.Namespace,
	}, &entities.Secret{
		Secret:      secret,
		AccountName: ctx.AccountName,
		ResourceMetadata: common.ResourceMetadata{
			DisplayName:   fmt.Sprintf("root credentials (%s)", mres.ManagedServiceName),
			CreatedBy:     mres.CreatedBy,
			LastUpdatedBy: mres.LastUpdatedBy,
		},
		SyncStatus: mres.SyncStatus,
		IsReadOnly: true,
	}); err != nil {
		return nil, errors.NewE(err)
	}

	return d.mresRepo.Upsert(ctx, repos.Filter{
		fields.AccountName:                   ctx.AccountName,
		fc.ManagedResourceManagedServiceName: mres.ManagedServiceName,
		fc.MetadataName:                      mres.Name,
		fc.ClusterName:                       mres.ClusterName,
	}, mres)
}

func (d *domain) CreateManagedResource(ctx ManagedResourceContext, mres entities.ManagedResource) (*entities.ManagedResource, error) {
	if err := d.canPerformActionInAccount(ctx.ConsoleContext, iamT.CreateManagedResource); err != nil {
		return nil, errors.NewE(err)
	}

	if ctx.ManagedServiceName == nil {
		return nil, errors.Newf("managed service name is required")
	}

	msvcOut, err := d.infraClient.GetClusterManagedService(ctx, &infra.GetClusterManagedServiceIn{
		UserId:      string(ctx.UserId),
		UserName:    ctx.UserName,
		UserEmail:   ctx.UserEmail,
		AccountName: ctx.AccountName,
		MsvcName:    *ctx.ManagedServiceName,
	})
	if err != nil {
		return nil, errors.NewE(err)
	}

	if mres.Spec.ResourceTemplate.TypeMeta.GroupVersionKind().GroupKind().Empty() {
		return nil, errors.New(".spec.resourceTemplate.apiVersion, and .spec.resourceTemplate.kind must be set")
	}

	mres.Namespace = msvcOut.TargetNamespace
	if mres.ManagedResource.Spec.ResourceTemplate.MsvcRef.ClusterName == nil {
		mres.ManagedResource.Spec.ResourceTemplate.MsvcRef.ClusterName = &msvcOut.ClusterName
	}

	mres.EnsureGVK()
	if err := d.k8sClient.ValidateObject(ctx, &mres.ManagedResource); err != nil {
		return nil, errors.NewE(err)
	}

	mres.IncrementRecordVersion()

	mres.CreatedBy = common.CreatedOrUpdatedBy{
		UserId:    ctx.UserId,
		UserName:  ctx.UserName,
		UserEmail: ctx.UserEmail,
	}
	mres.LastUpdatedBy = mres.CreatedBy

	mres.AccountName = ctx.AccountName
	mres.EnvironmentName = ""
	mres.ManagedServiceName = *ctx.ManagedServiceName
	mres.IsImported = false
	mres.ClusterName = msvcOut.ClusterName

	return d.createAndApplyManagedResource(ctx, msvcOut.ClusterName, &mres)
}

func genMresResourceName(envName string, mresName string) string {
	return fmt.Sprintf("env-%s-%s", envName, mresName)
}

func genMresCredentialsSecretName(name string) string {
	return fmt.Sprintf("mres-%s-creds", name)
}

func (d *domain) createAndApplyManagedResource(ctx ManagedResourceContext, clusterName string, mres *entities.ManagedResource) (*entities.ManagedResource, error) {
	mres.Spec.ResourceNamePrefix = nil
	mres.SyncStatus = t.GenSyncStatus(t.SyncActionApply, 0)

	m, err := d.mresRepo.Create(ctx, mres)
	if err != nil {
		if d.mresRepo.ErrAlreadyExists(err) {
			// TODO: better insights into error, when it is being caused by duplicated indexes
			return nil, errors.NewE(err)
		}
		return nil, errors.NewE(err)
	}

	d.resourceEventPublisher.PublishClusterManagedServiceEvent(ctx.ConsoleContext, m.ManagedServiceName, entities.ResourceTypeManagedResource, m.Name, PublishAdd)

	if err := d.applyK8sResourceOnCluster(ctx, clusterName, &m.ManagedResource, m.RecordVersion); err != nil {
		return m, errors.NewE(err)
	}

	return m, nil
}

func (d *domain) importAndApplyManagedResourceSecret(ctx ConsoleContext, envName string, mres *entities.ManagedResource) (*entities.ManagedResource, error) {
	mres.Spec.ResourceNamePrefix = fn.New(genMresResourceName(envName, mres.Name))
	mres.SyncStatus = t.GenSyncStatus(t.SyncActionApply, 0)

	if mres.SyncedOutputSecretRef == nil {
		return nil, errors.Newf("managed resource (%s), not ready yet, please try again", mres.Name)
	}

	ann := mres.SyncedOutputSecretRef.GetAnnotations()
	if ann == nil {
		ann = make(map[string]string, len(types.SecretWatchingAnnotation))
	}

	for k, v := range types.SecretWatchingAnnotation {
		ann[k] = v
	}

	nImpMres, err := d.mresRepo.Create(ctx, mres)
	if err != nil {
		return nil, errors.NewE(err)
	}

	d.resourceEventPublisher.PublishEnvironmentResourceEvent(ctx, envName, entities.ResourceTypeManagedResource, nImpMres.Name, PublishAdd)

	s, err := d.secretRepo.Create(ctx, &entities.Secret{
		Secret: corev1.Secret{
			TypeMeta: v1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
			ObjectMeta: v1.ObjectMeta{
				Name:      mres.SyncedOutputSecretRef.Name,
				Namespace: d.getEnvironmentTargetNamespace(envName),
				Labels: map[string]string{
					"kloudlite.io/mres.imported": "true",
				},
				Annotations: ann,
			},
			Data: mres.SyncedOutputSecretRef.Data,
		},
		AccountName:      mres.AccountName,
		EnvironmentName:  mres.EnvironmentName,
		ResourceMetadata: mres.ResourceMetadata,
		SyncStatus:       mres.SyncStatus,
		IsReadOnly:       true,
	})
	if err != nil {
		return nil, errors.NewE(err)
	}

	if err := d.applyK8sResource(ctx, envName, &s.Secret, s.RecordVersion); err != nil {
		return nImpMres, errors.NewE(err)
	}

	return nImpMres, nil
}

func (d *domain) applyMresSecretsToEnvironment(ctx ConsoleContext, envName string, mres *entities.ManagedResource) error {
	mres.SyncedOutputSecretRef.Namespace = d.getEnvironmentTargetNamespace(envName)
	return d.applyK8sResource(ctx, envName, mres.SyncedOutputSecretRef, mres.RecordVersion)
}

func (d *domain) applyMresSecrets(ctx ConsoleContext, mres *entities.ManagedResource) error {
	listImportedMres, err := d.listImportedMres(ctx, mres.Name)
	if err != nil {
		return nil
	}

	for _, mres := range listImportedMres {
		err := d.applyMresSecretsToEnvironment(ctx, mres.EnvironmentName, mres)
		if err != nil {
			return errors.NewE(err)
		}
	}
	return nil
}

func (d *domain) UpdateManagedResource(ctx ManagedResourceContext, mres entities.ManagedResource) (*entities.ManagedResource, error) {
	if err := d.canPerformActionInAccount(ctx.ConsoleContext, iamT.UpdateManagedResource); err != nil {
		return nil, errors.NewE(err)
	}

	mres.EnsureGVK()
	if err := d.k8sClient.ValidateObject(ctx, &mres.ManagedResource); err != nil {
		return nil, errors.NewE(err)
	}

	patchForUpdate := common.PatchForUpdate(
		ctx,
		&mres,
		common.PatchOpts{
			XPatch: repos.Document{
				fc.ManagedResourceSpecResourceTemplateSpec: mres.Spec.ResourceTemplate.Spec,
			},
		})

	f, err := ctx.MresDBFilters()
	if err != nil {
		return nil, errors.NewE(err)
	}

	upMres, err := d.mresRepo.Patch(
		ctx,
		f.Add(fields.MetadataName, mres.Name),
		patchForUpdate,
	)
	if err != nil {
		return nil, errors.NewE(err)
	}

	d.resourceEventPublisher.PublishClusterManagedServiceEvent(ctx.ConsoleContext, upMres.ManagedServiceName, entities.ResourceTypeManagedResource, upMres.Name, PublishUpdate)

	if err := d.applyK8sResourceOnCluster(ctx, upMres.ClusterName, &upMres.ManagedResource, upMres.RecordVersion); err != nil {
		return upMres, errors.NewE(err)
	}

	return upMres, nil
}

func (d *domain) DeleteManagedResource(ctx ManagedResourceContext, name string) error {
	if err := d.canPerformActionInAccount(ctx.ConsoleContext, iamT.DeleteManagedResource); err != nil {
		return errors.NewE(err)
	}

	f, err := ctx.MresDBFilters()
	if err != nil {
		return errors.NewE(err)
	}

	umres, err := d.mresRepo.Patch(
		ctx,
		f.Add(fields.MetadataName, name),
		common.PatchForMarkDeletion(),
	)
	if err != nil {
		return errors.NewE(err)
	}
	d.resourceEventPublisher.PublishClusterManagedServiceEvent(ctx.ConsoleContext, umres.ManagedServiceName, entities.ResourceTypeManagedResource, umres.Name, PublishUpdate)
	if err := d.deleteK8sResourceOfCluster(ctx, umres.ClusterName, &umres.ManagedResource); err != nil {
		if errors.Is(err, ErrNoClusterAttached) {
			return d.mresRepo.DeleteById(ctx, umres.Id)
		}
		return errors.NewE(err)
	}
	return nil
}
func (d *domain) OnManagedResourceDeleteMessage(ctx ConsoleContext, msvcName string, mres entities.ManagedResource) error {
	err := d.mresRepo.DeleteOne(
		ctx,
		repos.Filter{
			fields.AccountName:                   ctx.AccountName,
			fc.ManagedResourceManagedServiceName: msvcName,
			fields.MetadataName:                  mres.Name,
		},
	)
	if err != nil {
		return errors.NewE(err)
	}
	d.resourceEventPublisher.PublishClusterManagedServiceEvent(ctx, msvcName, entities.ResourceTypeManagedResource, mres.Name, PublishDelete)
	return nil
}

func (d *domain) OnManagedResourceUpdateMessage(ctx ConsoleContext, msvcName string, mres entities.ManagedResource, status types.ResourceStatus, opts UpdateAndDeleteOpts) error {
	xmres, err := d.findMRes(NewManagedResourceContext(ctx, msvcName), mres.Name)
	if err != nil {
		return errors.NewE(err)
	}

	if xmres == nil {
		return errors.Newf("no manage resource found")
	}

	recordVersion, err := d.MatchRecordVersion(mres.Annotations, xmres.RecordVersion)
	if err != nil {
		return d.resyncK8sResourceToCluster(ctx, xmres.ClusterName, mres.SyncStatus.Action, &mres.ManagedResource, mres.RecordVersion)
	}

	umres, err := d.mresRepo.PatchById(
		ctx,
		xmres.Id,
		common.PatchForSyncFromAgent(&mres, recordVersion, status, common.PatchOpts{
			MessageTimestamp: opts.MessageTimestamp,
			XPatch: repos.Document{
				fc.ManagedResourceSyncedOutputSecretRef: mres.SyncedOutputSecretRef,
			},
		}))
	if err != nil {
		return err
	}

	d.resourceEventPublisher.PublishClusterManagedServiceEvent(ctx, msvcName, umres.GetResourceType(), umres.GetName(), PublishUpdate)

	if mres.SyncedOutputSecretRef != nil {
		if mres.SyncedOutputSecretRef.Labels == nil {
			mres.SyncedOutputSecretRef.Labels = map[string]string{}
		}
		mres.SyncedOutputSecretRef.Labels["kloudlite.io/secret.synced-by"] = fmt.Sprintf("%s/%s", umres.GetNamespace(), umres.GetName())

		secretData := make(map[string]string, len(mres.SyncedOutputSecretRef.Data))

		for k, v := range mres.SyncedOutputSecretRef.Data {
			secretData[k] = string(v)
		}

		mres.SyncedOutputSecretRef.Data = nil
		mres.SyncedOutputSecretRef.StringData = secretData

		if _, err = d.secretRepo.Upsert(ctx, repos.Filter{
			fc.AccountName:       ctx.AccountName,
			fc.MetadataName:      mres.SyncedOutputSecretRef.GetName(),
			fc.MetadataNamespace: mres.SyncedOutputSecretRef.GetNamespace(),
		}, &entities.Secret{
			Secret:      *mres.SyncedOutputSecretRef,
			AccountName: ctx.AccountName,
			ResourceMetadata: common.ResourceMetadata{
				DisplayName:   umres.GetName(),
				CreatedBy:     common.CreatedOrUpdatedByResourceSync,
				LastUpdatedBy: common.CreatedOrUpdatedByResourceSync,
			},
			SyncStatus: t.SyncStatus{
				LastSyncedAt:  opts.MessageTimestamp,
				Action:        t.SyncActionApply,
				RecordVersion: recordVersion,
				State:         t.SyncStateUpdatedAtAgent,
				Error:         nil,
			},
			IsReadOnly: true,
		}); err != nil {
			return errors.NewE(err)
		}
	}

	err = d.applyMresSecrets(ctx, &mres)
	if err != nil {
		return errors.NewE(err)
	}

	return nil
}

func (d *domain) OnManagedResourceApplyError(ctx ConsoleContext, errMsg string, msvcName string, name string, opts UpdateAndDeleteOpts) error {
	umres, err := d.mresRepo.Patch(
		ctx,
		repos.Filter{
			fields.AccountName:                   ctx.AccountName,
			fc.ManagedResourceManagedServiceName: msvcName,
			fields.MetadataName:                  name,
		},
		common.PatchForErrorFromAgent(
			errMsg,
			common.PatchOpts{
				MessageTimestamp: opts.MessageTimestamp,
			},
		),
	)
	if err != nil {
		return errors.NewE(err)
	}
	d.resourceEventPublisher.PublishClusterManagedServiceEvent(ctx, msvcName, entities.ResourceTypeManagedResource, umres.Name, PublishDelete)
	return errors.NewE(err)
}

func (d *domain) ResyncManagedResource(ctx ConsoleContext, msvcName string, name string) error {
	if err := d.canPerformActionInAccount(ctx, iamT.CreateManagedResource); err != nil {
		return errors.NewE(err)
	}

	mres, err := d.findMRes(NewManagedResourceContext(ctx, msvcName), name)
	if err != nil {
		return errors.NewE(err)
	}
	return d.resyncK8sResourceToCluster(ctx, mres.ClusterName, mres.SyncStatus.Action, &mres.ManagedResource, mres.RecordVersion)
}
