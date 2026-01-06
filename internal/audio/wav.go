//go:build !wasm

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
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
	"github.com/ruanklein/synapseq/v3/internal/info"
)

// RenderWav renders the audio to a WAV file using go-audio/wav
func (r *AudioRenderer) RenderWav(outPath string) error {
	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	streamer := newRendererStreamer(r)
	format := beep.Format{
		SampleRate:  beep.SampleRate(r.SampleRate),
		NumChannels: audioChannels,
		Precision:   audioBitDepth / 8,
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
func WriteICMTChunkFromTextFile(wavPath string, metadata *info.Metadata) error {
	if metadata == nil {
		return fmt.Errorf("metadata is nil")
	}
	header := bytes.Buffer{}

	header.WriteString("SYNAPSEQ_META::ID=" + metadata.ID() + "\n")
	header.WriteString("VERSION=" + metadata.Version() + "\n")
	header.WriteString("GENERATED=" + metadata.Generated() + "\n")
	header.WriteString("PLATFORM=" + metadata.Platform() + "\n")
	header.WriteString("CONTENT=\n")
	header.WriteString(metadata.Content() + "\n")

	commentBytes := header.Bytes()
	paddedLen := (len(commentBytes) + 1) &^ 1 // padding to even
	icmtSize := uint32(paddedLen)
	totalSize := uint32(4 + 4 + 4 + paddedLen) // "INFO" + "ICMT" + size + data

	// Open the WAV file to append the chunk
	f, err := os.OpenFile(wavPath, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening wav: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	buf.WriteString("LIST")
	binary.Write(&buf, binary.LittleEndian, totalSize)
	buf.WriteString("INFO")
	buf.WriteString("ICMT")
	binary.Write(&buf, binary.LittleEndian, icmtSize)
	buf.Write(commentBytes)
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
func ExtractTextSequenceFromWAV(wavPath string) (string, error) {
	f, err := os.Open(wavPath)
	if err != nil {
		return "", fmt.Errorf("error opening WAV: %w", err)
	}
	defer f.Close()

	const chunkHeaderSize = 8

	buf := make([]byte, chunkHeaderSize)
	for {
		_, err := io.ReadFull(f, buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading chunk: %w", err)
		}

		chunkID := string(buf[:4])
		chunkSize := binary.LittleEndian.Uint32(buf[4:8])

		if chunkID == "LIST" {
			listData := make([]byte, chunkSize)
			_, err := io.ReadFull(f, listData)
			if err != nil {
				return "", fmt.Errorf("error reading LIST chunk: %w", err)
			}

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
						return "", errors.New("ICMT subchunk size exceeds LIST chunk size")
					}

					data := listData[offset : offset+int(subchunkSize)]
					data = bytes.TrimRight(data, "\x00") // remove padding null bytes

					// Parse lines
					lines := bytes.Split(data, []byte("\n"))
					readContent := false

					var (
						id, generated, version, platform string
						base64Content                    []byte
					)

					for _, line := range lines {
						if readContent {
							base64Content = append(base64Content, line...)
							base64Content = append(base64Content, '\n') // preserve line breaks
						} else if bytes.HasPrefix(line, []byte("SYNAPSEQ_META::ID=")) {
							id = string(bytes.TrimPrefix(line, []byte("SYNAPSEQ_META::ID=")))
						} else if bytes.HasPrefix(line, []byte("GENERATED=")) {
							generated = string(bytes.TrimPrefix(line, []byte("GENERATED=")))
						} else if bytes.HasPrefix(line, []byte("VERSION=")) {
							version = string(bytes.TrimPrefix(line, []byte("VERSION=")))
						} else if bytes.HasPrefix(line, []byte("PLATFORM=")) {
							platform = string(bytes.TrimPrefix(line, []byte("PLATFORM=")))
						} else if bytes.HasPrefix(line, []byte("CONTENT=")) {
							readContent = true
						}
					}

					decoded, err := base64.StdEncoding.DecodeString(string(base64Content))
					if err != nil {
						return "", fmt.Errorf("failed to decode base64 content: %w", err)
					}

					content := "# ================================================\n"
					content += "#  This sequence was exported from SynapSeq\n"
					content += fmt.Sprintf("#  ID       : %s\n", id)
					content += fmt.Sprintf("#  Date     : %s\n", generated)
					content += fmt.Sprintf("#  Version  : %s\n", version)
					content += fmt.Sprintf("#  Platform : %s\n", platform)
					content += "# ================================================\n\n\n"
					content += string(decoded)

					return content, nil
				}

				offset += int((subchunkSize + 1) &^ 1)
			}
		} else {
			if _, err := f.Seek(int64((chunkSize+1)&^1), io.SeekCurrent); err != nil {
				return "", fmt.Errorf("error skipping chunk: %w", err)
			}
		}
	}

	return "", errors.New("no ICMT metadata found in WAV")
}
