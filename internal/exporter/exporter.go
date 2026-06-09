package exporter

import (
	"errors"
	"fmt"
	"os"

	"github.com/ongyx/knap/internal/schema"
)

// Exporter exports an Obsidian vault to Outline format.
type Exporter struct {
	identity    schema.Identity
	vault       *Vault
	ftcSettings *FTCSettings
}

// Creates a new exporter with the given identity.
func New(identity schema.Identity, vaultPath string) (*Exporter, error) {
	v := NewVault(vaultPath)

	// Try to read settings from the Fast Text Color plugin.
	st, err := NewFTCSettings(vaultPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		// The settings file exists but it can't be read for some other reason.
		return nil, fmt.Errorf("fast text color settings exists in vault but parsing failed: %w", err)
	}

	return &Exporter{
		identity:    identity,
		vault:       v,
		ftcSettings: st,
	}, nil
}

// Returns the identity of the Outline user exporting the vault.
func (e *Exporter) Identity() schema.Identity {
	return e.identity
}

// Returns the vault being exported.
func (e *Exporter) Vault() *Vault {
	return e.vault
}

// Returns the Fast Text Color settings. If the plugin is not installed, this is nil.
func (e *Exporter) FTCSettings() *FTCSettings {
	return e.ftcSettings
}
