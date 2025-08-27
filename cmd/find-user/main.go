package main

import (
	"fmt"

	"github.com/yazidalg/lidm_backend/internal/config"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()

	fmt.Printf("Finding user with Module 2 and 3 answers...\n")
	
	// Query to find users with Module 2 answers (3 prequizzes + 1 video quiz)
	rows, err := db.Raw(`
		SELECT pua.user_id, 
			COUNT(CASE WHEN p.module_id = 2 THEN 1 END) as module2_prequizzes,
			COUNT(CASE WHEN p.module_id = 3 THEN 1 END) as module3_prequizzes
		FROM prequiz_user_answers pua 
		JOIN prequizzes p ON pua.prequiz_id = p.id
		GROUP BY pua.user_id
		HAVING module2_prequizzes = 3 OR module3_prequizzes = 3
	`).Rows()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var userID, module2Prequizzes, module3Prequizzes int
		rows.Scan(&userID, &module2Prequizzes, &module3Prequizzes)
		fmt.Printf("User %d: Module2(pre:%d) Module3(pre:%d)\n", 
			userID, module2Prequizzes, module3Prequizzes)
	}
	
	// Also check video quiz answers separately
	fmt.Printf("\nVideo quiz answers by user:\n")
	rows2, err := db.Raw(`
		SELECT vqua.user_id, vm.module_id, COUNT(*) as count
		FROM video_quiz_user_answers vqua
		JOIN video_quizzes vq ON vqua.video_quiz_id = vq.id
		JOIN video_materials vm ON vq.video_material_id = vm.id
		GROUP BY vqua.user_id, vm.module_id
		ORDER BY vqua.user_id, vm.module_id
	`).Rows()
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer rows2.Close()
	
	for rows2.Next() {
		var userID, moduleID, count int
		rows2.Scan(&userID, &moduleID, &count)
		fmt.Printf("User %d, Module %d: %d video quiz answers\n", userID, moduleID, count)
	}
}
