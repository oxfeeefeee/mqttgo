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


// This file implements subscribe related messages
// - MsgSubscribe
// - MsgSubAck
// - MsgUnsubscribe
// - MsgUnsubAck
package mqttgo

import (
    "io"
    "bytes"
    )

type MsgSubscribe struct {
    H       Header
    MsgId   uint16
    Topics  []struct{Topic string; QosLevel}
}

type MsgSubAck struct {
    H           Header
    MsgId       uint16
    GrantedQos  []QosLevel
}

type MsgUnsubscribe struct {
    H       Header
    MsgId   uint16
    Topics  []string
}

type MsgUnsubAck struct {
    msgSimpleAck
}

func (m *MsgSubscribe) MsgHeader() *Header {
    return &(m.H)
}

func (m *MsgSubscribe) Id() uint16 {
    return m.MsgId
}

func (m *MsgSubscribe) SetId(id uint16) {
    m.MsgId = id
}

func (m *MsgSubscribe) readFrom(r io.Reader, h Header, length uint32) error {
    m.H = h
    var err error
    lr := &io.LimitedReader{r, int64(length)}
    if m.MsgId, err = readUint16(lr); err != nil {
        return err
    }
    for lr.N > 0 {
        if topic, err := readStr(lr); err != nil {
            return err
        } else if qos, err := readUint8(lr); err != nil {
            return err
        } else {
            t := struct{Topic string; QosLevel}{topic, QosLevel(qos)}
            m.Topics = append(m.Topics, t)
        } 
    }
    if lr.N != 0 { // Could it be negative?
        return ErrWrongLength
    }
    return nil
}

func (m *MsgSubscribe) writeTo(w io.Writer) error {
    b := new(bytes.Buffer)
    if err := writeUint16(b, m.MsgId); err != nil {
        return err
    }
    for _, t := range m.Topics {
        if err := str(t.Topic).writeTo(b); err != nil {
            return err
        } else if err := writeUint8(b, byte(t.QosLevel)); err != nil {
            return err
        }
    }
    return writeMsgData(w, m.H, b.Bytes())
}

func (m *MsgSubAck) MsgHeader() *Header {
    return &(m.H)
}

func (m *MsgSubAck) readFrom(r io.Reader, h Header, length uint32) error {
    m.H = h
    var err error
    lr := &io.LimitedReader{r, int64(length)}
    if m.MsgId, err = readUint16(lr); err != nil {
        return err
    }
    for lr.N > 0 {
        if qos, err := readUint8(lr); err != nil {
            return err
        } else {
            m.GrantedQos = append(m.GrantedQos, QosLevel(qos))
        } 
    }
    if lr.N != 0 { // Could it be negative?
        return ErrWrongLength
    }
    return nil
}

func (m *MsgSubAck) writeTo(w io.Writer) error {
    b := new(bytes.Buffer)
    if err := writeUint16(b, m.MsgId); err != nil {
        return err
    }
    for _, t := range m.GrantedQos {
        if err := writeUint8(b, byte(t)); err != nil {
            return err
        }
    }
    return writeMsgData(w, m.H, b.Bytes())
}

func (m *MsgUnsubscribe) MsgHeader() *Header {
    return &(m.H)
}

func (m *MsgUnsubscribe) readFrom(r io.Reader, h Header, length uint32) error {
    m.H = h
    var err error
    lr := &io.LimitedReader{r, int64(length)}
    if m.MsgId, err = readUint16(lr); err != nil {
        return err
    }
    for lr.N > 0 {
        if topic, err := readStr(lr); err != nil {
            return err
        } else {
            m.Topics = append(m.Topics, topic)
        } 
    }
    if lr.N != 0 { // Could it be negative?
        return ErrWrongLength
    }
    return nil
}

func (m *MsgUnsubscribe) writeTo(w io.Writer) error {
    b := new(bytes.Buffer)
    if err := writeUint16(b, m.MsgId); err != nil {
        return err
    }
    for _, t := range m.Topics {
        if err := str(t).writeTo(b); err != nil {
            return err
        }
    }
    return writeMsgData(w, m.H, b.Bytes())
}
