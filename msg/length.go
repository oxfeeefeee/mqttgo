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

// Decode Len4
func ReadLen4(r io.Reader) (Len4, error) {
    var val Len4
    var shift uint
    for i := 0; i < 4; i++ {
        var buf [1]byte
        if _, err := io.ReadFull(r, buf[:]); err != nil {
            return val, err
        } else {
            b := buf[0]
            val |= Len4(uint32(b & 0x7f) << shift)
            if (b & 0x80) == 0 {
                return val, nil
            }
        }
        shift += 7
    }
    return val, errors.New("ReadLen4: Bad format reading Len4")
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