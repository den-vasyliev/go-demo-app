export ARG1=${1:-10}
/usr/local/bin/kubectl get svc front-v2 -o yaml --export -n demo|sed "s/weight:.[0-9][0-9]$/weight: $ARG1/"|/usr/local/bin/kubectl -n demo apply -f -
