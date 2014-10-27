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


// This file implements msgSimpleAck, i.e. messages contain only a header
// and a MsgId for for message types:
// - MsgPubAck
// - MsgPubRec
// - MsgPubRel
// - MsgPubComp
// - MsgUnsubAck
// and msgHeaderOnly for message types:
// - MsgPingReq
// - MsgPingResp
// - Disconnect
package mqttgo

import (
    "io"
    )

type msgSimpleAck struct {
    H       Header
    MsgId   uint16
}

func (m *msgSimpleAck) MsgHeader() *Header {
    return &(m.H)
}

func (m *msgSimpleAck) readFrom(r io.Reader, h Header, length uint32) error {
    m.H = h
    var err error
    lr := &io.LimitedReader{r, int64(length)}
    if m.MsgId, err = readUint16(lr); err != nil {
        return err
    } else if lr.N != 0 {
        return ErrWrongLength
    }
    return nil
}

func (m *msgSimpleAck) writeTo(w io.Writer) error {
    return writeMsgData(w, m.H, []byte{byte(m.MsgId >> 8), byte(m.MsgId & 0x00ff)})
}

type msgHeaderOnly struct {
    H   Header
}

func (m *msgHeaderOnly) MsgHeader() *Header {
    return &(m.H)
}

func (m *msgHeaderOnly) readFrom(r io.Reader, h Header, length uint32) error {
    m.H = h
    if length != 0 {
        return ErrWrongLength
    }
    return nil
}

func (m *msgHeaderOnly) writeTo(w io.Writer) error {
    return writeMsgData(w, m.H, nil)
}
