# Venue Search POS

## Overview
This project is a Go web server application that utilizes the Gin framework to provide endpoints for searching nearby stadiums and checking if a point is inside a stadium. It connects to a MongoDB database to retrieve venue information based on geographical coordinates.

## Files
- **main.go**: Contains the main application code for the web server.
- **Dockerfile**: Defines the Docker image for the application.
- **README.md**: Documentation for the project.

## Getting Started

### Prerequisites
- Go (version 1.16 or later)
- Docker (for containerization)
- MongoDB Atlas account (for database access)

### Setup

1. **Clone the repository**:
   ```
   git clone <repository-url>
   cd venue-searchpos
   ```

2. **Install dependencies**:
   Ensure you have the required Go packages by running:
   ```
   go mod tidy
   ```

3. **Configure MongoDB**:
   Update the MongoDB connection string in `main.go` with your credentials.

### Running the Application

#### Locally
To run the application locally, execute:
```
go run main.go
```
The server will start on `http://localhost:8080`.

#### Using Docker
To build and run the application using Docker, follow these steps:

1. **Build the Docker image**:
   ```
   docker build -t venue-searchpos .
   ```

2. **Run the Docker container**:
   ```
   docker run -p 8080:8080 venue-searchpos
   ```

The application will be accessible at `http://localhost:8080`.

### API Endpoints

- **GET /nearby**: Retrieve nearby stadiums based on longitude, latitude, and maximum distance.
  - Query parameters:
    - `long`: Longitude of the location.
    - `lat`: Latitude of the location.
    - `maxDistance`: Maximum distance in meters.

- **GET /inside**: Check if a point is inside a stadium.
  - Query parameters:
    - `long`: Longitude of the point.
    - `lat`: Latitude of the point.

### License
This project is licensed under the MIT License.

# Perimissions for SA 

> Gitaction
> The service account gitaction@lfs261-cicd-304112.iam.gserviceaccount.com needs additional permissions. You need to grant these roles to the service account in the Google Cloud Console:
> Required IAM Roles:
> Cloud Run Admin - To deploy services
> Cloud Build Editor - To build from source
> Storage Admin - For build artifacts
> Secret Manager Secret Accessor - To access MongoDB URI secret
> Service Account User - To use service accounts

# Get the Cloud Build service account
PROJECT_NUMBER=$(gcloud projects describe lfs261-cicd-304112 --format="value(projectNumber)")

# Grant Cloud Run Admin to Cloud Build service account
gcloud projects add-iam-policy-binding lfs261-cicd-304112 \
  --member="serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com" \
  --role="roles/run.admin"

# Get the Cloud Build service account
PROJECT_NUMBER=$(gcloud projects describe lfs261-cicd-304112 --format="value(projectNumber)")

# Grant Cloud Run Admin to Cloud Build service account
gcloud projects add-iam-policy-binding lfs261-cicd-304112 \
  --member="serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com" \
  --role="roles/run.admin"

```bash
roles/run.admin
roles/cloudbuild.builds.editor
roles/storage.admin
roles/secretmanager.secretAccessor
roles/iam.serviceAccountUser
```