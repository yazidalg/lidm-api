#!/bin/bash

echo "🌱 Fotosintesis Data Seeder"
echo "=========================="

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo "❌ Error: .env file not found!"
    echo "Please make sure you have .env file with database configuration"
    exit 1
fi

echo "📝 Loading environment variables..."
source .env

echo "🔄 Running Fotosintesis seeder..."
go run cmd/seed-fotosintesis/main.go

if [ $? -eq 0 ]; then
    echo "✅ Fotosintesis data seeding completed successfully!"
    echo ""
    echo "📚 Created content:"
    echo "   - 1 Module: Belajar Fotosintesis - Kelas 4 SD"
    echo "   - 5 SubMaterials dengan urutan pembelajaran"
    echo "   - 2 Video Materials dengan 5 Interactive Video Quizzes"
    echo "   - 1 AR Experience (AR Link: https://asblr.com/e4Z066)"
    echo "   - 6 Flashcards"
    echo "   - 30+ Prequizzes (10 per SubMaterial)"
    echo "   - 7 Module Quiz Questions"
    echo ""
    echo "🎯 Learning Path Coordinates:"
    echo "   - SubMaterial 1: Video Intro (120, 50)"
    echo "   - SubMaterial 2: Quiz Dasar (200, 180)"
    echo "   - SubMaterial 3: AR Lab (120, 310)"
    echo "   - SubMaterial 4: Video Proses (200, 440)"
    echo "   - SubMaterial 5: Flashcards & Quiz Final (120, 570)"
    echo ""
    echo "👥 Dummy Users Created:"
    echo "   - 5 users with role 'user' (Andi, Sari, Budi, Maya, Riko)"
    echo "   - All users have completed prequizzes and video quizzes"
    echo ""
    echo "� User Progress Stats:"
    echo "   - PreQuiz Answers: 150 total (30 per user)"
    echo "   - Video Quiz Answers: 25 total (5 per user)"
    echo "   - Accuracy rates: 70-90% on prequizzes, 100% on video quizzes"
    echo ""
    echo "�🌐 Test with these curl commands:"
    echo "   📚 Get all modules:"
    echo "      curl -X GET http://localhost:3000/module/all"
    echo ""
    echo "   👥 Get all users:"
    echo "      curl -X GET http://localhost:3000/user/all"
    echo ""
    echo "   📝 Get prequiz user answers:"
    echo "      curl -X GET 'http://localhost:3000/prequiz/user-answers?user_id=8'"
    echo ""
    echo "   🎥 Get video quiz user answers:"
    echo "      curl -X GET 'http://localhost:3000/video-quiz/user-answers?user_id=8'"
    echo ""
    echo "   📊 Get user progress by submaterial:"
    echo "      curl -X GET 'http://localhost:3000/progress/user/8'"
else
    echo "❌ Error: Seeding failed!"
    exit 1
fi
