# Technical Test

## Introduction

The assignment is divided into three tasks.

### Task 1: Create a DOCKERFILE

-   You are expected to put the application in a container so that it runs.

### Task 2: Create a Docker-compose

-   You are expected to write a docker-compose-yml
-   The compose file should bootstrap all dependecies of the service so that the service can reach them at "redis" and "kafka"
-   The entire infra should be running with a single \`doocker compose up\` statement.
-   The dependecies of the services are 
    -   **Kafka:** reachable at "kafka:9092"
    -   **Redis:** reachable at "redis:6379"
-   Mount the data directory of redis and kafka to a local directory(ies) in the docker-compose file.

### Task 3: Deploy service using CloudFormation OR Terraform

-   Write a CloudFormation/Terraform Temlplate so that we may deploy the application container using ECS
-   Expect redis and kafka to be running at \`redis\` and \`kafka\` endpoints.
-   Use your best judgement to set CPU and Memory requirements
-   Use your best judgement for LoadBalancer and Target Group configuration.
-   Assume that we have stored the Docker Image in a Container Repository and ECS can pull the Docker Image with no issues.


## Expectation

-   You are not expected to edit the code.
-   All the additional infrastructure that you need to deploy a service in ECS like VPC-id, Subnet-id, ClusterName etc would be already available and you can use placeholder values for them in your CloudFormation/Terraform template.


## Goals

-   The goal of this exercise is get a clear picture of how much you understand docker, docker-compose and aws-services.


## Notes

-   We believe that there is enough information in this document and the README.md file for you to successfully complete the tasks.
-   That said, if you have questions; please feel free to ask them.
-   If we have missed anything and you make assumptions for them, please mention the assumtions that you are making.


## How to Submit this assignment ?

-   Create a repository in gitlab/github and share the URL with us.
-   Write a README that shows how to deploy the application using the template.

