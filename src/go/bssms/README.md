# bootstrap secrets management service

How to provide secrets to a cloud-init booting server with minimal interception risk.

For instance:

- gitops: git/registry secret, sops secret
- basic deployment: web server TLS certificate private key, service account password 

## Design

### Actors

- Secrets Provisioner: runs aside IaC provisioner tool,
e.g. could be integrated as OpenTofu provider,
trustfully sends secrets to Secrets Installer and instructions for the latter to install them
- Secrets Installer: run by cloud init,
the installer receives secrets and instructions from Provisioner
- Proxy Service: securely connect both

### How security risks are addressed

MITM attacks exist and the actors are possibly dealing
with root secrets and cloud resources.
Bssms security is a kind of tradeoff between keeping it simple and remaining reasonably secure.

The Provisioner must not send secrets to an untrusted Installer.
The Installer must provide to the Provisioner various unique identity properties
that the Provisioner knows from another channel thanks to the IaC provisioning tool
(e.g. server uuid, mac and IP addresses).
This is what a human operator would do.

The Installer must not install secrets or execute instructions from an untrusted Provisioner.
This information is encrypted by the Provisioner using a temporary public key only known by it.
The corresponding private key is provided to the Installer with cloud-init.
The Installer is kept alive for a short period of time.

The Installer also communicates with the Provisioner using its temporary public key
provided with cloud-init.

The Proxy Service acts as a pass through between the Provisioners and the Installers.
The information exchanged is encrypted with secret keys that are unknown to the Proxy Service.
