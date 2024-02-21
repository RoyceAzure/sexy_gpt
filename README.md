# Application Setup Guide

This guide outlines the steps necessary to start the application using Docker and provides default user information for initial access.

## Start Instructions

1. **Initialization:**
    - Navigate to the root directory of the application.
    - Start the application by running the following command:
      ```sh
      docker-compose up -d
      ```

## Usage Instructions

- **Access the Application:**
    - The application is hosted locally and can be accessed using the following URL:
      ```
      http://localhost:8081
      ```
    - Utilize the HTTP protocol to connect to the host.

## Default Internal User Information

- **Credentials for internal use:**
    - **Email:** `roycewnag123@gmail.com`
    - **Password:** `Royce123456`



## PawAI Project
1. **Initialization:**
    - Navigate to the root directory of the paw_ai_service.
    - create virtual env:
      ```python
      python -m venv {env_name}
      ```
    - activate env:
      - windows:
        ```cmd
        {env_name}\Scripts\activate
        ```
      - MacOS/Linux:
        ```bash
        source {env_name}/bin/activate
        ```
    - Install dependencies:
      ```python
      pip install -r requirements.txt
      ```
    - Run the application:
      ```
      python app.py
      ```
---
