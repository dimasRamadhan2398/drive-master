# Backend Deployment Guide (VPS)

This project uses a microservices architecture containerized with Docker. The easiest way to deploy it to a VPS is using **Docker Compose**.

## Prerequisites

1.  **VPS** (Ubuntu 22.04+ recommended) with at least 4GB RAM (due to Kafka and multiple microservices).
2.  **Docker** and **Docker Compose** installed on the VPS.
3.  **Domain Name** (optional, but recommended for SSL).

## Step-by-Step Deployment

### 1. Prepare the VPS

Connect to your VPS:
```bash
ssh user@your-vps-ip
```

Install Docker (if not already installed):
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```

### 2. Clone the Repository

```bash
git clone <your-repo-url>
cd <repo-folder>/backend
```

### 3. Configure Environment Variables

Create a `.env` file in the `backend` directory:
```bash
cp .env.example .env
nano .env
```
Fill in the necessary credentials (Postgres, Redis, Kafka, Mailtrap, Midtrans, etc.).

### 4. Configure Kong API Gateway

The project uses Kong as the API Gateway. Ensure `backend/api-gateway/kong/kong.yml` is correctly configured with your VPS IP or domain name if necessary (though internal service names like `user-service` are used within the Docker network).

### 5. Deploy with Docker Compose

Run the following command to build and start all services in detached mode:
```bash
docker compose up -d --build
```

### 6. Verify the Deployment

Check the status of your containers:
```bash
docker compose ps
```

Monitor logs:
```bash
docker compose logs -f
```

Access the services:
- **API Gateway:** `http://your-vps-ip:8080`
- **Kafka UI:** `http://your-vps-ip:8085`
- **Redis Insight:** `http://your-vps-ip:5540`

### 7. Database Migrations and Seeders

The services are configured to run migrations on startup via environment variables:
- `RUN_MIGRATIONS=true`
- `RUN_SEEDERS=true`

If you need to run them manually, you can enter the container:
```bash
docker exec -it drive_master_user_service ./user-service serve --migrate --seed
```

## Security Recommendations

1.  **Firewall:** Use `ufw` to close all ports except `22` (SSH), `80` (HTTP), and `443` (HTTPS).
    ```bash
    sudo ufw allow 22
    sudo ufw allow 80
    sudo ufw allow 443
    sudo ufw enable
    ```
2.  **Reverse Proxy:** Use Nginx or Caddy on the host to handle SSL (Certbot/Let's Encrypt) and proxy traffic to Kong (`8080`).
3.  **Secrets:** Never commit your `.env` file to version control.

## Troubleshooting

- **Memory Issues:** If Kafka fails to start, ensure your VPS has enough memory. You can increase swap space if needed.
- **Connectivity:** Ensure all services are on the `drive_master_network`.
