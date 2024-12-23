# MobSF Static Scanning of App

## Overview

This Go project provides a core setup flow for performing static scanning of Android and iOS app packages using MobSF (Mobile Security Framework). It is designed to simplify the process of scanning mobile applications for security vulnerabilities.

## Prerequisites

- Go installed on your machine
- Kubernetes cluster set up
- `kubectl` configured to interact with your cluster
- MobSF installed and configured

## Getting Started

Follow these steps to set up and run the MobSF static scanning:

### Step 1: Navigate to the Project Directory

<!-- ```bash -->
cd mobsf

### Step 2: Update the API Key

Replace your 64-character alpha numeric API key with both lower and upper case in the appropriate YAML configuration file:

- For development: `build/mobsf-dev.yaml`
- For production: `build/mobsf-prod.yaml`

### Step 3: Update Namespace (if required)

If you need to change the namespace, do so in the YAML configuration files as needed.

### Step 4: Deploy the Configuration

Run the following command based on the environment type (dev or prod):

<!-- ```bash -->
kubectl apply -f build/mobsf-dev.yaml

Note: For this project, we will be using the dev setup in the mobsf namespace.

### Step 5: Get the NodePort

After deploying, pick the NodePort from the deployed service to access the MobSF server.

### Step 6: Update MobSF Server URL and App Path

Replace the MobSF server URL and the app path in the code with the correct details that correspond to your setup.

### Step 7: Clean Up Dependencies

Run the following command to tidy up your Go module dependencies:

<!-- ```bash -->
go mod tidy

### Step 8: Run the Application

Finally, execute the main Go program:

<!-- ```bash -->
go run main.go

Happy scanning! 🚀
