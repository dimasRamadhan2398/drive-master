set -e

export CONFIG_PATH=/app/config.yaml

echo "Running database migrations..."
./booking-service migrate

echo "Running database seeders..."
./booking-service seed

echo "Starting User Service..."
./booking-service serve