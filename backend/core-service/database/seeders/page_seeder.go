package seeders

import (
	"core-service/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func RunPageSeeder(db *gorm.DB) error {
	homeSections := `[
  {
    "id": "hero-home",
    "type": "hero",
    "data": {
      "heading": "Kuasai Jalan, Dorong Masa Depan",
      "subheading": "Akademi Mengemudi Premium pertama di Alam Sutera menggunakan 100% Kendaraan Listrik. Rasakan pembelajaran yang halus, tenang, dan berkelanjutan.",
      "ctaText": "Pesan Sesi Pertama",
      "ctaLink": "/auth/register",
      "secondaryCtaText": "Lihat Paket",
      "secondaryCtaLink": "/packages",
      "bgImage": "https://images.unsplash.com/photo-1593941707882-a5bba14938c7?w=800&auto=format&fit=crop&q=80",
      "features": [
        { "title": "Electric Vehicle - 100% Listrik", "icon": "i-lucide-battery-charging" },
        { "title": "Kontrol Ganda", "icon": "i-lucide-shield-check" }
      ]
    }
  },
  {
    "id": "material-home",
    "type": "course_material",
    "data": {
      "headline": "Materi kursus yang akan Anda pelajari",
      "title": "Materi Kursus",
      "description": "Materi kursus yang Anda pelajari akan membuat Anda lebih percaya diri dalam mengemudi.",
      "materials": [
        {
          "title": "Teori Materi",
          "icon": "i-lucide-book-open",
          "description": [
            "Pengenalan Kendaraan dan Kontrol Dasar",
            "Pelatihan kokpit (posisi duduk ergonomis, penyesuaian spion tengah dan samping, penggunaan sabuk pengaman)",
            "Pengenalan instrumen (pedal gas, rem, tuas transmisi, rem tangan, lampu indikator di dashboard)",
            "Cek keselamatan (periksa kondisi ban, oli, dan air radiator sebelum berkendara)"
          ]
        },
        {
          "title": "Kontrol Awal",
          "icon": "i-lucide-shield-check",
          "description": [
            "Menghidupkan & menghentikan mesin (prosedur standar untuk menghidupkan mesin dengan aman)",
            "Teknik pedal gas dengan aman (seimbang dan halus)",
            "Teknik rem dan pemberhentian (pengereman halus dan cara berhenti di titik tertentu dengan tepat)"
          ]
        },
        {
          "title": "Teknik Manuver Dasar",
          "icon": "i-lucide-radar",
          "description": [
            "Kontrol kemudi (teknik memutar kemudi saat belok cepat)",
            "Mundur (mengendalikan mobil mundur hanya dengan spion belakang)",
            "Belok di persimpangan (teknik mengambil sudut belok yang benar ke kiri atau kanan)"
          ]
        },
        {
          "title": "Teknik Mengemudi di Jalan Naik & Turun",
          "icon": "i-lucide-car",
          "description": [
            "Teknik start-stop di tanjakan",
            "Teknik start-stop di turunan"
          ]
        },
        {
          "title": "Teknik Parkir",
          "icon": "i-lucide-car",
          "description": [
            "Parkir mundur sudut atau lurus (masuk slot parkir dengan mobil mundur)",
            "Parkir paralel (teknik memasukkan mobil di antara dua mobil lain yang paralel)"
          ]
        },
        {
          "title": "Mengemudi di Jalan Tol",
          "icon": "i-lucide-car",
          "description": [
            "Rambu dan marka jalan (mematuhi rambu lalu lintas, rambu dilarang parkir dan marka jalan)",
            "Etika berkendara (menggunakan sein, menjaga jarak aman, dan cara menyalip kendaraan lain dengan benar)",
            "Blind spot (teknik memeriksa area yang tidak terlihat di spion sebelum berganti jalur)"
          ]
        }
      ]
    }
  }
]`

	pages := []models.Page{
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Title:     "Home",
			Slug:      "/",
			Status:    models.PageStatusPublished,
			Sections:  homeSections,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Title:     "About Us",
			Slug:      "/about",
			Status:    models.PageStatusDraft,
			Sections:  "[]",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, page := range pages {
		var count int64
		db.Model(&models.Page{}).Where("id = ? OR slug = ?", page.ID, page.Slug).Count(&count)
		if count == 0 {
			if err := db.Create(&page).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
