# gcloud projects add-iam-policy-binding lfs261-cicd-304112 \
# --member=serviceAccount:616245089133-compute@developer.gserviceaccount.com \
# --role=roles/cloudbuild.builds.builder

gcloud functions deploy venues-searchpos-service \
  --gen2 \
  --runtime=go122 \
  --entry-point=ignored \
  --source=. \
  --trigger-http \
  --region europe-southwest1 \
  --allow-unauthenticated \
  --set-secrets MONGO_URI=mongodb-supporters-uri:latest