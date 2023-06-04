# TikTok Tech Immersion Assignment 2023
[Assignment instruction](https://bytedance.sg.feishu.cn/docx/P9kQdDkh5oqG37xVm5slN1Mrgle)
- Design and implement a backend Instant Messaging system

# Tools used
- [Golang](https://go.dev/) programming language
- Kitex
- [Redis](https://redis.io/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) and Kubernetes
- Github Actions: to automate the testing
- [Postman](https://www.postman.com/downloads/): To test the api

# Setup
- You can either run the project using docker compose or kubernetes
- You may install [Docker Desktop](https://www.docker.com/products/docker-desktop/) and enable Kubernetes in the Docker Desktop
- Clone this repo

## Using docker compose to run the project
run the command `docker compose up --build` in the terminal

## Using Kubernetes
### 1. Build Docker Image
- Build the docker image by running the command `docker build -t {image tag} {path to dockerfile}` in the terminal
- For example, run the commands `docker build -t ernst1/http-server ./http-server` and `docker build -t ernst1/rpc-server ./rpc-server` in the terminal
- Note: If you change the image tag, remember to update the files in `k8` folder so that it gets the correct docker image

### 2. Push to Docker Hub
- Login to Docker Hub, if you haven't logged in yet
- Run the command to push your docker image to docker hub: `docker push {image tag}`
  - Eg. run the commands `docker push ernst1/rpc-server` and `docker push ernst1/http-server` in the terminal
- Note: When pushing to Docker Hub, the image tag format should be `{your docker username}/{image name}`

### 3. Run using Kubernetes
run the command `kubectl apply -f ./k8` in the terminal

### Change the number of pods
- You may change the number of pods running by changing the file in `k8` folder.
- For example, you may change `replicas: 2` to `replicas: 3` (to run 3 pods for that service) in `k8/http-server-depl.yaml` file

# JMeter
To use `HTTP Request.jmx`, install JMeter. For MacOS, `brew install jmeter`
