package processor

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

type TranscodeProfile string

const (
	ProfileMP3V0   TranscodeProfile = "V0"   // High Quality VBR (~240 kbps)
	ProfileMP3320K TranscodeProfile = "320k" // Constant Bitrate 320 kbps
	ProfileOpus160 TranscodeProfile = "160k" // Low Latency High-Fidelity Opus
)

type PipeTranscoder struct {
	FFmpegPath string
}

func NewPipeTranscoder() (*PipeTranscoder, error) {
	path, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, fmt.Errorf("ffmpeg binary not found in PATH: %w", err)
	}
	return &PipeTranscoder{FFmpegPath: path}, nil
}

// StreamTranscode pipes an input audio reader through FFmpeg using EBU R128 loudness normalization
func (pt *PipeTranscoder) StreamTranscode(ctx context.Context, input io.Reader, profile TranscodeProfile) (io.ReadCloser, error) {
	args := []string{
		"-i", "pipe:0",
		"-af", "loudnorm=I=-16:LRA=11:TP=-1.5", // EBU R128 audio normalization filter
		"-ar", "44100",
	}

	switch profile {
	case ProfileMP3V0:
		args = append(args, "-c:a", "libmp3lame", "-q:a", "0")
	case ProfileMP3320K:
		args = append(args, "-c:a", "libmp3lame", "-b:a", "320k")
	case ProfileOpus160:
		args = append(args, "-c:a", "libopus", "-b:a", "160k")
	}

	args = append(args, "-f", "mp3", "pipe:1")

	cmd := exec.CommandContext(ctx, pt.FFmpegPath, args...)
	cmd.Stdin = input

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return outPipe, nil
}

// MemoryPipedTranscoder handles transcoding with in-memory buffering.
type MemoryPipedTranscoder struct {
	transcoder *PipeTranscoder
}

func NewMemoryPipedTranscoder() *MemoryPipedTranscoder {
	pt, _ := NewPipeTranscoder() // Simplification for now, error handled elsewhere
	return &MemoryPipedTranscoder{transcoder: pt}
}

// StreamTranscode delegates to PipeTranscoder
func (mpt *MemoryPipedTranscoder) StreamTranscode(ctx context.Context, input io.Reader, profile TranscodeProfile) (io.ReadCloser, error) {
	return mpt.transcoder.StreamTranscode(ctx, input, profile)
}
