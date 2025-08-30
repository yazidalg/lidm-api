#!/bin/bash

# Script to seed video quizzes for LIDM API
# Usage: ./seed_video_quizzes.sh [command] [module_name]

set -e

echo "🎥 LIDM API - Video Quizzes Seeding Tool"
echo "======================================="

# Change to the project directory
cd "$(dirname "$0")"

# Function to show usage
show_usage() {
    echo "Usage:"
    echo "  ./seed_video_quizzes.sh all                              - Seed video quizzes for all modules"
    echo "  ./seed_video_quizzes.sh module \"Module Name\"             - Seed video quizzes for specific module"
    echo "  ./seed_video_quizzes.sh fotosintesis                     - Seed video quizzes for all fotosintesis modules"
    echo "  ./seed_video_quizzes.sh clear                            - Clear all video quizzes"
    echo "  ./seed_video_quizzes.sh summary                          - Show summary of current video quizzes"
    echo ""
    echo "Examples:"
    echo "  ./seed_video_quizzes.sh module \"Fotosintesis - Dasar\""
    echo "  ./seed_video_quizzes.sh module \"Fotosintesis - Lanjutan\""
    echo "  ./seed_video_quizzes.sh module \"Fotosintesis - Eksperimen\""
}

# Check if command is provided
if [ $# -eq 0 ]; then
    show_usage
    exit 1
fi

COMMAND=$1

case $COMMAND in
    "all")
        echo "🚀 Seeding video quizzes for all modules..."
        cd cmd/seed-video-quizzes && go run main.go all
        ;;
    "module")
        if [ $# -lt 2 ]; then
            echo "❌ Error: module name required"
            echo "Usage: ./seed_video_quizzes.sh module \"Module Name\""
            exit 1
        fi
        MODULE_NAME=$2
        echo "🚀 Seeding video quizzes for module: $MODULE_NAME"
        cd cmd/seed-video-quizzes && go run main.go module "$MODULE_NAME"
        ;;
    "fotosintesis")
        echo "🌱 Seeding video quizzes for all fotosintesis modules..."
        cd cmd/seed-video-quizzes && go run main.go fotosintesis
        ;;
    "clear")
        echo "🗑️  Clearing all video quizzes..."
        read -p "Are you sure you want to clear all video quizzes? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            cd cmd/seed-video-quizzes && go run main.go clear
        else
            echo "Operation cancelled."
        fi
        ;;
    "summary")
        echo "📊 Showing video quizzes summary..."
        cd cmd/seed-video-quizzes && go run main.go summary
        ;;
    *)
        echo "❌ Unknown command: $COMMAND"
        show_usage
        exit 1
        ;;
esac

echo "✅ Operation completed!"
