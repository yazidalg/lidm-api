package main

import (
	"log"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/config"
	"github.com/yazidalg/lidm_backend/internal/database"
)

func main() {
	// Load environment and connect to database
	config.LoadEnv()
	db := config.ConnectDB()
	database.Migrate(db)

	// Video data for all modules
	videoData := map[uint]struct {
		title string
		url   string
	}{
		1: {"Video Pembelajaran Fotosintesis", "https://www.youtube.com/watch?v=53v-bKx53Lo"},
		2: {"Video Bagian Tumbuhan dalam Fotosintesis", "https://youtu.be/eOTMazMILI4"},
		3: {"Video Faktor yang Mempengaruhi Fotosintesis", "https://www.youtube.com/watch?v=4rtlDVxQ7zA"},
		4: {"Video Proses Fotosintesis", "https://www.youtube.com/watch?v=nVoFOWTMQi4"},
		5: {"Video Produk Fotosintesis", "https://www.youtube.com/watch?v=4UknZXtQTGE"},
		6: {"Video Kaitan Fotosintesis dengan Makhluk Hidup", "https://www.youtube.com/watch?v=85UpVXMrZ2Q"},
	}

	for moduleID, video := range videoData {
		// Check if module exists
		var module models.Module
		result := db.First(&module, moduleID)
		if result.Error != nil {
			log.Printf("Module %d not found: %v", moduleID, result.Error)
			continue
		}

		// Check if video material already exists for this module
		var existingVideo models.VideoMaterial
		result = db.Where("module_id = ?", moduleID).First(&existingVideo)
		
		if result.Error == nil {
			// Video material exists, update it
			existingVideo.YoutubeLink = video.url
			existingVideo.Title = video.title
			existingVideo.Duration = 0 // Will be set later if needed

			err := db.Save(&existingVideo).Error
			if err != nil {
				log.Printf("Failed to update video material for Module %d: %v", moduleID, err)
				continue
			}

			log.Printf("Successfully updated video material for Module %d", moduleID)
			log.Printf("Video ID: %d", existingVideo.ID)
			log.Printf("YouTube Link: %s", existingVideo.YoutubeLink)
		} else {
			// Video material doesn't exist, create new one
			videoMaterial := models.VideoMaterial{
				ModuleID:    moduleID,
				Title:       video.title,
				YoutubeLink: video.url,
				Duration:    0, // Will be set later if needed
			}

			err := db.Create(&videoMaterial).Error
			if err != nil {
				log.Printf("Failed to create video material for Module %d: %v", moduleID, err)
				continue
			}

			log.Printf("Successfully created video material for Module %d", moduleID)
			log.Printf("Video ID: %d", videoMaterial.ID)
			log.Printf("YouTube Link: %s", videoMaterial.YoutubeLink)
		}
	}

	log.Println("All video materials setup completed!")
}
