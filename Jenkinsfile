//Define variables
def scmVars

//Start Pipeline
pipeline {
    
  //Configure Jenkins Slave
  agent {
    //Use Kubernetes as dynamic Jenkins Slave
    kubernetes {
      //Kubernetes Manifest File to spin up Pod to do build
      yaml """
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: docker
    image: docker:25.0.3-dind
    imagePullPolicy: IfNotPresent
    command:
    - dockerd
    - --iptables=false
    - --tls=false
    - --host=unix:///var/run/docker.sock
    - --host=tcp://0.0.0.0:2375
    - --storage-driver=overlay2
    tty: true
    securityContext:
      privileged: true
  - name: helm
    image: lachlanevenson/k8s-helm:v3.10.2
    imagePullPolicy: IfNotPresent
    command:
    - cat
    tty: true
"""
    }//End kubernetes
  }//End agent
  
    environment {
        ENV_NAME = "${BRANCH_NAME == "main" ? "prd" : "${BRANCH_NAME}"}"

        GIT_SSH = "git@github.com:Live-Quiz-Project/Central-Backend.git"
        APP_NAME = "oqp-backend"
        IMAGE = "ghcr.io/phurits/oqp-backend"
        // SCANNER_HOME = tool 'sonarqube-scanner'
        // PROJECT_KEY = "bookinfo-reviews"
        // PROJECT_NAME = "bookinfo-reviews"
        // GOOGLE_APPLICATION_CREDENTIALS = credentials('gke-sa-key-json')
    }
  
  //Start Pipeline
  stages {
      
      // ***** Stage Clone *****
      stage('Clone reviews source code') {
        // Steps to run build
        steps {
          // Run in Jenkins Slave container
          container('jnlp') {
            //Use script to run
            script {
              // Git clone repo and checkout branch as we put in parameter
              scmVars = git branch: "${BRANCH_NAME}",
                            credentialsId: 'git',
                            url: "${GIT_SSH}"
            }// End Script
          }// End Container
        }// End steps
      }//End stage

      // ***** Stage Build *****
      stage('Build reviews Docker Image and push') {
          steps {
              container('docker') {
                script {
                    
                  // Do docker login authentication
                  docker.withRegistry('https://ghcr.io','ghcr-registry') {
                      // Do docker build and docker push
                      docker.build("${IMAGE}:${ENV_NAME}").push()
                 }// End docker
                }//End script
              }//End container
          }//End steps
      }//End stage

      stage('Deploy reviews with Helm Chart') {
          steps {
              // Run on Helm container
              container('helm') {
                  script {
                      // Use kubeconfig from Jenkins Credential
                      withKubeConfig([credentialsId: 'kubeconfig']) {
                          // Use Google Service Account IAM for Kubernetes Authentication

                          // Run the helm command with the service account key JSON file
                          sh "helm upgrade -i ${APP_NAME} k8s/helm -f k8s/helm-values/apps-${ENV_NAME}.yaml \
                              --wait --namespace ${ENV_NAME}"

                      } // End withKubeConfig
                  } // End script
              } // End container
          } // End steps
      } // End stage

  }// End stages
}// End pipeline
