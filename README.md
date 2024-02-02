# DIA Temperature Scraper Deployment Infrastructure task

This repository houses the codebase for the the temperature scraper application infrastructure hosted on Heztner Cloud with proper storage, availability and security. The temperature scraper collects temperature data (in °C) for Zurich, London, Miami, Tokyo, and Singapore.

![ecs](images/ecs.jpg)

# Infrastructure Overview
The major components include:

- *Heztner Server:* This is the core of our infrastructure handling the minikube setup for all our services and applications.
- *Load Balancer:* Ensures even distribution of incoming requests across tasks, optimizing response times and resource utilization.
- *PostgreSQL Database:* Provides robust data storage for the Go application, ensuring data integrity and fast access.
- *Minikube:* Mini-Kubernetes setup in the cloud to automatically manage applications, scripts and services.
- *Prometheus & Grafana:* Integrated for comprehensive logging, enabling real-time monitoring and alerting for system health and performance.
- *Loki & Promtail:* Integrated for log aggregation for all services.
- *Ansible:* This is used to provision ArgoCD and manage all levels of configurations throughout the infrastructure.

## Continous Integration / Continous Deployment

Our CI/CD pipeline, facilitated by GitHub Actions, automates and streamlines our deployment process. For both the consumer app and the temeprature scraper, the pipeline is triggered by pushes to the main branch, it encompasses:

- *Dockerfile Linting:* Ensures adherence to best practices and standards in our Dockerfile configurations.
- *Docker Image Building:* Constructs a Docker image from the application and pushes it to DockerHub.

# How To Setup And Deploy The Infrastructure ?

## Requirements


- Infrastructure Initialization, Planning & Applying:

   - Run `cd ..` to return to the root directory.
   - Run `make init_all` to initialize Terraform scripts for the infrastructure.
   - Run `make plan_all` to review the infrastructure details before deployment.
   - Run `make apply_all` to apply the Terraform code and set up the infrastructure.

> **_Note:_** For some more information, the Node.js application code is built to expose prometheus metrics. You can view it with just an addition of /metrics to the load balancer URL.

## Author
- Samuel Arogbonlo - [GitHub](https://github.com/samuelarogbonlo)

## Collaborators
- [YOUR NAME HERE] - Feel free to contribute to the codebase by resolving any open issues, refactoring, adding new features, writing test cases or any other way to make the project better and helpful to the community. Please feel free to send pull requests.

## Hire me
Looking for a Senior DevOps Engineer to build your next infrastructure? Get in touch: [sbayo971@gmail.com](mailto:sbayo971@gmail.com)

## License

The MIT License (http://www.opensource.org/licenses/mit-license.php)