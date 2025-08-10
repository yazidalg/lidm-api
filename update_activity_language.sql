-- SQL script untuk update activity yang masih dalam bahasa Inggris ke Indonesia

-- Update ActivityType ke bahasa Indonesia
UPDATE user_activities SET activity_type = 'masuk' WHERE activity_type = 'login';
UPDATE user_activities SET activity_type = 'keluar' WHERE activity_type = 'logout';
UPDATE user_activities SET activity_type = 'gabung_kuis' WHERE activity_type = 'quiz_join';
UPDATE user_activities SET activity_type = 'selesai_kuis' WHERE activity_type = 'quiz_complete';
UPDATE user_activities SET activity_type = 'jawab_kuis' WHERE activity_type = 'quiz_answer';
UPDATE user_activities SET activity_type = 'lihat_pelajaran' WHERE activity_type = 'lesson_view';
UPDATE user_activities SET activity_type = 'selesai_pelajaran' WHERE activity_type = 'lesson_complete';
UPDATE user_activities SET activity_type = 'lihat_modul' WHERE activity_type = 'module_view';
UPDATE user_activities SET activity_type = 'selesai_modul' WHERE activity_type = 'module_complete';
UPDATE user_activities SET activity_type = 'update_profil' WHERE activity_type = 'profile_update';

-- Update Description ke bahasa Indonesia
UPDATE user_activities SET description = 'Pengguna berhasil masuk' WHERE description = 'User logged in successfully';
UPDATE user_activities SET description = 'Pengguna masuk dengan Google' WHERE description = 'User logged in with Google';
UPDATE user_activities SET description = 'Pengguna masuk dengan akun Belajar' WHERE description = 'User logged in with Belajar account';
UPDATE user_activities SET description = 'Pengguna keluar' WHERE description = 'User logged out';
UPDATE user_activities SET description = 'Bergabung dengan sesi kuis' WHERE description = 'Joined a quiz session';
UPDATE user_activities SET description = 'Menyelesaikan kuis' WHERE description = 'Completed a quiz';
UPDATE user_activities SET description = 'Menjawab pertanyaan kuis' WHERE description = 'Answered a quiz question';
UPDATE user_activities SET description = 'Melihat pelajaran' WHERE description = 'Viewed a lesson';
UPDATE user_activities SET description = 'Menyelesaikan pelajaran' WHERE description = 'Completed a lesson';
UPDATE user_activities SET description = 'Melihat modul' WHERE description = 'Viewed a module';
UPDATE user_activities SET description = 'Menyelesaikan modul' WHERE description = 'Completed a module';
UPDATE user_activities SET description = 'Memperbarui informasi profil' WHERE description = 'Updated profile information';

-- Verifikasi perubahan
SELECT DISTINCT activity_type, description FROM user_activities ORDER BY activity_type;
