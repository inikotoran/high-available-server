# this script is used to build docker image
# it doesn't push to any registry since we're going to use it locally
# if we want to deploy the application on cloud environment, it should be pushed to a remote registry

tag=$1

if [ "$(kubectl config current-context)" == minikube ]
then
  echo "Setup docker env for minikube"
  eval $(minikube -p minikube docker-env)
fi

if [[ -z $tag ]]
then
   tag=$(git rev-parse --short HEAD)
fi

docker build -t high-available-server:$tag .
