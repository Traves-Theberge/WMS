# Omarchy Integration

This folder contains all files needed to integrate WMS with Omarchy.

## Files:

- **`omarchy.toml`** - Main Omarchy configuration file
- **`wms-omarchy`** - Launcher script (executable)
- **`setup-omarchy.sh`** - Complete setup script
- **`OMARCHY_SETUP.md`** - Detailed documentation

## Quick Setup:

1. **Run setup script:**
   ```bash
   ./omarchy/setup-omarchy.sh
   ```

2. **Add to Omarchy:**
   - Path: `/home/traves/Development/WMS`
   - Config: Use `omarchy/omarchy.toml`
   - Launch: `/home/traves/Development/WMS/wms`

3. **Or use launcher script:**
   ```bash
   ./omarchy/wms-omarchy
   ```

## Controls in WMS:
- `1,2,3`: Switch tabs (Weather/Moon/Solar)
- `U`: Cycle units/time (Metric 24h → Metric 12h → Imperial 24h → Imperial 12h)
- `T`: Toggle time format only
- `R`: Refresh data
- `Q`: Quit

Your WMS TUI is ready for Omarchy! 🌦️🌙☀️