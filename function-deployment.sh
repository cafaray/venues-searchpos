gcloud functions deploy venue-searchpos-service \
  --gen2 \
  --runtime go122 \
  --entry-point EntryPoint \
  --source . \
  --trigger-http \
  --region europe-southwest1 \
  --allow-unauthenticated \
  --set-secrets MONGO_URI=mongodb-supporters-uri:latest