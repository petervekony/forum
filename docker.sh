# This is the bash script that helps create the docker image, and then runs it and prunes all dangling images
# The container is run using the --rm flag, meaning it the container is deleted after it stops.

# Building docker image
docker image build -f Dockerfile -t forum .

# Pruning dangling images and containers
docker image prune -f
docker container prune -f

# Running the container on port 443 named gritface
docker container run -p 443:443 --detach --rm --name gritface forum