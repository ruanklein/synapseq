//go:build wasm

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package audio

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gopxl/beep/v2"
	bwav "github.com/gopxl/beep/v2/wav"
)

// BufferWriteSeeker is a bytes.Buffer that supports writing and seeking
type BufferWriteSeeker struct {
	buf *bytes.Buffer
	pos int64
}

// NewBufferWriteSeeker creates a new BufferWriteSeeker
func NewBufferWriteSeeker() *BufferWriteSeeker {
	return &BufferWriteSeeker{
		buf: bytes.NewBuffer(nil),
		pos: 0,
	}
}

// Write writes data to the buffer at the current position
func (b *BufferWriteSeeker) Write(p []byte) (int, error) {
	if b.pos == int64(b.buf.Len()) {
		n, err := b.buf.Write(p)
		b.pos += int64(n)
		return n, err
	}

	total := int(b.pos) + len(p)
	data := b.buf.Bytes()

	if total > len(data) {
		newData := make([]byte, total)
		copy(newData, data)
		b.buf = bytes.NewBuffer(newData)
		data = b.buf.Bytes()
	}

	copy(data[b.pos:], p)
	b.pos += int64(len(p))

	return len(p), nil
}

// Seek sets the position for the next write
func (b *BufferWriteSeeker) Seek(offset int64, whence int) (int64, error) {
	var newPos int64
	switch whence {
	case io.SeekStart:
		newPos = offset
	case io.SeekCurrent:
		newPos = b.pos + offset
	case io.SeekEnd:
		newPos = int64(b.buf.Len()) + offset
	default:
		return 0, fmt.Errorf("invalid whence")
	}

	if newPos < 0 {
		return 0, fmt.Errorf("negative position")
	}

	b.pos = newPos
	return newPos, nil
}

// Bytes returns the contents of the buffer
func (b *BufferWriteSeeker) Bytes() []byte {
	return b.buf.Bytes()
}

// RenderWav renders the audio to a WAV byte slice
func (r *AudioRenderer) RenderWav() ([]byte, error) {
	buf := NewBufferWriteSeeker()

	streamer := newRendererStreamer(r)
	format := beep.Format{
		SampleRate:  beep.SampleRate(r.SampleRate),
		NumChannels: audioChannels,
		Precision:   audioBitDepth / 8,
	}

	if err := bwav.Encode(buf, streamer, format); err != nil {
		return nil, err
	}
	if streamer.err != nil {
		return nil, streamer.err
	}

	return buf.Bytes(), nil
}
