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


// This file implements type len4 and str,
// which are for encoding/decoding length value and string value.
// Also implements Read/Write uint8/uint16
package mqttgo

import (
    "io"
    "bytes"
    "errors"
)

// Taken from the standard:
// The variable length encoding scheme uses a single byte for messages up to 127 bytes long. 
// Longer messages are handled as follows. Seven bits of each byte encode the remaining 
// Length data, and the eighth bit indicates any following bytes in the representation.
type len4 uint32

// Read from io.Reader and decode as len4
func readLen4(r io.Reader) (uint32, error) {
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
    return val, errors.New("readLen4: Bad format reading len4")
}

// Write the encoded len4 to io.Writer
func (l len4) writeTo(w io.Writer) error {
    if p, err := l.bytes(); err != nil {
        return err
    } else {
        _, err := w.Write(p)
        return err
    }
}

// Encode len4
func (l len4) bytes() ([]byte, error) {
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
    return nil, errors.New("len4.bytes: len4 value out of range")
}

type str string

// Read from io.Reader and decode as str
func readStr(r io.Reader) (string, error) {
    if l, err := readLen4(r); err != nil {
        return "", err
    } else {
        p := make([]byte, l)
        if _, err := io.ReadFull(r, p); err != nil {
            return "", err
        } else {
            return string(p), nil
        }
    }
}

// Write the encoded str to io.Writer
func (s str) writeTo(w io.Writer) error {
    l := len4(len(s))
    if err := l.writeTo(w); err != nil {
        return err
    }
    if _, err := io.WriteString(w, string(s)); err != nil {
        return err
    }
    return nil
}

// Read uint8 from io.Reader
func readUint8(r io.Reader) (uint8, error) {
    var buf [1]byte
    if _, err := io.ReadFull(r, buf[:]); err != nil {
        return 0, err
    }
    return uint8(buf[0]), nil
}

// Write uint8 to io.Writer
func writeUint8(w io.Writer, val uint8) error {
    buf := [1]byte{val,}
    _, err := w.Write(buf[:])
    return err
}

// Read uint16 from io.Reader
func readUint16(r io.Reader) (uint16, error) {
    var buf [2]byte
    if _, err := io.ReadFull(r, buf[:]); err != nil {
        return 0, err
    }
    return (uint16(buf[0]) << 8) | uint16(buf[1]), nil
}

// Write uint16 to io.Write
func writeUint16(w io.Writer, val uint16) error {
    buf := [2]byte{byte(val >> 8), byte(val & 0x00ff)}
    _, err := w.Write(buf[:])
    return err
}

// Read the length of the rest of the message
func readMsgLen(r io.Reader) (uint32, error) {
    if l, err := readLen4(r); err != nil {
        return 0, err
    } else {
        return l, nil
    }
}

// Helper function for writing msg
func writeMsgData(w io.Writer, h Header, p []byte) error {
    if err := writeUint8(w, byte(h)); err != nil {
        return err
    } else if err := len4(len(p)).writeTo(w); err != nil {
        return err
    } else if _, err := w.Write(p); err != nil {
        return err
    } 
    return nil
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