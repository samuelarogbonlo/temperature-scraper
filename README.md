# DIA Temperature Scraper Deployment Infrastructure task

This repository houses the codebase for the the temperature scraper application infrastructure hosted on Heztner Cloud with proper storage, availability and security. The temperature scraper collects temperature data (in °C) for Zurich, London, Miami, Tokyo, and Singapore.

![DIA Infra drawio](https://github.com/samuelarogbonlo/temp-scraper/assets/47984109/753581bb-e8d1-4ba0-ab71-105ae7e64a8b)

# Infrastructure Overview
The major components include:

- *Heztner Server:* This is the core of our infrastructure handling the minikube setup for all our services and applications.
- *Minikube:* Mini-Kubernetes set up in the cloud to automatically manage applications, scripts and services.
- *Kubernetes Deployment:* In the Kubernetes folder in the root folder contains all the deployment files for setting up the deployment of the application. These applications include the temperature scraper app, Kafka cluster, consumer app, IPFS cluster, Postgres Database and logging & monitoring suite (Prometheus, Grafana and Loki).

## Continous Integration

Our CI is facilitated by GitHub Actions, which automates and integrates some processes in our infrastructure. For both the consumer app and the temperature scraper, the pipeline is triggered by pushes to the main branch, it encompasses:

- *Dockerfile Linting:* Ensures adherence to best practices and standards in our Dockerfile configurations.
- *Docker Image Building:* Constructs a Docker image from the application and pushes it to DockerHub.

# How To Setup And Deploy The Infrastructure?

## Requirements
- Heztner Server: 8GB RAM, 40GB ROM
- Helm
- Golang (probably the latest version)

## Steps
- SSH into the server
- Run `minikube addons enable metrics-server` to install the metrics server for autoscaling.
- Setup and install ArgoCD with this command `kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml.` and `kubectl create namespace argocd` For more information on installations, check [here](https://argo-cd.readthedocs.io/en/stable/getting_started/)
- Run this command on the server `kubectl port-forward --address 0.0.0.0 svc/argocd-server -n argocd 8081:443`.
- Go to the browser and access the ArgoCD UI with `http://<server-ip>:8081/`
- For the login details for the UI, use the default username `admin` and to get the password, run `kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d` on the server.
- After logging in, navigate to **Settings** then to **Repository** then click **Connect Repo** and fill all required details.
- Then create a new app by Navigating to **Applications** then click **New App** and fill in all details required for each app in the `Kubernetes` directory in the `temp-scraper` repository.
- Once all applications are deployed, synced and healthy, go back to the server and run `git clone` of the host repository.
- Depending on where you are in the server, run `cd IPFS-Upload` to get to the right directory then run `go build`.
- After that is completeed, run `./ipfs-upload -city <city e.g London> -date <date in YYYY-MM-DD format>` to get the script in operation. This should generate the URL that will be used to access the data.
- You can download the data with `wget http://<ipfs-gateway-url>:8080/ipfs/<generated-cid>`.

## Monitoring And Logging
Once you have installed and deployed the applications via ArgoCD, you can access grafana to view metrics and the in-built dashboard. Follow the steps here

- Access the server and get the Grafana login password with `kubectl get secret --namespace monitoring prometheus-stack-chart-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo` and username is `admin`.
- On Grafana, you may have to configure the data sources manually like Loki and Prometheus. However, for the current setup, it is all active.
- To view the current setup, you have first to access the server and run `kubectl port-forward --address 0.0.0.0 svc/prometheus-stack-chart-grafana -n monitoring 3001:80` then go to the browser and run this URL `http://<server-ip>:3001/`.

## Deployments And Rollbacks
Currently, the project is being managed with the GitOps approach using ArgoCD and this gives room for seamless deployments on changes to the main branch. Of course, in the future, we can set up different branches like `dev`, `staging` and `main` and still be able to manage rollbacks and updates from ArgoCD UI.

> **_Note:_**
- For ease of access to the services via the browser on my local machine, a simple NodePort configuration could’ve sufficed. Unfortunately, due to Docker’s networking and Minikube's reliance on Docker, the URL format `http://<server-ip>:<NodePort>` to access apps deployed on Kubernetes will not suffice outside of the Server, Minikube is running in. As a hotfix, port-forwarding had to be implemented to access the UIs for these apps.

- An additional information, we could have been able to set up HPA for Kafka Brokers but to do so effectively, there will have to be careful considerations of autoscaling implications for Kafka. These scaling decisions are based on more than just CPU and memory metrics, they involve operational tasks such as rebalancing partitions that HPA alone cannot handle as a result of having limited resources on the single-node Kubernetes cluster. HPA configurations such as offsets.topic.replication.factor etc. have been set to the barest minimum of 1 each due to limited resources available on the single-node Kubernetes Cluster.

- Also, there was a need to expand the server a little bit for us to accommodate all the services and deployments so we had to increase the server size.

- Lastly, we must also add that Minikube made things a little manual and slow because of the need to manage lots of things because of the small nature of the cluster. Of course, minikube is used for smaller setups and less production-based implementations.

## Author
- Samuel Arogbonlo - [GitHub](https://github.com/samuelarogbonlo)

## Collaborators
- [YOUR NAME HERE] - Feel free to contribute to the codebase by resolving any open issues, refactoring, adding new features, writing test cases or any other way to make the project better and helpful to the community. Please feel free to send pull requests.

## Hire me
Are you looking for a Senior DevOps Engineer to build your next infrastructure? Get in touch: [sbayo971@gmail.com](mailto:sbayo971@gmail.com)

## License

The MIT License (http://www.opensource.org/licenses/mit-license.php)
