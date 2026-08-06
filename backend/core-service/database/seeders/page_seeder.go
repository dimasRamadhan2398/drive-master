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
      "features": [
        { "title": "Pelopor Pendidikan Mengemudi Bebas Emisi", "icon": "i-lucide-leaf" }
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
    "id": "specifications-services",
    "type": "specifications",
    "data": {
      "headline": "",
      "title": "",
      "description": "",
      "items": [
        {
          "title": "Layanan Khusus",
          "subtitle": "Keuntungan memilih Kami",
          "icon": "i-lucide-star",
          "description": [
            "SIM",
            "Antar-jemput gratis (Alam Sutera, BSD, Gading Serpong)",
            "Materi teori termasuk",
            "Sertifikat penyelesaian"
          ]
        },
        {
          "title": "Waktu Sesi",
          "subtitle": "Waktu sesi asli kami",
          "icon": "i-lucide-car",
          "description": [
            "Sesi di hari kerja buka dari jam 08:00 sampai 17:00",
            "60 menit per sesi",
            "Sesi Malam: 18:00 - 20:00",
            "Sabtu - Minggu: 08:00 - 17:00"
          ]
        },
        {
          "title": "Transmisi Mobil",
          "subtitle": "Mobil kami menggunakan transmisi matic",
          "icon": "i-lucide-car",
          "description": [
            "Transmisi Matic"
          ]
        },
        {
          "title": "Sesi Malam",
          "subtitle": "Harga tambahan untuk sesi malam",
          "icon": "i-lucide-moon",
          "description": [
            "Sesi di malam hari buka dari jam 18:00 sampai 20:00",
            "Harga dasar 6x + Rp.100.000 untuk sesi malam",
            "Harga dasar 8x + Rp.100.000 untuk sesi malam",
            "Harga dasar 10x + Rp.100.000 untuk sesi malam",
            "Harga dasar 12x + Rp.100.000 untuk sesi malam"
          ]
        },
        {
          "title": "Sesi Akhir Pekan",
          "subtitle": "Harga tambahan untuk sesi akhir pekan",
          "icon": "i-lucide-clock",
          "description": [
            "Sesi di akhir pekan buka dari jam 08:00 sampai 17:00",
            "Harga dasar 6x + Rp.100.000 untuk sesi akhir pekan",
            "Harga dasar 8x + Rp.100.000 untuk sesi akhir pekan",
            "Harga dasar 10x + Rp.100.000 untuk sesi akhir pekan",
            "Harga dasar 12x + Rp.100.000 untuk sesi akhir pekan"
          ]
        }
      ]
    }
  },
  {
    "id": "material-services",
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
  },
  {
    "id": "coverage-services",
    "type": "service_areas",
    "data": {
      "headline": "Jangkauan Layanan",
      "title": "Area Operasional Kami",
      "description": "Kami melayani antar-jemput gratis di berbagai wilayah berikut",
      "footer": "Area Anda belum tertera? Hubungi tim kami untuk informasi lebih lanjut.",
      "areas": [
        "Alam Sutera & sekitarnya",
        "Serpong & BSD City",
        "Tangerang City Center",
        "Gading Serpong",
        "Lippo Karawaci",
        "Bintaro Jaya (terbatas)"
      ]
    }
  },
  {
    "id": "cta-services",
    "type": "cta",
    "data": {
      "heading": "Siap Memulai Perjalanan Mengemudi Anda?",
      "description": "Bergabunglah bersama ratusan siswa yang telah berhasil mendapatkan SIM & mahir mengemudi bersama instruktur bersertifikat.",
      "buttonText": "Lihat Paket Sekarang",
      "buttonLink": "/packages",
      "buttonIcon": "i-lucide-package",
      "secondaryButtonText": "Konsultasi WhatsApp",
      "secondaryButtonLink": "https://wa.me/628119124848",
      "secondaryButtonIcon": "i-simple-icons-whatsapp"
    }
  }
]`

	contactSections := `[
  {
    "id": "hero-contact",
    "type": "hero",
    "data": {
      "heading": "Kami di Sini untuk Membantu",
      "subheading": "Punya pertanyaan tentang paket mengemudi EV kami atau penjadwalan? Hubungi tim kami melalui metode di bawah ini atau isi formulir."
    }
  },
  {
    "id": "methods-contact",
    "type": "contact_methods",
    "data": {
      "methods": [
        {
          "title": "Pusat Pelatihan",
          "description": "The Smith Office, 9th Floor, Unit 0902 Jl. Jalur Sutera Timur, RT 002/003, Kunciran, Kec. Pinang, Kota Tangerang, Provinsi Banten 15144",
          "icon": "i-lucide-map-pin",
          "actionText": "Dapatkan Petunjuk Arah",
          "actionLink": "https://maps.app.goo.gl/qGngC2sF4G3jt8Vs8",
          "target": "_blank"
        },
        {
          "title": "Dukungan WhatsApp",
          "description": "+62 811-9124-848 (Available 08:00 - 18:00)",
          "icon": "i-simple-icons-whatsapp",
          "actionText": "Chat Sekarang",
          "actionLink": "https://wa.me/628119124848?text=Halo%20Drive%20Master%2C%20saya%20ingin%20bertanya%20tentang%20kursus%20mengemudi",
          "target": "_blank"
        },
        {
          "title": "Alamat Email",
          "description": "info@evdriveacademy.id",
          "icon": "i-lucide-mail",
          "actionText": "Kirim Email",
          "actionLink": "mailto:info@evdriveacademy.id",
          "target": "_self"
        },
        {
          "title": "Jam Operasional",
          "description": "Senin - Jumat: 08:00 - 17:00 | Sabtu - Minggu: 08:00 - 17:00 | Sesi Malam: 18:00 - 20:00",
          "icon": "i-lucide-clock",
          "actionText": "Lihat FAQ",
          "actionLink": "/#faq",
          "target": "_self"
        }
      ]
    }
  },
  {
    "id": "form-map-contact",
    "type": "contact_form_map",
    "data": {
      "headline": "Hubungi Kami",
      "title": "Kirim Pesan",
      "description": "Isi formulir di bawah ini dan tim sukses pelanggan kami akan segera menghubungi Anda.",
      "mapEmbedUrl": "https://maps.google.com/maps?q=-6.22369663061115,106.66409468196608&z=17&output=embed"
    }
  },
  {
    "id": "social-contact",
    "type": "social_media",
    "data": {
      "headline": "Bergabunglah dengan Komunitas Kami",
      "title": "Media Sosial",
      "description": "Ikuti kami di media sosial untuk tips mengemudi, berita EV, dan cerita sukses murid.",
      "links": [
        { "label": "TikTok", "icon": "i-simple-icons-tiktok", "to": "https://tiktok.com" },
        { "label": "Facebook", "icon": "i-simple-icons-facebook", "to": "https://facebook.com" },
        { "label": "Instagram", "icon": "i-simple-icons-instagram", "to": "https://instagram.com" },
        { "label": "YouTube", "icon": "i-simple-icons-youtube", "to": "https://youtube.com" }
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
