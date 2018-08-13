node('docker') {
    stage('Build') {
        println(env)
        sh "env"
        checkout scm
        sh "pwd; ls -l"
        String goPath = "/go/src/app/vendor/github.com/FINTprosjektet/fint-model"
        docker.image('golang').inside("-v /tmp:/tmp -v ${pwd()}:${goPath}") {            
            sh "pwd; ls -l"
            sh "go get github.com/mitchellh/gox && go install github.com/mitchellh/gox"
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
                    sh "for f in fint-model-*; do gothub upload --user FINTprosjektet --repo fint-model --tag v${commitVersion} --name \$f --file \$f; done"
                }
            }
        }        
    }
    stage('Cleanup') {
        deleteDir()
    }
}
