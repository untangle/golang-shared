void builddiscoverd(String libc, String buildDir) {
    // sh "docker pull untangleinc/discoverd:build-${libc}"
    sh "docker-compose -f ${buildDir}/build/docker-compose.build.yml -p discoverd_${libc} up --force-recreate --abort-on-container-exit --build ${libc}-local"
    sh "cp ${buildDir}/cmd/discoverd/discoverd cmd/discoverd/discoverd-${libc}"
}

void lintdiscoverd(String libc, String buildDir) {
  sh "docker-compose -f ${buildDir}/build/docker-compose.build.yml -p discoverd_${libc} up --force-recreate --abort-on-container-exit --build ${libc}-lint"
}

void archivediscoverd() {
    archiveArtifacts artifacts:'cmd/discoverd/discoverd*', fingerprint: true
}

pipeline {
    agent none

    stages {
        stage('Build') {
            parallel {
                stage('Build musl') {
                    agent { label 'docker' }

                    environment {
                        libc = 'musl'
                        buildDir = "${env.HOME}/build-discoverd-${env.BRANCH_NAME}-${libc}/go/src/github.com/untangle/discoverd"
                    }

                    stages {
                        stage('Prep WS musl') {
                            steps { dir(buildDir) { checkout scm } }
                        }

                        stage('Build discoverd musl') {
                            steps {
                                builddiscoverd(libc, buildDir)
                                stash(name:"discoverd-${libc}", includes:"cmd/discoverd/discoverd*")
                            }
                        }
                    }

                    post {
                        success { archivediscoverd() }
                    }
                }

                stage('Build glibc') {
                    agent { label 'docker' }

                    environment {
                        libc = 'glibc'
                        buildDir = "${env.HOME}/build-discoverd-${env.BRANCH_NAME}-${libc}/go/src/github.com/untangle/discoverd"
                    }

                    stages {
                        stage('Prep WS glibc') {
                            steps { dir(buildDir) { checkout scm } }
                        }

                        stage('Build discoverd glibc') {
                            steps {
                                builddiscoverd(libc, buildDir)
                                stash(name:"discoverd-${libc}", includes:'cmd/discoverd/discoverd*')
                            }
                        }
                    }

                    post {
                        success { archivediscoverd() }
                    }
                }
            }
        }

        stage('Lint') {

            parallel {
                stage('Lint musl') {
                    agent { label 'docker' }

                    environment {
                        libc = 'musl'
                        buildDir = "${env.HOME}/build-discoverd-${env.BRANCH_NAME}-${libc}/go/src/github.com/untangle/discoverd"
                    }

                    stages {
                        stage('Prep WS musl') {
                            steps { dir(buildDir) { checkout scm } }
                        }

                        stage('Lint discoverd musl') {
                            steps {
                                sshagent (credentials: ['buildbot']) {
                                    lintdiscoverd(libc, buildDir)
                                }
                            }
                        }
                    }
                }

                stage('Lint glibc') {
                    agent { label 'docker' }

                    environment {
                        libc = 'glibc'
                        buildDir = "${env.HOME}/build-discoverd-${env.BRANCH_NAME}-${libc}/go/src/github.com/untangle/discoverd"
                    }

                    stages {
                        stage('Prep WS glibc') {
                            steps { dir(buildDir) { checkout scm } }
                        }

                        stage('Lint discoverd glibc') {
                            steps {
                                sshagent (credentials: ['buildbot']) {
                                    lintdiscoverd(libc, buildDir)
                                }
                            }
                        }
                    }

                }
            }
        }

        stage('Test') {
            parallel {
                stage('Test musl') {
                    agent { label 'docker' }

                    environment {
                        libc = 'musl'
                        discoverd = "cmd/discoverd/discoverd-${libc}"
                    }

                    stages {
                        stage('Prep musl') {
                            steps {
                                unstash(name:"discoverd-${libc}")
                            }
                        }

                        stage('File testing for musl') {
                            steps {
                                sh "test -f ${discoverd} && file ${discoverd} | grep -q /ld-musl"
                            }
                        }
                    }
                }

                stage('Test libc') {
                    agent { label 'docker' }

                    environment {
                        libc = 'glibc'
                        discoverd = "cmd/discoverd/discoverd-${libc}"
                    }

                    stages {
                        stage('Prep libc') {
                            steps {
                                unstash(name:"discoverd-${libc}")
                            }
                        }

                        stage('File testing for libc') {
                            steps {
                                sh "test -f ${discoverd} && file ${discoverd} | grep -q /ld-linux"
                            }
                        }
                        
                    }
                }
            }

            post {
                changed {
                    script {
                        // set result before pipeline ends, so emailer sees it
                        currentBuild.result = currentBuild.currentResult
                    }
                    emailext(to:'nfgw-engineering@untangle.com', subject:"${env.JOB_NAME} #${env.BUILD_NUMBER}: ${currentBuild.result}", body:"${env.BUILD_URL}")
                    slackSend(channel:'#team_engineering', message:"${env.JOB_NAME} #${env.BUILD_NUMBER}: ${currentBuild.result} at ${env.BUILD_URL}")
                }
            }
        }
    }
}
