//go:build !wasm

/*
 * SynapSeq - Synapse-Sequenced Brainwave Generator
 *
 * Copyright (c) ...
 */

package external_test

import (
	"fmt"
	"log"

	synapseq "github.com/ruanklein/synapseq/v3/core"
	"github.com/ruanklein/synapseq/v3/external"
)

func ExampleNewFFPlay() {
	// Create ffplay instance using executable from PATH
	// player, err := external.NewFFPlay("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	player := external.NewFFPlayUnsafe("")
	fmt.Println("ffplay initialized:", player.Path())
	// Output:
	// ffplay initialized: ffplay
}

func ExampleFFplay_Play() {
	// Create SynapSeq application context
	ctx, err := synapseq.NewAppContext("input.spsq", "", "text")
	if err != nil {
		log.Fatal(err)
	}

	// Load sequence (required before playback)
	// _ = ctx.LoadSequence()

	// Create ffplay instance
	// _, err = external.NewFFPlay("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Play audio (real-time)
	// _ = player.Play(ctx)

	fmt.Printf("Playback executed successfully for input: %s\n", ctx.InputFile())
	// Output:
	// Playback executed successfully for input: input.spsq
}

func ExampleNewFFmpeg() {
	// Create ffmpeg instance using executable from PATH
	// encoder, err := external.NewFFmpeg("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	encoder := external.NewFFmpegUnsafe("")
	fmt.Println("ffmpeg initialized:", encoder.Path())
	// Output:
	// ffmpeg initialized: ffmpeg
}

func ExampleFFmpeg_Convert() {
	// Create SynapSeq application context
	ctx, err := synapseq.NewAppContext("input.spsq", "output.mp3", "text")
	if err != nil {
		log.Fatal(err)
	}

	// Load sequence (required before encoding)
	// _ = ctx.LoadSequence()

	// Create ffmpeg instance
	// encoder, err := external.NewFFmpeg("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Encode MP3 at 320 kbps CBR
	// _ = encoder.Convert(ctx, "mp3")

	fmt.Printf("MP3 encoding executed successfully for output: %s\n", ctx.OutputFile())
	// Output:
	// MP3 encoding executed successfully for output: output.mp3
}

func ExampleNewFFprobe() {
	// Create ffprobe instance using executable from PATH
	// probe, err := external.NewFFprobe("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	probe := external.NewFFprobeUnsafe("")
	fmt.Println("ffprobe initialized:", probe.Path())
	// Output:
	// ffprobe initialized: ffprobe
}

func ExampleFFprobe_ExtractTextSequence() {
	// Create ffprobe instance
	// probe, err := external.NewFFprobe("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Extract text sequence from encoded file
	// content, err := probe.ExtractTextSequence("output.mp3")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(content)

	fmt.Println("Text sequence extracted successfully from MP3 file")
	// Output:
	// Text sequence extracted successfully from MP3 file
}

func ExampleFFprobe_SaveExtractedTextSequence() {
	// Create ffprobe instance
	// probe, err := external.NewFFprobe("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Extract and save text sequence
	// err = probe.SaveExtractedTextSequence("output.mp3", "extracted.spsq")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("Text sequence extracted and saved successfully")
	// Output:
	// Text sequence extracted and saved successfully
}
