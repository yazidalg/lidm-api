package main

import (
	"log"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/config"
)

func main() {
	log.Println("Starting prequiz data seeding...")

	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()

	// Clear existing prequiz data
	if err := db.Exec("DELETE FROM prequizzes").Error; err != nil {
		log.Printf("Error clearing prequizzes: %v", err)
	}

	// Module 1 - Materi 1: Fotosintesis Dasar
	module1Questions := []models.Prequiz{
		{
			ModuleID: 1,
			Question: "Seorang siswa menaruh tanaman di dalam ruangan gelap selama seminggu. Daunnya menjadi pucat dan tanaman tampak layu. Apa penyebab utama kondisi tersebut berdasarkan proses fotosintesis?",
			Options: models.Options{
				OptionA: "Tanaman kekurangan air",
				OptionB: "Tidak ada cahaya matahari sebagai sumber energi",
				OptionC: "Tanaman tidak menyerap karbon dioksida",
				OptionD: "Tanaman terlalu banyak oksigen",
			},
			CorrectAnswer: "B",
			Explanation:   "Fotosintesis memerlukan cahaya matahari sebagai sumber energi. Tanpa cahaya, tumbuhan tidak dapat melakukan fotosintesis dan akan layu.",
		},
		{
			ModuleID: 1,
			Question: "Apa arti kata \"synthesis\" dalam istilah fotosintesis?",
			Options: models.Options{
				OptionA: "Cahaya",
				OptionB: "Menggabungkan atau menyusun",
				OptionC: "Energi",
				OptionD: "Warna hijau",
			},
			CorrectAnswer: "B",
			Explanation:   "Kata 'synthesis' berarti menggabungkan atau menyusun. Fotosintesis adalah proses menyusun makanan dengan bantuan cahaya.",
		},
		{
			ModuleID: 1,
			Question: "Mengapa tumbuhan disebut makhluk hidup autotrof?",
			Options: models.Options{
				OptionA: "Karena bisa bergerak mencari makanannya",
				OptionB: "Karena bisa membuat makanannya sendiri",
				OptionC: "Karena selalu membutuhkan hewan untuk hidup",
				OptionD: "Karena tidak membutuhkan air",
			},
			CorrectAnswer: "B",
			Explanation:   "Autotrof berarti dapat membuat makanan sendiri. Tumbuhan dapat menghasilkan makanan melalui proses fotosintesis.",
		},
	}

	// Module 2 - Materi 2: Bagian-bagian Tumbuhan dalam Fotosintesis
	module2Questions := []models.Prequiz{
		{
			ModuleID: 2,
			Question: "Bagian tumbuhan yang berfungsi utama dalam proses fotosintesis adalah…",
			Options: models.Options{
				OptionA: "Akar",
				OptionB: "Batang",
				OptionC: "Daun",
				OptionD: "Bunga",
			},
			CorrectAnswer: "C",
			Explanation:   "Daun adalah bagian utama tempat fotosintesis terjadi karena mengandung klorofil dan stomata.",
		},
		{
			ModuleID: 2,
			Question: "Mengapa daun disebut sebagai \"dapur makanan tumbuhan\"?",
			Options: models.Options{
				OptionA: "Karena warnanya hijau",
				OptionB: "Karena di dalamnya terjadi fotosintesis",
				OptionC: "Karena bentuknya lebar",
				OptionD: "Karena tempat menyimpan cadangan makanan",
			},
			CorrectAnswer: "B",
			Explanation:   "Daun disebut dapur makanan karena di sinilah proses pembuatan makanan (fotosintesis) berlangsung.",
		},
		{
			ModuleID: 2,
			Question: "Zat hijau daun disebut…",
			Options: models.Options{
				OptionA: "Stomata",
				OptionB: "Klorofil",
				OptionC: "Xilem",
				OptionD: "Floem",
			},
			CorrectAnswer: "B",
			Explanation:   "Klorofil adalah zat hijau yang terdapat dalam daun dan berfungsi menangkap cahaya matahari untuk fotosintesis.",
		},
	}

	// Module 3 - Materi 3: Faktor-faktor yang Mempengaruhi Fotosintesis
	module3Questions := []models.Prequiz{
		{
			ModuleID: 3,
			Question: "Faktor utama yang menjadi sumber energi dalam fotosintesis adalah…",
			Options: models.Options{
				OptionA: "Air",
				OptionB: "Cahaya matahari",
				OptionC: "Klorofil",
				OptionD: "Karbon dioksida",
			},
			CorrectAnswer: "B",
			Explanation:   "Cahaya matahari adalah sumber energi utama yang diperlukan untuk menggerakkan proses fotosintesis.",
		},
		{
			ModuleID: 3,
			Question: "Gas karbon dioksida masuk ke dalam daun melalui…",
			Options: models.Options{
				OptionA: "Akar",
				OptionB: "Stomata",
				OptionC: "Batang",
				OptionD: "Kloroplas",
			},
			CorrectAnswer: "B",
			Explanation:   "Stomata adalah pori-pori kecil pada daun yang menjadi jalan masuk karbon dioksida dari udara.",
		},
		{
			ModuleID: 3,
			Question: "Suhu yang terlalu tinggi dapat…",
			Options: models.Options{
				OptionA: "Mempercepat kerja enzim",
				OptionB: "Menghentikan fotosintesis",
				OptionC: "Mengubah klorofil",
				OptionD: "Membentuk lebih banyak air",
			},
			CorrectAnswer: "B",
			Explanation:   "Suhu yang terlalu tinggi dapat merusak enzim dan menghentikan proses fotosintesis.",
		},
	}

	// Module 4 - Materi 4: Proses dan Hasil Fotosintesis
	module4Questions := []models.Prequiz{
		{
			ModuleID: 4,
			Question: "Proses fotosintesis berlangsung di dalam organel bernama…",
			Options: models.Options{
				OptionA: "Nukleus",
				OptionB: "Kloroplas",
				OptionC: "Mitokondria",
				OptionD: "Sitoplasma",
			},
			CorrectAnswer: "B",
			Explanation:   "Kloroplas adalah organel dalam sel tumbuhan tempat fotosintesis berlangsung.",
		},
		{
			ModuleID: 4,
			Question: "Reaksi fotosintesis menghasilkan glukosa dan…",
			Options: models.Options{
				OptionA: "Karbon dioksida",
				OptionB: "Air",
				OptionC: "Oksigen",
				OptionD: "Nitrogen",
			},
			CorrectAnswer: "C",
			Explanation:   "Fotosintesis menghasilkan glukosa sebagai makanan dan oksigen sebagai produk sampingan.",
		},
		{
			ModuleID: 4,
			Question: "Rumus kimia glukosa adalah…",
			Options: models.Options{
				OptionA: "H₂O",
				OptionB: "O₂",
				OptionC: "CO₂",
				OptionD: "C₆H₁₂O₆",
			},
			CorrectAnswer: "D",
			Explanation:   "C₆H₁₂O₆ adalah rumus kimia glukosa, molekul gula yang dihasilkan fotosintesis.",
		},
	}

	// Module 5 - Materi 5: Manfaat Fotosintesis
	module5Questions := []models.Prequiz{
		{
			ModuleID: 5,
			Question: "Produk utama yang dihasilkan dari fotosintesis adalah…",
			Options: models.Options{
				OptionA: "Air dan karbondioksida",
				OptionB: "Glukosa dan oksigen",
				OptionC: "Nitrogen dan glukosa",
				OptionD: "Mineral dan oksigen",
			},
			CorrectAnswer: "B",
			Explanation:   "Fotosintesis menghasilkan glukosa sebagai makanan dan oksigen yang penting bagi kehidupan.",
		},
		{
			ModuleID: 5,
			Question: "Glukosa yang dihasilkan tumbuhan digunakan untuk…",
			Options: models.Options{
				OptionA: "Menyerap cahaya matahari",
				OptionB: "Menyerap air dari tanah",
				OptionC: "Tumbuh, berkembang, dan menyimpan cadangan makanan",
				OptionD: "Menghirup oksigen dari udara",
			},
			CorrectAnswer: "C",
			Explanation:   "Glukosa digunakan tumbuhan untuk pertumbuhan, perkembangan, dan disimpan sebagai cadangan energi.",
		},
		{
			ModuleID: 5,
			Question: "Fungsi utama oksigen hasil fotosintesis bagi manusia adalah…",
			Options: models.Options{
				OptionA: "Menyerap cahaya matahari",
				OptionB: "Bernapas",
				OptionC: "Membentuk cadangan makanan",
				OptionD: "Menyerap air",
			},
			CorrectAnswer: "B",
			Explanation:   "Oksigen hasil fotosintesis sangat penting bagi manusia dan hewan untuk proses pernapasan.",
		},
	}

	// Module 6 - Materi 6: Peran Fotosintesis dalam Kehidupan
	module6Questions := []models.Prequiz{
		{
			ModuleID: 6,
			Question: "Proses fotosintesis pada tumbuhan menghasilkan gas penting bagi manusia, yaitu…",
			Options: models.Options{
				OptionA: "Karbon dioksida (CO₂)",
				OptionB: "Oksigen (O₂)",
				OptionC: "Nitrogen (N₂)",
				OptionD: "Uap air (H₂O)",
			},
			CorrectAnswer: "B",
			Explanation:   "Fotosintesis menghasilkan oksigen yang sangat penting bagi pernapasan manusia dan hewan.",
		},
		{
			ModuleID: 6,
			Question: "Mengapa tumbuhan disebut produsen?",
			Options: models.Options{
				OptionA: "Karena menghasilkan oksigen",
				OptionB: "Karena menghasilkan makanan sendiri",
				OptionC: "Karena bisa berjalan mencari makanan",
				OptionD: "Karena tidak membutuhkan air",
			},
			CorrectAnswer: "B",
			Explanation:   "Tumbuhan disebut produsen karena dapat menghasilkan makanan sendiri melalui fotosintesis.",
		},
		{
			ModuleID: 6,
			Question: "Tanpa fotosintesis, kehidupan di bumi…",
			Options: models.Options{
				OptionA: "Akan tetap sama",
				OptionB: "Tidak akan ada seperti sekarang",
				OptionC: "Semakin mudah",
				OptionD: "Semakin banyak oksigen",
			},
			CorrectAnswer: "B",
			Explanation:   "Fotosintesis adalah dasar kehidupan di bumi karena menghasilkan oksigen dan makanan untuk semua makhluk hidup.",
		},
	}

	// Combine all questions
	allQuestions := append(module1Questions, module2Questions...)
	allQuestions = append(allQuestions, module3Questions...)
	allQuestions = append(allQuestions, module4Questions...)
	allQuestions = append(allQuestions, module5Questions...)
	allQuestions = append(allQuestions, module6Questions...)

	// Insert all questions
	for i, question := range allQuestions {
		if err := db.Create(&question).Error; err != nil {
			log.Printf("Error creating question %d: %v", i+1, err)
		} else {
			log.Printf("Created question %d for module %d", i+1, question.ModuleID)
		}
	}

	log.Printf("Prequiz seeding completed. Total questions created: %d", len(allQuestions))
}
