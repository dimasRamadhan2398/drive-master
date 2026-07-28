# Deployment Progress Report (VPS & Domain)

This report outlines the current status of the deployment for the **DriveMaster** application to the VPS (`203.194.114.20`) and the domain `drivemaster.id`.

---

## 📊 Summary of Current Status

* **DNS Configuration**:  **Completed**. All domains point to the correct VPS IP.
* **VPS Environment**:  **Active**. The VPS runs AlmaLinux 8.10 with cPanel/Apache on ports 80/443. Docker, Docker Compose, Node.js (`v20.20.2`), and PM2 are installed.
* **Code/Repository State**:  **Outdated/Incomplete**. The target deployment directory (`/var/www/drive`) on the VPS contains older code from June 26. Crucially, the `docker-compose.yml` on the server is the old version (3KB) and does not match the current local version (9.7KB) containing all microservices.
* **Application Execution**:  **Not Running**. No Docker containers or PM2 processes are currently running on the VPS. Visiting `drivemaster.id` returns Apache's default empty directory listing ("Index of /").

---

## 🔍 Detailed Checks & Findings

### 1. DNS Mapping
All DNS checks resolve correctly to your VPS IP `203.194.114.20`:
* `drivemaster.id` ➡️ `203.194.114.20`
* `api.drivemaster.id` ➡️ `203.194.114.20`
* `vps.drivemaster.id` ➡️ `203.194.114.20`

### 2. VPS System & Web Server Environment
* **Web Server**: Apache (`httpd`) is running on ports **80** (HTTP) and **443** (HTTPS) managed by cPanel.
* **Document Root**: The domain `drivemaster.id` points to `/home/drivemaster/public_html`. Currently, this folder only contains cPanel's default configuration files and is serving an empty index listing.
* **Subdomain (`api.drivemaster.id`)**: This domain is pointing to the VPS via DNS, but it is **not yet configured** in the Apache virtual hosts or cPanel userdata.

### 3. Docker & Containers
* The Docker daemon is running successfully on the VPS.
* Infrastructure images are pulled (`postgres`, `redis`, `redisinsight`, `kong-gateway`, `cp-kafka`, `cp-zookeeper`, `kafka-ui`).
* **No containers are currently running or built** for the backend microservices (`user-service`, `core-service`, `booking-service`, `payment-service`).

### 4. Code & Configuration Discrepancies
* The code on the server under `/var/www/drive` is missing the latest changes.
* **`docker-compose.yml` mismatch**:
  * Local version: **9.7 KB** (contains all microservices: booking, payment, kong-gateway, user, core, kafka, redis, etc.)
  * Remote version: **2.9 KB** (only contains postgres, redis, zookeeper, kafka, user-service, core-service, api-gateway)

---

## 🛠️ Action Plan to Complete Deployment

To finish deploying and make `drivemaster.id` work, the following steps must be executed:

### Step 1: Update Code on VPS
* Sync the latest local code (including the updated `docker-compose.yml` and backend services) to `/var/www/drive`.

### Step 2: Build & Start Microservices (Docker Compose)
* Run `docker compose up -d --build` on the remote server under `/var/www/drive` to compile the backend microservices and start all docker containers.

### Step 3: Run Nuxt Frontend via PM2
* Build the frontend locally or on the server.
* Start the frontend via PM2 on the server to listen on port `3000` (e.g. `pm2 start .output/server/index.mjs --name "drive-frontend"`).

### Step 4: Configure Apache Reverse Proxy
Since Apache is running on port 80/443, we must set up reverse proxies in Apache to route traffic from public ports to our internal services:
1. **Frontend Proxy (`drivemaster.id` to Port 3000)**:
   Add a reverse proxy rule in `/home/drivemaster/public_html/.htaccess` to route traffic to the Nuxt frontend running on port 3000:
   ```apache
   RewriteEngine On
   RewriteRule ^$ http://127.0.0.1:3000/ [P,L]
   RewriteCond %{REQUEST_FILENAME} !-f
   RewriteCond %{REQUEST_FILENAME} !-d
   RewriteRule ^(.*)$ http://127.0.0.1:3000/$1 [P,L]
   ```
2. **API Proxy (`api.drivemaster.id` to Port 8080)**:
   * Option A: Register `api.drivemaster.id` as a subdomain in cPanel, then use a similar `.htaccess` file in its document root to proxy to port `8080` (Kong Gateway).
   * Option B: Route all API requests through `drivemaster.id/api` (rather than a subdomain). We can update `.htaccess` to forward `/api` requests to `http://127.0.0.1:8080/api`. This avoids having to configure subdomains and manage multiple SSL certificates.
