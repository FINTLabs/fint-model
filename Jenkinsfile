pipeline {
    agent {
        label 'docker'
    }
    stages {
        stage('Build') {
            when {
                not { buildingTag() }
            }
            steps {
                sh "docker build --tag fint-model --build-arg VERSION=${BUILD_TAG} ."
                sh "docker run -i fint-model generate"
            }
        }
        stage('Deploy') {
            when {
                tag pattern: "v\\d+\\.\\d+\\.\\d+(-\\w+-\\d+)?", comparator: "REGEXP"
            }
            steps {
                script {
                    VERSION = TAG_NAME[1..-1]
                }
                sh "echo Version is ${VERSION}"
                sh "docker build --tag fint/fint-model:${VERSION} --build-arg VERSION=${VERSION} ."
                withDockerRegistry([credentialsId: 'asgeir-docker', url: '']) {
                    sh "docker push fint/fint-model:${VERSION}"
                    sh "docker tag fint/fint-model:${VERSION} fint/fint-model:latest"
                    sh "docker push fint/fint-model:latest"
                }
            }
        }
    }
}
