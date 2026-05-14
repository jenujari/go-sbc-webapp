pipeline {
    agent any

    triggers {
        githubPush()
    }

    environment {
        IMAGE_NAME = "docker.io/jhon5456/swewebapp"
        REGISTRY_CREDENTIALS = "docker-creds"
    }

    stages {

        stage('Checkout Code') {
            steps {
                checkout scm
            }
        }


        stage('Validate branch is main before build') {
            steps {
                script {
                    def branch = env.GIT_BRANCH?.replace("origin/", "")
                    env.BRANCH_NAME = branch

                    echo "Triggered branch: ${branch}"

                    if (branch != "main") {

                        currentBuild.result = 'NOT_BUILT'

                        echo """
                        Skipping futher build steps in pipeline.
                        Push detected on non-main branch:
                        ${branch}
                        """

                        return
                    }

                    echo "Main branch detected. Continuing pipeline..."
                }
            }
        }

        stage('Generate Commit SHA') {
            steps {
                script {
                    SHORT_SHA = sh(
                        script: 'git rev-parse --short HEAD',
                        returnStdout: true
                    ).trim()

                    env.IMAGE_TAG = SHORT_SHA

                    echo "Image tag set to: $IMAGE_TAG"
                }
            }
        }

        stage('Build Podman Image') {
            when {
                expression {
                    env.BRANCH_NAME == 'main'
                }
            }
            steps {
                sh """
                    podman build -t $IMAGE_NAME:$SHORT_SHA .
                    podman tag $IMAGE_NAME:$SHORT_SHA $IMAGE_NAME:latest
                """
            }
        }

        stage('Docker Login') {
            when {
                expression {
                    env.BRANCH_NAME == 'main'
                }
            }
            steps {
                withCredentials([
                    usernamePassword(
                        credentialsId: "${REGISTRY_CREDENTIALS}",
                        usernameVariable: 'DOCKER_USER',
                        passwordVariable: 'DOCKER_PASS'
                    )
                ]) {
                    sh 'echo $DOCKER_PASS | podman login docker.io -u $DOCKER_USER --password-stdin'
                }
            }
        }

        stage('Push Docker Image') {
            when {
                expression {
                    env.BRANCH_NAME == 'main'
                }
            }
            steps {
                sh """
                    podman push $IMAGE_NAME:$SHORT_SHA
                    podman push $IMAGE_NAME:latest
                """
            }
        }

        stage('Update Kubernetes Deployment') {
            when {
                expression {
                    env.BRANCH_NAME == 'main'
                }
            }
            steps {
                sh """
                    kubectl set image deployment/webapp \
                    webapp=$IMAGE_NAME:latest

                    kubectl rollout status deployment/webapp
                """
            }
        }
    }

    post {
        success {
            echo 'Deployment successful!'
        }

        failure {
            echo 'Pipeline failed!'
        }

        notBuilt {
            echo 'Pipeline skipped for non-main branch'
        }

        always {
            echo 'Cleaning workspace...'
            // upload coverage reports here
            cleanWs()
        }
    }
}
