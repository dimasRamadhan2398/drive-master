package cli

import (
	"log"
	"slices"

	"user-service/models"

	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run database migrations to create or update database schema.`,
	Run:   runMigrate,
}

var (
	migrateDown  bool
	migrateSteps int
	migrateDryRun bool
	migrateReset bool
)

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().BoolVar(&migrateReset, "reset", false, "Drop and recreate all tables")
	migrateCmd.Flags().BoolVar(&migrateDown, "down", false, "Roll back migrations")
	migrateCmd.Flags().IntVarP(&migrateSteps, "steps", "n", 0, "Number of migrations to run (positive for up, negative for down)")
	migrateCmd.Flags().BoolVar(&migrateDryRun, "dry-run", false, "Show what would be migrated without running")
}

func migrateTestimonialIDToUUID(db *gorm.DB) error {
    // 1. Add a new UUID column (nullable at first)
    if err := db.Exec(`
        ALTER TABLE testimonials
        ADD COLUMN IF NOT EXISTS new_id UUID
    `).Error; err != nil {
        return err
    }

    // 2. Generate UUIDs for existing rows
    if err := db.Exec(`
        UPDATE testimonials
        SET new_id = gen_random_uuid()
        WHERE new_id IS NULL
    `).Error; err != nil {
        return err
    }

    // 3. Handle dependent tables (example: testimonial_media)
    // Add a new FK column in the referencing table
    if err := db.Exec(`
        ALTER TABLE testimonial_media
        ADD COLUMN IF NOT EXISTS new_testimonial_id UUID
    `).Error; err != nil {
        return err
    }

    // Populate it based on the old integer id
    if err := db.Exec(`
        UPDATE testimonial_media tm
        SET new_testimonial_id = t.new_id
        FROM testimonials t
        WHERE tm.testimonial_id = t.id
    `).Error; err != nil {
        return err
    }

    // 4. Drop old foreign key constraints (if any)
    //    You'll need to know the actual constraint names – query them from
    //    information_schema or use a naming convention.
    //    Example:
    // db.Exec(`ALTER TABLE testimonial_media DROP CONSTRAINT fk_testimonial_media_testimonial`)

    // 5. Drop the old integer columns and rename the new ones
    //    Because this is risky, you may want to do it in a transaction.

    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Drop old FK column and rename new one in referencing table
    if err := tx.Exec(`ALTER TABLE testimonial_media DROP COLUMN testimonial_id`).Error; err != nil {
        tx.Rollback()
        return err
    }
    if err := tx.Exec(`ALTER TABLE testimonial_media RENAME COLUMN new_testimonial_id TO testimonial_id`).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Now swap the primary key on testimonials
    // Drop the old primary key constraint (this may cascade, be careful)
    if err := tx.Exec(`ALTER TABLE testimonials DROP CONSTRAINT testimonials_pkey`).Error; err != nil {
        tx.Rollback()
        return err
    }
    if err := tx.Exec(`ALTER TABLE testimonials DROP COLUMN id`).Error; err != nil {
        tx.Rollback()
        return err
    }
    if err := tx.Exec(`ALTER TABLE testimonials RENAME COLUMN new_id TO id`).Error; err != nil {
        tx.Rollback()
        return err
    }
    if err := tx.Exec(`ALTER TABLE testimonials ADD PRIMARY KEY (id)`).Error; err != nil {
        tx.Rollback()
        return err
    }

    // 6. Recreate foreign key constraints (adjust to your model)
    if err := tx.Exec(`
        ALTER TABLE testimonial_media
        ADD CONSTRAINT fk_testimonial_media_testimonial
        FOREIGN KEY (testimonial_id) REFERENCES testimonials(id)
    `).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}

func runMigrate(cmd *cobra.Command, args []string) {
	// Load config first
	LoadConfig()

	db, err := gorm.Open(postgres.Open(getDSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if migrateDryRun {
		log.Println("Dry run mode - no changes will be made")
		log.Println("Would run migrations for the following models:")
		modelsList := []string{
			"User", "Role", "MemberProfile", "InstructorProfile",
			"WorkExperience", "InstructorArea", "Certification", "Entitlement",
			"Testimonial", "TestimonialMedia",
		}
		for _, m := range modelsList {
			log.Printf("  - %s", m)
		}
		return
	}

	if migrateReset {
		log.Println("Resetting database - dropping all tables...")
		if err := dropAllTables(db); err != nil {
			log.Fatalf("Failed to drop tables: %v", err)
		}
		log.Println("All tables dropped successfully")
	}

	log.Println("Running database migrations...")

	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.MemberProfile{},
		&models.InstructorProfile{},
		&models.WorkExperience{},
		&models.InstructorArea{},
		&models.Certification{},
		&models.Entitlement{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}

	if err := migrateTestimonialIDToUUID(db); err != nil {
    	log.Fatal("migration failed: ", err)
	}
		// After the migration, AutoMigrate will recognise the new schema
	db.AutoMigrate(&models.Testimonial{}, &models.TestimonialMedia{})
	

	log.Println("Database migrations completed successfully")
}

// dropAllTables drops all tables in the database using GORM
func dropAllTables(db *gorm.DB) error {
	// Get all table names from the database
	var tables []string

	// Use the underlying SQL DB to get tables
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	rows, err := sqlDB.Query(`
		SELECT tablename FROM pg_tables
		WHERE schemaname = 'public'
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return err
		}
		tables = append(tables, tableName)
	}

	// Drop tables in reverse dependency order to avoid FK issues
	// Order matters for foreign key dependencies
	order := []string{
		// "users",                    // depends on nothing
		// "instructor_areas",         // depends on regions
		// "work_experiences",          // depends on users
		// "member_profiles",           // depends on users
		// "instructor_profiles",       // depends on users
		"roles",                    // depends on nothing
	}

	for _, table := range order {
		if err := db.Exec("DROP TABLE IF EXISTS " + table + " CASCADE").Error; err != nil {
			log.Printf("Warning: Failed to drop table %s: %v", table, err)
		} else {
			log.Printf("Dropped table: %s", table)
		}
	}

	// Also drop any tables not in our list
	for _, table := range tables {
		if !slices.Contains(order, table) {
			if err := db.Exec("DROP TABLE IF EXISTS " + table + " CASCADE").Error; err != nil {
				log.Printf("Warning: Failed to drop table %s: %v", table, err)
			} else {
				log.Printf("Dropped table: %s", table)
			}
		}
	}

	return nil
}