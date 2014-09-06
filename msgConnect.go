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


// This file implements type MsgConnect and MsgConnAct
package mqttgo

import (
    "io"
    "bytes"
    )

const (
    RCAccepted ReturnCode = iota
    RCBadVersion
    RCIdRejected
    RCServerUnavailable
    RCBadUserPassword
    )

type ReturnCode uint8 

func (c ReturnCode) Valid() bool {
    return c <= RCBadUserPassword
}

type MsgConnect struct {
    H           Header  // Fixed header
    ProtName    string  // Protocal name
    ProtVer     uint8   // Protocal version number
    flags       byte    // Connect flags
    KeepAlive   uint16  // Keep alive timer
    ClientId    string  // Client identifier
    WillTopic   string
    WillMsg     string
    UserName    string
    Password    string
}

type MsgConnAct struct {
    H   Header
    RC  ReturnCode
}

func (m *MsgConnect) readFrom(r io.Reader, h Header, length uint32) error {
    m.H = h
    var err error
    lr := &io.LimitedReader{r, int64(length)}
    if m.ProtName, err = readStr(lr); err != nil {
        return err
    } else if m.ProtVer, err = readUint8(lr); err != nil {
        return err
    } else if m.flags, err = readUint8(lr); err != nil {
        return err
    } else if m.KeepAlive, err = readUint16(lr); err != nil {
        return err
    } else if m.ClientId, err = readStr(lr); err != nil {
        return err
    } else if m.WillTopic, err = readStr(lr); err != nil {
        return err
    } else if m.WillMsg, err = readStr(lr); err != nil {
        return err
    } else if m.UserName, err = readStr(lr); err != nil {
        return err
    } else if m.Password, err = readStr(lr); err != nil {
        return err
    }
    if lr.N != 0 {
        return ErrWrongLength
    }
    return nil
}

func (m *MsgConnect) writeTo(w io.Writer) error {
    b := new(bytes.Buffer)
    if err := str(m.ProtName).writeTo(b); err != nil {
        return err 
    } else if err := writeUint8(b, m.ProtVer); err != nil {
        return err
    } else if err := writeUint8(b, m.flags); err != nil {
        return err
    } else if err := writeUint16(b, m.KeepAlive); err != nil {
        return err
    } else if err := str(m.ClientId).writeTo(b); err != nil {
        return err
    } else if err := str(m.WillTopic).writeTo(b); err != nil {
        return err
    } else if err := str(m.WillMsg).writeTo(b); err != nil {
        return err
    } else if err := str(m.UserName).writeTo(b); err != nil {
        return err
    } else if err := str(m.Password).writeTo(b); err != nil {
        return err
    }
    return writeMsgData(w, m.H, b.Bytes())
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
func (m *MsgConnect) WillQos() (QosLevel, error) {
    l := QosLevel(get2Bits(m.flags, 3))
    if l.Valid() {
        return l, nil
    } else {
        return l, ErrBadQosLevel
    }
}

// Setter of Will-Qos level
func (m *MsgConnect) SetWillQos(l QosLevel) error {
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

func (m *MsgConnAct) readFrom(r io.Reader, h Header, length uint32) error {
    m.H = h
    lr := &io.LimitedReader{r, int64(length)}
    if _, err := readUint8(lr); err != nil { // Reserved byte
        return err
    } else if rc, err := readUint8(lr); err != nil {
        return err
    } else {
        m.RC = ReturnCode(rc)
    }
    if !m.RC.Valid() {
        return ErrBadRC
    } else if lr.N != 0 {
        return ErrWrongLength
    }
    return nil
}

func (m *MsgConnAct) writeTo(w io.Writer) error {
    return writeMsgData(w, m.H, []byte{0, byte(m.RC)})
}