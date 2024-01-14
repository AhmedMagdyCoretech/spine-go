package spine

import (
	"testing"
	"time"

	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestDeviceLocalSuite(t *testing.T) {
	suite.Run(t, new(DeviceLocalTestSuite))
}

type DeviceLocalTestSuite struct {
	suite.Suite
}

var _ shipapi.SpineDataConnection = (*DeviceLocalTestSuite)(nil)

func (d *DeviceLocalTestSuite) WriteSpineMessage([]byte) {}

func (d *DeviceLocalTestSuite) Test_RemoveRemoteDevice() {
	sut := NewDeviceLocalImpl("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)

	ski := "test"
	_ = sut.SetupRemoteDevice(ski, d)
	rDevice := sut.RemoteDeviceForSki(ski)
	assert.NotNil(d.T(), rDevice)

	sut.RemoveRemoteDeviceConnection(ski)

	rDevice = sut.RemoteDeviceForSki(ski)
	assert.Nil(d.T(), rDevice)
}

func (d *DeviceLocalTestSuite) Test_RemoteDevice() {
	sut := NewDeviceLocalImpl("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)
	localEntity := NewEntityLocalImpl(sut, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}))
	sut.AddEntity(localEntity)

	f := NewFeatureLocalImpl(1, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = NewFeatureLocalImpl(2, localEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	localEntity.AddFeature(f)

	ski := "test"
	remote := sut.RemoteDeviceForSki(ski)
	assert.Nil(d.T(), remote)

	devices := sut.RemoteDevices()
	assert.Equal(d.T(), 0, len(devices))

	_ = sut.SetupRemoteDevice(ski, d)
	remote = sut.RemoteDeviceForSki(ski)
	assert.NotNil(d.T(), remote)

	devices = sut.RemoteDevices()
	assert.Equal(d.T(), 1, len(devices))

	entities := sut.Entities()
	assert.Equal(d.T(), 2, len(entities))

	entity1 := sut.Entity([]model.AddressEntityType{1})
	assert.NotNil(d.T(), entity1)

	entity2 := sut.Entity([]model.AddressEntityType{2})
	assert.Nil(d.T(), entity2)

	featureAddress := &model.FeatureAddressType{
		Entity:  []model.AddressEntityType{1},
		Feature: util.Ptr(model.AddressFeatureType(1)),
	}
	feature1 := sut.FeatureByAddress(featureAddress)
	assert.NotNil(d.T(), feature1)

	feature2 := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	assert.NotNil(d.T(), feature2)

	sut.RemoveEntity(entity1)
	entities = sut.Entities()
	assert.Equal(d.T(), 1, len(entities))

	sut.RemoveRemoteDevice(ski)
	remote = sut.RemoteDeviceForSki(ski)
	assert.Nil(d.T(), remote)
}

func (d *DeviceLocalTestSuite) Test_ProcessCmd_Errors() {
	sut := NewDeviceLocalImpl("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)
	localEntity := NewEntityLocalImpl(sut, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}))
	sut.AddEntity(localEntity)

	ski := "test"
	_ = sut.SetupRemoteDevice(ski, d)
	remote := sut.RemoteDeviceForSki(ski)
	assert.NotNil(d.T(), remote)

	datagram := model.DatagramType{
		Header: model.HeaderType{
			AddressSource: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("localdevice")),
			},
			AddressDestination: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("localdevice")),
			},
			MsgCounter:    util.Ptr(model.MsgCounterType(1)),
			CmdClassifier: util.Ptr(model.CmdClassifierTypeRead),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{},
		},
	}

	err := sut.ProcessCmd(datagram, remote)
	assert.NotNil(d.T(), err)

	datagram = model.DatagramType{
		Header: model.HeaderType{
			AddressSource: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("localdevice")),
			},
			AddressDestination: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("localdevice")),
			},
			MsgCounter:    util.Ptr(model.MsgCounterType(1)),
			CmdClassifier: util.Ptr(model.CmdClassifierTypeRead),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{
				{},
			},
		},
	}

	err = sut.ProcessCmd(datagram, remote)
	assert.NotNil(d.T(), err)
}

func (d *DeviceLocalTestSuite) Test_ProcessCmd() {
	sut := NewDeviceLocalImpl("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)
	localEntity := NewEntityLocalImpl(sut, model.EntityTypeTypeCEM, NewAddressEntityType([]uint{1}))
	sut.AddEntity(localEntity)

	f := NewFeatureLocalImpl(1, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = NewFeatureLocalImpl(2, localEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = NewFeatureLocalImpl(3, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	localEntity.AddFeature(f)

	ski := "test"
	remoteDeviceName := "remote"
	_ = sut.SetupRemoteDevice(ski, d)
	remote := sut.RemoteDeviceForSki(ski)
	assert.NotNil(d.T(), remote)

	detailedData := &model.NodeManagementDetailedDiscoveryDataType{
		DeviceInformation: &model.NodeManagementDetailedDiscoveryDeviceInformationType{
			Description: &model.NetworkManagementDeviceDescriptionDataType{
				DeviceAddress: &model.DeviceAddressType{
					Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
				},
			},
		},
		EntityInformation: []model.NodeManagementDetailedDiscoveryEntityInformationType{
			{
				Description: &model.NetworkManagementEntityDescriptionDataType{
					EntityAddress: &model.EntityAddressType{
						Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity: []model.AddressEntityType{1},
					},
					EntityType: util.Ptr(model.EntityTypeTypeEVSE),
				},
			},
		},
		FeatureInformation: []model.NodeManagementDetailedDiscoveryFeatureInformationType{
			{
				Description: &model.NetworkManagementFeatureDescriptionDataType{
					FeatureAddress: &model.FeatureAddressType{
						Device:  util.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity:  []model.AddressEntityType{1},
						Feature: util.Ptr(model.AddressFeatureType(1)),
					},
					FeatureType: util.Ptr(model.FeatureTypeTypeElectricalConnection),
					Role:        util.Ptr(model.RoleTypeServer),
				},
			},
		},
	}
	_, err := remote.AddEntityAndFeatures(true, detailedData)
	assert.Nil(d.T(), err)

	datagram := model.DatagramType{
		Header: model.HeaderType{
			AddressSource: &model.FeatureAddressType{
				Device:  util.Ptr(model.AddressDeviceType(remoteDeviceName)),
				Entity:  []model.AddressEntityType{1},
				Feature: util.Ptr(model.AddressFeatureType(1)),
			},
			AddressDestination: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("localdevice")),
				Entity: []model.AddressEntityType{1},
			},
			MsgCounter:    util.Ptr(model.MsgCounterType(1)),
			CmdClassifier: util.Ptr(model.CmdClassifierTypeRead),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{},
		},
	}

	err = sut.ProcessCmd(datagram, remote)
	assert.NotNil(d.T(), err)

	cmd := model.CmdType{
		ElectricalConnectionParameterDescriptionListData: &model.ElectricalConnectionParameterDescriptionListDataType{},
	}

	datagram.Payload.Cmd = append(datagram.Payload.Cmd, cmd)

	err = sut.ProcessCmd(datagram, remote)
	assert.NotNil(d.T(), err)

	datagram.Header.AddressDestination.Feature = util.Ptr(model.AddressFeatureType(1))

	err = sut.ProcessCmd(datagram, remote)
	assert.NotNil(d.T(), err)

	datagram.Header.AddressDestination.Feature = util.Ptr(model.AddressFeatureType(3))

	err = sut.ProcessCmd(datagram, remote)
	assert.Nil(d.T(), err)
}
