#!/bin/bash

# Start the server in the background
DB_TYPE=sqlite PORT=3006 JWT_SECRET=your-super-secret-key-123 go run . &
SERVER_PID=$!

# Wait for the server to start
sleep 5

# Create screenshots directory if it doesn't exist
mkdir -p assets

# Take screenshots using Chrome in headless mode
# Login screen
google-chrome --headless --screenshot=assets/login.png --window-size=1200,800 http://localhost:3006/client/

# Chat interface (you'll need to manually log in and navigate to this page)
# google-chrome --headless --screenshot=assets/chat.png --window-size=1200,800 http://localhost:3006/client/

# User list
# google-chrome --headless --screenshot=assets/users.png --window-size=1200,800 http://localhost:3006/client/

# Kill the server
kill $SERVER_PID

echo "Screenshots generated in assets directory" 