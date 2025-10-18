#!/bin/bash

# WMS Setup Script for Omarchy
echo "🌦️  WMS - Weather Management System Setup"
echo "========================================="

# Step 1: Check if .env file exists
if [[ ! -f ".env" ]]; then
    echo "📝 Setting up environment file..."
    cp .env.example .env
    echo "✅ Created .env file from template"
    echo ""
    echo "⚠️  IMPORTANT: You need to add your WeatherAPI key to the .env file"
    echo "   1. Get a free API key from: https://www.weatherapi.com/"
    echo "   2. Edit .env file and replace WEATHER_API_KEY=\"\" with your key"
    echo ""
    read -p "Press Enter when you've added your API key, or Ctrl+C to exit..."
else
    echo "✅ .env file already exists"
fi

# Step 2: Build the application
echo "🔨 Building WMS application..."
if go build -o wms ./cmd/wms; then
    echo "✅ Build successful!"
else
    echo "❌ Build failed. Please check for errors above."
    exit 1
fi

# Step 3: Test the application
echo "🧪 Testing application..."
if ./wms --help > /dev/null 2>&1; then
    echo "✅ Application test successful!"
else
    echo "❌ Application test failed."
    exit 1
fi

echo ""
echo "🎉 Setup complete! Your WMS TUI is ready for Omarchy!"
echo ""
echo "📋 Next steps for Omarchy:"
echo "   1. In Omarchy, add this directory as a TUI application:"
echo "      Path: $(pwd)"
echo "   2. Use the omarchy.toml configuration file provided"
echo "   3. Or run directly with: ./wms-omarchy"
echo ""
echo "🎮 Controls:"
echo "   • Tab/1-3: Switch between Weather/Moon/Solar views"
echo "   • q/Ctrl+C: Quit"
echo "   • Arrow keys: Navigate"
echo ""
echo "🔧 Configuration:"
echo "   • Config file: ~/.config/wms/wms.toml (created on first run)"
echo "   • Command options: ./wms --help"