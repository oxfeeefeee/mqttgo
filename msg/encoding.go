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


// This file implements type Len4 and Str,
// which are for encoding/decoding length value and string value.
// Also implements Read/Write uint8/uint16
package msg

import (
    "io"
    "bytes"
    "errors"
)

// Taken from the standard:
// The variable length encoding scheme uses a single byte for messages up to 127 bytes long. 
// Longer messages are handled as follows. Seven bits of each byte encode the remaining 
// Length data, and the eighth bit indicates any following bytes in the representation.
type Len4 uint32

// Read from io.Reader and decode as Len4
func ReadLen4(r io.Reader) (uint32, error) {
    var val uint32
    var shift uint
    for i := 0; i < 4; i++ {
        var buf [1]byte
        if _, err := io.ReadFull(r, buf[:]); err != nil {
            return val, err
        } else {
            b := buf[0]
            val |= (uint32(b & 0x7f) << shift)
            if (b & 0x80) == 0 {
                return val, nil
            }
        }
        shift += 7
    }
    return val, errors.New("ReadLen4: Bad format reading Len4")
}

// Write the encoded Len4 to io.Writer
func (l Len4) WriteTo(w io.Writer) error {
    if p, err := l.Bytes(); err != nil {
        return err
    } else {
        _, err := w.Write(p)
        return err
    }
}

// Encode Len4
func (l Len4) Bytes() ([]byte, error) {
    if l == 0 {
        return make([]byte, 1), nil
    }
    var p bytes.Buffer
    for i := 0; i < 4; i++ {
        digit := l & 0x7f
        if l >>= 7; l > 0 {
            digit = digit | 0x80
            p.WriteByte(byte(digit))
        } else {
            return p.Bytes(), nil
        }
    }
    return nil, errors.New("Len4.Bytes: Len4 value out of range")
}

// Validate Len4
func (l Len4) Validate() error {
    if l < 0 || l > 268435455 {
        errors.New("Len4.Validate: Len4 value out of range")
    }
    return nil
}

type Str string

// Read from io.Reader and decode as Str
func ReadStr(r io.Reader) (Str, error) {
    if l, err := ReadLen4(r); err != nil {
        return "", err
    } else {
        p := make([]byte, l)
        if _, err := io.ReadFull(r, p); err != nil {
            return "", err
        } else {
            return Str(p), nil
        }
    }
}

// Write the encoded Str to io.Writer
func (s Str) WriteTo(w io.Writer) error {
    l := Len4(len(s))
    if err := l.WriteTo(w); err != nil {
        return err
    }
    if _, err := io.WriteString(w, string(s)); err != nil {
        return err
    }
    return nil
}

// Read uint8 from io.Reader
func ReadUint8(r io.Reader) (uint8, error) {
    var buf [1]byte
    if _, err := io.ReadFull(r, buf[:]); err != nil {
        return 0, err
    }
    return uint8(buf[0]), nil
}

// Write uint8 to io.Writer
func WriteUint8(w io.Writer, val uint8) error {
    buf := [1]byte{val,}
    _, err := w.Write(buf[:])
    return err
}

// Read uint16 from io.Reader
func ReadUint16(r io.Reader) (uint16, error) {
    var buf [2]byte
    if _, err := io.ReadFull(r, buf[:]); err != nil {
        return 0, err
    }
    return (uint16(buf[0]) << 8) | uint16(buf[1]), nil
}

// Write uint16 to io.Write
func WriteUint16(w io.Writer, val uint16) error {
    buf := [2]byte{byte(val >> 8), byte(val & 0x00ff)}
    _, err := w.Write(buf[:])
    return err
}