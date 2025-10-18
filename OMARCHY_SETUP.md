# WMS Setup for Omarchy

## Prerequisites

1. **Go Language**: Make sure Go 1.24+ is installed
   ```bash
   go version
   ```

2. **WeatherAPI Key**: Get a free API key from https://www.weatherapi.com/

## Setup Steps

### 1. Configure API Key
```bash
# Copy the example environment file
cp .env.example .env

# Edit .env and add your WeatherAPI key
nano .env  # or your preferred editor
```

### 2. Test the Application
```bash
# Test run to make sure everything works
go run ./cmd/wms
```

### 3. Add to Omarchy

#### Option A: Add as local application
```bash
# In Omarchy, add this directory as a local TUI app
# Point to: /home/traves/Development/WMS
# Use the omarchy.toml configuration file
```

#### Option B: Use the launcher script directly
```bash
# Run the Omarchy-compatible launcher
./wms-omarchy
```

## Configuration

### Command Line Options
- `-location "City Name"` - Set specific location
- `-units metric|imperial` - Set units
- `-time 12|24` - Set time format  
- `-compact` - Use compact display mode
- `-refresh N` - Set refresh interval in minutes

### Configuration File
The app creates a config file at `~/.config/wms/wms.toml` on first run.

## Usage in Omarchy

1. **Navigation**: 
   - 1,2,3: Switch directly to Weather/Moon/Solar views
   - Tab/Shift+Tab: Navigate between views
   - Arrow keys: Navigate within views
   - q or Ctrl+C: Quit

2. **Controls**:
   - r: Refresh data manually
   - u: Toggle temperature units (metric/imperial)
   - t: Toggle time format (12/24 hour)
   - s: Open settings menu

3. **Features**:
   - Weather view: Detailed weather with ASCII art (IP-based location)
   - Moon view: Current phase and illumination
   - Solar view: Sunrise/sunset times

4. **Auto-refresh**: Data updates every 5 minutes (configurable)

## Troubleshooting

- **"No weather data"**: Check your API key in .env file
- **Build errors**: Ensure Go 1.24+ is installed
- **Permission denied**: Run `chmod +x wms-omarchy`
- **Module errors**: Run `go mod tidy` to update dependencies