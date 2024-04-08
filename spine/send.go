package spine

import (
	"encoding/json"
	"errors"
	"sync/atomic"

	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	lru "github.com/hashicorp/golang-lru/v2"
)

type Sender struct {
	msgNum uint64 // 64bit values need to be defined on top of the struct to make atomic commands work on 32bit systems

	// we cache the last 100 notify messages, so we can find the matching item for result errors being returned
	datagramNotifyCache *lru.Cache[model.MsgCounterType, model.DatagramType]

	writeHandler shipapi.ShipConnectionDataWriterInterface
}

var _ api.SenderInterface = (*Sender)(nil)

func NewSender(writeI shipapi.ShipConnectionDataWriterInterface) api.SenderInterface {
	cache, _ := lru.New[model.MsgCounterType, model.DatagramType](100)
	return &Sender{
		datagramNotifyCache: cache,
		writeHandler:        writeI,
	}
}

// return the datagram for a given msgCounter (only availbe for Notify messasges!), error if not found
func (c *Sender) DatagramForMsgCounter(msgCounter model.MsgCounterType) (model.DatagramType, error) {
	if datagram, ok := c.datagramNotifyCache.Get(msgCounter); ok {
		return datagram, nil
	}

	return model.DatagramType{}, errors.New("msgCounter not found")
}

func (c *Sender) sendSpineMessage(datagram model.DatagramType) error {
	// pack into datagram
	data := model.Datagram{
		Datagram: datagram,
	}

	// marshal
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if c.writeHandler == nil {
		return errors.New("outgoing interface implementation not set")
	}

	if msg == nil {
		return errors.New("message is nil")
	}

	logging.Log().Debug(datagram.PrintMessageOverview(true, "", ""))

	// write to channel
	c.writeHandler.WriteShipMessageWithPayload(msg)

	return nil
}

// Sends request
func (c *Sender) Request(cmdClassifier model.CmdClassifierType, senderAddress, destinationAddress *model.FeatureAddressType, ackRequest bool, cmd []model.CmdType) (*model.MsgCounterType, error) {
	msgCounter := c.getMsgCounter()

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        senderAddress,
			AddressDestination:   destinationAddress,
			MsgCounter:           msgCounter,
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: cmd,
		},
	}

	if ackRequest {
		datagram.Header.AckRequest = &ackRequest
	}

	return msgCounter, c.sendSpineMessage(datagram)
}

func (c *Sender) ResultSuccess(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType) error {
	return c.result(requestHeader, senderAddress, nil)
}

func (c *Sender) ResultError(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, err *model.ErrorType) error {
	return c.result(requestHeader, senderAddress, err)
}

// sends a result for a request
func (c *Sender) result(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, err *model.ErrorType) error {
	cmdClassifier := model.CmdClassifierTypeResult

	addressSource := *requestHeader.AddressDestination
	addressSource.Device = senderAddress.Device

	var resultData model.ResultDataType
	if err != nil {
		resultData = model.ResultDataType{
			ErrorNumber: &err.ErrorNumber,
			Description: err.Description,
		}
	} else {
		resultData = model.ResultDataType{
			ErrorNumber: util.Ptr(model.ErrorNumberTypeNoError),
		}
	}

	cmd := model.CmdType{
		ResultData: &resultData,
	}

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        &addressSource,
			AddressDestination:   requestHeader.AddressSource,
			MsgCounter:           c.getMsgCounter(),
			MsgCounterReference:  requestHeader.MsgCounter,
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	return c.sendSpineMessage(datagram)
}

// Reply sends reply to original sender
func (c *Sender) Reply(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, cmd model.CmdType) error {
	cmdClassifier := model.CmdClassifierTypeReply

	addressSource := *requestHeader.AddressDestination
	addressSource.Device = senderAddress.Device

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        &addressSource,
			AddressDestination:   requestHeader.AddressSource,
			MsgCounter:           c.getMsgCounter(),
			MsgCounterReference:  requestHeader.MsgCounter,
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	return c.sendSpineMessage(datagram)
}

// Notify sends notification to destination
func (c *Sender) Notify(senderAddress, destinationAddress *model.FeatureAddressType, cmd model.CmdType) (*model.MsgCounterType, error) {
	msgCounter := c.getMsgCounter()

	cmdClassifier := model.CmdClassifierTypeNotify

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        senderAddress,
			AddressDestination:   destinationAddress,
			MsgCounter:           msgCounter,
			CmdClassifier:        &cmdClassifier,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	c.datagramNotifyCache.Add(*msgCounter, datagram)

	return msgCounter, c.sendSpineMessage(datagram)
}

// Write sends notification to destination
func (c *Sender) Write(senderAddress, destinationAddress *model.FeatureAddressType, cmd model.CmdType) (*model.MsgCounterType, error) {
	msgCounter := c.getMsgCounter()

	cmdClassifier := model.CmdClassifierTypeWrite
	ackRequest := true

	datagram := model.DatagramType{
		Header: model.HeaderType{
			SpecificationVersion: &SpecificationVersion,
			AddressSource:        senderAddress,
			AddressDestination:   destinationAddress,
			MsgCounter:           msgCounter,
			CmdClassifier:        &cmdClassifier,
			AckRequest:           &ackRequest,
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{cmd},
		},
	}

	return msgCounter, c.sendSpineMessage(datagram)
}

// Send a subscription request to a remote server feature
func (c *Sender) Subscribe(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) (*model.MsgCounterType, error) {
	cmd := model.CmdType{
		NodeManagementSubscriptionRequestCall: NewNodeManagementSubscriptionRequestCallType(senderAddress, destinationAddress, serverFeatureType),
	}

	// we always send it to the remote NodeManagment feature, which always is at entity:[0],feature:0
	localAddress := NodeManagementAddress(senderAddress.Device)
	remoteAddress := NodeManagementAddress(destinationAddress.Device)

	return c.Request(model.CmdClassifierTypeCall, localAddress, remoteAddress, true, []model.CmdType{cmd})
}

// Send a subscription deletion request to a remote server feature
func (c *Sender) Unsubscribe(senderAddress, destinationAddress *model.FeatureAddressType) (*model.MsgCounterType, error) {
	cmd := model.CmdType{
		NodeManagementSubscriptionDeleteCall: NewNodeManagementSubscriptionDeleteCallType(senderAddress, destinationAddress),
	}

	// we always send it to the remote NodeManagment feature, which always is at entity:[0],feature:0
	localAddress := NodeManagementAddress(senderAddress.Device)
	remoteAddress := NodeManagementAddress(destinationAddress.Device)

	return c.Request(model.CmdClassifierTypeCall, localAddress, remoteAddress, true, []model.CmdType{cmd})
}

// Send a binding request to a remote server feature
func (c *Sender) Bind(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) (*model.MsgCounterType, error) {
	cmd := model.CmdType{
		NodeManagementBindingRequestCall: NewNodeManagementBindingRequestCallType(senderAddress, destinationAddress, serverFeatureType),
	}

	// we always send it to the remote NodeManagment feature, which always is at entity:[0],feature:0
	localAddress := NodeManagementAddress(senderAddress.Device)
	remoteAddress := NodeManagementAddress(destinationAddress.Device)

	return c.Request(model.CmdClassifierTypeCall, localAddress, remoteAddress, true, []model.CmdType{cmd})
}

// Send a binding request to a remote server feature
func (c *Sender) Unbind(senderAddress, destinationAddress *model.FeatureAddressType) (*model.MsgCounterType, error) {
	cmd := model.CmdType{
		NodeManagementBindingDeleteCall: NewNodeManagementBindingDeleteCallType(senderAddress, destinationAddress),
	}

	// we always send it to the remote NodeManagment feature, which always is at entity:[0],feature:0
	localAddress := NodeManagementAddress(senderAddress.Device)
	remoteAddress := NodeManagementAddress(destinationAddress.Device)

	return c.Request(model.CmdClassifierTypeCall, localAddress, remoteAddress, true, []model.CmdType{cmd})
}

func (c *Sender) getMsgCounter() *model.MsgCounterType {
	// TODO:  persistence
	i := model.MsgCounterType(atomic.AddUint64(&c.msgNum, 1))
	return &i
}
