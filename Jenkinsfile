pipeline {
    agent none
    stages {
        stage('Prepare') {
            agent { label 'master' }
            steps {
                sh 'git log --oneline | nl -nln | perl -lne \'if (/^(\\d+).*Version (\\d+\\.\\d+\\.\\d+)/) { print "$2-$1"; exit; }\' > version.txt'
                stash includes: 'version.txt', name: 'version'
            }
        }
        stage('Build') {
            agent { label 'docker' }
            steps {
                unstash 'version'
                script {
                    VERSION=readFile('version.txt').trim()
                    docker.build("fint-model:${env.BUILD_ID}", ".").inside {
                        dir '/go/src/app/vendor/github.com/FINTprosjektet/fint-model' {
                            sh "gox -output='./build/{{.Dir}}-{{.OS}}' -verbose -rebuild -osarch='darwin/amd64 windows/amd64' -ldflags='-X main.Version=${VERSION}'"
                            stash name: 'artifacts', includes: 'build/**'
                        }
                    }
                }
            }
        }
        stage('Publish') {
            agent { label 'docker' }
            when {
                branch 'master'
            }
            steps {
                unstash 'version'
                unstash 'artifacts'
                script {
                    VERSION=readFile('version.txt').trim()
                }
                archiveArtifacts 'build/**'
            }
        }
    }
}
