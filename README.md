# pbl-project-simpadu
Project Micro Service yang berkolaborasi dengan tim lain, service ini berfokus kepada:
- Design system login
- Super Admin
- Master Data (Admin Akademik)

Repository proyek PBL (Project Based Learning) - Sistem Informasi Manajemen Padu (SIMPADU).
Berikut adalah panduan untuk menjalankan proyek ini di lingkungan lokal menggunakan Docker.

## Requirement

- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [node.js](https://nodejs.org/)
- [npm](https://www.npmjs.com/)
- [flutter](https://flutter.dev/)

## Tutorial Penggunaan Repo

### 1. Clone Repository

    git clone https://github.com/Kar-Su/uas-mobile.git
    cd uas-mobile

### 2. Siapkan File Environment

Salin file `.env.example` yang ada di root proyek dan di dalam folder `backend/` menjadi `.env`.

**Linux / macOS:**

    cp .env.example .env
    cp ./backend/.env.example ./backend/.env

**Windows (Command Prompt):**

    copy .env.example .env
    copy .\backend\.env.example .\backend\.env

**Windows (PowerShell):**

    Copy-Item .env.example .env
    Copy-Item .\backend\.env.example .\backend\.env

> **Catatan:** Jika file `.env` sudah ada, Anda dapat menimpanya atau menyesuaikan isinya sesuai kebutuhan.

### 3. Jalankan Docker Containers

    docker compose up -d --build

Perintah ini akan membangun (jika diperlukan) dan menjalankan semua layanan dalam mode background. Tunggu hingga semua container siap.

### 4. Seed Database

Setelah container berjalan, isi basis data dengan data awal (seeder):

    docker exec -i golang_pbl go run cmd/seeder/main.go --seed

### 5. Jalankan Aplikasi

    cd frontend/
    npm install
    npm run dev

Terminal Baru

    cd mobile/
    flutter run
