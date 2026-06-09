package exporter

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const fastTextColorSettingsPath = ".obsidian/plugins/fast-text-color/data.json"

// The settings written by the Fast Text Color plugin.
// See https://github.com/Superschnizel/obsidian-fast-text-color/blob/23f7e97a4ffac7643506a8e422dca0f016331529/src/FTCSettings.ts#L59.
type FTCSettings struct {
	// The defined themes.
	Themes []FTCTheme `json:"themes"`
	// The default theme to use.
	ThemeIndex *int `json:"themeIndex"`
	// Theme lookup by name.
	ThemeMap map[string]*FTCTheme `json:"-"`
}

// Parses the settings file within the Obsidian vault.
func NewFTCSettings(vaultPath string) (*FTCSettings, error) {
	p := filepath.Join(vaultPath, fastTextColorSettingsPath)
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var s FTCSettings
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}

	s.ThemeMap = make(map[string]*FTCTheme, len(s.Themes))
	for i := range s.Themes {
		t := &s.Themes[i]
		s.ThemeMap[t.Name] = t

		t.ColorMap = make(map[string]*FTCColor, len(t.Colors))
		for i := range t.Colors {
			c := &t.Colors[i]
			t.ColorMap[c.ID] = c
		}
	}

	return &s, nil
}

// Returns the default theme in the settings.
func (st *FTCSettings) GetDefaultTheme() *FTCTheme {
	var idx int
	if st.ThemeIndex != nil {
		idx = *st.ThemeIndex
	}

	if idx < len(st.Themes) {
		return &st.Themes[idx]
	}

	return nil
}

// A Fast Text Color theme.
type FTCTheme struct {
	// The theme name.
	Name string `json:"name"`
	// The colors defined by the theme.
	Colors []FTCColor `json:"colors"`
	// Color lookup by name.
	ColorMap map[string]*FTCColor `json:"-"`
}

// A Fast Text Color color.
type FTCColor struct {
	// The color as a hexadecimal number (e.g., #FFFFFF).
	Color string `json:"color"`
	// The ID of the color.
	ID string `json:"id"`
	// Whether or not the colored text should be italic.
	Italic bool `json:"italic"`
	// Whether or not the colored text should be bold.
	Bold bool `json:"bold"`
}
