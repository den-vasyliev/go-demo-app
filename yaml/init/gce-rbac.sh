## Note that your GCP identity is case sensitive but gcloud info as of Google Cloud SDK 221.0.0 is not. 
## This means that if your IAM member contains capital letters, the above one-liner may not work for you. 
## If you have 403 forbidden responses after running the above command and kubectl apply -f kubernetes, 
## check the IAM member associated with your account at 
## https://console.cloud.google.com/iam-admin/iam?project=PROJECT_ID. 
## If it contains capital letters, you may need to set the --user flag in the command above to the 
## case-sensitive role listed at https://console.cloud.google.com/iam-admin/iam?project=PROJECT_ID.
##After running the above, if you see Clusterrolebinding "cluster-admin-binding" created, then you 
## are able to continue with the setup of this service.

kubectl create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin --user=$(gcloud config get-value core/account)
