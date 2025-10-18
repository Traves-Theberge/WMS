# Changelog

All notable changes to WMS (Weather Management System) will be documented in this file.

## [1.1.0] - 2025-10-17

### Added
- **Omarchy Integration Support**
  - Created `omarchy/` folder with all integration files
  - Added `omarchy.toml` configuration file
  - Added `wms-omarchy` launcher script
  - Added `setup-omarchy.sh` automated setup script
  - Added comprehensive `OMARCHY_SETUP.md` documentation

- **Enhanced Navigation**
  - Added number key navigation (1, 2, 3) for direct tab switching
  - Added Shift+Tab support for reverse tab navigation
  - Improved keyboard shortcuts with separate U and T keys

- **Improved Configuration**
  - Added 4-state unit/time cycling: Metric 24h → Metric 12h → Imperial 24h → Imperial 12h
  - Separated time format toggle (T key) from units toggle (U key)
  - Enhanced help text with complete keyboard shortcut documentation

- **Robust Location Detection**
  - Implemented proper IP-based location detection by default
  - Fixed location mode logic to respect "ip" vs "manual" settings
  - Automatic fallback to IP detection when location is empty

- **Moon Phase Resilience**
  - Added offline moon phase calculation as fallback
  - Implemented local astronomical calculations when API is unavailable
  - Reduced API timeout to 5 seconds for better responsiveness
  - Added graceful error handling for moon data API failures

### Fixed
- **Location Detection**: Fixed weather data to use IP-based location when configured
- **Moon Data**: Resolved "Moon data unavailable" issue with API fallback
- **Navigation**: Fixed tab switching to work with both Tab keys and number keys
- **Time Format**: Fixed 12-hour time format display and cycling
- **Configuration Loading**: Improved .env file loading and API key validation

### Changed
- **Default Configuration**: Changed default location to empty string to enable IP detection
- **API Timeouts**: Reduced moon API timeout from 10s to 5s for better UX
- **Build Process**: Enhanced build system with proper main entry point
- **Error Handling**: Improved error messages and user feedback
- **Code Organization**: Cleaned up project structure with dedicated Omarchy folder

### Technical Improvements
- **Main Entry Point**: Created proper `cmd/wms/main.go` with configuration support
- **Fallback Systems**: Implemented robust fallback mechanisms for external APIs
- **Path Resolution**: Fixed absolute path handling for Omarchy integration
- **Executable Management**: Optimized launcher to use compiled binary instead of `go run`

### Documentation
- **Setup Guides**: Created comprehensive setup documentation for Omarchy
- **Keyboard Shortcuts**: Updated all help text with current key bindings
- **Configuration**: Documented all configuration options and modes
- **Troubleshooting**: Added common issues and solutions

## [1.0.0] - Initial Release

### Features
- Weather dashboard with real-time data
- Moon phase tracking and display
- Solar information (sunrise/sunset)
- Tabbed TUI interface
- Configurable units and time formats
- ASCII art weather icons
- Automatic data refresh

---

**Note**: This changelog follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) format.