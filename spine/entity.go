package spine

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

const DeviceInformationEntityId uint = 0

var DeviceInformationAddressEntity = []model.AddressEntityType{model.AddressEntityType(DeviceInformationEntityId)}

type EntityImpl struct {
	eType        model.EntityTypeType
	address      *model.EntityAddressType
	description  *model.DescriptionType
	fIdGenerator func() uint
}

var _ api.Entity = (*EntityImpl)(nil)

func NewEntity(eType model.EntityTypeType, deviceAdress *model.AddressDeviceType, entityAddress []model.AddressEntityType) *EntityImpl {
	entity := &EntityImpl{
		eType: eType,
		address: &model.EntityAddressType{
			Device: deviceAdress,
			Entity: entityAddress,
		},
	}
	if entityAddress[0] == 0 {
		// Entity 0 Feature addresses start with 0
		entity.fIdGenerator = newFeatureIdGenerator(0)
	} else {
		// Entity 1 and further Feature addresses start with 1
		entity.fIdGenerator = newFeatureIdGenerator(1)
	}

	return entity
}

func (r *EntityImpl) EntityType() model.EntityTypeType {
	return r.eType
}

func (r *EntityImpl) Address() *model.EntityAddressType {
	return r.address
}

func (r *EntityImpl) Description() *model.DescriptionType {
	return r.description
}

func (r *EntityImpl) SetDescription(d *model.DescriptionType) {
	r.description = d
}

func (r *EntityImpl) NextFeatureId() uint {
	return r.fIdGenerator()
}

func EntityAddressType(deviceAdress *model.AddressDeviceType, entityAddress []model.AddressEntityType) *model.EntityAddressType {
	return &model.EntityAddressType{
		Device: deviceAdress,
		Entity: entityAddress,
	}
}

func NewEntityAddressType(deviceName string, entityIds []uint) *model.EntityAddressType {
	return &model.EntityAddressType{
		Device: util.Ptr(model.AddressDeviceType(deviceName)),
		Entity: NewAddressEntityType(entityIds),
	}
}

func NewAddressEntityType(entityIds []uint) []model.AddressEntityType {
	var addressEntity []model.AddressEntityType
	linq.From(entityIds).SelectT(func(i uint) model.AddressEntityType { return model.AddressEntityType(i) }).ToSlice(&addressEntity)
	return addressEntity
}

func newFeatureIdGenerator(id uint) func() uint {
	return func() uint {
		defer func() { id += 1 }()
		return id
	}
}
