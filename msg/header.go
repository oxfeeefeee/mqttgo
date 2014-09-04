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


//Message header format:
//-------------------------------------------------------------------
// bit    |  7  |  6  |  5  |  4  |     3     |  2  |  1  |    0    |
//-------------------------------------------------------------------
// byte1  |      Message Type     | DUP flag  | QoS level |RETAIN   |
//-------------------------------------------------------------------
// byte2+ |                  Remaining Length                       |
//-------------------------------------------------------------------

package msg

import (
    "io"
    "errors"
    )

var (
    ErrBadMsgType = errors.New("mqttgo/msg: Bad message type")
    ErrBadQosLevel = errors.New("mqttgo/msg: Bad Qos level")
    ErrBadRC = errors.New("mqttgo/msg: Bad return code")
    )

const (
    MsgConnect MsgType = iota
    MsgConnAck
    MsgPublish
    MsgPubAck
    MsgPubRec
    MsgPubRel
    MsgPubComp
    MsgSubscribe
    MsgSubAck
    MsgUnsubscribe
    MsgUnsubAck
    MsgPingReq
    MsgPingResp
    MsgDisconnect
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
    RCBadUserPasswd
    )

type MsgType uint8

func (t MsgType) Valid() bool {
    return t <= MsgDisconnect
}

// Qaulity of service
type Qos uint8

func (q Qos) Valid() bool {
    return q <= QosExactlyOnce
}

// Return code
type RC uint8 

func (c RC) Valid() bool {
    return c <= RCBadUserPasswd
}

// Header of all messages
type Header struct {
    byte0   byte
    length  Len4
}

// Getter of MsgType
func (h *Header) Type() (MsgType, error) {
    t := MsgType(h.byte0 >> 4)
    if t.Valid() {
        return t, nil
    } else {
        return t, ErrBadMsgType
    }
}

// Setter of MsgType
func (h *Header) SetType(t MsgType) error {
    if !t.Valid() {
        return ErrBadMsgType
    }
    other := h.byte0 & 0x0F
    h.byte0 = other | (byte(t) << 4)
    return nil
}

// Getter of Dup flag
func (h *Header) Dup() bool {
    return (h.byte0 & 0x08) != 0
}

// Setter of Dup flag
func (h *Header) SetDup(d bool) {
    h.byte0 = h.byte0 &^ 0x08
    if d {
        h.byte0 = h.byte0 | 0x08  
    }
}

// Getter of Qos level
func (h *Header) Qos() (Qos, error) {
    l := Qos(h.byte0 >> 1)
    if l.Valid() {
        return l, nil
    } else {
        return l, ErrBadQosLevel
    }
}

// Setter of Qos level
func (h *Header) SetQos(l Qos) error {
    if !l.Valid() {
        return ErrBadQosLevel
    }
    other := h.byte0 &^ 0x06
    h.byte0 = other | (byte(l) << 1)
    return nil
}

// Getter of Retain flag
func (h *Header) Retain() bool {
    return (h.byte0 & 0x01) != 0
}

// Setter of Retain flag
func (h *Header) SetRetain(d bool) {
    h.byte0 = h.byte0 &^ 0x01
    if d {
        h.byte0 = h.byte0 | 0x01
    }
}

// Getter of Remaining Length
func (h *Header) Len() uint32 {
    return uint32(h.length)
}

// Setter of Remaining Length
func (h *Header) SetLen(l uint32) {
    h.length = Len4(l) 
}

// Decode length from byte slice
func (h *Header) ReadLen(r io.Reader) error {
    if l, err := ReadLen4(r); err != nil {
        return err
    } else {
        h.length = l
        return nil
    }
}

// Decode header
func (h *Header) ReadFrom(r io.Reader) error {
    var buf [1]byte
    if _, err := io.ReadFull(r, buf[:]); err != nil {
        return err
    }
    h.byte0 = buf[0]
    if err := h.ReadLen(r); err != nil {
        return err
    }
    return h.Validate()
}

// Encode header
func (h *Header) Bytes() ([]byte, error) {
    p := make([]byte, 1, 5)
    p[0] = h.byte0
    if ret, err := h.length.Bytes(); err != nil {
        return nil, err
    } else {
        p = append(p, ret...)
    }
    return p, nil
}

// Validate header
func (h *Header) Validate() error {
    if _, err := h.Type(); err != nil {
        return err
    } else if _, err = h.Qos(); err != nil {
        return err
    } else if err = h.length.Validate(); err != nil {
        return err
    }
    return nil
}