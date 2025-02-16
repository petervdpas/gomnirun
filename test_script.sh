#!/bin/bash

# Colors
GREEN="\e[32m"
BLUE="\e[34m"
BOLD="\e[1m"
RESET="\e[0m"

echo -e "${GREEN}🚀 Test Script Executed!${RESET}"

# Default values
VAR1=""
VAR2=""

# Parse arguments
for arg in "$@"; do
    case $arg in
        -GreetingPartOne=*)
            VAR1="${arg#*=}" # Extract value after '='
            ;;
        -GreetingPartTwo=*)
            VAR2="${arg#*=}" # Extract value after '='
            ;;
    esac
done

echo -e "${BOLD}GreetingPartOne:${RESET} $VAR1"
echo -e "${BOLD}GreetingPartTwo:${RESET} $VAR2"

# Create a cool ASCII art message
echo -e "\n${BLUE}Creating a cool message...${RESET}"
sleep 1

MESSAGE="${VAR1} ${VAR2}"
echo -e "${GREEN}✨ Here is your message: ${BOLD}${MESSAGE}${RESET}"

# Simulate processing
echo -e "\nProcessing..."
sleep 1

# Generate a simple animated effect
for i in {1..3}; do
    echo -n "🔄 Loading."
    sleep 0.5
    echo -n "."
    sleep 0.5
    echo -n "."
    sleep 0.5
    echo -e "."
done

echo -e "\n${GREEN}✅ Done!${RESET}"
