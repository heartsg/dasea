//protobuf decoding/encoding supports (without 
// statically defined message definition)

package common

import (
	"errors"
	"io"
	"math"
	"time"
	"fmt"
)

var errOverflow = errors.New("protobuf: integer overflow")
var errBadWireType = errors.New("protobuf: bad wiretype")
var errBadTag = errors.New("protobuf: bad tag")

const (
	WireVarint uint64 = iota   	//int32, int64, uint32, uint64, sint32, sint64, bool, enum
	WireFixed64 uint64 = iota	//fixed64, sfixed64, double
	WireLengthDelimited uint64 = iota	//string, bytes, embedded messages, packed repeated fields
	WireStartGroup uint64 = iota
	WireEndGroup uint64 = iota
	WireFixed32 uint64 = iota //fixed32, sfixed32, float
)

//similar to proto.Buffer
type ProtoBuffer struct {
	buf   []byte // encode/decode byte stream
	index int    // write point
}

func NewProtoBuffer(e []byte) *ProtoBuffer {
	return &ProtoBuffer{buf: e}
}

func (p *ProtoBuffer) Reset() {
	p.buf = p.buf[0:0] // for reading/writing
	p.index = 0        // for reading
}

func (p *ProtoBuffer)DecodeKey() (wire uint64, tag uint64, err error) {
	u, err := p.DecodeVarint()
	if err != nil {
		return
	}
	tag = u >> 3
	wire = u - (tag << 3)
	if wire > WireFixed32 {
		err = errBadWireType
	}
	return
}

func (p *ProtoBuffer)DecodeCheckKey(wireCheck uint64, tagCheck uint64)  (err error) {
	wire, tag, err := p.DecodeKey()
	if err != nil {
		return
	}
	if wire != wireCheck {
		err = errBadWireType
		return
	}
	if tag != tagCheck {
		err = errBadTag
		return
	}
	return
}

func (p *ProtoBuffer) DecodeLength() (l uint64, err error) {
	l, err = p.DecodeVarint()
	return
}

func (p *ProtoBuffer) DecodeVarint() (x uint64, err error) {
	i := p.index
	l := len(p.buf)

	for shift := uint(0); shift < 64; shift += 7 {
		if i >= l {
			err = io.ErrUnexpectedEOF
			return
		}
		b := p.buf[i]
		i++
		x |= (uint64(b) & 0x7F) << shift
		if b < 0x80 {
			p.index = i
			return
		}
	}

	// The number is too large to represent in a 64-bit value.
	err = errOverflow
	return
}

func (p *ProtoBuffer) DecodeFixed64() (x uint64, err error) {
	// x, err already 0
	i := p.index + 8
	if i < 0 || i > len(p.buf) {
		err = io.ErrUnexpectedEOF
		return
	}
	p.index = i

	x = uint64(p.buf[i-8])
	x |= uint64(p.buf[i-7]) << 8
	x |= uint64(p.buf[i-6]) << 16
	x |= uint64(p.buf[i-5]) << 24
	x |= uint64(p.buf[i-4]) << 32
	x |= uint64(p.buf[i-3]) << 40
	x |= uint64(p.buf[i-2]) << 48
	x |= uint64(p.buf[i-1]) << 56
	return
}

func (p *ProtoBuffer) DecodeFloat64() (x float64, err error) {
	i, err := p.DecodeFixed64()
	if err != nil {
		return
	}
	x = math.Float64frombits(i)
	return
}

func (p *ProtoBuffer) DecodeFixed32() (x uint64, err error) {
	// x, err already 0
	i := p.index + 4
	if i < 0 || i > len(p.buf) {
		err = io.ErrUnexpectedEOF
		return
	}
	p.index = i

	x = uint64(p.buf[i-4])
	x |= uint64(p.buf[i-3]) << 8
	x |= uint64(p.buf[i-2]) << 16
	x |= uint64(p.buf[i-1]) << 24
	return
}

func (p *ProtoBuffer) DecodeFloat32() (x float32, err error) {
	i, err := p.DecodeFixed32()
	if err != nil {
		return
	}
	x = math.Float32frombits(uint32(i))
	return
}

func (p *ProtoBuffer) DecodeZigzag64() (x uint64, err error) {
	x, err = p.DecodeVarint()
	if err != nil {
		return
	}
	x = (x >> 1) ^ uint64((int64(x&1)<<63)>>63)
	return
}

func (p *ProtoBuffer) DecodeZigzag32() (x uint64, err error) {
	x, err = p.DecodeVarint()
	if err != nil {
		return
	}
	x = uint64((uint32(x) >> 1) ^ uint32((int32(x&1)<<31)>>31))
	return
}

func (p *ProtoBuffer) DecodeRawBytes(alloc bool) (buf []byte, err error) {
	n, err := p.DecodeVarint()
	if err != nil {
		return nil, err
	}

	nb := int(n)
	if nb < 0 {
		return nil, fmt.Errorf("proto: bad byte length %d", nb)
	}
	end := p.index + nb
	if end < p.index || end > len(p.buf) {
		return nil, io.ErrUnexpectedEOF
	}

	if !alloc {
		// todo: check if can get more uses of alloc=false
		buf = p.buf[p.index:end]
		p.index += nb
		return
	}

	buf = make([]byte, nb)
	copy(buf, p.buf[p.index:])
	p.index += nb
	return
}

func (p *ProtoBuffer) DecodeStringBytes() (s string, err error) {
	buf, err := p.DecodeRawBytes(false)
	if err != nil {
		return
	}
	return string(buf), nil
}


func (p *ProtoBuffer) DecodeTimestamp() (t time.Time, err error) {
	buf, err := p.DecodeRawBytes(false)
	if err != nil {
		return
	}
	if len(buf) == 0 {
		t = time.Unix(0, 0)
		return
	}
	protoBuf := NewProtoBuffer(buf)
	
	var seconds uint64
	var nanos uint64
	
	wire, tag, err := protoBuf.DecodeKey()
	if err != nil {
		return
	}
	if wire != WireVarint {
		err = errBadWireType
		return
	}
	if tag != 1 && tag != 2 {
		err = errBadTag
		return
	}
	
	if tag == 1 {
		seconds, err = protoBuf.DecodeVarint()
	}
	
	if tag == 2 {
		nanos, err = protoBuf.DecodeVarint()
	}
	if err != nil {
		return
	}

	if protoBuf.DecodeComplete() {
		t = time.Unix(int64(seconds), int64(nanos))
		return
	}
	if tag == 2 {
		err = errOverflow
		return
	}
	
	err = protoBuf.DecodeCheckKey(WireVarint, 2)
	if err != nil {
		return
	}
	
	nanos, err = protoBuf.DecodeVarint()
	if err != nil {
		return
	}
	
	t = time.Unix(int64(seconds), int64(nanos))
	return
}


func (p *ProtoBuffer) DecodeComplete() bool {
	return p.index >= len(p.buf)
}