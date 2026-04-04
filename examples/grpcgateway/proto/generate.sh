#!/bin/bash

# Generate Go code from proto files
echo "Generating Go code from proto files..."

# Create gen directory if it doesn't exist
if [ -d "gen" ]; then
    rm -rf gen
fi

# Generate code using buf
buf generate

echo "Code generation complete!"
