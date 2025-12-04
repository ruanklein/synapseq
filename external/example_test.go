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

func ExampleFFmpeg_MP3() {
	// Create SynapSeq application context
	ctx, err := synapseq.NewAppContext("input.spsq", "output.mp3", "text")
	if err != nil {
		log.Fatal(err)
	}

	// Load sequence (required before encoding)
	// _ = ctx.LoadSequence()

	// Create ffmpeg instance
	// _, err = external.NewFFmpeg("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Encode MP3 using default VBR mode (highest quality)
	// _ = encoder.MP3(ctx, nil)

	fmt.Printf("MP3 encoding (VBR) executed successfully for output: %s\n", ctx.OutputFile())
	// Output:
	// MP3 encoding (VBR) executed successfully for output: output.mp3
}

func ExampleFFmpeg_MP3_cbr() {
	// Create SynapSeq application context
	ctx, err := synapseq.NewAppContext("input.spsq", "output.mp3", "text")
	if err != nil {
		log.Fatal(err)
	}

	// Load sequence (required before encoding)
	// _ = ctx.LoadSequence()

	// Create ffmpeg instance
	// _, err = external.NewFFmpeg("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Encode MP3 using CBR mode at 320 kbps
	// _ = encoder.MP3(ctx, &external.MP3Options{Mode: external.MP3ModeCBR})

	fmt.Printf("MP3 encoding (CBR) executed successfully for output: %s\n", ctx.OutputFile())
	// Output:
	// MP3 encoding (CBR) executed successfully for output: output.mp3
}
