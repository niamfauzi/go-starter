#!/usr/bin/env bash

# Login
curl -X POST http://localhost:8080/api/v1/auth/login   -H "Content-Type: application/json"   -H "X-Tenant-Id: tenant-demo"   -d '{
    "email": "demo@example.com",
    "password": "password123"
  }'
