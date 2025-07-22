pipeline {
    agent any

    environment {
        JAVA_HOME = '/usr/lib/jvm/java-11-openjdk-amd64'
    }

    stages {
        stage('Checkout') {
            steps {
                // Descargar el código fuente
                checkout scm
            }
        }

        stage('Build') {
            steps {
                // Ejecutar la compilación, en este caso para un proyecto Java
                sh 'mvn clean install'
            }
        }

        stage('Test') {
            steps {
                // Ejecutar pruebas, por ejemplo con Maven
                sh 'mvn test'
            }
        }

        stage('Deploy') {
            steps {
                // Desplegar, por ejemplo, en un servidor de producción o staging
                echo 'Deploying application...'
            }
        }
    }

    post {
        success {
            echo 'Pipeline finished successfully!'
        }
        failure {
            echo 'Pipeline failed!'
        }
    }
}