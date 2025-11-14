#!/bin/bash
# WMS Testing Script

echo "=== WMS Testing Helper ==="
echo ""

# Check if binary exists
if [ ! -f "./wms" ]; then
    echo "❌ Binary not found. Building..."
    go build -o wms ./cmd/wms
    if [ $? -eq 0 ]; then
        echo "✅ Build successful!"
    else
        echo "❌ Build failed!"
        exit 1
    fi
else
    echo "✅ Binary found: ./wms"
fi

echo ""
echo "=== Test Options ==="
echo "1. Run WMS normally"
echo "2. Run with API key from .env.example"
echo "3. Check config/env file locations"
echo "4. Clean config directory (reset)"
echo ""
read -p "Select option (1-4): " option

case $option in
    1)
        echo ""
        echo "Starting WMS..."
        echo "Press 'S' to open Settings and test API key input"
        echo "Try resizing your terminal to test responsive scaling!"
        echo ""
        ./wms
        ;;
    2)
        echo ""
        echo "Using API key from .env.example..."
        export WEATHER_API_KEY="fe5293644b554a78b4b03141250410"
        ./wms
        ;;
    3)
        CONFIG_DIR="$HOME/.config/wms"
        echo ""
        echo "Configuration locations:"
        echo "  Config file: $CONFIG_DIR/wms.toml"
        echo "  Env file:    $CONFIG_DIR/.env"
        echo ""

        if [ -f "$CONFIG_DIR/wms.toml" ]; then
            echo "✅ Config file exists"
            ls -lh "$CONFIG_DIR/wms.toml"
        else
            echo "❌ Config file not found (will be created on first run)"
        fi

        echo ""

        if [ -f "$CONFIG_DIR/.env" ]; then
            echo "✅ .env file exists"
            ls -lh "$CONFIG_DIR/.env"
            echo ""
            echo "Contents (API key should be masked):"
            cat "$CONFIG_DIR/.env"
        else
            echo "❌ .env file not found (will be created when you set API key)"
        fi
        ;;
    4)
        CONFIG_DIR="$HOME/.config/wms"
        echo ""
        read -p "⚠️  This will delete $CONFIG_DIR. Continue? (y/N): " confirm
        if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
            rm -rf "$CONFIG_DIR"
            echo "✅ Config directory cleaned!"
        else
            echo "Cancelled."
        fi
        ;;
    *)
        echo "Invalid option"
        exit 1
        ;;
esac
