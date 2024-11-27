# Restaurants Finder

## Overview

Restaurants Finder is a Go-based application that provides APIs for managing and searching restaurants. It uses AWS DynamoDB for storage and is orchestrated with Kubernetes. The infrastructure is managed via Terraform, and CI/CD pipelines are implemented using GitHub Actions.

---

## Features

- Add, update, delete, and search for restaurants.
- Audit logs for API requests.
- Deployment to AWS EKS using Kubernetes.
- Secure API with secrets stored in Kubernetes.
- Readiness and liveness probes for health checks.

---

## Prerequisites

### Tools Required

1. **Install Terraform**:  
   [Install Terraform](https://developer.hashicorp.com/terraform/downloads)

   # Example for Linux
    ```
    wget https://releases.hashicorp.com/terraform/1.x.x/terraform_1.x.x_linux_amd64.zip
    unzip terraform_1.x.x_linux_amd64.zip
    sudo mv terraform /usr/local/bin/
    ```
2.	**Initialize and Apply Terraform**:
    ```
    cd infra
    terraform init
    terraform apply
    ```

## Kubernetes Setup

	•	Ensure you have kubectl and helm installed and configured on your local machine.
	•	Make sure your AWS credentials are properly configured to interact with EKS and ECR.
	•	Confirm the restaurant-finder namespace is created or allow Helm to create it.

# Configure the kubeconfig file

`aws eks update-kubeconfig --region <AWS_REGION> --name <EKS_CLUSTER_NAME>`

# Helm Deployment

1.	**Install the Application**:

2.	Update Secrets Locally:
For local testing, create the admin-secret Kubernetes secret:
```
kubectl create secret generic admin-secret \
  --namespace restaurant-finder \
  --from-literal=ADMIN_PASSWORD=your-admin-password
```

Use Helm to deploy the application:
```
helm upgrade --install restaurant-finder ./k8s/ \
  --namespace restaurant-finder \
  --create-namespace \
  --set serviceAccount.name=restaurant-finder-sa \
  --set serviceAccount.roleArn=arn:aws:iam::${AWS_ACCOUNT_ID}:role/restaurants-finder-eks-cluster-irsa \
  --set image.repository=${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/restaurant-finder \
  --set image.tag=latest \
  --set adminPassword=${ADMIN_PASSWORD}
```

## Interacting with the API

Example curl Commands

	1.	Health Check:
    `curl http://<load-balancer-endpoint>/healthz`

    2.	Search Restaurants:
    `curl "http://<load-balancer-endpoint>/restaurants/search?cuisine=Italian&is_kosher=true&is_open=true"`
    3.	Admin Actions:
        Replace <admin-password> with your ADMIN_PASSWORD.
	•	Add a Restaurant:
    ```
    curl -X POST -H "Authorization: <admin-password>" -H "Content-Type: application/json" \
    -d '{"restaurant_name":"New Place","address":"123 Main St","cuisine_type":"Italian","is_kosher":true}' \
    http://<load-balancer-endpoint>/admin/restaurants
    ``` 
    •	Fetch Audit Logs:
    ```
    curl -X GET -H "Authorization: <admin-password>" \
    http://<load-balancer-endpoint>/admin/logs?minutes=60
    ```

## CI/CD Pipelines

    Build and Push to ECR

    The pipeline builds the Docker image and pushes it to ECR on any changes to the server/ folder.

    Deploy to Kubernetes

    The pipeline deploys to the Kubernetes cluster on changes to the k8s/ folder.

    GitHub Actions

    Required Secrets

    Set the following secrets in your repository:
        •	AWS_ACCESS_KEY_ID: Your AWS access key.
        •	AWS_SECRET_ACCESS_KEY: Your AWS secret key.
        •	AWS_REGION: The AWS region for your resources (e.g., us-east-1).
        •	AWS_ACCOUNT_ID: The AWS account id for your resources (e.g., 123456789).
        •	ADMIN_PASSWORD: The admin password for the application.

    Required Vars

    Set the following vars in your repository:
        •	ECR_REPOSITORY: restaurant-finder


## Prometheus installation

1.	Install the Helm Chart:
```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
```

2.	Deploy the Stack:
```
helm upgrade --install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --set grafana.admin.user=admin \
  -f resources/prometheus.yaml \
  --set grafana.admin.password=admin
```

3.	Verify the Installation:
```
kubectl get pods -n monitoring
kubectl get svc -n monitoring
```

