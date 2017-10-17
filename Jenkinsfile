node('master') {
    stage('Prepare') {
        sh 'git log --oneline | nl -nln | perl -lne \'if (/^(\\d+).*Version (\\d+\\.\\d+\\.\\d+)/) { print "$2-$1"; exit; }\' > version.txt'
        stash includes: 'version.txt', name: 'version'
    }
}
node('docker') {
    stage('Build') {
        docker.image('golang').inside("-v /tmp:/tmp -v ${pwd()}:/go/src/app/vendor/github.com/FINTprosjektet/fint-model") {
            unstash 'version'
            VERSION=readFile('version.txt').trim()
            sh "cd /go/src/app/vendor/github.com/FINTprosjektet/fint-model; GOARCH=amd64; for GOOS in darwin windows; do go build -v -ldflags='-X main.Version=${VERSION}' -o fint-model-\$GOOS; done"
            stash name: 'artifacts', includes: 'fint-model-*'
        }
    }
    stage('Publish') {
        unstash 'version'
        unstash 'artifacts'
        VERSION=readFile('version.txt').trim()
        archiveArtifacts 'build/**'
    }
}
