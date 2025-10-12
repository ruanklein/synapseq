#!/usr/bin/env python3
"""
SynapSeq JSON Generator - Example Integration
==============================================

This script demonstrates how to programmatically generate a JSON sequence
for SynapSeq and invoke the CLI to create audio output.

Use case: Automate brainwave entrainment session creation based on parameters
"""

import json
import subprocess
import sys
from pathlib import Path


def create_synapseq_sequence(duration_minutes: int, start_freq: float, end_freq: float):
    """
    Generate a progressive sequence in JSON format.
    
    Args:
        duration_minutes: Total duration in minutes
        start_freq: Starting resonance frequency (Hz)
        end_freq: Ending resonance frequency (Hz)
    
    Returns:
        dict: SynapSeq JSON sequence structure
    """
    
    # Calculate time points
    fade_in = 15000  # 15 seconds fade-in
    fade_out = 30000  # 30 seconds fade-out
    total_time = duration_minutes * 60 * 1000
    main_duration = total_time - fade_in - fade_out
    
    # Calculate frequency transition points (every 2 minutes)
    interval = 120000  # 2 minutes in milliseconds
    num_steps = int(main_duration / interval)
    freq_step = (end_freq - start_freq) / num_steps if num_steps > 0 else 0
    
    sequence = {
        "description": [
            f"Progressive Session - {duration_minutes} minutes",
            f"Frequency range: {start_freq}Hz -> {end_freq}Hz",
            "Generated programmatically using Python"
        ],
        "options": {
            "samplerate": 44100,
            "volume": 100
        },
        "sequence": []
    }
    
    # Add fade-in (silence)
    sequence["sequence"].append({
        "time": 0,
        "transition": "smooth",
        "track": {
            "tones": [{
                "mode": "binaural",
                "carrier": 250,
                "resonance": start_freq,
                "amplitude": 0,
                "waveform": "sine"
            }],
            "noises": [{
                "mode": "brown",
                "amplitude": 0
            }]
        }
    })
    
    # Main session start
    sequence["sequence"].append({
        "time": fade_in,
        "transition": "smooth",
        "track": {
            "tones": [{
                "mode": "binaural",
                "carrier": 250,
                "resonance": start_freq,
                "amplitude": 15,
                "waveform": "sine"
            }],
            "noises": [{
                "mode": "brown",
                "amplitude": 30
            }]
        }
    })
    
    # Progressive frequency transitions
    current_time = fade_in
    current_freq = start_freq
    
    for step in range(num_steps):
        current_time += interval
        current_freq += freq_step
        
        sequence["sequence"].append({
            "time": current_time,
            "transition": "smooth",
            "track": {
                "tones": [{
                    "mode": "binaural",
                    "carrier": 250,
                    "resonance": round(current_freq, 2),
                    "amplitude": 15,
                    "waveform": "sine"
                }],
                "noises": [{
                    "mode": "brown",
                    "amplitude": 30
                }]
            }
        })
    
    # Fade-out
    sequence["sequence"].append({
        "time": total_time - fade_out,
        "transition": "smooth",
        "track": {
            "tones": [{
                "mode": "binaural",
                "carrier": 250,
                "resonance": end_freq,
                "amplitude": 0,
                "waveform": "sine"
            }],
            "noises": [{
                "mode": "brown",
                "amplitude": 0
            }]
        }
    })
    
    return sequence


def main():
    # Parse command-line arguments
    if len(sys.argv) < 4:
        print("Usage: generate_json_sequence.py <start_freq> <end_freq> <duration_minutes>")
        print("\nExamples:")
        print("  python3 generate_json_sequence.py 10.0 6.0 10    # Alpha to Theta, 10 minutes")
        print("  python3 generate_json_sequence.py 14.0 10.0 15   # Beta to Alpha, 15 minutes")
        print("  python3 generate_json_sequence.py 8.0 4.0 20     # Alpha to Theta, 20 minutes")
        print("\nFrequency ranges:")
        print("  Delta: 0.5-4 Hz (deep sleep)")
        print("  Theta: 4-8 Hz (meditation)")
        print("  Alpha: 8-13 Hz (relaxation)")
        print("  Beta: 13-30 Hz (focus)")
        sys.exit(1)
    
    try:
        start_frequency = float(sys.argv[1])
        end_frequency = float(sys.argv[2])
        duration = int(sys.argv[3])
    except ValueError:
        print("Error: Invalid parameters. Please provide numeric values.", file=sys.stderr)
        sys.exit(1)
    
    # Validate inputs
    if start_frequency <= 0 or end_frequency <= 0:
        print("Error: Frequencies must be positive numbers.", file=sys.stderr)
        sys.exit(1)
    
    if duration <= 0:
        print("Error: Duration must be a positive number of minutes.", file=sys.stderr)
        sys.exit(1)
    
    # Create output directory
    generated_dir = Path("generated")
    generated_dir.mkdir(exist_ok=True)
    
    # Generate output filenames based on parameters
    filename = f"progressive_{start_frequency}Hz-{end_frequency}Hz_{duration}min"
    output_json = generated_dir / f"{filename}.json"
    output_wav = generated_dir / f"{filename}.wav"

    print(f"Generating {duration}-minute progressive sequence...")
    print(f"   Frequency: {start_frequency}Hz -> {end_frequency}Hz")
    
    # Generate sequence
    sequence = create_synapseq_sequence(duration, start_frequency, end_frequency)
    
    # Save JSON file
    with open(output_json, 'w', encoding='utf-8') as f:
        json.dump(sequence, f, indent=2)
    
    print(f"JSON sequence saved to: {output_json}")
    
    # Invoke SynapSeq CLI
    print(f"Generating audio with SynapSeq...")
    
    try:
        result = subprocess.run(
            ['synapseq', '-json', str(output_json), str(output_wav)],
            capture_output=True,
            text=True,
            check=True
        )

        print(result.stderr)
        print(f"Audio file created: {output_wav}")
        
    except subprocess.CalledProcessError as e:
        print(f"Error running SynapSeq:", file=sys.stderr)
        print(e.stderr, file=sys.stderr)
        sys.exit(1)
    
    except FileNotFoundError:
        print("SynapSeq not found in PATH. Please install it first.", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
