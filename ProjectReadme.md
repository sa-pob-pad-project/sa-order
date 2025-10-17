# How to Run the Project

## System Architecture
Microservices

This project consists of three services:
- **sa-user**: Backend (User Service)
- **sa-appointment**: Backend (Appointment Service)
- **sa-frontend**: Frontend Application

## Installation
- Install Docker Compose
- Install Go 1.24.4
- Install Node.js 15.5.2

## Running the Services

1. Start the database using Docker Compose (from any service root):
   ```cmd
   docker-compose up -d
   ```

2. Run the backend services
   ```cmd
   cd sa-user
   go run main.go
   ```
   Open a new terminal and run:
   ```cmd
   cd sa-appointment
   go run main.go
   ```

3. Run the frontend application
   ```cmd
   cd sa-frontend
   npm install
   npm start
   ```

## Access URLs
- Frontend: http://localhost:3000
- User Service API: http://localhost:8000
- Appointment Service API: http://localhost:8001

## Contribution
  1. นพณัช สาทิพย์พงษ์ besterOz
  2. พงศธร รักงาน prukngan
  3. ภวัต เลิศตระกูลชัย Phawat Loedtrakunchai
  4. ธฤต จันทร์ดี tharitpr



