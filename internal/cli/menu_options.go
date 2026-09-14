package cli

// Option constants for type-safe menu definitions
const (
    // Main menu options
    MainMostWanted = iota
    MainDiscography
    MainSingleTrack
    MainSystemAdmin
    MainExit
)

// Format selection constants
const (
    FormatFLAC = iota
    FormatMP3
    FormatAAC
    FormatOGG
)

// Predefined menu items with all metadata
var MainMenuItems = []MenuItem[int]{
    {
        ID:     MainMostWanted,
        Icon:   "🔥",
        Label:  "Most Wanted",
        Description: "Top downloads & trending tracks",
    },
    {
        ID:     MainDiscography,
        Icon:   "💿",
        Label:  "Full Discography",
        Description: "Album collections & artist catalogs",
    },
    {
        ID:     MainSingleTrack,
        Icon:   "🎵",
        Label:  "Single Track",
        Description: "Custom search & specific tracks",
    },
    {
        ID:     MainSystemAdmin,
        Icon:   "⚙️",
        Label:  "System Maintenance",
        Description: "Storage admin & cleanup",
    },
    {
        ID:     MainExit,
        Icon:   "🚪",
        Label:  "Exit",
        Description: "Close application",
    },
}

var FormatMenuItems = []MenuItem[int]{
    {
        ID:     FormatFLAC,
        Icon:   "📀",
        Label:  "FLAC",
        Description: "Lossless 24-bit / 192kHz",
    },
    {
        ID:     FormatMP3,
        Icon:   "🎧",
        Label:  "MP3",
        Description: "Extreme 320kbps CBR",
    },
    {
        ID:     FormatAAC,
        Icon:   "🎶",
        Label:  "AAC",
        Description: "VBR High Efficiency ~256kbps",
    },
    {
        ID:     FormatOGG,
        Icon:   "📻",
        Label:  "OGG",
        Description: "Vorbis Quality 9 ~320kbps",
    },
}
