-- Database trigger untuk auto-unlock module berikutnya
-- Trigger ini akan dijalankan setiap kali module_progresses diupdate

DELIMITER $$

DROP TRIGGER IF EXISTS auto_unlock_next_module$$

CREATE TRIGGER auto_unlock_next_module
    AFTER UPDATE ON module_progresses
    FOR EACH ROW
BEGIN
    -- Hanya jalankan jika module baru saja complete (dari false ke true)
    IF NEW.is_complete = 1 AND OLD.is_complete = 0 THEN
        
        -- Cari module selanjutnya berdasarkan ID
        SET @next_module_id = (
            SELECT id 
            FROM modules 
            WHERE id > NEW.module_id 
              AND deleted_at IS NULL 
            ORDER BY id ASC 
            LIMIT 1
        );
        
        -- Jika ada module selanjutnya
        IF @next_module_id IS NOT NULL THEN
            
            -- Check apakah progress record sudah ada
            SET @existing_progress = (
                SELECT id 
                FROM module_progresses 
                WHERE user_id = NEW.user_id 
                  AND module_id = @next_module_id 
                  AND deleted_at IS NULL
            );
            
            -- Jika belum ada, buat record baru
            IF @existing_progress IS NULL THEN
                INSERT INTO module_progresses (
                    user_id, 
                    module_id, 
                    is_unlocked, 
                    is_complete, 
                    progress, 
                    created_at, 
                    updated_at
                ) VALUES (
                    NEW.user_id, 
                    @next_module_id, 
                    1, 
                    0, 
                    0, 
                    NOW(), 
                    NOW()
                );
            ELSE
                -- Jika sudah ada, unlock saja
                UPDATE module_progresses 
                SET is_unlocked = 1, updated_at = NOW() 
                WHERE id = @existing_progress;
            END IF;
            
        END IF;
        
    END IF;
END$$

DELIMITER ;
