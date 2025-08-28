package database

import (
	"fmt"
	"log"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

// SeedVideoQuizzes seeds video quizzes for all existing video materials
// This function will create comprehensive video quizzes for all modules
func SeedVideoQuizzes(db *gorm.DB) {
	log.Println("Starting Video Quizzes seeding...")

	// Get all video materials with their modules
	var videoMaterials []struct {
		models.VideoMaterial
		ModuleTitle string `gorm:"column:module_title"`
	}
	
	if err := db.Table("video_materials vm").
		Select("vm.*, m.title as module_title").
		Joins("JOIN modules m ON vm.module_id = m.id").
		Find(&videoMaterials).Error; err != nil {
		log.Printf("Error fetching video materials: %v", err)
		return
	}

	if len(videoMaterials) == 0 {
		log.Println("No video materials found. Please seed modules and video materials first.")
		return
	}

	for _, vmData := range videoMaterials {
		vm := vmData.VideoMaterial
		moduleTitle := vmData.ModuleTitle
		
		// Check if video quizzes already exist for this video material
		var existingCount int64
		db.Model(&models.VideoQuiz{}).Where("video_material_id = ?", vm.ID).Count(&existingCount)
		
		if existingCount > 0 {
			log.Printf("Video quizzes already exist for video material '%s' (Module: %s), skipping...", 
				vm.Title, moduleTitle)
			continue
		}

		// Create video quizzes based on module type
		switch moduleTitle {
		case "Fotosintesis - Dasar", "Materi Pokok 1: Pengenalan Fotosintesis":
			seedFotosintesisDasarQuizzes(db, vm, moduleTitle)
		case "Fotosintesis - Lanjutan", "Materi Pokok 2: Bagian Tumbuhan yang Berperan dalam Proses Fotosintesis":
			seedFotosintesisLanjutanQuizzes(db, vm, moduleTitle)
		case "Fotosintesis - Eksperimen", "Materi Pokok 3: Faktor yang Mempengaruhi Fotosintesis":
			seedFotosintesisEksperimenQuizzes(db, vm, moduleTitle)
		case "Materi Pokok 4: Proses Fotosintesis":
			seedFotosintesisProsesQuizzes(db, vm, moduleTitle)
		default:
			// For other modules, create generic quizzes
			seedGenericVideoQuizzes(db, vm, moduleTitle)
		}
	}

	log.Println("Video Quizzes seeding completed!")
}

// seedFotosintesisDasarQuizzes creates quizzes for basic photosynthesis videos (Module 1)
func seedFotosintesisDasarQuizzes(db *gorm.DB, vm models.VideoMaterial, moduleTitle string) {
	log.Printf("Creating quizzes for basic photosynthesis video: %s", vm.Title)

	quizzes := []models.VideoQuiz{
		{
			VideoMaterialID: vm.ID,
			Question:        "Kata fotosintesis berasal dari bahasa Yunani. Apa arti dari kata \"photo\" dan \"synthesis\"?",
			TimestampStart:  50, // Muncul di detik ke-50 (0:50)
			TimestampEnd:    65, // Berakhir di detik ke-65
			Options: models.Options{
				OptionA: "Photo berarti cahaya, synthesis berarti menggabungkan atau menyusun",
				OptionB: "Photo berarti air, synthesis berarti tumbuhan",
				OptionC: "Photo berarti udara, synthesis berarti makanan",
				OptionD: "Photo berarti daun, synthesis berarti energi",
			},
			CorrectAnswer: "A",
			Explanation:   "Benar! Photo berarti cahaya dan synthesis berarti menggabungkan atau menyusun. Jadi fotosintesis artinya menggabungkan dengan bantuan cahaya!",
			Order:         1,
		},
		{
			VideoMaterialID: vm.ID,
			Question:        "Agar proses fotosintesis dapat berlangsung, tumbuhan memerlukan beberapa bahan. Manakah kombinasi yang benar?",
			TimestampStart:  60, // Muncul di detik ke-60 (1:00)
			TimestampEnd:    75, // Berakhir di detik ke-75
			Options: models.Options{
				OptionA: "Cahaya matahari, Air (H₂O), Karbon dioksida (CO₂), dan Klorofil",
				OptionB: "Air (H₂O), Nitrogen, Karbon dioksida (CO₂), dan Cahaya bulan",
				OptionC: "Cahaya matahari, Oksigen, Air (H₂O), dan Tanah",
				OptionD: "Cahaya matahari, Protein, Oksigen, dan Kloroplas",
			},
			CorrectAnswer: "A",
			Explanation:   "Tepat! Fotosintesis membutuhkan 4 komponen utama: cahaya matahari sebagai energi, air dari akar, karbon dioksida dari udara, dan klorofil untuk menangkap cahaya!",
			Order:         2,
		},
		{
			VideoMaterialID: vm.ID,
			Question:        "Tumbuhan membuat makanannya sendiri di bagian dalam daun yang tidak terlihat oleh mata kita. Bagian tersebut disebut",
			TimestampStart:  90,  // Muncul di detik ke-90 (1:30)
			TimestampEnd:    105, // Berakhir di detik ke-105
			Options: models.Options{
				OptionA: "Klorofil",
				OptionB: "Kloroplas",
				OptionC: "Stomata",
				OptionD: "Xilem",
			},
			CorrectAnswer: "B",
			Explanation:   "Benar! Kloroplas adalah organel kecil di dalam sel daun tempat fotosintesis terjadi. Di sanalah klorofil berada!",
			Order:         3,
		},
	}

	createVideoQuizzes(db, quizzes, fmt.Sprintf("Fotosintesis Dasar - %s", vm.Title))
}

// seedFotosintesisLanjutanQuizzes creates quizzes for advanced photosynthesis videos (Module 2 - Bagian Tumbuhan)
func seedFotosintesisLanjutanQuizzes(db *gorm.DB, vm models.VideoMaterial, moduleTitle string) {
	log.Printf("Creating quizzes for advanced photosynthesis video: %s", vm.Title)

	quizzes := []models.VideoQuiz{
		{
			VideoMaterialID: vm.ID,
			Question:        "Apa fungsi utama akar pada tumbuhan?",
			TimestampStart:  127, // Muncul di detik ke-127 (2:07)
			TimestampEnd:    142, // Berakhir di detik ke-142
			Options: models.Options{
				OptionA: "Menyerap air dan nutrisi dari tanah",
				OptionB: "Menghasilkan cahaya",
				OptionC: "Menyimpan oksigen",
				OptionD: "Menghindari hama",
			},
			CorrectAnswer: "A",
			Explanation:   "Benar! Akar berfungsi menyerap air dan nutrisi dari tanah yang dibutuhkan untuk fotosintesis dan pertumbuhan tumbuhan!",
			Order:         1,
		},
		{
			VideoMaterialID: vm.ID,
			Question:        "Bagian tumbuhan manakah yang berfungsi mengangkut air dan zat makanan dari akar ke daun?",
			TimestampStart:  152, // Muncul di detik ke-152 (2:32)
			TimestampEnd:    167, // Berakhir di detik ke-167
			Options: models.Options{
				OptionA: "Akar",
				OptionB: "Batang",
				OptionC: "Daun",
				OptionD: "Biji",
			},
			CorrectAnswer: "B",
			Explanation:   "Tepat! Batang berfungsi sebagai jalan raya bagi tumbuhan untuk mengangkut air dari akar ke daun dan makanan dari daun ke seluruh tubuh tumbuhan!",
			Order:         2,
		},
		{
			VideoMaterialID: vm.ID,
			Question:        "Apa fungsi daun dalam tumbuhan?",
			TimestampStart:  196, // Muncul di detik ke-196 (3:16)
			TimestampEnd:    211, // Berakhir di detik ke-211
			Options: models.Options{
				OptionA: "Menyerap air dan nutrisi",
				OptionB: "Tempat fotosintesis untuk membuat makanan",
				OptionC: "Melindungi tanah",
				OptionD: "Menyimpan cadangan makanan",
			},
			CorrectAnswer: "B",
			Explanation:   "Benar! Daun adalah tempat fotosintesis terjadi. Di daun inilah tumbuhan membuat makanannya sendiri menggunakan cahaya matahari!",
			Order:         3,
		},
	}

	createVideoQuizzes(db, quizzes, fmt.Sprintf("Fotosintesis Lanjutan - %s", vm.Title))
}

// seedFotosintesisProsesQuizzes creates quizzes for photosynthesis process videos (Module 4)
func seedFotosintesisProsesQuizzes(db *gorm.DB, vm models.VideoMaterial, moduleTitle string) {
	log.Printf("Creating quizzes for photosynthesis process video: %s", vm.Title)

	quizzes := []models.VideoQuiz{
		{
			VideoMaterialID: vm.ID,
			Question:        "Apa saja yang dibutuhkan tumbuhan untuk melakukan fotosintesis?",
			TimestampStart:  227, // Muncul di detik ke-227 (3:47)
			TimestampEnd:    242, // Berakhir di detik ke-242
			Options: models.Options{
				OptionA: "Air, Karbondioksida, Matahari, Klorofil",
				OptionB: "Oksigen, Air, Karbondioksida, Tanah",
				OptionC: "Air, Matahari, Tanah, Oksigen",
				OptionD: "Air, Klorofil, Nitrogen, Cahaya Bulan",
			},
			CorrectAnswer: "A",
			Explanation:   "Benar! Tumbuhan membutuhkan air dari akar, karbondioksida dari udara, cahaya matahari sebagai energi, dan klorofil untuk menangkap cahaya!",
			Order:         1,
		},
		{
			VideoMaterialID: vm.ID,
			Question:        "Apa hasil utama dari proses fotosintesis?",
			TimestampStart:  248, // Muncul di detik ke-248 (4:08)
			TimestampEnd:    263, // Berakhir di detik ke-263
			Options: models.Options{
				OptionA: "Nitrogen dan Air",
				OptionB: "Oksigen dan Karbohidrat (makanan/glukosa)",
				OptionC: "Karbondioksida dan Oksigen",
				OptionD: "Cahaya Matahari dan Air",
			},
			CorrectAnswer: "B",
			Explanation:   "Tepat! Fotosintesis menghasilkan oksigen yang kita hirup dan karbohidrat (glukosa) sebagai makanan untuk tumbuhan!",
			Order:         2,
		},
		{
			VideoMaterialID: vm.ID,
			Question:        "Mengapa fotosintesis penting bagi makhluk hidup di Bumi?",
			TimestampStart:  284, // Muncul di detik ke-284 (4:44)
			TimestampEnd:    299, // Berakhir di detik ke-299
			Options: models.Options{
				OptionA: "Karena menghasilkan tanah yang subur",
				OptionB: "Karena menghasilkan oksigen dan makanan untuk rantai makanan",
				OptionC: "Karena membuat daun selalu hijau",
				OptionD: "Karena menyerap semua air hujan",
			},
			CorrectAnswer: "B",
			Explanation:   "Benar! Fotosintesis menghasilkan oksigen untuk bernapas dan menjadi dasar rantai makanan bagi semua makhluk hidup di Bumi!",
			Order:         3,
		},
	}

	createVideoQuizzes(db, quizzes, fmt.Sprintf("Fotosintesis Proses - %s", vm.Title))
}

// seedFotosintesisEksperimenQuizzes creates quizzes for photosynthesis experiment videos (Module 3)
func seedFotosintesisEksperimenQuizzes(db *gorm.DB, vm models.VideoMaterial, moduleTitle string) {
	log.Printf("Creating quizzes for photosynthesis experiment video: %s", vm.Title)

	quizzes := []models.VideoQuiz{
		{
			VideoMaterialID: vm.ID,
			Question:        "Dalam eksperimen Ingenhousz, apa yang terjadi pada tumbuhan air ketika terkena cahaya?",
			TimestampStart:  50,
			TimestampEnd:    65,
			Options: models.Options{
				OptionA: "Tumbuhan menjadi layu",
				OptionB: "Tumbuhan mengeluarkan gelembung gas oksigen",
				OptionC: "Tumbuhan menyerap semua air",
				OptionD: "Tumbuhan berubah warna",
			},
			CorrectAnswer: "B",
			Explanation:   "Dalam eksperimen Ingenhousz, tumbuhan air yang terkena cahaya akan mengeluarkan gelembung gas oksigen sebagai hasil fotosintesis.",
			Order:         1,
		},
		{
			VideoMaterialID: vm.ID,
			Question:        "Mengapa dalam eksperimen fotosintesis, tumbuhan diletakkan dalam air?",
			TimestampStart:  80,
			TimestampEnd:    95,
			Options: models.Options{
				OptionA: "Agar tumbuhan tidak mati",
				OptionB: "Untuk melihat gelembung oksigen yang dihasilkan",
				OptionC: "Agar mudah diamati",
				OptionD: "Untuk memberikan nutrisi tambahan",
			},
			CorrectAnswer: "B",
			Explanation:   "Tumbuhan diletakkan dalam air agar kita dapat melihat gelembung gas oksigen yang dihasilkan dari proses fotosintesis dengan jelas.",
			Order:         2,
		},
		{
			VideoMaterialID: vm.ID,
			Question:        "Apa yang terjadi jika eksperimen fotosintesis dilakukan di tempat gelap?",
			TimestampStart:  110,
			TimestampEnd:    125,
			Options: models.Options{
				OptionA: "Gelembung gas tetap dihasilkan dengan jumlah sama",
				OptionB: "Gelembung gas dihasilkan lebih banyak",
				OptionC: "Tidak ada atau sangat sedikit gelembung gas yang dihasilkan",
				OptionD: "Tumbuhan akan menghasilkan karbon dioksida",
			},
			CorrectAnswer: "C",
			Explanation:   "Tanpa cahaya, fotosintesis tidak dapat berlangsung optimal, sehingga tidak ada atau sangat sedikit oksigen yang dihasilkan.",
			Order:         3,
		},
	}

	createVideoQuizzes(db, quizzes, fmt.Sprintf("Fotosintesis Eksperimen - %s", vm.Title))
}

// seedGenericVideoQuizzes creates generic quizzes for other video materials
func seedGenericVideoQuizzes(db *gorm.DB, vm models.VideoMaterial, moduleTitle string) {
	log.Printf("Creating generic quizzes for video: %s (Module: %s)", vm.Title, moduleTitle)

	quizzes := []models.VideoQuiz{
		{
			VideoMaterialID: vm.ID,
			Question:        fmt.Sprintf("Apa topik utama yang dibahas dalam video '%s'?", vm.Title),
			TimestampStart:  30,
			TimestampEnd:    45,
			Options: models.Options{
				OptionA: fmt.Sprintf("Konsep dasar %s", moduleTitle),
				OptionB: "Teori umum biologi",
				OptionC: "Matematika dasar",
				OptionD: "Sejarah sains",
			},
			CorrectAnswer: "A",
			Explanation:   fmt.Sprintf("Video ini membahas konsep dasar dari %s yang penting untuk dipahami.", moduleTitle),
			Order:         1,
		},
		{
			VideoMaterialID: vm.ID,
			Question:        "Mengapa materi ini penting untuk dipelajari?",
			TimestampStart:  60,
			TimestampEnd:    75,
			Options: models.Options{
				OptionA: "Untuk ujian saja",
				OptionB: "Membantu memahami konsep sains dalam kehidupan sehari-hari",
				OptionC: "Tidak ada manfaatnya",
				OptionD: "Hanya untuk akademisi",
			},
			CorrectAnswer: "B",
			Explanation:   "Mempelajari sains membantu kita memahami fenomena alam dan menerapkan pengetahuan dalam kehidupan sehari-hari.",
			Order:         2,
		},
	}

	createVideoQuizzes(db, quizzes, fmt.Sprintf("Generic - %s", vm.Title))
}

// createVideoQuizzes helper function to create video quizzes in database
func createVideoQuizzes(db *gorm.DB, quizzes []models.VideoQuiz, context string) {
	for _, quiz := range quizzes {
		if err := db.Create(&quiz).Error; err != nil {
			log.Printf("Error creating video quiz for %s: %v", context, err)
		} else {
			questionPreview := quiz.Question
			if len(questionPreview) > 50 {
				questionPreview = questionPreview[:50] + "..."
			}
			log.Printf("✓ Created video quiz: %s (Order: %d)", questionPreview, quiz.Order)
		}
	}
}

// SeedVideoQuizzesForModule seeds video quizzes for specific module
func SeedVideoQuizzesForModule(db *gorm.DB, moduleTitle string) {
	log.Printf("Seeding video quizzes for module: %s", moduleTitle)

	var module models.Module
	if err := db.Where("title = ?", moduleTitle).First(&module).Error; err != nil {
		log.Printf("Module '%s' not found: %v", moduleTitle, err)
		return
	}

	var videoMaterials []models.VideoMaterial
	if err := db.Where("module_id = ?", module.ID).Find(&videoMaterials).Error; err != nil {
		log.Printf("Error fetching video materials for module '%s': %v", moduleTitle, err)
		return
	}

	for _, vm := range videoMaterials {
		// Check if video quizzes already exist
		var existingCount int64
		db.Model(&models.VideoQuiz{}).Where("video_material_id = ?", vm.ID).Count(&existingCount)
		
		if existingCount > 0 {
			log.Printf("Video quizzes already exist for video material '%s', skipping...", vm.Title)
			continue
		}

		// Create quizzes based on module type
		switch moduleTitle {
		case "Fotosintesis - Dasar", "Materi Pokok 1: Pengenalan Fotosintesis":
			seedFotosintesisDasarQuizzes(db, vm, moduleTitle)
		case "Fotosintesis - Lanjutan", "Materi Pokok 2: Bagian Tumbuhan yang Berperan dalam Proses Fotosintesis":
			seedFotosintesisLanjutanQuizzes(db, vm, moduleTitle)
		case "Fotosintesis - Eksperimen", "Materi Pokok 3: Faktor yang Mempengaruhi Fotosintesis":
			seedFotosintesisEksperimenQuizzes(db, vm, moduleTitle)
		case "Materi Pokok 4: Proses Fotosintesis":
			seedFotosintesisProsesQuizzes(db, vm, moduleTitle)
		default:
			seedGenericVideoQuizzes(db, vm, moduleTitle)
		}
	}

	log.Printf("Completed seeding video quizzes for module: %s", moduleTitle)
}

// ClearVideoQuizzes removes all video quizzes and their answers
func ClearVideoQuizzes(db *gorm.DB) {
	log.Println("Clearing all video quizzes...")
	
	// Delete user answers first (foreign key constraint)
	if err := db.Exec("DELETE FROM video_quiz_user_answers").Error; err != nil {
		log.Printf("Error clearing video quiz user answers: %v", err)
	}
	
	// Delete video quizzes
	if err := db.Exec("DELETE FROM video_quizzes").Error; err != nil {
		log.Printf("Error clearing video quizzes: %v", err)
	} else {
		log.Println("✓ All video quizzes cleared successfully")
	}
}

// GetVideoQuizzesSummary returns a summary of video quizzes in the database
func GetVideoQuizzesSummary(db *gorm.DB) {
	log.Println("=== VIDEO QUIZZES SUMMARY ===")
	
	var totalQuizzes int64
	db.Model(&models.VideoQuiz{}).Count(&totalQuizzes)
	log.Printf("Total Video Quizzes: %d", totalQuizzes)
	
	// Group by module
	var results []struct {
		ModuleTitle string
		VideoTitle  string
		QuizCount   int64
	}
	
	db.Table("video_quizzes vq").
		Select("m.title as module_title, vm.title as video_title, COUNT(vq.id) as quiz_count").
		Joins("JOIN video_materials vm ON vq.video_material_id = vm.id").
		Joins("JOIN modules m ON vm.module_id = m.id").
		Group("m.id, vm.id").
		Order("m.title, vm.title").
		Find(&results)
	
	currentModule := ""
	for _, result := range results {
		if result.ModuleTitle != currentModule {
			log.Printf("\n📚 Module: %s", result.ModuleTitle)
			currentModule = result.ModuleTitle
		}
		log.Printf("  🎥 %s: %d quizzes", result.VideoTitle, result.QuizCount)
	}
	
	log.Println("==============================")
}
