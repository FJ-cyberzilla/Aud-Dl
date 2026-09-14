package ui

type Theme struct {
	Name      string
	Primary   string
	Secondary string
	Accent    string
	Success   string
	Warning   string
	Error     string
	Border    string
	TextDim   string
	Subtext   string
	Highlight string
	Text      string
}

var CyberDarkTheme = Theme{
	Name:      "Cyber Dark",
	Primary:   "#00F0FF", // Bright Cyan
	Secondary: "#FF007F", // Neon Magenta
	Accent:    "#FFE600", // Bright Yellow
	Success:   "#00FF66", // Electric Green
	Warning:   "#FF9900", // Amber
	Error:     "#FF0033", // Bright Red
	Border:    "#333344", // Dark Purple Gray
	TextDim:   "#777788",
	Subtext:   "#777788",
	Highlight: "#FFE600",
	Text:      "#FFFFFF",
}

var MidnightNeonTheme = Theme{
	Name:      "Midnight Neon",
	Primary:   "#8A2BE2", // Blue Violet
	Secondary: "#00FFFF", // Cyan
	Accent:    "#FF1493", // Deep Pink
	Success:   "#32CD32", // Lime Green
	Warning:   "#FF8C00", // Dark Orange
	Error:     "#DC143C", // Crimson
	Border:    "#2E0854",
	TextDim:   "#6A5ACD",
	Subtext:   "#6A5ACD",
	Highlight: "#00FFFF",
	Text:      "#FFFFFF",
}
