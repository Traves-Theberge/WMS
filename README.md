# WMS - Weather, Moon, & Solar

A comprehensive, terminal-based weather dashboard inspired by the aesthetics and functionality of [Stormy](https://github.com/dpr-1/stormy'). WMS provides detailed weather, moon phase, and solar information in a clean, elegant, and highly configurable TUI (Text-based User Interface).

![WMS Screenshot](https://user-images.githubusercontent.com/12345/67890.png) <!-- Placeholder -->

## Features

- **Tabbed Interface**: Switch between three distinct views:
    - **Weather**: A detailed, Stormy-style weather display with ASCII art icons.
    - **Moon**: Information about the current moon phase, illumination, and next phase.
    - **Solar**: Sunrise, sunset, and daylight duration information.
- **Dynamic ASCII Art**: Weather icons change based on the conditions, and the solar tab shows a sun during the day and a moon at night.
- **Highly Configurable**: Customize units, time format, and more using a simple TOML configuration file or command-line flags.
- **Automatic Location Detection**: If no location is specified, WMS will attempt to determine your location automatically based on your IP address.
- **Real-time Updates**: Weather and time information updates automatically.

## Installation

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/your-username/wms.git
    cd wms
    ```

2.  **Build the application**:
    ```bash
    go build -o wms ./cmd/wms
    ```

## Usage

Run the application from your terminal:

```bash
./wms
```

### Command-Line Flags

You can override the default configuration using command-line flags:

| Flag          | Description                  | Default |
|---------------|------------------------------|---------|
| `-location`   | Location to get weather for  | ""      |
| `-units`      | Units (metric, imperial)     | "metric"|
| `-time`       | Time format (12, 24)         | "24"    |
| `-compact`    | Compact display mode         | false   |
| `-refresh`    | Refresh interval in minutes  | 5       |
| `-help`       | Show help                    | false   |

**Example**:

```bash
./wms -location "New York" -units "imperial"
```

## Configuration

WMS uses a `.env` file to manage API keys and a `wms.toml` file for other settings.

### `.env` File

Create a `.env` file in the root of the project and add your API key:

```
# WMS Environment Variables

# You can get a free API key from https://www.weatherapi.com/
WEATHER_API_KEY="your-weather-api-key"
```

### `wms.toml`

The `wms.toml` file is located at:
- **Linux/macOS**: `~/.config/wms/wms.toml`
- **Windows**: `%APPDATA%\wms\wms.toml`

The application will create a default configuration file on the first run.

```toml
# Weather settings
weather_provider = "WeatherAPI" #
location = "" # e.g., "London"

# Display settings
units = "metric" # "metric" or "imperial"
time_format = "24" # "12" or "24"
use_colors = true
compact = false
show_city_name = true

# Update settings
refresh_interval = 5 # minutes
```

## Keyboard Shortcuts

| Key      | Action                       |
|----------|------------------------------|
| `1` / `w`| Switch to Weather Tab        |
| `2` / `m`| Switch to Moon Tab           |
| `3`      | Switch to Solar Tab          |
| `Tab`    | Cycle through tabs           |
| `u`      | Toggle units and time format |
| `s`      | Open settings menu           |
| `r`      | Refresh all data             |
| `q`      | Quit the application         |

## Inspiration

WMS draws inspiration from several fantastic open-source projects. A special thanks to the creators and maintainers of:

- **[chubin/wttr.in](https://github.com/chubin/wttr.in)**: The original console-based weather service that set the standard for terminal weather reports.
- **[dpr-1/stormy](https://github.com/dpr-1/stormy)**: Another excellent Go-based weather tool that provided valuable insights and ideas.
- **[liveslol/rainy](https://github.com/liveslol/rainy)**: A beautiful terminal-based weather application that inspired the UI design.

## Dependencies
- [bubbletea](https://github.com/charmbracelet/bubbletea)
- [lipgloss](https://github.com/charmbracelet/lipgloss)
- [toml](https://github.com/BurntSushi/toml)
- [godotenv](https://github.com/joho/godotenv) 