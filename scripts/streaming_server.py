#!/usr/bin/env python3
"""
SynapSeq Streaming Server - HTTP Audio Streaming Example
=========================================================

This script demonstrates real-time audio streaming using SynapSeq's RAW output.

Use case: Build a web service that streams brainwave entrainment audio
to clients in real-time, allowing consumption via curl, browsers with HTML5
audio players, or custom clients.

Architecture:
-------------
1. HTTP server receives request with parameters (frequency, duration, etc.)
2. Server generates JSON sequence on-the-fly
3. SynapSeq generates RAW audio and streams to stdout
4. Server chunks the output and streams to client
5. Client receives chunks and plays/processes in real-time

Usage:
------
Start server:
    python3 streaming_server.py

Test with curl + sox:
    curl -N "http://localhost:8000/stream?freq=10&duration=5" | \
        play -t raw -r 44100 -e signed-integer -b 24 -c 2 -

The server supports query parameters:
    - freq: Resonance frequency (Hz) - default: 10
    - duration: Duration in minutes - default: 5
    - noise: Noise type (white, pink, brown) - default: pink
    - mode: Tone mode (binaural, monaural, isochronic) - default: binaural
    - carrier: Carrier frequency (Hz) - default: 300
    - amplitude: Amplitude (0-100) - default: 15
"""

import json
import subprocess
import tempfile
from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs
import sys
import os


def create_streaming_sequence(params: dict) -> dict:
    """
    Create a SynapSeq JSON sequence based on request parameters.
    
    Args:
        params: Dictionary with sequence parameters
    
    Returns:
        dict: SynapSeq JSON sequence
    """
    freq = float(params.get('freq', [10])[0])
    duration = int(params.get('duration', [5])[0])
    noise_type = params.get('noise', ['pink'])[0]
    tone_mode = params.get('mode', ['binaural'])[0]
    carrier = int(params.get('carrier', [300])[0])
    amplitude = int(params.get('amplitude', [15])[0])
    
    # Duration in milliseconds
    total_ms = duration * 60 * 1000
    fade_in = 15000  # 15 seconds
    fade_out = 30000  # 30 seconds
    
    sequence = {
        "description": [
            f"Streaming Session - {duration} minutes",
            f"Resonance: {freq}Hz, Carrier: {carrier}Hz",
            f"Mode: {tone_mode}"
        ],
        "options": {
            "samplerate": 44100,
            "volume": 100
        },
        "sequence": [
            {
                "time": 0,
                "transition": "steady",
                "track": {
                    "noises": [
                        {"mode": noise_type, "amplitude": 0}
                    ],
                    "tones": [
                        {
                            "mode": tone_mode,
                            "carrier": carrier,
                            "resonance": freq,
                            "amplitude": 0,
                            "waveform": "sine"
                        }
                    ]
                }
            },
            {
                "time": fade_in,
                "transition": "steady",
                "track": {
                    "noises": [
                        {"mode": noise_type, "amplitude": 30}
                    ],
                    "tones": [
                        {
                            "mode": tone_mode,
                            "carrier": carrier,
                            "resonance": freq,
                            "amplitude": amplitude,
                            "waveform": "sine"
                        }
                    ]
                }
            },
            {
                "time": total_ms - fade_out,
                "transition": "steady",
                "track": {
                    "noises": [
                        {"mode": noise_type, "amplitude": 30}
                    ],
                    "tones": [
                        {
                            "mode": tone_mode,
                            "carrier": carrier,
                            "resonance": freq,
                            "amplitude": amplitude,
                            "waveform": "sine"
                        }
                    ]
                }
            },
            {
                "time": total_ms,
                "transition": "steady",
                "track": {
                    "noises": [
                        {"mode": noise_type, "amplitude": 0}
                    ],
                    "tones": [
                        {
                            "mode": tone_mode,
                            "carrier": carrier,
                            "resonance": freq,
                            "amplitude": 0,
                            "waveform": "sine"
                        }
                    ]
                }
            }
        ]
    }
    
    return sequence


class StreamingHandler(BaseHTTPRequestHandler):
    """HTTP request handler for audio streaming."""
    
    def log_message(self, format, *args):
        """Override to add custom logging."""
        sys.stderr.write(f"[{self.log_date_time_string()}] {format % args}\n")
    
    def do_GET(self):
        """Handle GET requests."""
        parsed_path = urlparse(self.path)
        
        if parsed_path.path == '/stream':
            self.stream_audio(parsed_path)
        elif parsed_path.path == '/':
            self.send_info()
        else:
            self.send_error(404, "Not Found")
    
    def stream_audio(self, parsed_path):
        """
        Stream RAW audio generated by SynapSeq.
        
        Args:
            parsed_path: Parsed URL with query parameters
        """
        process = None
        tmp_path = None
        
        try:
            # Parse query parameters
            params = parse_qs(parsed_path.query)
            
            # Log request
            freq = params.get('freq', ['10'])[0]
            duration = params.get('duration', ['5'])[0]
            mode = params.get('mode', ['binaural'])[0]
            sys.stderr.write(f"\n📻 Streaming request: freq={freq}Hz, duration={duration}min, mode={mode}\n")
            
            # Generate sequence
            sequence = create_streaming_sequence(params)
            
            # Create temporary file for sequence
            with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as tmp:
                json.dump(sequence, tmp, indent=2)
                tmp_path = tmp.name
            
            # Start SynapSeq process with RAW output to stdout
            # Output goes to stdout (pipe) instead of a file
            process = subprocess.Popen(
                ['synapseq', '-quiet', '-json', tmp_path, '-'],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                bufsize=8192  # Use buffered I/O for better performance
            )
            
            # Send HTTP headers
            self.send_response(200)
            self.send_header('Content-Type', 'audio/x-raw')
            self.send_header('Cache-Control', 'no-cache')
            self.send_header('X-Audio-Format', 'raw')
            self.send_header('X-Audio-Samplerate', '44100')
            self.send_header('X-Audio-Bitdepth', '24')
            self.send_header('X-Audio-Channels', '2')
            self.send_header('X-Audio-Encoding', 'signed-integer')
            self.end_headers()
            
            # Stream audio chunks
            chunk_size = 8192  # 8KB chunks for better performance
            total_bytes = 0
            
            sys.stderr.write("Streaming started...\n")
            
            try:
                while True:
                    chunk = process.stdout.read(chunk_size)
                    if not chunk:
                        break
                    
                    try:
                        # Send chunk to client
                        self.wfile.write(chunk)
                        self.wfile.flush()  # Ensure data is sent immediately
                        total_bytes += len(chunk)
                        
                        # Log progress every 1MB
                        if total_bytes % (1024 * 1024) == 0:
                            sys.stderr.write(f"   Streamed: {total_bytes / (1024 * 1024):.1f} MB\n")
                    
                    except BrokenPipeError:
                        # Client disconnected - this is normal, not an error
                        sys.stderr.write(f"Client disconnected after {total_bytes / (1024 * 1024):.2f} MB\n")
                        break
                
                # Wait for process to complete
                process.wait()
                
                # Log completion
                mb_streamed = total_bytes / (1024 * 1024)
                sys.stderr.write(f"Stream completed: {mb_streamed:.2f} MB sent\n\n")
                
                # Check for errors
                if process.returncode != 0:
                    stderr_output = process.stderr.read().decode('utf-8')
                    sys.stderr.write(f"SynapSeq error: {stderr_output}\n")
            
            except BrokenPipeError:
                # Client disconnected during streaming - this is normal
                mb_streamed = total_bytes / (1024 * 1024)
                sys.stderr.write(f"Client disconnected: {mb_streamed:.2f} MB sent\n\n")
        
        except BrokenPipeError:
            # Client disconnected before we could start - ignore
            sys.stderr.write("Client disconnected early\n\n")
        
        except Exception as e:
            sys.stderr.write(f"Error: {str(e)}\n\n")
            # Only try to send error if connection is still alive
            try:
                self.send_error(500, f"Internal Server Error: {str(e)}")
            except:
                pass  # Ignore if we can't send the error
        
        finally:
            # Cleanup
            if process and process.poll() is None:
                process.terminate()
                process.wait()
            
            if tmp_path and os.path.exists(tmp_path):
                os.unlink(tmp_path)
    
    def send_info(self):
        """Send information page about the streaming server."""
        info = """<!DOCTYPE html>
<html>
<head>
    <title>SynapSeq Streaming Server</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        * { box-sizing: border-box; }
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', monospace; 
            max-width: 900px; 
            margin: 40px auto; 
            padding: 20px;
            background: #f5f5f5;
            color: #333;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            margin-bottom: 20px;
        }
        h1 { 
            color: #ffffff;
            margin-top: 0;
            border-bottom: 3px solid rgba(255,255,255,0.5);
            padding-bottom: 10px;
        }
        h2 { 
            color: #34495e;
            margin-top: 30px;
            border-bottom: 1px solid #ecf0f1;
            padding-bottom: 5px;
        }
        
        /* Player Section */
        .player-section {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            border-radius: 10px;
            margin-bottom: 20px;
        }
        .player-section h2 {
            color: white;
            border-bottom: 2px solid rgba(255,255,255,0.3);
        }
        
        /* Form Styles */
        .form-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin: 20px 0;
        }
        .form-group {
            display: flex;
            flex-direction: column;
        }
        label {
            font-weight: 600;
            margin-bottom: 5px;
            font-size: 0.9em;
            color: rgba(255,255,255,0.9);
        }
        input, select {
            padding: 10px;
            border: 2px solid rgba(255,255,255,0.3);
            border-radius: 5px;
            font-family: monospace;
            font-size: 14px;
            background: rgba(255,255,255,0.1);
            color: white;
            transition: border-color 0.3s;
        }
        input:focus, select:focus {
            outline: none;
            border-color: rgba(255,255,255,0.8);
            background: rgba(255,255,255,0.15);
        }
        input::placeholder {
            color: rgba(255,255,255,0.5);
        }
        
        /* Button Styles */
        .button-group {
            display: flex;
            gap: 10px;
            margin-top: 20px;
        }
        button {
            flex: 1;
            padding: 15px 30px;
            font-size: 16px;
            font-weight: 600;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            transition: all 0.3s;
            font-family: monospace;
        }
        button:disabled {
            opacity: 0.5;
            cursor: not-allowed;
        }
        .btn-play {
            background: #2ecc71;
            color: white;
        }
        .btn-play:hover:not(:disabled) {
            background: #27ae60;
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(46,204,113,0.4);
        }
        .btn-stop {
            background: #e74c3c;
            color: white;
        }
        .btn-stop:hover:not(:disabled) {
            background: #c0392b;
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(231,76,60,0.4);
        }
        
        /* Audio Player */
        audio {
            width: 100%;
            margin: 20px 0;
            border-radius: 5px;
        }
        
        /* Status Messages */
        .status {
            padding: 15px;
            border-radius: 5px;
            margin: 15px 0;
            font-weight: 500;
            display: none;
        }
        .status.info {
            background: rgba(52,152,219,0.2);
            border-left: 4px solid #3498db;
            color: #2c3e50;
        }
        .status.success {
            background: rgba(46,204,113,0.2);
            border-left: 4px solid #2ecc71;
            color: #27ae60;
        }
        .status.error {
            background: rgba(231,76,60,0.2);
            border-left: 4px solid #e74c3c;
            color: #c0392b;
        }
        .status.warning {
            background: rgba(241,196,15,0.2);
            border-left: 4px solid #f1c40f;
            color: #f39c12;
        }
        
        /* Documentation Styles */
        pre {
            background: #2c3e50;
            color: #ecf0f1;
            padding: 15px;
            border-radius: 5px;
            overflow-x: auto;
            font-size: 13px;
        }
        code {
            background: #ecf0f1;
            padding: 2px 6px;
            border-radius: 3px;
            font-size: 0.9em;
            color: #e74c3c;
        }
        .param { color: #3498db; font-weight: 600; }
        .example { margin: 20px 0; }
        
        ul { line-height: 1.8; }
        
        .info-note {
            background: rgba(255, 255, 255, 0.2);
            border-left: 4px solid #ffd700;
            padding: 15px;
            margin: 15px 0;
            border-radius: 5px;
            color: #ffffff;
        }
        
        /* Navbar */
        .navbar {
            background: #2c3e50;
            padding: 15px 30px;
            border-radius: 10px;
            margin-bottom: 20px;
            display: flex;
            justify-content: space-between;
            align-items: center;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .navbar-brand {
            color: #ecf0f1;
            font-size: 1.2em;
            font-weight: 700;
            text-decoration: none;
        }
        .navbar-links {
            display: flex;
            gap: 20px;
        }
        .navbar-links a {
            color: #3498db;
            text-decoration: none;
            font-weight: 600;
            transition: color 0.3s;
        }
        .navbar-links a:hover {
            color: #2ecc71;
        }
        
        /* Footer */
        .footer {
            background: #2c3e50;
            color: #ecf0f1;
            padding: 20px 30px;
            border-radius: 10px;
            margin-top: 20px;
            text-align: center;
            box-shadow: 0 -2px 10px rgba(0,0,0,0.1);
        }
        .footer a {
            color: #3498db;
            text-decoration: none;
            transition: color 0.3s;
        }
        .footer a:hover {
            color: #2ecc71;
        }
    </style>
</head>
<body>
    <nav class="navbar">
        <div class="navbar-brand">SynapSeq Streaming</div>
        <div class="navbar-links">
            <a href="https://github.com/ruanklein/synapseq" target="_blank">📦 GitHub Repository</a>
            <a href="https://github.com/ruanklein/synapseq/blob/main/docs/USAGE.md" target="_blank">📖 Documentation</a>
        </div>
    </nav>
    <div class="player-section">
        <h1>SynapSeq Streaming Server</h1>
        <p>Real-time brainwave entrainment audio streaming using SynapSeq RAW output.</p>
        
        <h2>Web Player</h2>
        
        <form id="streamForm">
            <div class="form-grid">
                <div class="form-group">
                    <label for="freq">Frequency (Hz)</label>
                    <input type="number" id="freq" name="freq" value="10" min="0.5" max="30" step="0.5" required>
                </div>
                
                <div class="form-group">
                    <label for="duration">Duration (min)</label>
                    <input type="number" id="duration" name="duration" value="5" min="1" max="60" required>
                </div>
                
                <div class="form-group">
                    <label for="mode">Tone Mode</label>
                    <select id="mode" name="mode">
                        <option value="binaural" selected>Binaural</option>
                        <option value="monaural">Monaural</option>
                        <option value="isochronic">Isochronic</option>
                    </select>
                </div>
                
                <div class="form-group">
                    <label for="noise">Noise Type</label>
                    <select id="noise" name="noise">
                        <option value="pink" selected>Pink</option>
                        <option value="white">White</option>
                        <option value="brown">Brown</option>
                    </select>
                </div>
                
                <div class="form-group">
                    <label for="carrier">Carrier (Hz)</label>
                    <input type="number" id="carrier" name="carrier" value="300" min="100" max="1000" required>
                </div>
                
                <div class="form-group">
                    <label for="amplitude">Amplitude (0-100)</label>
                    <input type="number" id="amplitude" name="amplitude" value="15" min="0" max="100" required>
                </div>
            </div>
            
            <div class="button-group">
                <button type="submit" class="btn-play" id="playBtn">Generate & Play</button>
                <button type="button" class="btn-stop" id="stopBtn" disabled>Stop</button>
            </div>
        </form>
        
        <div id="status" class="status"></div>
        
        <audio id="audioPlayer" controls style="display:none;"></audio>
        
        <div class="info-note">
            <strong>Note:</strong> For binaural beats, use headphones for optimal effect. 
            Audio generation may take a few seconds to start.
        </div>
    </div>
    
    <div class="container">
        <h2>API Endpoint</h2>
        <pre>GET /stream?freq=10&duration=5&mode=binaural</pre>
        
        <h2>Query Parameters</h2>
        <ul>
            <li><code class="param">freq</code> - Resonance frequency in Hz (default: 10)</li>
            <li><code class="param">duration</code> - Duration in minutes (default: 5)</li>
            <li><code class="param">noise</code> - Noise type: white, pink, brown (default: pink)</li>
            <li><code class="param">mode</code> - Tone mode: binaural, monaural, isochronic (default: binaural)</li>
            <li><code class="param">carrier</code> - Carrier frequency (Hz, default: 300)</li>
            <li><code class="param">amplitude</code> - Amplitude 0-100 (default: 15)</li>
        </ul>
        
        <h2>Command-Line Usage</h2>
        
        <div class="example">
            <strong>Play with curl + sox:</strong>
            <pre>curl -N "http://localhost:8000/stream?freq=10&duration=5" | \\
    play -t raw -r 44100 -e signed-integer -b 24 -c 2 -</pre>
        </div>
        
        <div class="example">
            <strong>Save to file:</strong>
            <pre>curl -N "http://localhost:8000/stream?freq=8&duration=10&mode=isochronic" > output.raw</pre>
        </div>
        
        <div class="example">
            <strong>Convert to WAV with ffmpeg:</strong>
            <pre>curl -N "http://localhost:8000/stream?freq=6&duration=15" | \\
    ffmpeg -f s24le -ar 44100 -ac 2 -i - output.wav</pre>
        </div>

        <h2>Audio Format</h2>
        <ul>
            <li><strong>Type:</strong> RAW (PCM)</li>
            <li><strong>Sample Rate:</strong> 44100 Hz</li>
            <li><strong>Bit Depth:</strong> 24 bits</li>
            <li><strong>Channels:</strong> 2 (stereo)</li>
            <li><strong>Encoding:</strong> Signed Integer (little-endian)</li>
        </ul>

        <h2>Integration Patterns</h2>
        <ul>
            <li><strong>Command-line clients:</strong> curl + audio players (sox, ffplay)</li>
            <li><strong>Web applications:</strong> Fetch API + WAV conversion</li>
            <li><strong>Mobile apps:</strong> HTTP streaming + native audio APIs</li>
            <li><strong>Processing pipelines:</strong> Stream to audio processing tools</li>
        </ul>
    </div>
    
    <footer class="footer">
        <p>Copyright © 2025 Ruan • <a href="https://ruan.sh" target="_blank">ruan.sh</a></p>
        <p style="margin-top: 10px; font-size: 0.9em; color: #95a5a6;">
            Licensed under GPL v2
        </p>
    </footer>

    <script>
        const form = document.getElementById('streamForm');
        const playBtn = document.getElementById('playBtn');
        const stopBtn = document.getElementById('stopBtn');
        const audioPlayer = document.getElementById('audioPlayer');
        const status = document.getElementById('status');
        
        let audioContext = null;
        let sourceNode = null;
        let currentRequest = null;
        
        function showStatus(message, type = 'info') {
            status.textContent = message;
            status.className = 'status ' + type;
            status.style.display = 'block';
        }
        
        function hideStatus() {
            status.style.display = 'none';
        }
        
        async function convertRawToWav(rawData, sampleRate = 44100, numChannels = 2, bitDepth = 24) {
            // Calculate sizes
            const bytesPerSample = bitDepth / 8;
            const dataSize = rawData.byteLength;
            
            // Convert 24-bit to 16-bit for browser compatibility
            const samples16 = new Int16Array(dataSize / bytesPerSample);
            const dataView = new DataView(rawData);
            
            for (let i = 0; i < samples16.length; i++) {
                // Read 24-bit sample (little-endian)
                const byte1 = dataView.getUint8(i * 3);
                const byte2 = dataView.getUint8(i * 3 + 1);
                const byte3 = dataView.getUint8(i * 3 + 2);
                
                // Combine into 24-bit signed integer
                let sample24 = (byte3 << 16) | (byte2 << 8) | byte1;
                
                // Sign extend if negative
                if (sample24 & 0x800000) {
                    sample24 |= 0xFF000000;
                }
                
                // Convert to 16-bit
                samples16[i] = sample24 >> 8;
            }
            
            // Create WAV file
            const wavSize = 44 + samples16.length * 2;
            const wavBuffer = new ArrayBuffer(wavSize);
            const view = new DataView(wavBuffer);
            
            // Write WAV header
            const writeString = (offset, string) => {
                for (let i = 0; i < string.length; i++) {
                    view.setUint8(offset + i, string.charCodeAt(i));
                }
            };
            
            writeString(0, 'RIFF');
            view.setUint32(4, wavSize - 8, true);
            writeString(8, 'WAVE');
            writeString(12, 'fmt ');
            view.setUint32(16, 16, true); // fmt chunk size
            view.setUint16(20, 1, true); // PCM format
            view.setUint16(22, numChannels, true);
            view.setUint32(24, sampleRate, true);
            view.setUint32(28, sampleRate * numChannels * 2, true); // byte rate
            view.setUint16(32, numChannels * 2, true); // block align
            view.setUint16(34, 16, true); // bits per sample (converted to 16-bit)
            writeString(36, 'data');
            view.setUint32(40, samples16.length * 2, true);
            
            // Copy audio data
            const samples = new Int16Array(wavBuffer, 44);
            samples.set(samples16);
            
            return wavBuffer;
        }
        
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            // Build query string
            const formData = new FormData(form);
            const params = new URLSearchParams(formData);
            const streamUrl = '/stream?' + params.toString();
            
            try {
                playBtn.disabled = true;
                stopBtn.disabled = false;
                showStatus('Generating audio stream... Please wait.', 'info');
                
                // Fetch the RAW audio stream
                currentRequest = fetch(streamUrl);
                const response = await currentRequest;
                
                if (!response.ok) {
                    throw new Error('Failed to fetch audio stream');
                }
                
                showStatus('Downloading audio data...', 'info');
                
                // Read the entire stream
                const rawData = await response.arrayBuffer();
                
                showStatus('Converting to WAV format...', 'info');
                
                // Convert RAW to WAV
                const wavData = await convertRawToWav(rawData);
                
                // Create blob and URL
                const blob = new Blob([wavData], { type: 'audio/wav' });
                const audioUrl = URL.createObjectURL(blob);
                
                // Set up audio player
                audioPlayer.src = audioUrl;
                audioPlayer.style.display = 'block';
                
                showStatus('Audio ready! Playing...', 'success');
                
                // Play audio
                try {
                    await audioPlayer.play();
                } catch (playError) {
                    showStatus('Audio ready, but autoplay blocked. Click play button.', 'warning');
                }
                
                // Clean up when audio ends
                audioPlayer.onended = () => {
                    URL.revokeObjectURL(audioUrl);
                    playBtn.disabled = false;
                    stopBtn.disabled = true;
                    showStatus('Playback completed.', 'success');
                };
                
            } catch (error) {
                console.error('Error:', error);
                showStatus('Error: ' + error.message, 'error');
                playBtn.disabled = false;
                stopBtn.disabled = true;
            }
        });
        
        stopBtn.addEventListener('click', () => {
            if (audioPlayer) {
                audioPlayer.pause();
                audioPlayer.currentTime = 0;
                audioPlayer.style.display = 'none';
            }
            
            if (currentRequest) {
                currentRequest = null;
            }
            
            playBtn.disabled = false;
            stopBtn.disabled = true;
            showStatus('Playback stopped.', 'info');
        });
        
        // Frequency range helper
        const freqInput = document.getElementById('freq');
        freqInput.addEventListener('change', (e) => {
            const freq = parseFloat(e.target.value);
            let range = '';
            
            if (freq >= 0.5 && freq < 4) range = 'Delta (Deep sleep)';
            else if (freq >= 4 && freq < 8) range = 'Theta (Meditation)';
            else if (freq >= 8 && freq < 13) range = 'Alpha (Relaxation)';
            else if (freq >= 13 && freq < 30) range = 'Beta (Focus)';
            else if (freq >= 30) range = 'Gamma (Peak focus)';
            
            if (range) {
                showStatus(range, 'info');
                setTimeout(hideStatus, 3000);
            }
        });
    </script>
</body>
</html>"""
        
        self.send_response(200)
        self.send_header('Content-Type', 'text/html; charset=utf-8')
        self.end_headers()
        self.wfile.write(info.encode('utf-8'))


def main():
    """Start the streaming server."""
    host = 'localhost'
    port = 8000
    
    server = HTTPServer((host, port), StreamingHandler)
    
    print(f"SynapSeq Streaming Server")
    print(f"{'=' * 50}")
    print(f"Server running at: http://{host}:{port}")
    print(f"Info page: http://{host}:{port}/")
    print(f"Stream endpoint: http://{host}:{port}/stream")
    print(f"\n{'=' * 50}")
    print(f"Test with curl + sox:")
    print(f'  curl -N "http://{host}:{port}/stream?freq=10&duration=5" | \\')
    print(f'      play -t raw -r 44100 -e signed-integer -b 24 -c 2 -')
    print(f"\n{'=' * 50}")
    print(f"Press Ctrl+C to stop the server\n")
    
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\n\nServer stopped")
        sys.exit(0)


if __name__ == '__main__':
    main()
