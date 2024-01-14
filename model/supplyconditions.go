package model

type ConditionIdType uint

type SupplyConditionEventTypeType string

const (
	SupplyConditionEventTypeTypeThesholdExceeded     SupplyConditionEventTypeType = "thesholdExceeded"
	SupplyConditionEventTypeTypeFallenBelowThreshold SupplyConditionEventTypeType = "fallenBelowThreshold"
	SupplyConditionEventTypeTypeSupplyInterrupt      SupplyConditionEventTypeType = "supplyInterrupt"
	SupplyConditionEventTypeTypeReleaseOfLimitations SupplyConditionEventTypeType = "releaseOfLimitations"
	SupplyConditionEventTypeTypeOtherProblem         SupplyConditionEventTypeType = "otherProblem"
	SupplyConditionEventTypeTypeGridConditionUpdate  SupplyConditionEventTypeType = "gridConditionUpdate"
)

type SupplyConditionOriginatorType string

const (
	SupplyConditionOriginatorTypeExternDSO       SupplyConditionOriginatorType = "externDSO"
	SupplyConditionOriginatorTypeExternSupplier  SupplyConditionOriginatorType = "externSupplier"
	SupplyConditionOriginatorTypeInternalLimit   SupplyConditionOriginatorType = "internalLimit"
	SupplyConditionOriginatorTypeInternalService SupplyConditionOriginatorType = "internalService"
	SupplyConditionOriginatorTypeInternalUser    SupplyConditionOriginatorType = "internalUser"
)

type GridConditionType string

const (
	GridConditionTypeConsumptionRed    GridConditionType = "consumptionRed"
	GridConditionTypeConsumptionYellow GridConditionType = "consumptionYellow"
	GridConditionTypeGood              GridConditionType = "good"
	GridConditionTypeProductionRed     GridConditionType = "productionRed"
	GridConditionTypeProductionYellow  GridConditionType = "productionYellow"
)

type SupplyConditionDataType struct {
	ConditionId         *ConditionIdType               `json:"conditionId,omitempty" eebus:"key"`
	Timestamp           *AbsoluteOrRelativeTimeType    `json:"timestamp,omitempty"`
	EventType           *SupplyConditionEventTypeType  `json:"eventType,omitempty"`
	Originator          *SupplyConditionOriginatorType `json:"originator,omitempty"`
	ThresholdId         *ThresholdIdType               `json:"thresholdId,omitempty"`
	ThresholdPercentage *ScaledNumberType              `json:"thresholdPercentage,omitempty"`
	RelevantPeriod      *TimePeriodType                `json:"relevantPeriod,omitempty"`
	Description         *DescriptionType               `json:"description,omitempty"`
	GridCondition       *GridConditionType             `json:"gridCondition,omitempty"`
}

type SupplyConditionDataElementsType struct {
	ConditionId         *ElementTagType           `json:"conditionId,omitempty"`
	Timestamp           *ElementTagType           `json:"timestamp,omitempty"`
	EventType           *ElementTagType           `json:"eventType,omitempty"`
	Originator          *ElementTagType           `json:"originator,omitempty"`
	ThresholdId         *ElementTagType           `json:"thresholdId,omitempty"`
	ThresholdPercentage *ScaledNumberElementsType `json:"thresholdPercentage,omitempty"`
	RelevantPeriod      *TimePeriodElementsType   `json:"relevantPeriod,omitempty"`
	Description         *ElementTagType           `json:"description,omitempty"`
	GridCondition       *ElementTagType           `json:"gridCondition,omitempty"`
}

type SupplyConditionListDataType struct {
	SupplyConditionData []SupplyConditionDataType `json:"supplyConditionData,omitempty"`
}

type SupplyConditionListDataSelectorsType struct {
	ConditionId       *ConditionIdType               `json:"conditionId,omitempty"`
	TimestampInterval *TimestampIntervalType         `json:"timestampInterval,omitempty"`
	EventType         *SupplyConditionEventTypeType  `json:"eventType,omitempty"`
	Originator        *SupplyConditionOriginatorType `json:"originator,omitempty"`
}

type SupplyConditionDescriptionDataType struct {
	ConditionId             *ConditionIdType     `json:"conditionId,omitempty" eebus:"key"`
	CommodityType           *CommodityTypeType   `json:"commodityType,omitempty"`
	PositiveEnergyDirection *EnergyDirectionType `json:"positiveEnergyDirection,omitempty"`
	Label                   *LabelType           `json:"label,omitempty"`
	Description             *DescriptionType     `json:"description,omitempty"`
}

type SupplyConditionDescriptionDataElementsType struct {
	ConditionId             *ElementTagType `json:"conditionId,omitempty"`
	CommodityType           *ElementTagType `json:"commodityType,omitempty"`
	PositiveEnergyDirection *ElementTagType `json:"positiveEnergyDirection,omitempty"`
	Label                   *ElementTagType `json:"label,omitempty"`
	Description             *ElementTagType `json:"description,omitempty"`
}

type SupplyConditionDescriptionListDataType struct {
	SupplyConditionDescriptionData []SupplyConditionDescriptionDataType `json:"supplyConditionDescriptionData,omitempty"`
}

type SupplyConditionDescriptionListDataSelectorsType struct {
	ConditionId *ConditionIdType `json:"conditionId,omitempty"`
}

type SupplyConditionThresholdRelationDataType struct {
	ConditionId *ConditionIdType  `json:"conditionId,omitempty" eebus:"key"`
	ThresholdId []ThresholdIdType `json:"thresholdId,omitempty"`
}

type SupplyConditionThresholdRelationDataElementsType struct {
	ConditionId *ElementTagType `json:"conditionId,omitempty"`
	ThresholdId *ElementTagType `json:"thresholdId,omitempty"`
}

type SupplyConditionThresholdRelationListDataType struct {
	SupplyConditionThresholdRelationData []SupplyConditionThresholdRelationDataType `json:"SupplyConditionThresholdRelationDataType,omitempty"`
}

type SupplyConditionThresholdRelationListDataSelectorsType struct {
	ConditionId *ConditionIdType `json:"conditionId,omitempty"`
	ThresholdId *ThresholdIdType `json:"thresholdId,omitempty"`
}
