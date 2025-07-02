package config

import (
	"fmt"
	"runtime"
)

// BinaryConfig holds configuration for external binaries
var BinaryConfig = map[string]BinaryInfo{
	"libreoffice": {
		Windows: BinaryPlatformConfig{
			URL:        "https://download.documentfoundation.org/libreoffice/stable/7.5.6/win/x86_64/LibreOffice_7.5.6_Win_x86-64.msi",
			BinaryPath: "program/soffice.exe",
		},
		MacOS: BinaryPlatformConfig{
			URL:        "https://download.documentfoundation.org/libreoffice/stable/7.5.6/mac/x86_64/LibreOffice_7.5.6_MacOS_x86-64.dmg",
			BinaryPath: "LibreOffice.app/Contents/MacOS/soffice",
		},
		Linux: BinaryPlatformConfig{
			URL:        "https://download.documentfoundation.org/libreoffice/stable/7.5.6/deb/x86_64/LibreOffice_7.5.6_Linux_x86-64_deb.tar.gz",
			BinaryPath: "LibreOffice_7.5.6.1_Linux_x86-64_deb/DEBS/desktop-integration/libreoffice7.5-debian-menus_7.5.6-1_all.deb",
		},
	},
	"imagemagick": {
		Windows: BinaryPlatformConfig{
			URL:        "https://imagemagick.org/archive/binaries/ImageMagick-7.1.1-15-portable-Q16-x64.zip",
			BinaryPath: "convert.exe",
		},
		MacOS: BinaryPlatformConfig{
			URL:        "https://imagemagick.org/archive/binaries/ImageMagick-x86_64-apple-darwin20.1.0.tar.gz",
			BinaryPath: "bin/convert",
		},
		Linux: BinaryPlatformConfig{
			URL:        "https://imagemagick.org/archive/binaries/ImageMagick-x86_64-pc-linux-gnu.tar.gz",
			BinaryPath: "bin/convert",
		},
	},
}

type BinaryInfo struct {
	Windows BinaryPlatformConfig
	MacOS   BinaryPlatformConfig
	Linux   BinaryPlatformConfig
}

type BinaryPlatformConfig struct {
	URL        string
	BinaryPath string
}

// GetBinaryConfig returns the appropriate binary configuration for the current platform
func GetBinaryConfig(name string) (BinaryPlatformConfig, error) {
	binary, exists := BinaryConfig[name]
	if !exists {
		return BinaryPlatformConfig{}, fmt.Errorf("binary %s not found", name)
	}

	switch runtime.GOOS {
	case "windows":
		return binary.Windows, nil
	case "darwin":
		return binary.MacOS, nil
	case "linux":
		return binary.Linux, nil
	default:
		return BinaryPlatformConfig{}, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}
