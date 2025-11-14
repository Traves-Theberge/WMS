# Changelog

All notable changes to WMS (Weather Management System) will be documented in this file.

## [1.1.0] - 2025-11-13

### Added
- **API Key Management in Settings**
  - In-app API key input with masked display (shows asterisks)
  - Secure storage in `~/.config/wms/.env` with 0600 permissions
  - Paste support (Ctrl+Shift+V) for easy key entry
  - Real-time connection testing after saving key
  - Visual display of masked API key in settings (shows last 4 chars)
  - Help text with link to WeatherAPI.com signup

- **Responsive UI Scaling**
  - Dynamic content scaling to fit terminal size
  - Card borders and content adapt to available space
  - Centered content layout (horizontal and vertical)
  - Minimum dimensions to maintain readability
  - Works across all terminal sizes

- **Enhanced Input Fields**
  - Paste support for API key and location inputs
  - Keyboard shortcuts: Ctrl+U (clear line), Ctrl+W (delete word)
  - Multi-character input handling for pasted content
  - ASCII character filtering (32-126) for clean input

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
- **API Key Persistence**: API key now properly persists across sessions
- **API Key Cleaning**: Automatically removes quotes, brackets, and whitespace from pasted keys
- **Content Centering**: Weather/Moon/Solar content now properly centered in cards
- **Location Detection**: Fixed weather data to use IP-based location when configured
- **Moon Data**: Resolved "Moon data unavailable" issue with API fallback
- **Navigation**: Fixed tab switching to work with both Tab keys and number keys
- **Time Format**: Fixed 12-hour time format display and cycling
- **Configuration Loading**: Improved .env file loading and API key validation

### Changed
- **Settings Menu**: Now has 4 options (added "Set API Key")
- **API Key Storage**: Moved from .env in project root to ~/.config/wms/.env
- **Default Configuration**: Changed default location to empty string to enable IP detection
- **API Timeouts**: Reduced moon API timeout from 10s to 5s for better UX
- **Build Process**: Enhanced build system with proper main entry point
- **Error Handling**: Improved error messages and user feedback with detailed error descriptions
- **Code Organization**: Cleaned up project structure
- **Card Padding**: Increased from (1,2) to (2,4) for better spacing

### Security
- **API Key Security**: Keys stored with 0600 permissions (owner read/write only)
- **Masked Input**: API keys shown as asterisks during entry
- **Masked Display**: Settings shows only last 4 characters of key
- **Warning Headers**: .env file includes warning about not committing to VCS

### Technical Improvements
- **Main Entry Point**: Created proper `cmd/wms/main.go` with configuration support
- **Fallback Systems**: Implemented robust fallback mechanisms for external APIs
- **SaveAPIKey()**: Returns cleaned API key for immediate use
- **Input Handling**: Improved paste detection and multi-character input processing
- **Config Loading**: .env loaded from config directory first, then current directory

### Documentation
- **Keyboard Shortcuts**: Updated all help text with current key bindings
- **Configuration**: Documented all configuration options and modes
- **Troubleshooting**: Added common issues and solutions
- **Testing Guide**: Created TESTING.md with comprehensive test checklist

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