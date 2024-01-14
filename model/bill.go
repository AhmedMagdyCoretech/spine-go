package model

type BillIdType uint

type BillTypeType string

const (
	BillTypeTypeChargingSummary BillTypeType = "chargingSummary"
)

type BillPositionIdType uint

type BillPositionCountType BillPositionIdType

type BillPositionTypeType string

const (
	BillPositionTypeTypeGridElectricEnergy         BillPositionTypeType = "gridElectricEnergy"
	BillPositionTypeTypeSelfProducedElectricEnergy BillPositionTypeType = "selfProducedElectricEnergy"
)

type BillValueIdType uint

type BillCostIdType uint

type BillCostTypeType string

const (
	BillCostTypeTypeAbsolutePrice    BillCostTypeType = "absolutePrice"
	BillCostTypeTypeRelativePrice    BillCostTypeType = "relativePrice"
	BillCostTypeTypeCo2Emission      BillCostTypeType = "co2Emission"
	BillCostTypeTypeRenewableEnergy  BillCostTypeType = "renewableEnergy"
	BillCostTypeTypeRadioactiveWaste BillCostTypeType = "radioactiveWaste"
)

type BillValueType struct {
	ValueId         *BillValueIdType       `json:"valueId,omitempty"`
	Unit            *UnitOfMeasurementType `json:"unit,omitempty"`
	Value           *ScaledNumberType      `json:"value,omitempty"`
	ValuePercentage *ScaledNumberType      `json:"valuePercentage,omitempty"`
}

type BillValueElementsType struct {
	ValueId         *ElementTagType `json:"valueId,omitempty"`
	Unit            *ElementTagType `json:"unit,omitempty"`
	Value           *ElementTagType `json:"value,omitempty"`
	ValuePercentage *ElementTagType `json:"valuePercentage,omitempty"`
}

type BillCostType struct {
	CostId         *BillCostIdType        `json:"costId,omitempty"`
	CostType       *BillCostTypeType      `json:"costType,omitempty"`
	ValueId        *BillValueIdType       `json:"valueId,omitempty"`
	Unit           *UnitOfMeasurementType `json:"unit,omitempty"`
	Currency       *CurrencyType          `json:"currency,omitempty"`
	Cost           *ScaledNumberType      `json:"cost,omitempty"`
	CostPercentage *ScaledNumberType      `json:"costPercentage,omitempty"`
}

type BillCostElementsType struct {
	CostId         *ElementTagType           `json:"costId,omitempty"`
	CostType       *ElementTagType           `json:"costType,omitempty"`
	ValueId        *ElementTagType           `json:"valueId,omitempty"`
	Unit           *ElementTagType           `json:"unit,omitempty"`
	Currency       *ElementTagType           `json:"currency,omitempty"`
	Cost           *ScaledNumberElementsType `json:"cost,omitempty"`
	CostPercentage *ScaledNumberElementsType `json:"costPercentage,omitempty"`
}

type BillPositionType struct {
	PositionId   *BillPositionIdType   `json:"positionId,omitempty"`
	PositionType *BillPositionTypeType `json:"positionType,omitempty"`
	TimePeriod   *TimePeriodType       `json:"timePeriod,omitempty"`
	Value        *BillValueType        `json:"value,omitempty"`
	Cost         *BillCostType         `json:"cost,omitempty"`
	Label        *LabelType            `json:"label,omitempty"`
	Description  *DescriptionType      `json:"description,omitempty"`
}

type BillPositionElementsType struct {
	PositionId   *ElementTagType         `json:"positionId,omitempty"`
	PositionType *ElementTagType         `json:"positionType,omitempty"`
	TimePeriod   *TimePeriodElementsType `json:"timePeriod,omitempty"`
	Value        *BillValueElementsType  `json:"value,omitempty"`
	Cost         *BillCostElementsType   `json:"cost,omitempty"`
	Label        *ElementTagType         `json:"label,omitempty"`
	Description  *ElementTagType         `json:"description,omitempty"`
}

type BillDataType struct {
	BillId    *BillIdType        `json:"billId,omitempty" eebus:"key"`
	BillType  *BillTypeType      `json:"billType,omitempty"`
	ScopeType *ScopeTypeType     `json:"scopeType,omitempty"`
	Total     *BillPositionType  `json:"total,omitempty"`
	Position  []BillPositionType `json:"position,omitempty"`
}

type BillDataElementsType struct {
	BillId    *ElementTagType           `json:"billId,omitempty"`
	BillType  *ElementTagType           `json:"billType,omitempty"`
	ScopeType *ElementTagType           `json:"scopeType,omitempty"`
	Total     *BillPositionElementsType `json:"total,omitempty"`
	Position  *BillPositionElementsType `json:"position,omitempty"`
}

type BillListDataType struct {
	BillData []BillDataType `json:"billData,omitempty"`
}

type BillListDataSelectorsType struct {
	BillId    *BillIdType    `json:"billId,omitempty"`
	ScopeType *ScopeTypeType `json:"scopeType,omitempty"`
}

type BillConstraintsDataType struct {
	BillId           *BillIdType            `json:"billId,omitempty" eebus:"key"`
	PositionCountMin *BillPositionCountType `json:"positionCountMin,omitempty"`
	PositionCountMax *BillPositionCountType `json:"positionCountMax,omitempty"`
}

type BillConstraintsDataElementsType struct {
	BillId           *ElementTagType `json:"billId,omitempty"`
	PositionCountMin *ElementTagType `json:"positionCountMin,omitempty"`
	PositionCountMax *ElementTagType `json:"positionCountMax,omitempty"`
}

type BillConstraintsListDataType struct {
	BillConstraintsData []BillConstraintsDataType `json:"billConstraintsData,omitempty"`
}

type BillConstraintsListDataSelectorsType struct {
	BillId *BillIdType `json:"billId,omitempty"`
}

type BillDescriptionDataType struct {
	BillId            *BillIdType    `json:"billId,omitempty" eebus:"key"`
	BillWriteable     *bool          `json:"billWriteable,omitempty"`
	UpdateRequired    *bool          `json:"updateRequired,omitempty"`
	SupportedBillType []BillTypeType `json:"supportedBillType,omitempty"`
	SessionId         *SessionIdType `json:"sessionId,omitempty"`
}

type BillDescriptionDataElementsType struct {
	BillId            *ElementTagType `json:"billId,omitempty"`
	BillWriteable     *ElementTagType `json:"billWriteable,omitempty"`
	UpdateRequired    *ElementTagType `json:"updateRequired,omitempty"`
	SupportedBillType *ElementTagType `json:"supportedBillType,omitempty"`
	SessionId         *ElementTagType `json:"sessionId,omitempty"`
}

type BillDescriptionListDataType struct {
	BillDescriptionData []BillDescriptionDataType `json:"billDescriptionData,omitempty"`
}

type BillDescriptionListDataSelectorsType struct {
	BillId *BillIdType `json:"billId,omitempty"`
}
