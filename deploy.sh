# this is the service account used by Cloud Functions and Cloud Run
# gcloud projects add-iam-policy-binding lfs261-cicd-304112 \
# --member=serviceAccount:616245089133-compute@developer.gserviceaccount.com \
# --role=roles/cloudbuild.builds.builder

# Bind the service account to the Secret Manager role
# gcloud secrets add-iam-policy-binding mongodb-supporters-uri \
#   --member="serviceAccount:616245089133-compute@developer.gserviceaccount.com" \
#   --role="roles/secretmanager.secretAccessor"

IMAGE_NAME="gcr.io/lfs261-cicd-304112/venue-searchpos-service:v3"

docker build --platform=linux/amd64 -t $IMAGE_NAME .
# if [[ "$?" -ne 0 ]] ; then
#     echo 'could not perform build'; exit $rc
# fi
docker push $IMAGE_NAME
# if [[ "$?" -ne 0 ]] ; then
#     echo 'could not perform image push'; exit $rc
# fi

gcloud run deploy venues-searchpos-service \
  --image=$IMAGE_NAME \
  --platform=managed \
  --region europe-southwest1 \
  --timeout=120s \
  --allow-unauthenticated \
  --memory=512Mi \
  --cpu=1 \
  --concurrency=80  
  --set-secrets MONGO_URI=mongodb-supporters-uri:latest