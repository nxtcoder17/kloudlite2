package domain

import (
	iamT "github.com/kloudlite/api/apps/iam/types"
	"github.com/kloudlite/api/apps/infra/internal/entities"
	"github.com/kloudlite/api/pkg/errors"
	fn "github.com/kloudlite/api/pkg/functions"

	"github.com/kloudlite/api/common"
	"github.com/kloudlite/api/pkg/repos"
	t "github.com/kloudlite/api/pkg/types"
	redpandaMsvcv1 "github.com/kloudlite/operator/apis/redpanda.msvc/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (d *domain) findBYOCCluster(ctx InfraContext, clusterName string) (*entities.BYOCCluster, error) {
	accNs, err := d.getAccNamespace(ctx, ctx.AccountName)
	if err != nil {
		return nil, errors.NewE(err)
	}

	cluster, err := d.byocClusterRepo.FindOne(ctx, repos.Filter{
		"spec.accountName":   ctx.AccountName,
		"metadata.name":      clusterName,
		"metadata.namespace": accNs,
	})
	if err != nil {
		return nil, errors.NewE(err)
	}
	if cluster == nil {
		return nil, errors.Newf("BYOC cluster with name %q not found", clusterName)
	}
	return cluster, nil
}

func (d *domain) CreateBYOCCluster(ctx InfraContext, cluster entities.BYOCCluster) (*entities.BYOCCluster, error) {
	if err := d.canPerformActionInAccount(ctx, iamT.CreateCluster); err != nil {
		return nil, errors.NewE(err)
	}

	cluster.EnsureGVK()
	cluster.IncomingKafkaTopicName = common.GetKafkaTopicName(ctx.AccountName, cluster.Name)

	if err := d.k8sClient.ValidateObject(ctx, &cluster.BYOC); err != nil {
		return nil, errors.NewE(err)
	}

	cluster.IncrementRecordVersion()
	cluster.CreatedBy = common.CreatedOrUpdatedBy{
		UserId:    ctx.UserId,
		UserName:  ctx.UserName,
		UserEmail: ctx.UserEmail,
	}
	cluster.LastUpdatedBy = cluster.CreatedBy

	cluster.IsConnected = false
	cluster.Spec.AccountName = ctx.AccountName
	cluster.SyncStatus = t.GenSyncStatus(t.SyncActionApply, cluster.RecordVersion)

	nCluster, err := d.byocClusterRepo.Create(ctx, &cluster)
	if err != nil {
		if d.clusterRepo.ErrAlreadyExists(err) {
			return nil, errors.NewE(err)
		}
	}

	if err := d.applyK8sResource(ctx, &nCluster.BYOC, nCluster.RecordVersion); err != nil {
		return nil, errors.NewE(err)
	}

	redpandaTopic := redpandaMsvcv1.Topic{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Name: cluster.IncomingKafkaTopicName, Namespace: d.env.ProviderSecretNamespace},
	}

	redpandaTopic.EnsureGVK()

	if err := d.applyK8sResource(ctx, &redpandaTopic, nCluster.RecordVersion); err != nil {
		return nil, errors.NewE(err)
	}

	return nCluster, nil
}

func (d *domain) ListBYOCClusters(ctx InfraContext, filters map[string]repos.MatchFilter, pagination repos.CursorPagination) (*repos.PaginatedRecord[*entities.BYOCCluster], error) {
	if err := d.canPerformActionInAccount(ctx, iamT.ListClusters); err != nil {
		return nil, errors.NewE(err)
	}

	accNs, err := d.getAccNamespace(ctx, ctx.AccountName)
	if err != nil {
		return nil, errors.NewE(err)
	}

	f := repos.Filter{
		"accountName":        ctx.AccountName,
		"metadata.namespace": accNs,
	}
	return d.byocClusterRepo.FindPaginated(ctx, d.byocClusterRepo.MergeMatchFilters(f, filters), pagination)
}

func (d *domain) GetBYOCCluster(ctx InfraContext, name string) (*entities.BYOCCluster, error) {
	if err := d.canPerformActionInAccount(ctx, iamT.GetCluster); err != nil {
		return nil, errors.NewE(err)
	}
	return d.findBYOCCluster(ctx, name)
}

func (d *domain) UpdateBYOCCluster(ctx InfraContext, cluster entities.BYOCCluster) (*entities.BYOCCluster, error) {
	if err := d.canPerformActionInAccount(ctx, iamT.UpdateCluster); err != nil {
		return nil, errors.NewE(err)
	}

	cluster.EnsureGVK()
	if err := d.k8sClient.ValidateObject(ctx, &cluster.BYOC); err != nil {
		return nil, errors.NewE(err)
	}

	c, err := d.findBYOCCluster(ctx, cluster.Name)
	if err != nil {
		return nil, errors.NewE(err)
	}

	c.IncrementRecordVersion()
	c.LastUpdatedBy = common.CreatedOrUpdatedBy{
		UserId:    ctx.UserId,
		UserName:  ctx.UserName,
		UserEmail: ctx.UserEmail,
	}

	c.BYOC = cluster.BYOC
	c.SyncStatus = t.GenSyncStatus(t.SyncActionApply, c.RecordVersion)

	uCluster, err := d.byocClusterRepo.UpdateById(ctx, c.Id, c)
	if err != nil {
		return nil, errors.NewE(err)
	}

	if err := d.applyK8sResource(ctx, &uCluster.BYOC, uCluster.RecordVersion); err != nil {
		return nil, errors.NewE(err)
	}

	return uCluster, nil
}

func (d *domain) DeleteBYOCCluster(ctx InfraContext, name string) error {
	if err := d.canPerformActionInAccount(ctx, iamT.DeleteCluster); err != nil {
		return errors.NewE(err)
	}

	clus, err := d.findBYOCCluster(ctx, name)
	if err != nil {
		return errors.NewE(err)
	}

	if clus.IsMarkedForDeletion() {
		return errors.Newf("BYOC cluster %q is already marked for deletion", name)
	}

	clus.MarkedForDeletion = fn.New(true)
	clus.SyncStatus = t.GetSyncStatusForDeletion(clus.Generation)
	upC, err := d.byocClusterRepo.UpdateById(ctx, clus.Id, clus)
	if err != nil {
		return errors.NewE(err)
	}
	return d.deleteK8sResource(ctx, &upC.BYOC)
}

func (d *domain) ResyncBYOCCluster(ctx InfraContext, name string) error {
	clus, err := d.findBYOCCluster(ctx, name)
	if err != nil {
		return errors.NewE(err)
	}

	if err := d.applyK8sResource(ctx, &clus.BYOC, clus.RecordVersion); err != nil {
		return errors.NewE(err)
	}

	redpandaTopic := redpandaMsvcv1.Topic{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      clus.IncomingKafkaTopicName,
			Namespace: d.env.ProviderSecretNamespace,
		},
	}

	redpandaTopic.EnsureGVK()
	return d.applyK8sResource(ctx, &redpandaTopic, clus.RecordVersion)
}

func (d *domain) OnDeleteBYOCClusterMessage(ctx InfraContext, cluster entities.BYOCCluster) error {
	accNs, err := d.getAccNamespace(ctx, ctx.AccountName)
	if err != nil {
		return errors.NewE(err)
	}

	return d.clusterRepo.DeleteOne(ctx, repos.Filter{
		"accountName":        ctx.AccountName,
		"metadata.name":      cluster.Name,
		"metadata.namespace": accNs,
	})
}

func (d *domain) OnUpdateBYOCClusterMessage(ctx InfraContext, cluster entities.BYOCCluster) error {
	c, err := d.findBYOCCluster(ctx, cluster.Name)
	if err != nil {
		return errors.NewE(err)
	}

	c.SyncStatus.State = t.SyncStateReceivedUpdateFromAgent

	_, err = d.byocClusterRepo.UpdateById(ctx, c.Id, &cluster)
	if err != nil {
		return errors.NewE(err)
	}
	return nil
}
