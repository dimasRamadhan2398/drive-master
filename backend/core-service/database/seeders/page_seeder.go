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

	aboutSections := `[
  {
    "id": "about-hero",
    "type": "hero",
    "data": {
      "heading": "Langkah Maju, Belajar dari Masa Depan",
      "subheading": "Drive Master bukan hanya tentang mengajar cara mengemudi, ini tentang mendefinisikan ulang standar pendidikan mengemudi di Indonesia. Sebagai pelopor sekolah mengemudi Kendaraan Listrik, kami percaya bahwa pengemudi masa depan harus lahir dari teknologi masa depan—modern, cerdas, dan ramah lingkungan.",
      "ctaText": "Mulai Perjalanan Anda",
      "ctaLink": "/auth/register",
      "secondaryCtaText": "Hubungi Kami",
      "secondaryCtaLink": "/contact",
      "features": [
        { "title": "Pelopor EV Bebas Emisi", "icon": "i-lucide-leaf" },
        { "title": "Instruktur Bersertifikat", "icon": "i-lucide-award" }
      ]
    }
  },
  {
    "id": "about-safety",
    "type": "specifications",
    "data": {
      "headline": "Keselamatan Utama",
      "title": "Prioritas Kami adalah Keselamatan Anda",
      "description": "Dalam bisnis kursus mengemudi, keselamatan bukan hanya sebuah fitur; itu adalah fondasi inti kami.",
      "items": [
        {
          "title": "Instruktur Bersertifikat",
          "subtitle": "Instruktur kami adalah profesional berlisensi yang bersertifikat khusus untuk mengoperasikan kendaraan listrik premium.",
          "icon": "i-lucide-award",
          "description": []
        },
        {
          "title": "Teknologi Keselamatan Aktif",
          "subtitle": "Memanfaatkan fitur keselamatan bawaan EV seperti Collision Avoidance dan Blind Spot Monitoring untuk meminimalkan risiko.",
          "icon": "i-lucide-radar",
          "description": []
        }
      ]
    }
  },
  {
    "id": "about-quote",
    "type": "quote",
    "data": {
      "quote": "Visi kami bukan hanya untuk menghasilkan pengemudi yang bisa memutar kemudi, tetapi untuk membina pengemudi yang cerdas dan aman yang siap merangkul era elektrifikasi.",
      "description": "Di Drive Master Indonesia, kami percaya bahwa cara kita belajar mengemudi harus berevolusi seiring dengan evolusi teknologi otomotif. Kami berkomitmen untuk menjadi standar baru dalam pendidikan mengemudi yang ramah lingkungan, memastikan bahwa setiap lulusan memiliki keterampilan mengemudi tingkat tinggi serta kesadaran akan masa depan mobilitas yang berkelanjutan.",
      "ctaText": "Mulai Perjalanan Anda",
      "ctaLink": "/auth/register",
      "secondaryCtaText": "Chat WhatsApp",
      "secondaryCtaLink": "https://wa.me/628119124848"
    }
  }
]`

	servicesSections := `[
  {
    "id": "hero-services",
    "type": "hero",
    "data": {
      "heading": "Layanan Drive Master Academy",
      "subheading": "Kursus mengemudi komprehensif yang dirancang untuk masa depan listrik. Dari pemula hingga pengemudi tingkat lanjut, kami memiliki program yang sempurna untuk Anda.",
      "ctaText": "Lihat Paket",
      "ctaLink": "/packages",
      "secondaryCtaText": "Konsultasi WA",
      "secondaryCtaLink": "https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi",
      "bgImage": "https://images.unsplash.com/photo-1593941707882-a5bba14938c7?w=800&auto=format&fit=crop&q=80",
      "features": [
        { "title": "100% Kendaraan Listrik", "icon": "i-lucide-battery-charging" },
        { "title": "Instruktur Bersertifikat", "icon": "i-lucide-award" }
      ]
    }
  },
  {
    "id": "material-services",
    "type": "course_material",
    "data": {
      "headline": "Materi kursus yang akan Anda pelajari",
      "title": "Materi Kursus Komprehensif",
      "description": "Kurikulum pelatihan terstruktur kami mencakup dasar-dasar hingga manuver tingkat lanjut.",
      "materials": [
        {
          "title": "Teori & Pengenalan",
          "icon": "i-lucide-book-open",
          "description": [
            "Pengenalan instrumen & kontrol kendaraan",
            "Posisi berkendara ergonomis & keselamatan dasar",
            "Pemeriksaan keselamatan kendaraan (pre-drive checks)"
          ]
        },
        {
          "title": "Kontrol Kendaraan",
          "icon": "i-lucide-shield-check",
          "description": [
            "Akselerasi & deselerasi halus",
            "Pengereman presisi & darurat",
            "Penggunaan gigi & transmisi matic"
          ]
        },
        {
          "title": "Manuver Jalan Raya",
          "icon": "i-lucide-radar",
          "description": [
            "Teknik belok & persimpangan",
            "Mengemudi di jalan tanjakan & turunan",
            "Teknik mundur dengan cermin"
          ]
        },
        {
          "title": "Parkir & Tol",
          "icon": "i-lucide-car",
          "description": [
            "Parkir mundur (slot parkir)",
            "Parkir paralel di antara kendaraan",
            "Navigasi jalan tol & pindah jalur aman"
          ]
        }
      ]
    }
  },
  {
    "id": "cta-services",
    "type": "cta",
    "data": {
      "heading": "Siap untuk menguasai jalanan bersama kami?",
      "buttonText": "Daftar Sekarang",
      "buttonLink": "/auth/register"
    }
  }
]`

	contactSections := `[
  {
    "id": "contact-hero",
    "type": "hero",
    "data": {
      "heading": "Kami di Sini untuk Membantu",
      "subheading": "Punya pertanyaan tentang paket mengemudi EV kami atau penjadwalan? Hubungi tim kami melalui metode di bawah ini.",
      "ctaText": "Chat WhatsApp",
      "ctaLink": "https://wa.me/628119124848",
      "secondaryCtaText": "Lihat FAQ",
      "secondaryCtaLink": "/#faq"
    }
  },
  {
    "id": "contact-cta",
    "type": "cta",
    "data": {
      "heading": "Hubungi Customer Service",
      "description": "Tim kami siap memberikan informasi detail mengenai jadwal, paket, dan lokasi pelatihan.",
      "buttonText": "Lihat Paket",
      "buttonLink": "/packages",
      "buttonIcon": "i-lucide-package",
      "secondaryButtonText": "Chat WhatsApp Sekarang",
      "secondaryButtonLink": "https://wa.me/628119124848",
      "secondaryButtonIcon": "i-simple-icons-whatsapp"
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
			Status:    models.PageStatusPublished,
			Sections:  aboutSections,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			Title:     "Services",
			Slug:      "/services",
			Status:    models.PageStatusPublished,
			Sections:  servicesSections,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			Title:     "Contact Us",
			Slug:      "/contact",
			Status:    models.PageStatusPublished,
			Sections:  contactSections,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, page := range pages {
		var existing models.Page
		err := db.Where("id = ? OR slug = ?", page.ID, page.Slug).First(&existing).Error
		if err != nil {
			if err := db.Create(&page).Error; err != nil {
				return err
			}
		} else {
			if existing.Sections == "" || existing.Sections == "[]" || existing.Status == models.PageStatusDraft {
				db.Model(&existing).Updates(map[string]interface{}{
					"title":    page.Title,
					"status":   page.Status,
					"sections": page.Sections,
				})
			}
		}
	}

	return nil
}
