#!/bin/bash

# Get the directory where the script is located
SCRIPT_DIR=$(dirname "$(readlink -f "\$0")")

# Service name
SERVICE_NAME="relayapi-server"

# Path to the executable (assuming relayapi-server is in the same directory as the script)
EXECUTABLE="$SCRIPT_DIR/relayapi-server"

# Automatically get the current user
USER=$(whoami)

# Path to the service file
SERVICE_FILE="/etc/systemd/system/$SERVICE_NAME.service"

# Check if the executable exists and has execute permissions
echo "Checking if executable '$EXECUTABLE' exists and is executable..."
if [ ! -x "$EXECUTABLE" ]; then
  echo "Error: Executable '$EXECUTABLE' does not exist or does not have execute permissions."
  exit 1
fi
echo "Executable check passed."

# Create the systemd service file (requires sudo)
echo "Creating systemd service file: '$SERVICE_FILE'..."
SERVICE_CONTENT="
[Unit]
Description=Relay API Server
After=network.target

[Service]
WorkingDirectory=$SCRIPT_DIR
User=$USER
ExecStart=$EXECUTABLE
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
"

if sudo sh -c "echo \"$SERVICE_CONTENT\" > '$SERVICE_FILE'"; then
  echo "Service file created successfully."
  if [ ! -f "$SERVICE_FILE" ]; then
    echo "Error: Service file was not created at '$SERVICE_FILE'."
    exit 1
  fi
else
  echo "Error: Failed to create service file."
  exit 1
fi

# Set permissions for the service file (requires sudo)
echo "Setting permissions for service file..."
if sudo chmod 644 "$SERVICE_FILE"; then
  echo "Service file permissions set successfully."
else
  echo "Error: Failed to set service file permissions."
  exit 1
fi

# Reload systemd configuration (requires sudo)
echo "Reloading systemd configuration..."
if sudo systemctl daemon-reload; then
  echo "Systemd configuration reloaded successfully."
else
  echo "Error: Failed to reload systemd configuration."
  exit 1
fi

# Enable the service to start on boot (requires sudo)
echo "Enabling service to start on boot..."
if sudo systemctl enable "$SERVICE_NAME"; then
  echo "Service enabled successfully."
else
  echo "Error: Failed to enable service."
  exit 1
fi

# Start the service (requires sudo)
echo "Starting the service..."
if sudo systemctl start "$SERVICE_NAME"; then
  echo "Service started successfully."
else
  echo "Error: Failed to start service."
  exit 1
fi

echo "Service '$SERVICE_NAME' has been registered and started."
