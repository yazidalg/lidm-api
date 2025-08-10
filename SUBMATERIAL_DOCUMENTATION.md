# SubMaterial Documentation - Complete Guide

## Overview

SubMaterials adalah komponen pembelajaran yang berada di dalam Module. Berbeda dengan Lesson yang merupakan struktur lama, SubMaterial adalah struktur konten yang lebih modern dan fleksibel yang mendukung berbagai jenis konten pembelajaran.

## SubMaterial Structure

### Model Structure
```go
type SubMaterial struct {
    ID          uint      `json:"id"`
    ModuleID    uint      `json:"module_id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Order       int       `json:"order"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    // Relationships - Different content types
    VideoMaterial  VideoMaterial  `json:"video_material,omitempty"`
    QuizQuestions  []QuizQuestion `json:"quiz_questions,omitempty"`
    ARExperiment   *ARExperiment  `json:"ar_experiment,omitempty"`
    Flashcards     []Flashcard    `json:"flashcards,omitempty"`
}
```

### Content Types

1. **Video Material**: Video pembelajaran dengan durasi dan URL
2. **Quiz Questions**: Pertanyaan kuis interaktif
3. **AR Experiment**: Pengalaman Augmented Reality
4. **Flashcards**: Kartu belajar untuk menghafal

## How to Access SubMaterials

### 1. Get All Modules with SubMaterials
```bash
curl -X GET http://localhost:3000/module/all \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response includes SubMaterials:**
```json
{
  "message": "Modules retrieved successfully",
  "data": [
    {
      "ID": 1,
      "title": "Introduction to Programming",
      "description": "Basic programming concepts",
      "icon": "/uploads/icons/module1.png",
      "thumbnail": "/uploads/thumbnails/module1.jpg",
      "sub_materials": [
        {
          "ID": 1,
          "module_id": 1,
          "title": "Welcome Video",
          "description": "Introduction video for the module",
          "order": 1,
          "video_material": {
            "url": "https://example.com/video1.mp4",
            "duration": 300,
            "thumbnail": "https://example.com/thumb1.jpg"
          }
        },
        {
          "ID": 2,
          "module_id": 1,
          "title": "Basic Concepts Quiz",
          "description": "Test your understanding",
          "order": 2,
          "quiz_questions": [
            {
              "ID": 1,
              "sub_material_id": 2,
              "question": "What is a variable?",
              "question_type": "multiple_choice",
              "options": ["A storage location", "A function", "A loop", "A condition"],
              "correct_answer": "A storage location"
            }
          ]
        },
        {
          "ID": 3,
          "module_id": 1,
          "title": "AR Lab Experience",
          "description": "Virtual laboratory",
          "order": 3,
          "ar_experiment": {
            "ID": 1,
            "name": "Chemistry Lab",
            "description": "Virtual chemistry experiments",
            "scene_url": "https://example.com/ar-scene.json"
          }
        },
        {
          "ID": 4,
          "module_id": 1,
          "title": "Key Terms",
          "description": "Important vocabulary",
          "order": 4,
          "flashcards": [
            {
              "ID": 1,
              "sub_material_id": 4,
              "front": "Variable",
              "back": "A storage location with a name and value"
            },
            {
              "ID": 2,
              "sub_material_id": 4,
              "front": "Function",
              "back": "A reusable block of code"
            }
          ]
        }
      ]
    }
  ]
}
```

### 2. Get Specific Module with SubMaterials
```bash
curl -X GET http://localhost:3000/module/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Returns detailed module with all SubMaterials and their complete content.**

## SubMaterial Types in Detail

### 1. Video Material
```json
{
  "ID": 1,
  "title": "Introduction Video",
  "description": "Welcome to the course",
  "order": 1,
  "video_material": {
    "url": "https://example.com/video.mp4",
    "duration": 300,
    "thumbnail": "https://example.com/thumbnail.jpg",
    "quality": "1080p",
    "size": 50000000
  }
}
```

### 2. Quiz Questions
```json
{
  "ID": 2,
  "title": "Knowledge Check",
  "description": "Test your understanding",
  "order": 2,
  "quiz_questions": [
    {
      "ID": 1,
      "question": "What is object-oriented programming?",
      "question_type": "multiple_choice",
      "options": [
        "A programming paradigm",
        "A data structure", 
        "An algorithm",
        "A database"
      ],
      "correct_answer": "A programming paradigm",
      "explanation": "OOP is a programming paradigm based on objects and classes"
    }
  ]
}
```

### 3. AR Experiment
```json
{
  "ID": 3,
  "title": "Virtual Laboratory",
  "description": "Hands-on AR experience",
  "order": 3,
  "ar_experiment": {
    "ID": 1,
    "name": "Chemistry Lab",
    "description": "Virtual chemistry experiments",
    "scene_url": "https://example.com/ar-scene.json",
    "instructions": "Point camera at the marker",
    "duration": 600
  }
}
```

### 4. Flashcards
```json
{
  "ID": 4,
  "title": "Key Vocabulary",
  "description": "Important terms to remember",
  "order": 4,
  "flashcards": [
    {
      "ID": 1,
      "front": "API",
      "back": "Application Programming Interface",
      "difficulty": "easy",
      "category": "programming"
    },
    {
      "ID": 2,
      "front": "JSON",
      "back": "JavaScript Object Notation",
      "difficulty": "easy",
      "category": "data-format"
    }
  ]
}
```

## RAG Integration for SubMaterials

When accessing modules with SubMaterials, the system tracks enhanced metadata for AI/RAG analysis:

### Activity Tracking
```json
{
  "activity_type": "lihat_modul",
  "description": "Melihat modul spesifik",
  "metadata": {
    "path": "/module/1",
    "method": "GET",
    "module_id": "1",
    "action": "view_specific_module",
    "content_type": "module_detail",
    "learning_activity": true,
    "knowledge_area": "module_content",
    "user_intent": "study_module",
    "session_context": "structured_learning",
    "engagement_type": "curriculum_following",
    "includes_sub_materials": true,
    "content_structure": "module_with_sub_materials"
  }
}
```

### RAG Endpoint Enhancement
The `/user-activity/for-rag` endpoint now includes SubMaterial information:

```bash
curl -X GET http://localhost:3000/user-activity/for-rag?limit=10
```

**Enhanced Response for Module Activities:**
```json
{
  "data": {
    "activities": [
      {
        "id": 25,
        "activity_type": "lihat_modul",
        "description": "Melihat modul spesifik",
        "learning_category": "curriculum_exploration",
        "module_details": {
          "id": 1,
          "title": "Introduction to Programming",
          "description": "Basic programming concepts",
          "lessons_count": 5,
          "sub_materials_count": 8,
          "sub_materials": [
            {
              "id": 1,
              "title": "Welcome Video",
              "description": "Introduction video",
              "order": 1
            },
            {
              "id": 2,
              "title": "Basic Quiz",
              "description": "Knowledge check",
              "order": 2
            }
          ]
        },
        "content_structure": "module_with_sub_materials",
        "includes_sub_materials": true
      }
    ]
  }
}
```

## Learning Path Structure

### Modules vs Lessons vs SubMaterials

1. **Modules**: Main learning units (current structure)
   - Contains SubMaterials (modern approach)
   - Contains Lessons (legacy, for backward compatibility)

2. **SubMaterials**: Individual learning components within modules
   - Video content
   - Interactive quizzes
   - AR experiences
   - Flashcards
   - Ordered sequence (1, 2, 3, ...)

3. **Lessons**: Legacy structure
   - Simple text-based content
   - Being phased out in favor of SubMaterials

### Recommended Content Structure
```
Module: "Introduction to Programming"
├── SubMaterial 1: Welcome Video (order: 1)
├── SubMaterial 2: Reading Material (order: 2)
├── SubMaterial 3: Interactive Quiz (order: 3)
├── SubMaterial 4: AR Lab Experience (order: 4)
├── SubMaterial 5: Flashcards Review (order: 5)
└── SubMaterial 6: Final Assessment (order: 6)
```

## Future SubMaterial Endpoints

**Note**: These endpoints would need to be implemented for direct SubMaterial management:

### Get SubMaterials for a Module
```bash
curl -X GET http://localhost:3000/module/1/sub-materials \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Create SubMaterial (Admin)
```bash
curl -X POST http://localhost:3000/module/1/sub-materials \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "New Video Lesson",
    "description": "Advanced concepts",
    "order": 5,
    "type": "video",
    "video_material": {
      "url": "https://example.com/video.mp4",
      "duration": 480
    }
  }'
```

### Update SubMaterial (Admin)
```bash
curl -X PUT http://localhost:3000/sub-material/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Video Title",
    "description": "Updated description",
    "order": 3
  }'
```

### Delete SubMaterial (Admin)
```bash
curl -X DELETE http://localhost:3000/sub-material/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Add Video to SubMaterial (Admin)
```bash
curl -X POST http://localhost:3000/sub-material/1/video \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com/video.mp4",
    "duration": 300,
    "thumbnail": "https://example.com/thumb.jpg"
  }'
```

### Add Quiz Questions to SubMaterial (Admin)
```bash
curl -X POST http://localhost:3000/sub-material/1/quiz-questions \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "questions": [
      {
        "question": "What is a variable?",
        "question_type": "multiple_choice",
        "options": ["Storage", "Function", "Loop", "Condition"],
        "correct_answer": "Storage",
        "explanation": "A variable stores data"
      }
    ]
  }'
```

### Add AR Experiment to SubMaterial (Admin)
```bash
curl -X POST http://localhost:3000/sub-material/1/ar-experiment \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Virtual Lab",
    "description": "Chemistry experiments",
    "scene_url": "https://example.com/ar-scene.json",
    "instructions": "Point camera at marker"
  }'
```

### Add Flashcards to SubMaterial (Admin)
```bash
curl -X POST http://localhost:3000/sub-material/1/flashcards \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "flashcards": [
      {
        "front": "API",
        "back": "Application Programming Interface",
        "difficulty": "easy"
      },
      {
        "front": "REST",
        "back": "Representational State Transfer",
        "difficulty": "medium"
      }
    ]
  }'
```

## Best Practices

### 1. Content Organization
- Use `order` field to sequence SubMaterials logically
- Start with video introduction
- Follow with reading or interactive content
- Include knowledge checks (quizzes)
- End with practical exercises (AR/flashcards)

### 2. Content Types Selection
- **Videos**: For explanations and demonstrations
- **Quizzes**: For knowledge assessment
- **AR Experiments**: For hands-on practice
- **Flashcards**: For memorization and review

### 3. Activity Tracking
- All module access automatically tracks SubMaterial information
- RAG system gets enriched data about content structure
- Learning patterns can be analyzed across different content types

### 4. Progressive Learning
```
Order 1: Introduction Video (overview)
Order 2: Core Concept Video (detailed)
Order 3: Knowledge Check Quiz
Order 4: Hands-on AR Lab
Order 5: Review Flashcards
Order 6: Final Assessment Quiz
```

## Error Responses

### Module Not Found
```json
{
  "error": "Module not found"
}
```

### SubMaterial Not Found
```json
{
  "error": "SubMaterial not found"
}
```

### Invalid Order
```json
{
  "error": "Order must be a positive integer"
}
```

## Summary

SubMaterials provide a modern, flexible approach to learning content within modules. They support multiple content types and are fully integrated with the activity tracking and RAG systems for comprehensive learning analytics.

**Key Points:**
- Access through Module endpoints (`/module/all` and `/module/:id`)
- Supports videos, quizzes, AR, and flashcards
- Automatically tracked for learning analytics
- Enriched metadata for RAG/AI systems
- Ordered sequence for structured learning
- Future direct CRUD endpoints planned
