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
            agent { 
                docker {
                    label 'docker'
                    image 'golang'
                    args "-v /tmp:/tmp -v .:/go/src/app/vendor/github.com/FINTprosjektet/fint-model -w /go/src/app/vendor/github.com/FINTprosjektet/fint-model bash"
                }
            }
            steps {
                unstash 'version'
                script {
                    VERSION=readFile('version.txt').trim()
                    sh "GOARCH=amd64; for GOOS in darwin windows; do go build -v -ldflags='-X main.Version=${VERSION}' -o fint-model-\$GOOS; done"
                    stash name: 'artifacts', includes: 'fint-model-*'
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
