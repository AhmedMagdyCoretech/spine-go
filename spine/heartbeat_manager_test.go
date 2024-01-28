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

func TestHeartbeatManagerSuite(t *testing.T) {
	suite.Run(t, new(HeartBeatManagerSuite))
}

type HeartBeatManagerSuite struct {
	suite.Suite

	localDevice  api.DeviceLocalInterface
	remoteDevice api.DeviceRemoteInterface
	sut          api.HeartbeatManagerInterface
}

func (suite *HeartBeatManagerSuite) WriteShipMessageWithPayload([]byte) {}

func (suite *HeartBeatManagerSuite) SetupSuite() {
	suite.localDevice = NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)

	ski := "test"
	sender := NewSender(suite)
	suite.remoteDevice = NewDeviceRemote(suite.localDevice, ski, sender)

	_ = suite.localDevice.SetupRemoteDevice(ski, suite)

	suite.sut = suite.localDevice.HeartbeatManager()
}

func (suite *HeartBeatManagerSuite) Test_HeartbeatFailure() {
	err := suite.sut.StartHeartbeat()
	assert.NotNil(suite.T(), err)
}

func (suite *HeartBeatManagerSuite) Test_HeartbeatSuccess() {
	entity := NewEntityLocal(suite.localDevice, model.EntityTypeTypeCEM, []model.AddressEntityType{1})
	suite.localDevice.AddEntity(entity)

	localFeature := entity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	localFeature.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, false, false)
	entity.AddFeature(localFeature)

	remoteEntity := NewEntityRemote(suite.remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})
	suite.remoteDevice.AddEntity(remoteEntity)

	remoteFeature := NewFeatureRemote(remoteEntity.NextFeatureId(), remoteEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	remoteEntity.AddFeature(remoteFeature)

	subscrRequest := &model.SubscriptionManagementRequestCallType{
		ClientAddress:     remoteFeature.Address(),
		ServerAddress:     localFeature.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeDeviceDiagnosis),
	}

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource: &model.FeatureAddressType{
				Device:  suite.remoteDevice.Address(),
				Entity:  []model.AddressEntityType{0},
				Feature: util.Ptr(model.AddressFeatureType(0)),
			},
			AddressDestination: &model.FeatureAddressType{
				Device:  suite.localDevice.Address(),
				Entity:  []model.AddressEntityType{0},
				Feature: util.Ptr(model.AddressFeatureType(0)),
			},
			MsgCounter:    util.Ptr(model.MsgCounterType(1000)),
			CmdClassifier: util.Ptr(model.CmdClassifierTypeCall),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{
				{
					NodeManagementSubscriptionRequestCall: &model.NodeManagementSubscriptionRequestCallType{
						SubscriptionRequest: subscrRequest,
					},
				},
			},
		},
	}
	err := suite.localDevice.ProcessCmd(datagram, suite.remoteDevice)
	assert.Nil(suite.T(), err)

	data := localFeature.DataCopy(model.FunctionTypeDeviceDiagnosisHeartbeatData)
	assert.Nil(suite.T(), data)

	running := suite.sut.IsHeartbeatRunning()
	assert.Equal(suite.T(), false, running)

	suite.localDevice.RemoveEntity(entity)
	entity = NewEntityLocal(suite.localDevice, model.EntityTypeTypeCEM, []model.AddressEntityType{1})
	suite.localDevice.AddEntity(entity)

	localFeature = entity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	localFeature.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)
	entity.AddFeature(localFeature)

	err = suite.localDevice.ProcessCmd(datagram, suite.remoteDevice)
	assert.Nil(suite.T(), err)

	err = suite.sut.StartHeartbeat()
	assert.Nil(suite.T(), err)

	time.Sleep(time.Second * 5)

	running = suite.sut.IsHeartbeatRunning()
	assert.Equal(suite.T(), true, running)

	data = localFeature.DataCopy(model.FunctionTypeDeviceDiagnosisHeartbeatData)
	assert.NotNil(suite.T(), data)

	fctData := data.(*model.DeviceDiagnosisHeartbeatDataType)
	var resultCounter uint64 = 1
	assert.Equal(suite.T(), resultCounter, *fctData.HeartbeatCounter)
	resultTimeout, err := fctData.HeartbeatTimeout.GetTimeDuration()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), time.Second*4, resultTimeout)

	subscrDelRequest := &model.SubscriptionManagementDeleteCallType{
		ClientAddress: remoteFeature.Address(),
		ServerAddress: localFeature.Address(),
	}

	datagram.Payload = model.PayloadType{
		Cmd: []model.CmdType{
			{
				NodeManagementSubscriptionDeleteCall: &model.NodeManagementSubscriptionDeleteCallType{
					SubscriptionDelete: subscrDelRequest,
				},
			},
		},
	}

	err = suite.localDevice.ProcessCmd(datagram, suite.remoteDevice)
	assert.Nil(suite.T(), err)

	isHeartbeatRunning := suite.sut.IsHeartbeatRunning()
	assert.Equal(suite.T(), false, isHeartbeatRunning)

	suite.sut.StopHeartbeat()
}
