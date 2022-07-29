package domain

import (
	"context"
	"fmt"
	"kloudlite.io/apps/console/internal/domain/entities"
	internal_crds "kloudlite.io/apps/console/internal/domain/op-crds/internal-crds"
	"kloudlite.io/pkg/repos"
)

func (d *domain) GetDevice(ctx context.Context, id repos.ID) (*entities.Device, error) {
	return d.deviceRepo.FindById(ctx, id)
}
func (d *domain) ListAccountDevices(ctx context.Context, accountId repos.ID) ([]*entities.Device, error) {
	q := make(repos.Filter)
	q["account_id"] = accountId
	return d.deviceRepo.Find(ctx, repos.Query{
		Filter: q,
	})
}
func (d *domain) ListUserDevices(ctx context.Context, userId repos.ID) ([]*entities.Device, error) {
	q := make(repos.Filter)
	q["user_id"] = userId
	return d.deviceRepo.Find(ctx, repos.Query{
		Filter: q,
	})
}

func (d *domain) GetDeviceConfig(ctx context.Context, deviceId repos.ID) (map[string]any, error) {
	dev, err := d.deviceRepo.FindById(ctx, deviceId)
	if err != nil {
		return nil, err
	}
	secret, err := d.kubeCli.GetSecret(ctx, fmt.Sprint("wg-", dev.AccountId), fmt.Sprint("wg-device-config-", dev.Id))
	if err != nil {
		return nil, err
	}
	parsedSec := make(map[string]any)
	for k, v := range secret.Data {
		parsedSec[k] = string(v)
	}
	return parsedSec, nil
}

func (d *domain) DeviceByNameExists(ctx context.Context, accountId repos.ID, name string) (bool, error) {
	one, err := d.deviceRepo.FindOne(ctx, repos.Filter{
		"account_id": accountId,
		"name":       name,
	})
	if err != nil {
		return false, err
	}
	return one != nil, nil
}

func (d *domain) AddDevice(ctx context.Context, deviceName string, accountId repos.ID, userId repos.ID) (*entities.Device, error) {
	devices, err := d.deviceRepo.Find(ctx, repos.Query{
		Filter: repos.Filter{
			"account_id": accountId,
		},
		Sort: map[string]any{
			"index": 1,
		},
	})
	if err != nil {
		return nil, err
	}
	index := -1
	count := 0
	for i, d := range devices {
		count++
		if d.Index != i {
			index = i
			break
		}
	}
	if index == -1 {
		index = count
	}
	device, e := d.deviceRepo.Create(ctx, &entities.Device{
		Name:      deviceName,
		AccountId: accountId,
		UserId:    userId,
		Status:    entities.DeviceStateSyncing,
		Index:     index,
	})
	if e != nil {
		return nil, fmt.Errorf("unable to persist in db %v", e)
	}
	err = d.workloadMessenger.SendAction("apply", string(device.Id), &internal_crds.Device{
		APIVersion: internal_crds.DeviceAPIVersion,
		Kind:       internal_crds.DeviceKind,
		Metadata: internal_crds.DeviceMetadata{
			Name: string(device.Id),
			Annotations: map[string]string{
				"kloudlite.io/account-id": string(device.AccountId),
			},
			Labels: map[string]string{
				"kloudlite.io/account-id": string(device.AccountId),
				"kloudlite.io/device-id":  string(device.Id),
			},
		},
		Spec: internal_crds.DeviceSpec{
			Account: string(device.AccountId),
			ActiveRegion: func() string {
				if device.ActiveRegion != nil {
					return *device.ActiveRegion
				}
				return ""
			}(),
			DeviceName: deviceName,
			Offset:     device.Index,
			DeviceId:   string(device.Id),
			Ports:      device.ExposedPorts,
		},
	})
	if err != nil {
		return device, err
	}
	return device, e
}
func (d *domain) RemoveDevice(ctx context.Context, deviceId repos.ID) error {
	device, err := d.deviceRepo.FindById(ctx, deviceId)
	if err != nil {
		return err
	}
	device.Status = entities.DeviceStateSyncing
	_, err = d.deviceRepo.UpdateById(ctx, deviceId, device)
	if err != nil {
		return err
	}
	err = d.workloadMessenger.SendAction("delete", string(device.Id), &internal_crds.Device{
		APIVersion: internal_crds.DeviceAPIVersion,
		Kind:       internal_crds.DeviceKind,
		Metadata: internal_crds.DeviceMetadata{
			Name: string(device.Id),
		},
	})
	return err
}
func (d *domain) UpdateDevice(ctx context.Context, deviceId repos.ID, deviceName string, region *string, ports []int32) (done bool, e error) {
	device, e := d.deviceRepo.FindById(ctx, deviceId)
	if region != nil {
		device.ActiveRegion = region
	}
	device.Name = deviceName
	device.ExposedPorts = ports
	_, err := d.deviceRepo.UpdateById(ctx, deviceId, device)
	if err != nil {
		return false, e
	}
	err = d.workloadMessenger.SendAction("apply", string(device.Id), &internal_crds.Device{
		APIVersion: internal_crds.DeviceAPIVersion,
		Kind:       internal_crds.DeviceKind,
		Metadata: internal_crds.DeviceMetadata{
			Name: string(device.Id),
			Annotations: map[string]string{
				"kloudlite.io/account-id": string(device.AccountId),
			},
			Labels: map[string]string{
				"kloudlite.io/account-id": string(device.AccountId),
				"kloudlite.io/device-id":  string(device.Id),
			},
		},
		Spec: internal_crds.DeviceSpec{
			Account: string(device.AccountId),
			ActiveRegion: func() string {
				if device.ActiveRegion != nil {
					return *device.ActiveRegion
				}
				return ""
			}(),
			Offset:   device.Index,
			DeviceId: string(device.Id),
			Ports:    device.ExposedPorts,
		},
	})
	if err != nil {
		return false, err
	}
	return true, nil

}
