package domain

import (
	"time"

	"github.com/kloudlite/api/apps/console/internal/entities"
	fc "github.com/kloudlite/api/apps/console/internal/entities/field-constants"
	"github.com/kloudlite/api/common"
	"github.com/kloudlite/api/pkg/errors"
	fn "github.com/kloudlite/api/pkg/functions"
	"github.com/kloudlite/api/pkg/repos"
	t "github.com/kloudlite/api/pkg/types"
	"github.com/kloudlite/operator/operators/resource-watcher/types"
)

func (d *domain) ListConfigs(ctx ResourceContext, search map[string]repos.MatchFilter, pq repos.CursorPagination) (*repos.PaginatedRecord[*entities.Config], error) {
	if err := d.canReadResourcesInEnvironment(ctx); err != nil {
		return nil, errors.NewE(err)
	}
	filter := ctx.DBFilters()
	return d.configRepo.FindPaginated(ctx, d.configRepo.MergeMatchFilters(filter, search), pq)
}

func (d *domain) findConfig(ctx ResourceContext, name string) (*entities.Config, error) {
	filters := ctx.DBFilters()
	filters.Add(fc.MetadataName, name)

	cfg, err := d.configRepo.FindOne(ctx, filters)
	if err != nil {
		return nil, errors.NewE(err)
	}
	if cfg == nil {
		return nil, errors.Newf("no config with name (%q)", name)
	}
	return cfg, nil
}

func (d *domain) GetConfig(ctx ResourceContext, name string) (*entities.Config, error) {
	if err := d.canReadResourcesInEnvironment(ctx); err != nil {
		return nil, errors.NewE(err)
	}
	return d.findConfig(ctx, name)
}

// GetConfigEntries implements Domain.
func (d *domain) GetConfigEntries(ctx ResourceContext, keyrefs []ConfigKeyRef) ([]*ConfigKeyValueRef, error) {
	filters := ctx.DBFilters()

	names := make([]any, 0, len(keyrefs))
	for i := range keyrefs {
		names = append(names, keyrefs[i].ConfigName)
	}

	filters = d.configRepo.MergeMatchFilters(filters, map[string]repos.MatchFilter{
		fc.MetadataName: {
			MatchType: repos.MatchTypeArray,
			Array:     names,
		},
	})

	configs, err := d.configRepo.Find(ctx, repos.Query{Filter: filters})
	if err != nil {
		return nil, errors.NewE(err)
	}

	results := make([]*ConfigKeyValueRef, 0, len(configs))

	data := make(map[string]map[string]string)

	for i := range configs {
		data[configs[i].Name] = configs[i].Data
	}

	for i := range keyrefs {
		results = append(results, &ConfigKeyValueRef{
			ConfigName: keyrefs[i].ConfigName,
			Key:        keyrefs[i].Key,
			Value:      data[keyrefs[i].ConfigName][keyrefs[i].Key],
		})
	}

	return results, nil
}

func (d *domain) CreateConfig(ctx ResourceContext, config entities.Config) (*entities.Config, error) {
	if err := d.canMutateResourcesInEnvironment(ctx); err != nil {
		return nil, errors.NewE(err)
	}

	config.SetGroupVersionKind(fn.GVK("v1", "ConfigMap"))

	var err error
	config.Namespace, err = d.envTargetNamespace(ctx.ConsoleContext, ctx.ProjectName, ctx.EnvironmentName)
	if err != nil {
		return nil, err
	}

	config.IncrementRecordVersion()
	config.CreatedBy = common.CreatedOrUpdatedBy{
		UserId:    ctx.UserId,
		UserName:  ctx.UserName,
		UserEmail: ctx.UserEmail,
	}
	config.LastUpdatedBy = config.CreatedBy

	config.AccountName = ctx.AccountName
	config.ProjectName = ctx.ProjectName
	config.EnvironmentName = ctx.EnvironmentName

	config.SyncStatus = t.GenSyncStatus(t.SyncActionApply, config.RecordVersion)

	if _, err := d.upsertEnvironmentResourceMapping(ctx, &config); err != nil {
		return nil, errors.NewE(err)
	}

	cfg, err := d.configRepo.Create(ctx, &config)
	if err != nil {
		if d.configRepo.ErrAlreadyExists(err) {
			// TODO: better insights into error, when it is being caused by duplicated indexes
			return nil, errors.NewE(err)
		}
		return nil, errors.NewE(err)
	}

	if cfg.Annotations == nil {
		cfg.Annotations = make(map[string]string)
	}

	for k, v := range types.ConfigWatchingAnnotation {
		cfg.Annotations[k] = v
	}

	if err := d.applyK8sResource(ctx, config.ProjectName, &cfg.ConfigMap, cfg.RecordVersion); err != nil {
		return cfg, errors.NewE(err)
	}

	return cfg, nil
}

func (d *domain) UpdateConfig(ctx ResourceContext, config entities.Config) (*entities.Config, error) {
	if err := d.canMutateResourcesInEnvironment(ctx); err != nil {
		return nil, errors.NewE(err)
	}

	config.SetGroupVersionKind(fn.GVK("v1", "ConfigMap"))
	xconfig, err := d.findConfig(ctx, config.Name)
	if err != nil {
		return nil, errors.NewE(err)
	}

	if xconfig == nil {
		return nil, errors.Newf("no config found")
	}

	upConfig, err := d.configRepo.PatchById(ctx, xconfig.Id, repos.Document{
		fc.RecordVersion: xconfig.RecordVersion + 1,
		fc.DisplayName:   config.DisplayName,
		fc.LastUpdatedBy: common.CreatedOrUpdatedBy{
			UserId:    ctx.UserId,
			UserName:  ctx.UserName,
			UserEmail: ctx.UserEmail,
		},

		fc.MetadataLabels:      config.Labels,
		fc.MetadataAnnotations: config.Annotations,

		fc.ConfigData: config.Data,

		fc.SyncStatusState:           t.SyncStateInQueue,
		fc.SyncStatusSyncScheduledAt: time.Now(),
		fc.SyncStatusAction:          t.SyncActionApply,
	})
	if err != nil {
		return nil, errors.NewE(err)
	}

	if upConfig.Annotations == nil {
		upConfig.Annotations = make(map[string]string)
	}

	for k, v := range types.ConfigWatchingAnnotation {
		upConfig.Annotations[k] = v
	}

	if err := d.applyK8sResource(ctx, xconfig.ProjectName, &upConfig.ConfigMap, upConfig.RecordVersion); err != nil {
		return upConfig, errors.NewE(err)
	}

	return upConfig, nil
}

func (d *domain) DeleteConfig(ctx ResourceContext, name string) error {
	if err := d.canMutateResourcesInEnvironment(ctx); err != nil {
		return errors.NewE(err)
	}

	c, err := d.findConfig(ctx, name)
	if err != nil {
		return errors.NewE(err)
	}

	if _, err := d.configRepo.PatchById(ctx, c.Id, repos.Document{
		fc.MarkedForDeletion:         true,
		fc.SyncStatusSyncScheduledAt: time.Now(),
		fc.SyncStatusAction:          t.SyncActionDelete,
		fc.SyncStatusState:           t.SyncStateInQueue,
	}); err != nil {
		return errors.NewE(err)
	}

	if err := d.deleteK8sResource(ctx, c.ProjectName, &c.ConfigMap); err != nil {
		if errors.Is(err, ErrNoClusterAttached) {
			return d.configRepo.DeleteById(ctx, c.Id)
		}
		return errors.NewE(err)
	}
	return nil
}

func (d *domain) OnConfigApplyError(ctx ResourceContext, errMsg, name string, opts UpdateAndDeleteOpts) error {
	c, err := d.findConfig(ctx, name)
	if err != nil {
		return errors.NewE(err)
	}

	if c == nil {
		return errors.Newf("no config found")
	}

	_, err = d.configRepo.PatchById(ctx, c.Id, repos.Document{
		fc.SyncStatusState:        t.SyncStateErroredAtAgent,
		fc.SyncStatusLastSyncedAt: opts.MessageTimestamp,
		fc.SyncStatusError:        errMsg,
	})
	return errors.NewE(err)
}

func (d *domain) OnConfigDeleteMessage(ctx ResourceContext, config entities.Config) error {
	c, err := d.findConfig(ctx, config.Name)
	if err != nil {
		return errors.NewE(err)
	}

	return d.configRepo.DeleteById(ctx, c.Id)
}

func (d *domain) OnConfigUpdateMessage(ctx ResourceContext, configIn entities.Config, status types.ResourceStatus, opts UpdateAndDeleteOpts) error {
	xconfig, err := d.findConfig(ctx, configIn.Name)
	if err != nil {
		return errors.NewE(err)
	}

	if xconfig == nil {
		return errors.Newf("no config found")
	}

	if err := d.MatchRecordVersion(configIn.Annotations, xconfig.RecordVersion); err != nil {
		return d.resyncK8sResource(ctx, xconfig.ProjectName, xconfig.SyncStatus.Action, &xconfig.ConfigMap, xconfig.RecordVersion)
	}

	_, err = d.configRepo.PatchById(ctx, xconfig.Id, repos.Document{
		fc.MetadataCreationTimestamp: configIn.CreationTimestamp,
		fc.MetadataLabels:            configIn.Labels,
		fc.MetadataAnnotations:       configIn.Annotations,
		fc.MetadataGeneration:        configIn.Generation,

		fc.SyncStatusState: func() t.SyncState {
			if status == types.ResourceStatusDeleting {
				return t.SyncStateDeletingAtAgent
			}
			return t.SyncStateUpdatedAtAgent
		}(),
		fc.SyncStatusRecordVersion: xconfig.RecordVersion,
		fc.SyncStatusLastSyncedAt:  opts.MessageTimestamp,
		fc.SyncStatusError:         nil,
	})
	return errors.NewE(err)
}

func (d *domain) ResyncConfig(ctx ResourceContext, name string) error {
	if err := d.canMutateResourcesInEnvironment(ctx); err != nil {
		return errors.NewE(err)
	}

	cfg, err := d.findConfig(ctx, name)
	if err != nil {
		return errors.NewE(err)
	}

	return d.resyncK8sResource(ctx, cfg.ProjectName, cfg.SyncStatus.Action, &cfg.ConfigMap, cfg.RecordVersion)
}
