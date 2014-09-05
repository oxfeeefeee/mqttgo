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


// This file contains commonly used types in msg package
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
package msg

import (
    "io"
    "errors"
    )

var (
    ErrBadMsgType = errors.New("mqttgo/msg: Bad message type")
    ErrBadQosLevel = errors.New("mqttgo/msg: Bad Qos level")
    ErrTooLong = errors.New("mqttgo/msg: Message too long")
    ErrBadRC = errors.New("mqttgo/msg: Bad return code")
    ErrWrongLength = errors.New("mqttgo/msg: Message length doesn't match with content")
    )

// Max length of the content of MsgPublish
const PublishMaxLen = 1024 * 1024

// Max length of the content of other Messages
const DefaultMaxLen = 1024 * 10

const (
    MsgTypeConnect MsgType = iota
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
)

const (
    QosAtMostOnce Qos = iota
    QosAtLeastOnce
    QosExactlyOnce
    )

const (
    RCAccepted RC = iota
    RCBadVersion
    RCIdRejected
    RCServerUnavailable
    RCBadUserPassword
    )

// All MQTT messages implement this interface
type Msg interface {
    // Decode Msg from r, the fixed header and the length is already read
    ReadFrom(r io.Reader, h Header, length uint32) error
    // Encode Msg
    WriteTo(w io.Writer) error
}

type MsgType uint8

func (t MsgType) Valid() bool {
    return t <= MsgTypeDisconnect
}

// Qaulity of service
type Qos uint8

func (q Qos) Valid() bool {
    return q <= QosExactlyOnce
}

// Return code
type RC uint8 

func (c RC) Valid() bool {
    return c <= RCBadUserPassword
}

// Read the length of the rest of the message
func readMsgLen(r io.Reader) (uint32, error) {
    if l, err := ReadLen4(r); err != nil {
        return 0, err
    } else {
        return l, nil
    }
}

// Gets one bit in a byte and returns as bool
func get1Bit(f byte, mask byte) bool {
    return (byte(f) & mask) != 0
}

// Sets one bit in a byte
func set1Bit(f *byte, v bool, mask byte) {
    *f = *f &^ 0x08
    if v {
        *f = *f | 0x08  
    }
}

// Gets two bits in a byte and returns as byte
func get2Bits(f byte, from uint) byte {
    return (f << (6 - from)) >> 6
}

// Sets two bits in a byte
func set2Bits(f *byte, val byte, from uint) {
    mask := ^((byte(0xff) << (6 - from)) >> 6)
    *f = (byte(*f) & mask) | (val << from)
}


