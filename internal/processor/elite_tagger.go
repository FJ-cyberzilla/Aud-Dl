package processor

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bogem/id3v2/v2"
)

type TagPayload struct {
	Title        string
	Artist       string
	Album        string
	CoverArtURL  string
	AppName      string // Your Brand / App Name
}

type EliteTagger struct{}

func NewEliteTagger() *EliteTagger {
	return &EliteTagger{}
}

func (et *EliteTagger) InjectMetadataAndBrand(filePath string, payload TagPayload) error {
	tag, err := id3v2.Open(filePath, id3v2.Options{Parse: true})
	if err != nil {
		return fmt.Errorf("failed to open mp3 for tagging: %w", err)
	}
	defer tag.Close()

	// 1. Standard Metadata
	tag.SetTitle(payload.Title)
	tag.SetArtist(payload.Artist)
	if payload.Album != "" {
		tag.SetAlbum(payload.Album)
	}

	// 2. Non-Destructive App Branding
	// Inject Brand Name into 'EncodedBy' (TENC) and 'Publisher' (TPUB) frames
	tag.AddTextFrame(tag.CommonID("Encoded by"), id3v2.EncodingUTF8, payload.AppName)
	tag.AddTextFrame(tag.CommonID("Publisher"), id3v2.EncodingUTF8, payload.AppName)

	// Inject custom user text frame (TXXX) for application signature
	customFrame := id3v2.UserDefinedTextFrame{
		Encoding:    id3v2.EncodingUTF8,
		Description: "SOFTWARE_SIGNATURE",
		Value:       fmt.Sprintf("Downloaded via %s", payload.AppName),
	}
	tag.AddUserDefinedTextFrame(customFrame)

	// 3. High-Resolution Cover Art Injection (APIC Frame)
	if payload.CoverArtURL != "" {
		imageData, mimeType, err := et.fetchCoverArtBytes(payload.CoverArtURL)
		if err == nil && len(imageData) > 0 {
			pic := id3v2.PictureFrame{
				Encoding:    id3v2.EncodingUTF8,
				MimeType:    mimeType,
				PictureType: id3v2.PTFrontCover,
				Description: "Cover Art",
				Picture:     imageData,
			}
			tag.AddAttachedPicture(pic)
		}
	}

	if err := tag.Save(); err != nil {
		return fmt.Errorf("failed to write ID3v2.4 tags: %w", err)
	}

	return nil
}

func (et *EliteTagger) fetchCoverArtBytes(coverURL string) ([]byte, string, error) {
	resp, err := http.Get(coverURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download cover art")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	return data, mimeType, nil
}
