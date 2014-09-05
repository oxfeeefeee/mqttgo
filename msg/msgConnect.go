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


// This file implements type MsgConnect
package msg

import (
    "io"
    "bytes"
    )

type MsgConnect struct {
    H           Header  // Fixed header
    ProtName    Str     // Protocal name
    ProtVer     uint8   // Protocal version number
    flags       byte    // Connect flags
    KeepAlive   uint16  // Keep alive timer
    ClientId    Str     // Client identifier
    WillTopic   Str
    WillMsg     Str
    UserName    Str
    Password    Str
}

// Decode Msg from r, the fixed header and the length is already read
func (m *MsgConnect) ReadFrom(r io.Reader, h Header, length uint32) error {
    m.H = h
    var err error
    lr := &io.LimitedReader{r, int64(length)}
    if m.ProtName, err = ReadStr(lr); err != nil {
        return err
    } else if m.ProtVer, err = ReadUint8(lr); err != nil {
        return err
    } else if m.flags, err = ReadUint8(lr); err != nil {
        return err
    } else if m.KeepAlive, err = ReadUint16(lr); err != nil {
        return err
    } else if m.ClientId, err = ReadStr(lr); err != nil {
        return err
    } else if m.WillTopic, err = ReadStr(lr); err != nil {
        return err
    } else if m.WillMsg, err = ReadStr(lr); err != nil {
        return err
    } else if m.UserName, err = ReadStr(lr); err != nil {
        return err
    } else if m.Password, err = ReadStr(lr); err != nil {
        return err
    }
    if lr.N != 0 {
        return ErrWrongLength
    }
    return nil
}

// Encode Msg
func (m *MsgConnect) WriteTo(w io.Writer) error {
    b := new(bytes.Buffer)
    if err := m.ProtName.WriteTo(b); err != nil {
        return err 
    } else if err := WriteUint8(b, m.ProtVer); err != nil {
        return err
    } else if err := WriteUint8(b, m.flags); err != nil {
        return err
    } else if err := WriteUint16(b, m.KeepAlive); err != nil {
        return err
    } else if err := m.ClientId.WriteTo(b); err != nil {
        return err
    } else if err := m.WillTopic.WriteTo(b); err != nil {
        return err
    } else if err := m.WillMsg.WriteTo(b); err != nil {
        return err
    } else if err := m.UserName.WriteTo(b); err != nil {
        return err
    } else if err := m.Password.WriteTo(b); err != nil {
        return err
    }
    p := b.Bytes()
    if err := WriteUint8(w, byte(m.H)); err != nil {
        return err
    } else if err := Len4(len(p)).WriteTo(w); err != nil {
        return err
    } else if _, err := w.Write(p); err != nil {
        return err
    } 
    return nil
}

// Getter of Clean Session flag
func (m *MsgConnect) CleanSession() bool {
    return get1Bit(m.flags, 0x02) 
}

// Setter of Clean Session flag
func (m *MsgConnect) SetCleanSession(v bool) {
    set1Bit(&m.flags, v, 0x02)
}

// Getter of Will-Flag flag
func (m *MsgConnect) WillFlag() bool {
    return get1Bit(m.flags, 0x04) 
}

// Setter of Will-Flag flag
func (m *MsgConnect) SetWillFlag(v bool) {
    set1Bit(&m.flags, v, 0x04)
}

// Getter of Will-Qos level
func (m *MsgConnect) WillQos() (Qos, error) {
    l := Qos(get2Bits(m.flags, 3))
    if l.Valid() {
        return l, nil
    } else {
        return l, ErrBadQosLevel
    }
}

// Setter of Will-Qos level
func (m *MsgConnect) SetWillQos(l Qos) error {
    if !l.Valid() {
        return ErrBadQosLevel
    }
    set2Bits(&m.flags, byte(l), 3)
    return nil
}

// Getter of Will-Retain flag
func (m *MsgConnect) WillRetain() bool {
    return get1Bit(m.flags, 0x20)
}

// Setter of Will-Retain flag
func (m *MsgConnect) SetWillRetain(d bool) {
    set1Bit(&m.flags, d, 0x20)
}

// Getter of Password flag
func (m *MsgConnect) PasswordFlag() bool {
    return get1Bit(m.flags, 0x40)
}

// Setter of Password flag
func (m *MsgConnect) SetPasswordFlag(d bool) {
    set1Bit(&m.flags, d, 0x40)
}

// Getter of User Name flag
func (m *MsgConnect) UserNameFlag() bool {
    return get1Bit(m.flags, 0x80)
}

// Setter of User Name flag
func (m *MsgConnect) SetUserNameFlag(d bool) {
    set1Bit(&m.flags, d, 0x80)
}