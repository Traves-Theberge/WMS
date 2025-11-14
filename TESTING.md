# WMS Testing Checklist

## Setup
1. Build the application: `./test-wms.sh` or `go build -o wms ./cmd/wms`

## Test 1: Responsive UI Scaling ✅

### Small Terminal (80x24)
- [ ] Resize terminal to 80 columns x 24 rows
- [ ] Run `./wms`
- [ ] Verify border fits within terminal
- [ ] Verify no content is cut off
- [ ] Switch between tabs (1, 2, 3)
- [ ] All tabs should display correctly

### Large Terminal (150x40)
- [ ] Resize terminal to 150 columns x 40 rows
- [ ] Verify content is centered
- [ ] Verify border scales appropriately
- [ ] Content should still be readable

### Dynamic Resize
- [ ] Start WMS in normal terminal
- [ ] Resize terminal while running (make smaller)
- [ ] Verify UI adapts dynamically
- [ ] Resize terminal larger
- [ ] Verify UI scales back up

## Test 2: API Key Management ✅

### Setting API Key via Settings Menu
1. [ ] Run `./wms`
2. [ ] Press `S` to open Settings
3. [ ] Navigate to "Set API Key" using arrow keys (↓)
4. [ ] Press `Enter` to start editing
5. [ ] Verify input screen shows:
   - Prompt: "Enter your WeatherAPI key:"
   - Help: Link to weatherapi.com/signup.aspx
   - Input field with cursor: `> █`
6. [ ] Type test key: `test123456789`
7. [ ] Verify input is masked (shows `***********`)
8. [ ] Press `Enter` to save
9. [ ] Verify status message: "API key saved!"
10. [ ] Return to Settings (`S`)
11. [ ] Verify API key shows as `****6789` (last 4 chars)

### Cancel API Key Input
1. [ ] Press `S` to open Settings
2. [ ] Navigate to "Set API Key"
3. [ ] Press `Enter`
4. [ ] Start typing
5. [ ] Press `Esc` to cancel
6. [ ] Verify status message: "Cancelled"
7. [ ] Verify no changes were saved

### Verify Secure Storage
1. [ ] Run: `./test-wms.sh` and select option 3
2. [ ] Verify `.env` file exists at `~/.config/wms/.env`
3. [ ] Check file permissions: `ls -l ~/.config/wms/.env`
4. [ ] Verify permissions are `-rw-------` (0600)
5. [ ] Check file contents: `cat ~/.config/wms/.env`
6. [ ] Verify it contains:
   ```
   # WMS Environment Variables
   # This file contains sensitive API keys - do not commit to version control

   WEATHER_API_KEY=your_key_here
   ```

## Test 3: API Key Loading ✅

### From .env in Config Directory
1. [ ] Set API key via settings (previous test)
2. [ ] Quit WMS (`Q`)
3. [ ] Restart WMS: `./wms`
4. [ ] Press `S` to check Settings
5. [ ] Verify API key is loaded (shows `****XXXX`)
6. [ ] Weather data should load if key is valid

### From Environment Variable
1. [ ] Run: `export WEATHER_API_KEY="fe5293644b554a78b4b03141250410"`
2. [ ] Run: `./wms`
3. [ ] Verify weather data loads
4. [ ] Weather tab should show current conditions

### Priority Test (Config dir vs Current dir)
1. [ ] Create `.env` in current directory with different key
2. [ ] Create `.env` in `~/.config/wms/` with valid key
3. [ ] Run WMS
4. [ ] Verify config directory `.env` takes priority

## Test 4: Settings Menu Navigation ✅

### Menu Options
- [ ] Press `S` to open Settings
- [ ] Verify 4 options are visible:
  1. Location Mode: ip
  2. Set Location: [location]
  3. Set API Key: [masked or "Not Set"]
  4. Save and Exit
- [ ] Use `↑` and `↓` to navigate
- [ ] Verify cursor (`>`) moves correctly
- [ ] Verify wrapping (down from option 4 goes to option 1)
- [ ] Press `Esc` to cancel without saving
- [ ] Verify returns to Weather view

### Help Text
- [ ] In Settings menu, verify help text is visible:
  - "Get your free API key at: https://www.weatherapi.com/signup.aspx"
  - "(Use ↑/↓ to navigate, Enter to select, Esc to cancel)"

## Test 5: Complete User Flow ✅

### New User Setup
1. [ ] Clean config: `./test-wms.sh` and select option 4
2. [ ] Run WMS: `./wms`
3. [ ] Verify warning: "weather_api_key is required"
4. [ ] Press `S` for Settings
5. [ ] Navigate to "Set API Key"
6. [ ] Enter key from `.env.example`: `fe5293644b554a78b4b03141250410`
7. [ ] Navigate to "Save and Exit"
8. [ ] Press `Enter`
9. [ ] Verify weather data loads
10. [ ] Quit and restart
11. [ ] Verify API key persists

## Test 6: Edge Cases ✅

### Empty API Key
- [ ] Try setting empty API key (just press Enter without typing)
- [ ] Verify it saves empty value
- [ ] Settings should show "Not Set"

### Very Long API Key
- [ ] Try entering 100+ character string
- [ ] Verify it's accepted
- [ ] Verify masking works correctly
- [ ] Verify display shows last 4 chars

### Special Characters
- [ ] Try API key with special chars: `abc-123_DEF!@#$`
- [ ] Verify it's stored correctly
- [ ] Check `.env` file contents

### Terminal Too Small
- [ ] Resize to very small (40x10)
- [ ] Verify WMS handles gracefully
- [ ] No crashes or rendering errors

## Success Criteria
All checkboxes should be checked (✅) with no errors or unexpected behavior.

## Known Issues / Notes
- Document any issues found during testing here
