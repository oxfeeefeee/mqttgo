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


// This file implements publish related messages
// - MsgPublish
// - MsgPubAck
// - MsgPubRec
// - MsgPubRel
// - MsgPubComp
package mqttgo

import (
    "io"
    "bytes"
    )

type MsgPublish struct {
    H       Header
    Topic   string
    MsgId   uint16
    Content []byte
}

type MsgPubAck struct {
    msgSimpleAck
}

type MsgPubRec struct {
    msgSimpleAck
}

type MsgPubRel struct {
    msgSimpleAck
}

type MsgPubComp struct {
    msgSimpleAck
}

func (m *MsgPublish) MsgHeader() *Header {
    return &(m.H)
}

func (m *MsgPublish) Id() uint16 {
    return m.MsgId
}

func (m *MsgPublish) SetId(id uint16) {
    m.MsgId = id
}

func (m *MsgPublish) readFrom(r io.Reader, h Header, length uint32) error {
    m.H = h
    var err error
    lr := &io.LimitedReader{r, int64(length)}
    if m.Topic, err = readStr(lr); err != nil {
        return err
    }
    if qos, err := h.Qos(); err != nil {
        return err
    } else if qos >= QosAtLeastOnce {
        if m.MsgId, err = readUint16(lr); err != nil {
            return err
        }
    }
    m.Content = make([]byte, lr.N)
    for lr.N > 0 {
        has := len(m.Content) - int(lr.N)
        if _, err := lr.Read(m.Content[has:]); err != nil {
            return err
        }
    }
    return nil
}

func (m *MsgPublish) writeTo(w io.Writer) error {
    b := new(bytes.Buffer)
    if err := str(m.Topic).writeTo(b); err != nil {
        return err 
    }
    if qos, err := m.H.Qos(); err != nil {
        return err
    } else if qos >= QosAtLeastOnce {
        if err := writeUint16(b, m.MsgId); err != nil {
            return err
        }
    }
    p := b.Bytes()
    // Do not use writeMsgData becasue we don't want to merge two silces beforehand
    if err := writeUint8(w, byte(m.H)); err != nil {
        return err
    } else if err := len4(len(p) + len(m.Content)).writeTo(w); err != nil {
        return err
    } else if _, err := w.Write(p); err != nil {
        return err
    } else if _, err := w.Write(m.Content); err != nil {
        return err
    } 
    return nil
}
