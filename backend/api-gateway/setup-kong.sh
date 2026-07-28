#!/bin/bash

# Kong API Gateway Setup Script
# This script helps set up and manage Kong API Gateway

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
KONG_CONFIG_FILE="./kong/kong.yml"
KONG_ADMIN_URL="http://localhost:8001"
KONG_PROXY_URL="http://localhost:8000"
KONG_MANAGER_URL="http://localhost:8002"

echo -e "${BLUE}Kong API Gateway Setup Script${NC}"
echo "=================================="

# Function to check if Kong is running
check_kong_health() {
    echo -e "${YELLOW}Checking Kong health...${NC}"
    if curl -s -f "${KONG_ADMIN_URL}/health" > /dev/null; then
        echo -e "${GREEN}Kong is healthy and running${NC}"
        return 0
    else
        echo -e "${RED}Kong is not responding${NC}"
        return 1
    fi
}

# Function to validate Kong configuration
validate_config() {
    echo -e "${YELLOW}Validating Kong configuration...${NC}"
    if [ -f "$KONG_CONFIG_FILE" ]; then
        echo -e "${GREEN}Kong configuration file found: ${KONG_CONFIG_FILE}${NC}"
        # You can add more validation here
    else
        echo -e "${RED}Kong configuration file not found: ${KONG_CONFIG_FILE}${NC}"
        exit 1
    fi
}

# Function to show Kong status
show_status() {
    echo -e "${BLUE}Kong Gateway Status${NC}"
    echo "======================"

    echo -e "${YELLOW}Admin API:${NC} ${KONG_ADMIN_URL}"
    echo -e "${YELLOW}Proxy:${NC} ${KONG_PROXY_URL}"
    echo -e "${YELLOW}Manager UI:${NC} ${KONG_MANAGER_URL}"

    if check_kong_health; then
        echo -e "\n${GREEN}=== Kong Routes ===${NC}"
        curl -s "${KONG_ADMIN_URL}/routes" | jq '.data[] | {name: .name, paths: .paths, methods: .methods}' 2>/dev/null || echo "Unable to fetch routes"

        echo -e "\n${GREEN}=== Kong Services ===${NC}"
        curl -s "${KONG_ADMIN_URL}/services" | jq '.data[] | {name: .name, url: .url}' 2>/dev/null || echo "Unable to fetch services"
    fi
}

# Function to test connectivity to microservices
test_services() {
    echo -e "\n${BLUE}Testing Service Connectivity${NC}"
    echo "============================="

    services=("user-service" "core-service" "booking-service")
    for service in "${services[@]}"; do
        echo -e "${YELLOW}Testing $service...${NC}"
        if docker-compose exec -T kong-gateway curl -s "http://${service}:800${service#user-service}" > /dev/null; then
            echo -e "${GREEN}✓ $service is reachable${NC}"
        else
            echo -e "${RED}✗ $service is not reachable${NC}"
        fi
    done
}

# Function to start/stop Kong
manage_kong() {
    local action=$1
    echo -e "${YELLOW}Running: docker-compose ${action} kong-gateway${NC}"
    docker-compose ${action} kong-gateway
}

# Function to show Kong logs
show_logs() {
    echo -e "${BLUE}Kong Gateway Logs${NC}"
    echo "==================="
    docker-compose logs -f kong-gateway
}

# Function to open Kong Manager UI
open_manager() {
    echo -e "${YELLOW}Opening Kong Manager UI in browser...${NC}"
    open "${KONG_MANAGER_URL}" || echo "Please open ${KONG_MANAGER_URL} manually"
}

# Main menu
case "${1:-status}" in
    "health")
        check_kong_health
        ;;
    "validate")
        validate_config
        ;;
    "status")
        show_status
        ;;
    "test")
        test_services
        ;;
    "start")
        manage_kong "up -d"
        ;;
    "stop")
        manage_kong "stop"
        ;;
    "restart")
        manage_kong "restart"
        ;;
    "logs")
        show_logs
        ;;
    "manager")
        open_manager
        ;;
    "reload")
        echo -e "${YELLOW}Reloading Kong configuration...${NC}"
        docker-compose exec -T kong-gateway kong reload
        ;;
    *)
        echo "Usage: $0 {health|validate|status|test|start|stop|restart|logs|manager|reload}"
        echo ""
        echo "Commands:"
        echo "  health     - Check if Kong is running and healthy"
        echo "  validate   - Validate Kong configuration file"
        echo "  status     - Show Kong status and configuration"
        echo "  test       - Test connectivity to microservices"
        echo "  start      - Start Kong gateway"
        echo "  stop       - Stop Kong gateway"
        echo "  restart    - Restart Kong gateway"
        echo "  logs       - Show Kong logs"
        echo "  manager    - Open Kong Manager UI in browser"
        echo "  reload     - Reload Kong configuration"
        exit 1
        ;;
esac