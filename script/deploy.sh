cd deployment

tag=$1

if [[ -z $tag ]]
then
  tag=$(git rev-parse --short HEAD)
fi

# install is for first time deployment. If it's failed, try to upgrade
helm --set server.image.tag=$tag install high-available-server . || helm --set server.image.tag=$tag upgrade high-available-server .
