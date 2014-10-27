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

package mqttgo

import (
    "io"
    )

// Header of all messages
type Header byte

// Getter of MsgType
func (h Header) Type() MsgType {
    t := MsgType(h >> 4)
    if t.Valid() {
        return t
    } else {
        return MsgTypeInvaild
    }
}

// Setter of MsgType
func (h *Header) SetType(t MsgType) error {
    if !t.Valid() {
        return ErrBadMsgType
    }
    other := byte(*h) & 0x0F
    *h = Header(other | (byte(t) << 4))
    return nil
}

// Getter of Dup flag
func (h Header) Dup() bool {
    return get1Bit(byte(h), 0x08)
}

// Setter of Dup flag
func (h *Header) SetDup(d bool) {
    set1Bit((*byte)(h), d, 0x08)
}

// Getter of Qos level
func (h Header) Qos() (QosLevel, error) {
    l := QosLevel(get2Bits(byte(h), 1))
    if l.Valid() {
        return l, nil
    } else {
        return l, ErrBadQosLevel
    }
}

// Setter of Qos level
func (h *Header) SetQos(l QosLevel) error {
    if !l.Valid() {
        return ErrBadQosLevel
    }
    set2Bits((*byte)(h), byte(l), 1)
    return nil
}

// Getter of Retain flag
func (h Header) Retain() bool {
    return get1Bit(byte(h), 0x01)
}

// Setter of Retain flag
func (h *Header) SetRetain(d bool) {
    set1Bit((*byte)(h), d, 0x01)
}

// Decode header
func (h *Header) readFrom(r io.Reader) error {
    var buf [1]byte
    if _, err := io.ReadFull(r, buf[:]); err != nil {
        return err
    }
    *h = Header(buf[0])
    return nil
}

// Encode header, along with the length
func (h *Header) bytes(length uint32) ([]byte, error) {
    p := make([]byte, 1, 5)
    p[0] = byte(*h)
    if ret, err := len4(length).bytes(); err != nil {
        return nil, err
    } else {
        p = append(p, ret...)
    }
    return p, nil
}

// Validate header and the length of the message
func (h *Header) Validate(length uint32) error {
    t := h.Type()
    if t == MsgTypeInvaild {
        return ErrBadMsgType
    } else if _, err := h.Qos(); err != nil {
        return err
    }
    switch t { // TODO better validation
    case MsgTypePublish:
        if length > PublishMaxLen {
            return ErrTooLong
        }
    default:
        if length > DefaultMaxLen {
            return ErrTooLong
        }
    }
    return nil
}
