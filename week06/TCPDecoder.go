package week06

import (
	"errors"
	"encoding/binary"
)

/*
goim协议结构如下：
https://github.com/Terry-Mao/goim/blob/master/api/protocol/protocol.go
http://www.tony.wiki/development/2016/09/04/goim-protocol.html

PacketLen	HeaderLen	Version	Operation	Sequence	Body
4bytes	2bytes	2bytes	4bytes	4bytes	PacketLen - HeaderLen
协议字段说明：
1. PacketLen 包长度，在数据流传输过程中，先写入整个包的长度，方便整个包的数据读取。
2. HeaderLen 头长度，在处理数据时，会先解析头部，可以知道具体业务操作。
3. Version 协议版本号，主要用于上行和下行数据包按版本号进行解析。
4. Operation 业务操作码，可以按操作码进行分发数据包到具体业务当中。
5. Sequence 序列号，数据包的唯一标记，可以做具体业务处理，或者数据包去重。
6. Body 实际业务数据，在业务层中会进行数据解码和编码。
*/
const (
	// MaxBodySize max proto body size
	MaxBodySize = uint32(1 << 12)
)

const (
	// size
	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
	_maxPackSize   = MaxBodySize + uint32(_rawHeaderSize)
	// offset
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
	_msgOffset	  = _seqOffset + _seqSize

)

var (
	// ErrProtoPackLen proto packet len error
	ErrProtoPackLen = errors.New("default server codec pack length error")
	// ErrProtoHeaderLen proto header len error
	ErrProtoHeaderLen = errors.New("default server codec header length error")
)

type Proto struct {
	Ver uint16
	Op  uint32
	Seq uint32
	Body []byte
}

func (p *Proto) TCPDecoder(buf []byte) (err error) {

	if len(buf) < _rawHeaderSize {
		return ErrProtoPackLen
	}

	packLen := binary.BigEndian.Uint32(buf[:_packSize])
	headerLen := binary.BigEndian.Uint16(buf[_headerOffset:_verOffset])

	p.Ver = binary.BigEndian.Uint16(buf[_verOffset:_opOffset])
	p.Op = binary.BigEndian.Uint32(buf[_opOffset:_seqOffset])
	p.Seq = binary.BigEndian.Uint32(buf[_seqOffset:_msgOffset])

	if packLen < 0 || packLen > _maxPackSize {
		return ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}
	if bodyLen := int(packLen - uint32(headerLen)); bodyLen > 0 {
		p.Body = buf[headerLen:packLen]
	} else {
		p.Body = nil
	}
	return

}