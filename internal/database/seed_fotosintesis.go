package database

import (
	"fmt"
	"log"
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

func SeedFotosintesisData(db *gorm.DB) {
	log.Println("Starting Fotosintesis seeding...")

	// Auto migrate all models first
	err := db.AutoMigrate(
		&models.Module{},
		&models.SubMaterial{},
		&models.VideoMaterial{},
		&models.VideoQuiz{},
		&models.VideoQuizUserAnswer{},
		&models.Question{},
		&models.ARExperiment{},
		&models.Flashcard{},
		&models.Prequiz{},
		&models.PrequizUserAnswer{},
		&models.Role{},
		&models.User{},
	)
	if err != nil {
		log.Printf("Error migrating models: %v", err)
		return
	}

	// Clear existing data
	db.Exec("DELETE FROM video_quiz_user_answers")
	db.Exec("DELETE FROM prequiz_user_answers")
	db.Exec("DELETE FROM users WHERE email LIKE '%@student.com'")
	db.Exec("DELETE FROM flashcards")
	db.Exec("DELETE FROM questions WHERE module_id IN (SELECT id FROM modules WHERE title LIKE '%Fotosintesis%')")
	db.Exec("DELETE FROM video_materials")
	db.Exec("DELETE FROM ar_experiments")
	db.Exec("DELETE FROM sub_materials")
	db.Exec("DELETE FROM modules WHERE title LIKE '%Fotosintesis%'")

	// Create Module: Fotosintesis untuk Kelas 4 SD
	module := models.Module{
		Title:       "Belajar Fotosintesis - Kelas 4 SD",
		Description: "Mari belajar tentang fotosintesis dengan cara yang menyenangkan! Kita akan memahami bagaimana tumbuhan membuat makanannya sendiri.",
		OffsetX:     0,
		OffsetY:     0,
		Icon:        "🌱",
		Thumbnail:   "https://i.pinimg.com/736x/6b/45/da/6b45da50aa5726a3b36e68d2859a9232.jpg",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(&module).Error; err != nil {
		log.Printf("Error creating module: %v", err)
		return
	}

	log.Printf("Created module: %s (ID: %d)", module.Title, module.ID)

	// Create AR Experiment first (will be used in SubMaterial 3)
	arExperiment := models.ARExperiment{
		Title:     "Laboratorium Virtual Fotosintesis",
		LinkAR:    "https://asblr.com/e4Z066",
		LinkEmbed: "https://asblr.com/e4Z066",
		OffsetX:   120,
		OffsetY:   310,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&arExperiment).Error; err != nil {
		log.Printf("Error creating AR experiment: %v", err)
		return
	}

	log.Printf("Created AR experiment: %s (ID: %d)", arExperiment.Title, arExperiment.ID)

	// SubMaterial 1: Video Pengenalan (Offset 120, 50)
	subMaterial1 := models.SubMaterial{
		ModuleID:    module.ID,
		Title:       "Apa itu Fotosintesis?",
		Description: "Video pengenalan tentang fotosintesis untuk anak-anak. Mari mengenal bagaimana tumbuhan bisa hidup!",
		Order:       1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(&subMaterial1).Error; err != nil {
		log.Printf("Error creating sub material 1: %v", err)
		return
	}

	// Prequizzes untuk SubMaterial 1 (minimal 10)
	prequizzes1 := []models.Prequiz{
		{
			SubMaterialID: subMaterial1.ID,
			Question:      "Apakah kamu pernah melihat tumbuhan?",
			Options: models.Options{
				OptionA: "Ya, setiap hari",
				OptionB: "Kadang-kadang",
				OptionC: "Jarang sekali",
				OptionD: "Tidak pernah",
			},
			CorrectAnswer: "A",
			Explanation:   "Bagus! Tumbuhan ada di mana-mana di sekitar kita. Mari kita pelajari lebih lanjut!",
		},
		{
			SubMaterialID: subMaterial1.ID,
			Question:      "Menurutmu, bagaimana tumbuhan mendapatkan makanan?",
			Options: models.Options{
				OptionA: "Membeli di toko",
				OptionB: "Meminta dari tumbuhan lain",
				OptionC: "Membuat sendiri",
				OptionD: "Tidak perlu makan",
			},
			CorrectAnswer: "C",
			Explanation:   "Benar! Tumbuhan bisa membuat makanan sendiri. Mari kita pelajari caranya!",
		},
		{
			SubMaterialID: subMaterial1.ID,
			Question:      "Apa warna daun yang paling sering kamu lihat?",
			Options: models.Options{
				OptionA: "Hijau",
				OptionB: "Merah",
				OptionC: "Kuning",
				OptionD: "Biru",
			},
			CorrectAnswer: "A",
			Explanation:   "Ya! Daun biasanya berwarna hijau. Warna hijau ini sangat penting untuk fotosintesis!",
		},
		{
			SubMaterialID: subMaterial1.ID,
			Question:      "Kapan tumbuhan terlihat paling segar?",
			Options: models.Options{
				OptionA: "Saat malam",
				OptionB: "Saat pagi hari",
				OptionC: "Saat siang hari",
				OptionD: "Tidak ada bedanya",
			},
			CorrectAnswer: "B",
			Explanation:   "Benar! Pagi hari tumbuhan terlihat segar karena sudah siap menerima sinar matahari!",
		},
		{
			SubMaterialID: subMaterial1.ID,
			Question:      "Apa yang terjadi jika tumbuhan tidak disiram?",
			Options: models.Options{
				OptionA: "Akan tumbuh lebih cepat",
				OptionB: "Tidak ada perubahan",
				OptionC: "Akan layu dan kering",
				OptionD: "Akan berubah warna menjadi ungu",
			},
			CorrectAnswer: "C",
			Explanation:   "Ya! Tumbuhan membutuhkan air untuk hidup. Tanpa air, mereka akan layu dan kering.",
		},
		{
			SubMaterialID: subMaterial1.ID,
			Question:      "Di mana biasanya kamu melihat banyak tumbuhan?",
			Options: models.Options{
				OptionA: "Di taman atau kebun",
				OptionB: "Di dalam lemari es",
				OptionC: "Di dalam mobil",
				OptionD: "Di dalam tas",
			},
			CorrectAnswer: "A",
			Explanation:   "Benar! Taman dan kebun adalah tempat yang bagus untuk melihat berbagai jenis tumbuhan!",
		},
		{
			SubMaterialID: subMaterial1.ID,
			Question:      "Menurutmu, apakah tumbuhan bisa bernapas?",
			Options: models.Options{
				OptionA: "Ya, seperti manusia",
				OptionB: "Tidak, mereka tidak hidup",
				OptionC: "Ya, tapi caranya berbeda",
				OptionD: "Hanya saat malam hari",
			},
			CorrectAnswer: "C",
			Explanation:   "Pintar! Tumbuhan bisa bernapas, tapi caranya berbeda dari manusia. Mari kita pelajari!",
		},
		{
			SubMaterialID: subMaterial1.ID,
			Question:      "Bagian tumbuhan mana yang biasanya berwarna hijau?",
			Options: models.Options{
				OptionA: "Akar",
				OptionB: "Daun",
				OptionC: "Batang",
				OptionD: "Bunga",
			},
			CorrectAnswer: "B",
			Explanation:   "Benar! Daun biasanya berwarna hijau karena mengandung zat khusus yang disebut klorofil!",
		},
		{
			SubMaterialID: subMaterial1.ID,
			Question:      "Apa yang kamu rasakan saat berada di dekat banyak tumbuhan?",
			Options: models.Options{
				OptionA: "Udara terasa lebih segar",
				OptionB: "Tidak ada perbedaan",
				OptionC: "Udara terasa lebih panas",
				OptionD: "Sulit bernapas",
			},
			CorrectAnswer: "A",
			Explanation:   "Ya! Tumbuhan membantu membuat udara lebih segar dengan menghasilkan oksigen!",
		},
		{
			SubMaterialID: subMaterial1.ID,
			Question:      "Menurutmu, apakah matahari penting untuk tumbuhan?",
			Options: models.Options{
				OptionA: "Ya, sangat penting",
				OptionB: "Tidak penting",
				OptionC: "Hanya sedikit penting",
				OptionD: "Malah berbahaya",
			},
			CorrectAnswer: "A",
			Explanation:   "Tepat sekali! Matahari sangat penting untuk tumbuhan dalam proses fotosintesis!",
		},
	}

	for _, prequiz := range prequizzes1 {
		if err := db.Create(&prequiz).Error; err != nil {
			log.Printf("Error creating prequiz for sub material 1: %v", err)
			continue
		}
	}

	// Video Material untuk SubMaterial 1
	videoMaterial1 := models.VideoMaterial{
		SubMaterialID: subMaterial1.ID,
		Title:         "Video Pengenalan Fotosintesis",
		YoutubeLink:   "https://www.youtube.com/watch?v=example-fotosintesis-intro",
		Duration:      180, // 3 menit
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := db.Create(&videoMaterial1).Error; err != nil {
		log.Printf("Error creating video material 1: %v", err)
		return
	}

	// Video Quizzes untuk Video Material 1 (muncul di tengah video)
	videoQuizzes1 := []models.VideoQuiz{
		{
			VideoMaterialID: videoMaterial1.ID,
			Question:        "Berdasarkan video, apa yang membuat daun berwarna hijau?",
			TimestampStart:  60,  // Muncul di detik ke-60 (1 menit)
			TimestampEnd:    75,  // Berakhir di detik ke-75
			Options: models.Options{
				OptionA: "Cat hijau",
				OptionB: "Klorofil",
				OptionC: "Air",
				OptionD: "Tanah",
			},
			CorrectAnswer: "B",
			Explanation:   "Benar! Klorofil adalah zat hijau yang membuat daun berwarna hijau dan berperan penting dalam fotosintesis!",
			Order:         1,
		},
		{
			VideoMaterialID: videoMaterial1.ID,
			Question:        "Menurut video, apa saja yang dibutuhkan tumbuhan untuk fotosintesis?",
			TimestampStart:  120, // Muncul di detik ke-120 (2 menit)
			TimestampEnd:    135, // Berakhir di detik ke-135
			Options: models.Options{
				OptionA: "Hanya air",
				OptionB: "Hanya sinar matahari",
				OptionC: "Air, sinar matahari, dan karbon dioksida",
				OptionD: "Pupuk dan pestisida",
			},
			CorrectAnswer: "C",
			Explanation:   "Tepat! Fotosintesis membutuhkan 3 bahan utama: air, sinar matahari, dan karbon dioksida dari udara!",
			Order:         2,
		},
	}

	for _, videoQuiz := range videoQuizzes1 {
		if err := db.Create(&videoQuiz).Error; err != nil {
			log.Printf("Error creating video quiz for video material 1: %v", err)
			continue
		}
	}

	// SubMaterial 2: Quiz Pengetahuan Dasar (Offset 200, 180)
	subMaterial2 := models.SubMaterial{
		ModuleID:    module.ID,
		Title:       "Kuis: Apa yang Kamu Tahu?",
		Description: "Ayo uji pengetahuan dasar kamu tentang tumbuhan dan fotosintesis!",
		Order:       2,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(&subMaterial2).Error; err != nil {
		log.Printf("Error creating sub material 2: %v", err)
		return
	}

	// Prequizzes untuk SubMaterial 2 (minimal 10)
	prequizzes2 := []models.Prequiz{
		{
			SubMaterialID: subMaterial2.ID,
			Question:      "Apa yang kamu ketahui tentang daun hijau?",
			Options: models.Options{
				OptionA: "Hanya untuk hiasan",
				OptionB: "Bisa membuat makanan",
				OptionC: "Tidak berguna",
				OptionD: "Hanya untuk berteduh",
			},
			CorrectAnswer: "B",
			Explanation:   "Daun hijau mengandung klorofil yang bisa membuat makanan melalui fotosintesis!",
		},
		{
			SubMaterialID: subMaterial2.ID,
			Question:      "Menurutmu, apa yang terjadi jika tidak ada matahari?",
			Options: models.Options{
				OptionA: "Tumbuhan tetap sehat",
				OptionB: "Tidak ada perubahan",
				OptionC: "Tumbuhan tidak bisa membuat makanan",
				OptionD: "Tumbuhan tumbuh lebih cepat",
			},
			CorrectAnswer: "C",
			Explanation:   "Tanpa matahari, tumbuhan tidak bisa melakukan fotosintesis untuk membuat makanan!",
		},
		{
			SubMaterialID: subMaterial2.ID,
			Question:      "Apa yang dihasilkan tumbuhan yang berguna untuk kita?",
			Options: models.Options{
				OptionA: "Karbon dioksida",
				OptionB: "Oksigen",
				OptionC: "Racun",
				OptionD: "Debu",
			},
			CorrectAnswer: "B",
			Explanation:   "Tumbuhan menghasilkan oksigen yang kita butuhkan untuk bernapas!",
		},
		{
			SubMaterialID: subMaterial2.ID,
			Question:      "Bagian tumbuhan mana yang menyerap air?",
			Options: models.Options{
				OptionA: "Daun",
				OptionB: "Bunga",
				OptionC: "Akar",
				OptionD: "Batang",
			},
			CorrectAnswer: "C",
			Explanation:   "Akar berfungsi menyerap air dan mineral dari tanah untuk tumbuhan!",
		},
		{
			SubMaterialID: subMaterial2.ID,
			Question:      "Apa nama proses tumbuhan membuat makanan?",
			Options: models.Options{
				OptionA: "Fotosintesis",
				OptionB: "Bernapas",
				OptionC: "Makan",
				OptionD: "Tidur",
			},
			CorrectAnswer: "A",
			Explanation:   "Fotosintesis adalah proses tumbuhan membuat makanan menggunakan sinar matahari!",
		},
		{
			SubMaterialID: subMaterial2.ID,
			Question:      "Gas apa yang diambil tumbuhan dari udara?",
			Options: models.Options{
				OptionA: "Oksigen",
				OptionB: "Nitrogen",
				OptionC: "Karbon dioksida",
				OptionD: "Hidrogen",
			},
			CorrectAnswer: "C",
			Explanation:   "Tumbuhan mengambil karbon dioksida dari udara untuk fotosintesis!",
		},
		{
			SubMaterialID: subMaterial2.ID,
			Question:      "Kapan waktu terbaik tumbuhan melakukan fotosintesis?",
			Options: models.Options{
				OptionA: "Malam hari",
				OptionB: "Saat hujan",
				OptionC: "Saat ada sinar matahari",
				OptionD: "Saat angin kencang",
			},
			CorrectAnswer: "C",
			Explanation:   "Fotosintesis membutuhkan sinar matahari sebagai sumber energi!",
		},
		{
			SubMaterialID: subMaterial2.ID,
			Question:      "Mengapa daun berwarna hijau?",
			Options: models.Options{
				OptionA: "Karena dicat",
				OptionB: "Karena klorofil",
				OptionC: "Karena air",
				OptionD: "Karena tanah",
			},
			CorrectAnswer: "B",
			Explanation:   "Klorofil adalah zat hijau yang membantu tumbuhan menangkap sinar matahari!",
		},
		{
			SubMaterialID: subMaterial2.ID,
			Question:      "Apa yang terjadi jika tumbuhan tidak mendapat air?",
			Options: models.Options{
				OptionA: "Tumbuh lebih tinggi",
				OptionB: "Berbunga lebih banyak",
				OptionC: "Layu dan mati",
				OptionD: "Berubah warna menjadi biru",
			},
			CorrectAnswer: "C",
			Explanation:   "Air sangat penting untuk tumbuhan, tanpa air mereka akan layu dan mati!",
		},
		{
			SubMaterialID: subMaterial2.ID,
			Question:      "Selain oksigen, apa lagi yang dihasilkan fotosintesis?",
			Options: models.Options{
				OptionA: "Air",
				OptionB: "Glukosa (gula)",
				OptionC: "Tanah",
				OptionD: "Batu",
			},
			CorrectAnswer: "B",
			Explanation:   "Fotosintesis menghasilkan glukosa (gula) sebagai makanan tumbuhan dan oksigen!",
		},
	}

	for _, prequiz := range prequizzes2 {
		if err := db.Create(&prequiz).Error; err != nil {
			log.Printf("Error creating prequiz for sub material 2: %v", err)
			continue
		}
	}

	// Quiz Questions untuk SubMaterial 2 (sebagai Questions dengan ModuleID)
	quizQuestions2 := []models.Question{
		{
			ModuleID:      &module.ID,
			Question:      "Apa warna daun yang bisa membuat makanan sendiri?",
			AnswerTime:    30,
			ReadTime:      15,
			Options: models.Options{
				OptionA: "Merah",
				OptionB: "Hijau",
				OptionC: "Kuning",
				OptionD: "Ungu",
			},
			CorrectAnswer: "B",
			Explanation:   "Daun hijau mengandung klorofil yang bisa menangkap sinar matahari untuk fotosintesis!",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ModuleID:      &module.ID,
			Question:      "Apa yang dibutuhkan tumbuhan untuk hidup?",
			AnswerTime:    45,
			ReadTime:      20,
			Options: models.Options{
				OptionA: "Air saja",
				OptionB: "Sinar matahari saja",
				OptionC: "Air, sinar matahari, dan udara",
				OptionD: "Tanah saja",
			},
			CorrectAnswer: "C",
			Explanation:   "Tumbuhan membutuhkan air, sinar matahari, dan udara (karbon dioksida) untuk fotosintesis!",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	for _, quiz := range quizQuestions2 {
		if err := db.Create(&quiz).Error; err != nil {
			log.Printf("Error creating quiz question: %v", err)
			continue
		}
	}

	// SubMaterial 3: AR Experience (Offset 120, 310)
	subMaterial3 := models.SubMaterial{
		ModuleID:       module.ID,
		Title:          "Laboratorium Virtual AR",
		Description:    "Masuk ke dunia AR dan lihat langsung bagaimana fotosintesis terjadi di dalam daun!",
		Order:          3,
		ARExperimentID: &arExperiment.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := db.Create(&subMaterial3).Error; err != nil {
		log.Printf("Error creating sub material 3: %v", err)
		return
	}

	// Prequizzes untuk SubMaterial 3 (AR Lab - minimal 10)
	prequizzes3 := []models.Prequiz{
		{
			SubMaterialID: subMaterial3.ID,
			Question:      "Pernahkah kamu menggunakan teknologi AR (Augmented Reality)?",
			Options: models.Options{
				OptionA: "Sering sekali",
				OptionB: "Pernah beberapa kali",
				OptionC: "Jarang sekali",
				OptionD: "Belum pernah",
			},
			CorrectAnswer: "A",
			Explanation:   "AR akan membantu kita melihat fotosintesis secara virtual dan interaktif!",
		},
		{
			SubMaterialID: subMaterial3.ID,
			Question:      "Apa yang ingin kamu lihat di dalam daun?",
			Options: models.Options{
				OptionA: "Warna hijau saja",
				OptionB: "Proses fotosintesis",
				OptionC: "Tidak ada yang menarik",
				OptionD: "Serangga kecil",
			},
			CorrectAnswer: "B",
			Explanation:   "Dengan AR, kita bisa melihat bagaimana fotosintesis terjadi di dalam daun!",
		},
		{
			SubMaterialID: subMaterial3.ID,
			Question:      "Bagaimana menurutmu cara terbaik memahami sains?",
			Options: models.Options{
				OptionA: "Membaca buku saja",
				OptionB: "Melihat langsung dengan teknologi",
				OptionC: "Mendengar penjelasan",
				OptionD: "Menghapal rumus",
			},
			CorrectAnswer: "B",
			Explanation:   "Teknologi AR membantu kita melihat dan memahami proses yang tidak bisa dilihat mata!",
		},
		{
			SubMaterialID: subMaterial3.ID,
			Question:      "Apa yang paling kamu ingin ketahui tentang fotosintesis?",
			Options: models.Options{
				OptionA: "Bagaimana prosesnya terjadi",
				OptionB: "Kenapa daun berwarna hijau",
				OptionC: "Apa yang dihasilkan",
				OptionD: "Semua hal di atas",
			},
			CorrectAnswer: "D",
			Explanation:   "AR Lab akan menunjukkan semua aspek fotosintesis secara visual dan interaktif!",
		},
		{
			SubMaterialID: subMaterial3.ID,
			Question:      "Menurutmu, di bagian mana fotosintesis terjadi?",
			Options: models.Options{
				OptionA: "Di seluruh tumbuhan",
				OptionB: "Hanya di daun",
				OptionC: "Di akar saja",
				OptionD: "Di batang saja",
			},
			CorrectAnswer: "B",
			Explanation:   "Fotosintesis terutama terjadi di daun, mari kita lihat di AR Lab!",
		},
		{
			SubMaterialID: subMaterial3.ID,
			Question:      "Apa yang kamu bayangkan ada di dalam sel daun?",
			Options: models.Options{
				OptionA: "Ruang kosong",
				OptionB: "Kloroplas hijau",
				OptionC: "Air saja",
				OptionD: "Tidak tahu",
			},
			CorrectAnswer: "B",
			Explanation:   "Kloroplas adalah tempat fotosintesis terjadi, kita akan melihatnya di AR!",
		},
		{
			SubMaterialID: subMaterial3.ID,
			Question:      "Bagaimana menurutmu sinar matahari masuk ke daun?",
			Options: models.Options{
				OptionA: "Langsung masuk ke dalam",
				OptionB: "Ditangkap oleh klorofil",
				OptionC: "Dipantulkan kembali",
				OptionD: "Tidak bisa masuk",
			},
			CorrectAnswer: "B",
			Explanation:   "Klorofil menangkap sinar matahari untuk digunakan dalam fotosintesis!",
		},
		{
			SubMaterialID: subMaterial3.ID,
			Question:      "Apa yang terjadi dengan karbon dioksida di dalam daun?",
			Options: models.Options{
				OptionA: "Dikeluarkan kembali",
				OptionB: "Disimpan saja",
				OptionC: "Diubah menjadi gula",
				OptionD: "Dibiarkan mengapung",
			},
			CorrectAnswer: "C",
			Explanation:   "Karbon dioksida diubah menjadi gula (glukosa) dalam proses fotosintesis!",
		},
		{
			SubMaterialID: subMaterial3.ID,
			Question:      "Menurutmu, berapa lama proses fotosintesis berlangsung?",
			Options: models.Options{
				OptionA: "Hanya sebentar",
				OptionB: "Sepanjang ada sinar matahari",
				OptionC: "Hanya saat pagi",
				OptionD: "Hanya saat siang",
			},
			CorrectAnswer: "B",
			Explanation:   "Fotosintesis berlangsung terus menerus selama ada sinar matahari!",
		},
		{
			SubMaterialID: subMaterial3.ID,
			Question:      "Apa yang paling kamu harapkan dari pengalaman AR ini?",
			Options: models.Options{
				OptionA: "Memahami fotosintesis lebih baik",
				OptionB: "Melihat teknologi canggih",
				OptionC: "Bermain-main saja",
				OptionD: "Semua hal di atas",
			},
			CorrectAnswer: "A",
			Explanation:   "Tujuan utama AR Lab adalah membantu memahami fotosintesis dengan cara yang menyenangkan!",
		},
	}

	for _, prequiz := range prequizzes3 {
		if err := db.Create(&prequiz).Error; err != nil {
			log.Printf("Error creating prequiz for sub material 3: %v", err)
			continue
		}
	}

	// SubMaterial 4: Video Proses Fotosintesis (Offset 200, 440)
	subMaterial4 := models.SubMaterial{
		ModuleID:    module.ID,
		Title:       "Bagaimana Fotosintesis Terjadi?",
		Description: "Mari pelajari langkah-langkah fotosintesis dengan animasi yang mudah dipahami!",
		Order:       4,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(&subMaterial4).Error; err != nil {
		log.Printf("Error creating sub material 4: %v", err)
		return
	}

	// Video Material untuk SubMaterial 4
	videoMaterial4 := models.VideoMaterial{
		SubMaterialID: subMaterial4.ID,
		Title:         "Video Proses Fotosintesis",
		YoutubeLink:   "https://www.youtube.com/watch?v=example-fotosintesis-proses",
		Duration:      240, // 4 menit
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := db.Create(&videoMaterial4).Error; err != nil {
		log.Printf("Error creating video material 4: %v", err)
		return
	}

	// Video Quizzes untuk Video Material 4 (proses fotosintesis)
	videoQuizzes4 := []models.VideoQuiz{
		{
			VideoMaterialID: videoMaterial4.ID,
			Question:        "Berdasarkan video, di bagian mana fotosintesis terjadi?",
			TimestampStart:  45,  // Muncul di detik ke-45
			TimestampEnd:    60,  // Berakhir di detik ke-60
			Options: models.Options{
				OptionA: "Di akar",
				OptionB: "Di batang",
				OptionC: "Di kloroplas dalam daun",
				OptionD: "Di bunga",
			},
			CorrectAnswer: "C",
			Explanation:   "Benar! Fotosintesis terjadi di kloroplas yang terdapat dalam sel-sel daun!",
			Order:         1,
		},
		{
			VideoMaterialID: videoMaterial4.ID,
			Question:        "Apa hasil utama dari proses fotosintesis yang dijelaskan dalam video?",
			TimestampStart:  120, // Muncul di detik ke-120 (2 menit)
			TimestampEnd:    135, // Berakhir di detik ke-135
			Options: models.Options{
				OptionA: "Air dan tanah",
				OptionB: "Glukosa dan oksigen",
				OptionC: "Karbon dioksida dan nitrogen",
				OptionD: "Protein dan lemak",
			},
			CorrectAnswer: "B",
			Explanation:   "Tepat! Fotosintesis menghasilkan glukosa (makanan tumbuhan) dan oksigen yang kita hirup!",
			Order:         2,
		},
		{
			VideoMaterialID: videoMaterial4.ID,
			Question:        "Menurut video, mengapa fotosintesis penting bagi kehidupan?",
			TimestampStart:  200, // Muncul di detik ke-200
			TimestampEnd:    215, // Berakhir di detik ke-215
			Options: models.Options{
				OptionA: "Menghasilkan oksigen untuk bernapas",
				OptionB: "Membuat tumbuhan terlihat cantik",
				OptionC: "Menghasilkan suara",
				OptionD: "Membuat udara panas",
			},
			CorrectAnswer: "A",
			Explanation:   "Benar! Fotosintesis menghasilkan oksigen yang sangat penting untuk semua makhluk hidup bernapas!",
			Order:         3,
		},
	}

	for _, videoQuiz := range videoQuizzes4 {
		if err := db.Create(&videoQuiz).Error; err != nil {
			log.Printf("Error creating video quiz for video material 4: %v", err)
			continue
		}
	}

	// SubMaterial 5: Flashcards & Quiz Final (Offset 120, 570)
	subMaterial5 := models.SubMaterial{
		ModuleID:    module.ID,
		Title:       "Kartu Belajar & Kuis Akhir",
		Description: "Hafal kata-kata penting dan uji semua yang sudah kamu pelajari!",
		Order:       5,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(&subMaterial5).Error; err != nil {
		log.Printf("Error creating sub material 5: %v", err)
		return
	}

	// Flashcards untuk SubMaterial 5
	flashcards := []models.Flashcard{
		{
			SubMaterialID: subMaterial5.ID,
			FrontText:     "Fotosintesis",
			BackText:      "Proses tumbuhan membuat makanan menggunakan sinar matahari, air, dan udara",
			Order:         1,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			SubMaterialID: subMaterial5.ID,
			FrontText:     "Klorofil",
			BackText:      "Zat hijau di dalam daun yang menangkap sinar matahari",
			Order:         2,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			SubMaterialID: subMaterial5.ID,
			FrontText:     "Karbon Dioksida",
			BackText:      "Gas dari udara yang dibutuhkan tumbuhan untuk fotosintesis",
			Order:         3,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			SubMaterialID: subMaterial5.ID,
			FrontText:     "Oksigen",
			BackText:      "Gas yang dihasilkan tumbuhan saat fotosintesis, berguna untuk kita bernapas",
			Order:         4,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			SubMaterialID: subMaterial5.ID,
			FrontText:     "Stomata",
			BackText:      "Lubang kecil di daun tempat keluar masuknya gas",
			Order:         5,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			SubMaterialID: subMaterial5.ID,
			FrontText:     "Glukosa",
			BackText:      "Gula yang dibuat tumbuhan sebagai makanannya",
			Order:         6,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	for _, flashcard := range flashcards {
		if err := db.Create(&flashcard).Error; err != nil {
			log.Printf("Error creating flashcard: %v", err)
			continue
		}
	}

	// Quiz Questions Final untuk SubMaterial 5 (sebagai Questions dengan ModuleID)
	finalQuizQuestions := []models.Question{
		{
			ModuleID:   &module.ID,
			Question:   "Bahan apa saja yang dibutuhkan untuk fotosintesis?",
			AnswerTime: 60,
			ReadTime:   30,
			Options: models.Options{
				OptionA: "Air dan tanah",
				OptionB: "Sinar matahari, air, dan karbon dioksida",
				OptionC: "Pupuk dan air",
				OptionD: "Oksigen dan nitrogen",
			},
			CorrectAnswer: "B",
			Explanation:   "Fotosintesis membutuhkan 3 bahan utama: sinar matahari sebagai energi, air dari akar, dan karbon dioksida dari udara!",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ModuleID:   &module.ID,
			Question:   "Apa yang dihasilkan dari fotosintesis?",
			AnswerTime: 45,
			ReadTime:   25,
			Options: models.Options{
				OptionA: "Air dan tanah",
				OptionB: "Glukosa dan oksigen",
				OptionC: "Karbon dioksida dan air",
				OptionD: "Nitrogen dan fosfor",
			},
			CorrectAnswer: "B",
			Explanation:   "Fotosintesis menghasilkan glukosa (makanan tumbuhan) dan oksigen (untuk kita bernapas)!",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ModuleID:   &module.ID,
			Question:   "Di bagian mana fotosintesis terjadi?",
			AnswerTime: 30,
			ReadTime:   15,
			Options: models.Options{
				OptionA: "Akar",
				OptionB: "Batang",
				OptionC: "Daun",
				OptionD: "Bunga",
			},
			CorrectAnswer: "C",
			Explanation:   "Fotosintesis terjadi di daun karena daun mengandung klorofil yang bisa menangkap sinar matahari!",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ModuleID:   &module.ID,
			Question:   "Kapan fotosintesis bisa terjadi?",
			AnswerTime: 35,
			ReadTime:   20,
			Options: models.Options{
				OptionA: "Malam hari",
				OptionB: "Siang hari ada matahari",
				OptionC: "Kapan saja",
				OptionD: "Saat hujan",
			},
			CorrectAnswer: "B",
			Explanation:   "Fotosintesis hanya bisa terjadi saat ada sinar matahari, karena matahari adalah sumber energinya!",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ModuleID:   &module.ID,
			Question:   "Mengapa tumbuhan penting bagi kita?",
			AnswerTime: 40,
			ReadTime:   25,
			Options: models.Options{
				OptionA: "Karena cantik",
				OptionB: "Menghasilkan oksigen untuk bernapas",
				OptionC: "Membuat suara",
				OptionD: "Tidak penting",
			},
			CorrectAnswer: "B",
			Explanation:   "Tumbuhan sangat penting karena menghasilkan oksigen yang kita butuhkan untuk bernapas setiap hari!",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	for _, quiz := range finalQuizQuestions {
		if err := db.Create(&quiz).Error; err != nil {
			log.Printf("Error creating final quiz question: %v", err)
			continue
		}
	}

	log.Println("✅ Fotosintesis seeding completed successfully!")
	log.Printf("Created:")
	log.Printf("- 1 Module: %s", module.Title)
	log.Printf("- 5 SubMaterials with sequential order")
	log.Printf("- 2 Video Materials")
	log.Printf("- 1 AR Experiment (Link: %s)", arExperiment.LinkAR)
	log.Printf("- 6 Flashcards")
	log.Printf("- 7 Questions for the module")
	// Create dummy users and simulate their quiz answers
	if err := createDummyUsersAndAnswers(db); err != nil {
		log.Fatalf("Failed to create dummy users and answers: %v", err)
	}

	log.Println("Offset coordinates applied for learning path visualization:")
	log.Println("- SubMaterial 1: Video Intro (implicitly positioned)")
	log.Println("- SubMaterial 2: Quiz Dasar (implicitly positioned)")
	log.Printf("- SubMaterial 3: AR Lab (%g, %g)", arExperiment.OffsetX, arExperiment.OffsetY)
	log.Println("- SubMaterial 4: Video Proses (implicitly positioned)")
	log.Println("- SubMaterial 5: Flashcards & Quiz Final (implicitly positioned)")
}

func createDummyUsersAndAnswers(db *gorm.DB) error {
	log.Println("Creating dummy users and their quiz answers...")

	// Get user role
	var userRole models.Role
	if err := db.Where("name = ?", "user").First(&userRole).Error; err != nil {
		return fmt.Errorf("failed to find user role: %v", err)
	}

	// Create dummy users
	dummyUsers := []models.User{
		{
			Name:      "Andi Pratama",
			Email:     "andi.pratama@student.com",
			Password:  "$2a$14$hashedpassword1", // This would be properly hashed in real implementation
			RoleID:    userRole.ID,
			IsVerified: true,
		},
		{
			Name:      "Sari Dewi",
			Email:     "sari.dewi@student.com",
			Password:  "$2a$14$hashedpassword2",
			RoleID:    userRole.ID,
			IsVerified: true,
		},
		{
			Name:      "Budi Santoso",
			Email:     "budi.santoso@student.com",
			Password:  "$2a$14$hashedpassword3",
			RoleID:    userRole.ID,
			IsVerified: true,
		},
		{
			Name:      "Maya Putri",
			Email:     "maya.putri@student.com",
			Password:  "$2a$14$hashedpassword4",
			RoleID:    userRole.ID,
			IsVerified: true,
		},
		{
			Name:      "Riko Firmansyah",
			Email:     "riko.firmansyah@student.com",
			Password:  "$2a$14$hashedpassword5",
			RoleID:    userRole.ID,
			IsVerified: true,
		},
	}

	// Create users
	for i := range dummyUsers {
		if err := db.Create(&dummyUsers[i]).Error; err != nil {
			return fmt.Errorf("failed to create dummy user %s: %v", dummyUsers[i].Name, err)
		}
		log.Printf("Created dummy user: %s (ID: %d)", dummyUsers[i].Name, dummyUsers[i].ID)
	}

	// Get all prequizzes for fotosintesis module
	var prequizzes []models.Prequiz
	if err := db.Joins("JOIN sub_materials ON prequizzes.sub_material_id = sub_materials.id").
		Joins("JOIN modules ON sub_materials.module_id = modules.id").
		Where("modules.title LIKE ?", "%Fotosintesis%").
		Find(&prequizzes).Error; err != nil {
		return fmt.Errorf("failed to get prequizzes: %v", err)
	}

	// Get all video quizzes for fotosintesis videos
	var videoQuizzes []models.VideoQuiz
	if err := db.Joins("JOIN video_materials ON video_quizzes.video_material_id = video_materials.id").
		Joins("JOIN sub_materials ON video_materials.sub_material_id = sub_materials.id").
		Joins("JOIN modules ON sub_materials.module_id = modules.id").
		Where("modules.title LIKE ?", "%Fotosintesis%").
		Find(&videoQuizzes).Error; err != nil {
		return fmt.Errorf("failed to get video quizzes: %v", err)
	}

	// Simulate prequiz answers for each user
	for _, user := range dummyUsers {
		// Each user answers 70-90% of prequizzes correctly
		correctAnswerRate := 0.7 + (float64(user.ID%3) * 0.1) // 70%, 80%, or 90%
		
		for i, prequiz := range prequizzes {
			// Simulate whether this user answers correctly based on their rate
			isCorrect := float64(i%10) < (correctAnswerRate * 10)
			
			selectedAnswer := prequiz.CorrectAnswer
			if !isCorrect {
				// Pick a random wrong answer
				wrongAnswers := []string{"A", "B", "C", "D"}
				// Remove correct answer from options
				for j, ans := range wrongAnswers {
					if ans == prequiz.CorrectAnswer {
						wrongAnswers = append(wrongAnswers[:j], wrongAnswers[j+1:]...)
						break
					}
				}
				selectedAnswer = wrongAnswers[int(user.ID)%len(wrongAnswers)]
			}

			// Create prequiz user answer
			prequizAnswer := models.PrequizUserAnswer{
				UserID:     user.ID,
				PrequizID:  prequiz.ID,
				Answer:     selectedAnswer,
				IsCorrect:  isCorrect,
				AnsweredAt: time.Now().Add(-time.Duration(len(prequizzes)-i) * time.Hour).Unix(),
			}

			if err := db.Create(&prequizAnswer).Error; err != nil {
				log.Printf("Warning: Could not create prequiz answer for user %s: %v", user.Name, err)
			}
		}
		log.Printf("Created prequiz answers for user: %s", user.Name)
	}

	// Simulate video quiz answers for each user
	for _, user := range dummyUsers {
		// Each user answers 60-85% of video quizzes correctly (slightly harder)
		correctAnswerRate := 0.6 + (float64(user.ID%4) * 0.0625) // 60%, 66.25%, 72.5%, 78.75%, 85%
		
		for i, videoQuiz := range videoQuizzes {
			// Simulate whether this user answers correctly
			isCorrect := float64(i%8) < (correctAnswerRate * 8)
			
			selectedAnswer := videoQuiz.CorrectAnswer
			if !isCorrect {
				// Pick a random wrong answer
				wrongAnswers := []string{"A", "B", "C", "D"}
				for j, ans := range wrongAnswers {
					if ans == videoQuiz.CorrectAnswer {
						wrongAnswers = append(wrongAnswers[:j], wrongAnswers[j+1:]...)
						break
					}
				}
				selectedAnswer = wrongAnswers[int(user.ID)%len(wrongAnswers)]
			}

			// Create video quiz user answer
			videoQuizAnswer := models.VideoQuizUserAnswer{
				UserID:        user.ID,
				VideoQuizID:   videoQuiz.ID,
				SelectedAnswer: selectedAnswer,
				IsCorrect:     isCorrect,
				AnsweredAt:    time.Now().Add(-time.Duration(len(videoQuizzes)-i) * time.Hour * 2).Unix(),
				ResponseTime:  5 + int(user.ID%10), // Simulate 5-15 second response time
			}

			if err := db.Create(&videoQuizAnswer).Error; err != nil {
				log.Printf("Warning: Could not create video quiz answer for user %s: %v", user.Name, err)
			}
		}
		log.Printf("Created video quiz answers for user: %s", user.Name)
	}

	log.Printf("✅ Created %d dummy users with simulated quiz answers", len(dummyUsers))
	log.Println("Users created:")
	for _, user := range dummyUsers {
		log.Printf("- %s (%s)", user.Name, user.Email)
	}

	return nil
}
