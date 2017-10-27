String commitMessage = "unknown"
def isRelease = false
String commitVersion = "0.0.0"
node('master') {
    stage('Prepare') {
        checkout scm
        sh 'git log --oneline | nl -nln | perl -lne \'if (/^(\\d+).*Version (\\d+\\.\\d+\\.\\d+)/) { print "$2-$1"; exit; }\' > version.txt'
        stash includes: 'version.txt', name: 'version'
        commitMessage = sh(script: "git log -1 --pretty=%B", returnStdout: true).trim()
        def m = commitMessage =~ /Version (\d+\.\d+\.\d+)/
        if (m) {
            commitVersion = m.group(1)
            isRelease = true
        }
    }
}
node('docker') {
    stage('Build') {
        String goPath = "/go/src/app/vendor/github.com/FINTprosjektet/fint-model"
        docker.image('golang').inside("-v /tmp:/tmp -v ${pwd()}:${goPath}") {
            sh "go-wrapper download github.com/mitchellh/gox && go-wrapper install github.com/mitchellh/gox"
            unstash 'version'
            VERSION=readFile('version.txt').trim()
            sh "cd ${goPath}; gox -output=\"./{{.Dir}}-${VERSION}-{{.OS}}\" -rebuild -osarch=\"darwin/amd64 windows/amd64\" -ldflags='-X main.Version=${VERSION}'"
            stash name: 'artifacts', includes: 'fint-model-*'
        }
    }
    stage('Publish') {
        unstash 'artifacts'
        archiveArtifacts 'fint-model-*'
    }
    if (isRelease) {
        stage('Release') {
            docker.image('asgeir/gothub').inside {
                unstash 'artifacts'
                sh "pwd; ls -l"
                echo "This is the commit message: <<${commitMessage}>>"
                echo "This is the version: <<${commitVersion}>>"
                withCredentials([string(credentialsId: 'github_fint_jenkins', variable: 'GITHUB_TOKEN')]) {
                    sh "gothub release --user FINTprosjektet --repo fint-model --tag v${commitVersion} --name v${commitVersion} --description '${commitMessage}'"
                    sh "for f in fint-model-*; do gothub upload --user FINTprosjektet --repo fint-model --tag v${commitVersion} -fint-model-* \$f --file \$f; done"
                }
            }
        }        
    }
}
