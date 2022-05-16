package main

import (
        "os"
        "os/exec"
        "context"
        "time"
        log "github.com/sirupsen/logrus"

        "k8s.io/client-go/kubernetes"
        "k8s.io/client-go/rest"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        v1 "k8s.io/api/core/v1"

        certv1alpha2 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha2"
        cmclient "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"

)


var exit = make(chan bool)


func main() {

        log.Info("Starting Application")
        certname := os.Getenv("CERT_NAME")
        namespace := os.Getenv("CERT_NAMESPACE")
	pidfile := os.Getenv("PID_FILE")

        // creates the in-cluster config
        config, err := rest.InClusterConfig()
        if err != nil {
            panic(err.Error())
        }
        // creates the clientset
        clientset, err := kubernetes.NewForConfig(config)
        if err != nil {
            panic(err.Error())
        }

        crtClient, err := cmclient.NewForConfig(config)
        if err != nil {
                panic(err)
        }



        go func () {
                for {

                        watch, err := clientset.CoreV1().Secrets(namespace).Watch(context.TODO(), metav1.ListOptions{})

                        if err != nil {
                                log.Fatal(err.Error())
                        }

                        for event := range watch.ResultChan() {
                                log.Info("Secret: Type: " + event.Type)
                                p, ok := event.Object.(*v1.Secret)
                                if !ok {
                                        log.Fatal("Secret: unexpected type")
                                }
                                log.Info("Secret: " + p.Name)
                        }
                        log.Info("Secret: Breaking out of loop")
                        time.Sleep(5 * time.Second)
                }
        }()





         go func () {
                for {
                        watchcert, err := crtClient.CertmanagerV1alpha2().Certificates(namespace).Watch(context.TODO(), metav1.ListOptions{})

                        if err != nil {
                                log.Fatal(err.Error())
                        }

                        for event := range watchcert.ResultChan() {
                                log.Info("Cert: Type: " + event.Type)
                                p, ok := event.Object.(*certv1alpha2.Certificate)
                                if !ok {
                                        log.Fatal("Cert: unexpected type")
                                }
                                log.Info("Cert: " + p.Name)

                                if p.Name == certname && event.Type == "MODIFIED" {
                                        // Process Restart
                                        log.Info("Process Restart")
                                        cmd := exec.Command("uwsgi", "--reload", pidfile)
                                        stdout, err := cmd.StdoutPipe()
                                        if err != nil {
                                                log.Fatal(err)
                                        }

                                        if err := cmd.Start(); err != nil {
                                                log.Fatal(err)
                                        }
                                        _ = stdout

                                }

                        }
                log.Info("Cert: Breaking out of loop")
                time.Sleep(5 * time.Second)
                }
        }()


        <-exit

}


