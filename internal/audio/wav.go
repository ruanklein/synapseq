/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) 2025 Ruan <https://ruan.sh/>
 * Licensed under GNU GPL v2. See COPYING.txt for details.
 */

package audio

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/gopxl/beep/v2"
	bwav "github.com/gopxl/beep/v2/wav"
)

// RenderWav renders the audio to a WAV file using go-audio/wav
func (r *AudioRenderer) RenderWav(outPath string) error {
	// If in debug mode, render to stdout (nil writer)
	if r.Debug {
		return r.Render(nil)
	}

	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	streamer := newRendererStreamer(r)
	format := beep.Format{
		SampleRate:  beep.SampleRate(r.SampleRate),
		NumChannels: audioChannels,
		Precision:   audioBitDepth / 8, // 24-bit -> 3 bytes
	}

	if err := bwav.Encode(out, streamer, format); err != nil {
		return err
	}
	if streamer.err != nil {
		return streamer.err
	}

	return nil
}

// WriteICMTChunkFromTextFile appends an ICMT chunk with base64-encoded content from the specified text file
func WriteICMTChunkFromTextFile(wavPath, filePath string) error {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	comment := base64.StdEncoding.EncodeToString(raw)

	f, err := os.OpenFile(wavPath, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening wav: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer

	// "LIST" header
	buf.WriteString("LIST")

	commentBytes := []byte(comment)
	paddedLen := (len(commentBytes) + 1) &^ 1 // padding para par
	icmtSize := uint32(paddedLen)

	totalSize := uint32(4 + 4 + 4 + paddedLen) // "INFO" + "ICMT" + size + data

	binary.Write(&buf, binary.LittleEndian, totalSize) // LIST chunk size
	buf.WriteString("INFO")                            // LIST type
	buf.WriteString("ICMT")                            // ICMT field
	binary.Write(&buf, binary.LittleEndian, icmtSize)  // ICMT data size
	buf.Write(commentBytes)                            // ICMT content
	if len(commentBytes)%2 != 0 {
		buf.WriteByte(0) // padding
	}

	_, err = f.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("error writing metadata: %w", err)
	}

	return nil
}

// ExtractTextSequenceFromWAV extracts the text sequence metadata from the ICMT chunk of a WAV file
func ExtractTextSequenceFromWAV(wavPath, outTextPath string) error {
	f, err := os.Open(wavPath)
	if err != nil {
		return fmt.Errorf("error opening WAV: %w", err)
	}
	defer f.Close()

	const (
		chunkHeaderSize = 8 // 4 bytes ID + 4 bytes size
	)

	buf := make([]byte, chunkHeaderSize)
	for {
		_, err := io.ReadFull(f, buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading chunk: %w", err)
		}

		chunkID := string(buf[:4])
		chunkSize := binary.LittleEndian.Uint32(buf[4:8])

		if chunkID == "LIST" {
			listData := make([]byte, chunkSize)
			_, err := io.ReadFull(f, listData)
			if err != nil {
				return fmt.Errorf("error reading LIST chunk: %w", err)
			}

			// Check if it's of type "INFO"
			if !bytes.HasPrefix(listData, []byte("INFO")) {
				continue
			}

			offset := 4
			for offset+8 <= len(listData) {
				subchunkID := string(listData[offset : offset+4])
				subchunkSize := binary.LittleEndian.Uint32(listData[offset+4 : offset+8])
				offset += 8

				if subchunkID == "ICMT" {
					if offset+int(subchunkSize) > len(listData) {
						return errors.New("ICMT subchunk size exceeds LIST chunk size")
					}

					data := listData[offset : offset+int(subchunkSize)]
					data = bytes.TrimRight(data, "\x00") // remove padding null bytes
					decoded, err := base64.StdEncoding.DecodeString(string(data))
					if err != nil {
						return fmt.Errorf("failed to decode base64: %w", err)
					}

					err = os.WriteFile(outTextPath, decoded, 0644)
					if err != nil {
						return fmt.Errorf("error saving text sequence: %w", err)
					}

					return nil // success
				}

				offset += int((subchunkSize + 1) &^ 1) // align to even
			}
		} else {
			// Skip to the next chunk
			if _, err := f.Seek(int64((chunkSize+1)&^1), io.SeekCurrent); err != nil {
				return fmt.Errorf("error skipping chunk: %w", err)
			}
		}
	}

	return errors.New("no ICMT metadata found in WAV")
}
