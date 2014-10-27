// Copyright 2014 mqttgo author
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


// This file contains commonly used types in mqttgo package
//
// Message format:
// |------------------------------------------------------------------|
// | bit    |  7  |  6  |  5  |  4  |     3     |  2  |  1  |    0    |
// |------------------------------------------------------------------|
// | byte1  |      Message Type     | DUP flag  | QoS level |RETAIN   |
// |------------------------------------------------------------------|
// | byte2+ |                  Remaining Length                       |
// |------------------------------------------------------------------|
// |                              Payload                             |
// |------------------------------------------------------------------|
package mqttgo

import (
    "io"
    "log"
    "errors"
    )

// Max length of the content of MsgPublish
const PublishMaxLen = 1024 * 1024

// Max length of the content of other Messages
const DefaultMaxLen = 1024 * 10

const (
    _               MsgType = iota
    MsgTypeConnect
    MsgTypeConnAck
    MsgTypePublish
    MsgTypePubAck
    MsgTypePubRec
    MsgTypePubRel
    MsgTypePubComp
    MsgTypeSubscribe
    MsgTypeSubAck
    MsgTypeUnsubscribe
    MsgTypeUnsubAck
    MsgTypePingReq
    MsgTypePingResp
    MsgTypeDisconnect
    MsgTypeInvaild
)

const (
    QosAtMostOnce QosLevel = iota
    QosAtLeastOnce
    QosExactlyOnce
    )

// Errors could happen when reading Msg 
var (
    ErrBadMsgType = errors.New("mqttgo/msg: Bad message type")
    ErrBadQosLevel = errors.New("mqttgo/msg: Bad Qos level")
    ErrTooLong = errors.New("mqttgo/msg: Message too long")
    ErrBadRC = errors.New("mqttgo/msg: Bad return code")
    ErrWrongLength = errors.New("mqttgo/msg: Message length doesn't match with content")
    )

// A registry for creating Msg objects
var msgRegistry map[MsgType]func() Msg = map[MsgType]func() Msg {
    MsgTypeConnect:     func() Msg { return new(MsgConnect) },
    MsgTypeConnAck:     func() Msg { return new(MsgConnAct) },
    MsgTypePublish:     func() Msg { return new(MsgPublish) },
    MsgTypePubAck:      func() Msg { return new(MsgPubAck) },
    MsgTypePubRec:      func() Msg { return new(MsgPubRec) },
    MsgTypePubRel:      func() Msg { return new(MsgPubRel) },
    MsgTypePubComp:     func() Msg { return new(MsgPubComp) },
    MsgTypeSubscribe:   func() Msg { return new(MsgSubscribe) },
    MsgTypeSubAck:      func() Msg { return new(MsgSubAck) },
    MsgTypeUnsubscribe: func() Msg { return new(MsgUnsubscribe) },
    MsgTypeUnsubAck:    func() Msg { return new(MsgUnsubAck) },
    MsgTypePingReq:     func() Msg { return new(MsgPingReq) },
    MsgTypePingResp:    func() Msg { return new(MsgPingResp) },
    MsgTypeDisconnect:  func() Msg { return new(MsgDisconnect) },
}

// All MQTT messages implement this interface
type Msg interface {
    // Returns the type of Msg
    MsgHeader() *Header
    // Decode Msg from r, the fixed header and the length is already read
    readFrom(r io.Reader, h Header, length uint32) error
    // Encode Msg
    writeTo(w io.Writer) error
}

type MsgType uint8

func (t MsgType) Valid() bool {
    return t <= MsgTypeDisconnect
}

// Qaulity of service
type QosLevel uint8

func (q QosLevel) Valid() bool {
    return q <= QosExactlyOnce
}

// Read a Msg from an io.Reader
func Read(r io.Reader) (Msg, error) {
    var h Header
    if err := h.readFrom(r); err != nil {
        return nil, err
    } else if l, err := readMsgLen(r); err != nil {
        return nil, err
    } else if err := h.Validate(l); err != nil {
        return nil, err
    } else {
        if t := h.Type(); t == MsgTypeInvaild {
            return nil, ErrBadMsgType
        } else {
            msg := msgRegistry[t]()
            if err := msg.readFrom(r, h, l); err != nil {
                return nil, err
            }
            log.Printf("READ message type: %d", t)
            return msg, nil
        }
    }
}

// Write a Msg to io.Writer
func Write(w io.Writer, m Msg) error {
    return m.writeTo(w)
}

// Write MsgTypeConnAck
func WriteConnAck(w io.Writer, rc ReturnCode) error {
    var m MsgConnAct
    m.H.SetType(MsgTypeConnAck)
    m.RC = rc
    return m.writeTo(w)
}

type MsgPingReq struct {
    msgHeaderOnly
}

type MsgPingResp struct {
    msgHeaderOnly
}

type MsgDisconnect struct {
    msgHeaderOnly
}
