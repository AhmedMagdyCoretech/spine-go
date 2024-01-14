package spine

import (
	"testing"

	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestFunctionData_UpdateData(t *testing.T) {
	newData := &model.DeviceClassificationManufacturerDataType{
		DeviceName: util.Ptr(model.DeviceClassificationStringType("device name")),
	}
	functionType := model.FunctionTypeDeviceClassificationManufacturerData
	sut := NewFunctionData[model.DeviceClassificationManufacturerDataType](functionType)
	sut.UpdateData(newData, nil, nil)
	getData := sut.DataCopy()

	assert.Equal(t, newData.DeviceName, getData.DeviceName)
	assert.Equal(t, functionType, sut.Function())

	// another update should not be reflected in the first dataset
	newData = &model.DeviceClassificationManufacturerDataType{
		DeviceName: util.Ptr(model.DeviceClassificationStringType("new device name")),
	}
	sut.UpdateData(newData, nil, nil)
	getNewData := sut.DataCopy()

	assert.Equal(t, newData.DeviceName, getNewData.DeviceName)
	assert.NotEqual(t, getData.DeviceName, getNewData.DeviceName)
	assert.Equal(t, functionType, sut.Function())
}

func TestFunctionData_UpdateDataPartial(t *testing.T) {
	newData := &model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(1)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Range: []model.ScaledNumberRangeType{
							{
								Min: &model.ScaledNumberType{
									Number: util.Ptr(model.NumberType(6)),
									Scale:  util.Ptr(model.ScaleType(0)),
								},
							},
						},
					},
				},
			},
		},
	}
	functionType := model.FunctionTypeElectricalConnectionPermittedValueSetListData
	sut := NewFunctionData[model.ElectricalConnectionPermittedValueSetListDataType](functionType)

	err := sut.UpdateData(newData, &model.FilterType{CmdControl: &model.CmdControlType{Partial: &model.ElementTagType{}}}, nil)
	if assert.Nil(t, err) {
		getData := sut.DataCopy()
		assert.Equal(t, 1, len(getData.ElectricalConnectionPermittedValueSetData))
	}
}

func TestFunctionData_UpdateDataPartial_Supported(t *testing.T) {
	newData := &model.HvacOverrunListDataType{
		HvacOverrunData: []model.HvacOverrunDataType{
			{
				OverrunId: util.Ptr(model.HvacOverrunIdType(1)),
			},
		},
	}
	functionType := model.FunctionTypeHvacOverrunListData
	sut := NewFunctionData[model.HvacOverrunListDataType](functionType)

	err := sut.UpdateData(newData, &model.FilterType{CmdControl: &model.CmdControlType{Partial: &model.ElementTagType{}}}, nil)
	assert.Nil(t, err)
}
