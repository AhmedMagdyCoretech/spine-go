package spine

import (
	"testing"
	"time"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestBindingManagerSuite(t *testing.T) {
	suite.Run(t, new(BindingManagerSuite))
}

type BindingManagerSuite struct {
	suite.Suite

	localDevice  api.DeviceLocal
	writeHandler *WriteMessageHandler
	remoteDevice api.DeviceRemote

	sut api.BindingManager
}

func (s *BindingManagerSuite) BeforeTest(suiteName, testName string) {
	s.localDevice = NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)
	remoteSki := "TestRemoteSki"

	s.writeHandler = &WriteMessageHandler{}

	sender := NewSender(s.writeHandler)
	s.remoteDevice = NewDeviceRemoteImpl(s.localDevice, remoteSki, sender)

	s.sut = NewBindingManager(s.localDevice)
}

func (suite *BindingManagerSuite) Test_Bindings() {
	entity := NewEntityLocalImpl(suite.localDevice, model.EntityTypeTypeCEM, []model.AddressEntityType{1})
	suite.localDevice.AddEntity(entity)

	localFeature := entity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)

	remoteEntity := NewEntityRemoteImpl(suite.remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})
	suite.remoteDevice.AddEntity(remoteEntity)

	remoteFeature := NewFeatureRemoteImpl(remoteEntity.NextFeatureId(), remoteEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	remoteFeature.Address().Device = util.Ptr(model.AddressDeviceType("remoteDevice"))
	remoteEntity.AddFeature(remoteFeature)

	bindingRequest := model.BindingManagementRequestCallType{
		ClientAddress:     remoteFeature.Address(),
		ServerAddress:     localFeature.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeDeviceDiagnosis),
	}

	bindingMgr := suite.localDevice.BindingManager()
	err := bindingMgr.AddBinding(suite.remoteDevice, bindingRequest)
	assert.Nil(suite.T(), err)

	subs := bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 1, len(subs))

	err = bindingMgr.AddBinding(suite.remoteDevice, bindingRequest)
	assert.NotNil(suite.T(), err)

	subs = bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 1, len(subs))

	address := model.FeatureAddressType{
		Device:  entity.Device().Address(),
		Entity:  entity.Address().Entity,
		Feature: util.Ptr(model.AddressFeatureType(10)),
	}
	entries := bindingMgr.BindingsOnFeature(address)
	assert.Equal(suite.T(), 0, len(entries))

	address.Feature = localFeature.Address().Feature
	entries = bindingMgr.BindingsOnFeature(address)
	assert.Equal(suite.T(), 1, len(entries))

	bindingDelete := model.BindingManagementDeleteCallType{
		ClientAddress: remoteFeature.Address(),
		ServerAddress: localFeature.Address(),
	}

	err = bindingMgr.RemoveBinding(bindingDelete, suite.remoteDevice)
	assert.Nil(suite.T(), err)

	subs = bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 0, len(subs))

	err = bindingMgr.RemoveBinding(bindingDelete, suite.remoteDevice)
	assert.NotNil(suite.T(), err)

	err = bindingMgr.AddBinding(suite.remoteDevice, bindingRequest)
	assert.Nil(suite.T(), err)

	subs = bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 1, len(subs))

	bindingMgr.RemoveBindingsForDevice(suite.remoteDevice)

	subs = bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 0, len(subs))
}
