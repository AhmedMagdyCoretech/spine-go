package model

// DeviceConfigurationKeyValueListDataType

var _ Updater = (*DeviceConfigurationKeyValueListDataType)(nil)

func (r *DeviceConfigurationKeyValueListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []DeviceConfigurationKeyValueDataType
	if newList != nil {
		newData = newList.(*DeviceConfigurationKeyValueListDataType).DeviceConfigurationKeyValueData
	}

	r.DeviceConfigurationKeyValueData = UpdateList(remoteWrite, r.DeviceConfigurationKeyValueData, newData, filterPartial, filterDelete)
}

// DeviceConfigurationKeyValueDescriptionListDataType

var _ Updater = (*DeviceConfigurationKeyValueDescriptionListDataType)(nil)

func (r *DeviceConfigurationKeyValueDescriptionListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []DeviceConfigurationKeyValueDescriptionDataType
	if newList != nil {
		newData = newList.(*DeviceConfigurationKeyValueDescriptionListDataType).DeviceConfigurationKeyValueDescriptionData
	}

	r.DeviceConfigurationKeyValueDescriptionData = UpdateList(remoteWrite, r.DeviceConfigurationKeyValueDescriptionData, newData, filterPartial, filterDelete)
}

// DeviceConfigurationKeyValueConstraintsListDataType

var _ Updater = (*DeviceConfigurationKeyValueConstraintsListDataType)(nil)

func (r *DeviceConfigurationKeyValueConstraintsListDataType) UpdateList(remoteWrite bool, newList any, filterPartial, filterDelete *FilterType) {
	var newData []DeviceConfigurationKeyValueConstraintsDataType
	if newList != nil {
		newData = newList.(*DeviceConfigurationKeyValueConstraintsListDataType).DeviceConfigurationKeyValueConstraintsData
	}

	r.DeviceConfigurationKeyValueConstraintsData = UpdateList(remoteWrite, r.DeviceConfigurationKeyValueConstraintsData, newData, filterPartial, filterDelete)
}
