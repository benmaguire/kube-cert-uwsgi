# kube-cert-uwsgi
Watcher for certificate changes in Kubernetes to gracefully reload UWSGI

Built to run as a binary inside a container running in kubernetes to listen for changes to certificates.
If a certificate change is detected based on the environment inputs, then UWSGI will reload based on the PID file identifier


