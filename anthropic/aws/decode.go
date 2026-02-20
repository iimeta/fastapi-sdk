package aws

import (
	"bytes"
	"encoding/binary"
	"hash"
	"hash/crc32"
	"io"
)

func DecodeMessage(reader io.Reader, payloadBuf []byte) (m Messages, err error) {

	crc := crc32.New(crc32IEEETable)
	hashReader := io.TeeReader(reader, crc)

	prelude, err := decodePrelude(hashReader, crc)
	if err != nil {
		return Messages{}, err
	}

	if prelude.HeadersLen > 0 {
		lr := io.LimitReader(hashReader, int64(prelude.HeadersLen))
		if m.Headers, err = decodeHeaders(lr); err != nil {
			return Messages{}, err
		}
	}

	if payloadLen := prelude.PayloadLen(); payloadLen > 0 {
		buf, err := decodePayload(payloadBuf, io.LimitReader(hashReader, int64(payloadLen)))
		if err != nil {
			return Messages{}, err
		}
		m.Payload = buf
	}

	if _, err := decodeUint32(reader); err != nil {
		return Messages{}, err
	}

	return m, nil
}

func decodePrelude(r io.Reader, crc hash.Hash32) (messagePrelude, error) {

	var p messagePrelude

	var err error
	p.Length, err = decodeUint32(r)
	if err != nil {
		return messagePrelude{}, err
	}

	p.HeadersLen, err = decodeUint32(r)
	if err != nil {
		return messagePrelude{}, err
	}

	if err := p.ValidateLens(); err != nil {
		return messagePrelude{}, err
	}

	preludeCRC := crc.Sum32()
	if _, err := decodeUint32(r); err != nil {
		return messagePrelude{}, err
	}

	p.PreludeCRC = preludeCRC

	return p, nil
}

func decodePayload(buf []byte, r io.Reader) ([]byte, error) {
	w := bytes.NewBuffer(buf[0:0])

	_, err := io.Copy(w, r)

	return w.Bytes(), err
}

func decodeUint8(r io.Reader) (uint8, error) {

	type byteReader interface {
		ReadByte() (byte, error)
	}

	if br, ok := r.(byteReader); ok {
		v, err := br.ReadByte()
		return v, err
	}

	var b [1]byte

	_, err := io.ReadFull(r, b[:])

	return b[0], err
}

func decodeUint16(r io.Reader) (uint16, error) {

	var b [2]byte
	bs := b[:]

	if _, err := io.ReadFull(r, bs); err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint16(bs), nil
}

func decodeUint32(r io.Reader) (uint32, error) {

	var b [4]byte
	bs := b[:]

	if _, err := io.ReadFull(r, bs); err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(bs), nil
}

func decodeUint64(r io.Reader) (uint64, error) {

	var b [8]byte
	bs := b[:]

	if _, err := io.ReadFull(r, bs); err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(bs), nil
}
