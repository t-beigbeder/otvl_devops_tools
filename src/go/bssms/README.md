# bootstrap secrets management service

How to provide secrets to a cloud-init booting server with minimal interception risk.

For instance:

- gitops: git/registry secret, sops secret

## Protocol:

- provider (QUIC):
  - P1 - initiate temporary pub/pri key-pair for each booting server's uuid/IP, 
  send temporary pri key/IP to service, keep pub secret
  - P2 - encrypt secrets with temporary pub key, send to service,
  not before B1 occurred (see below)

- booting server (HTTPS):
  - B1 - looping get temporary pri key
  - B2 - looping get encrypted secrets
  
- service:
  - on P1: keep state(uuid/IP) = temporary pri key
  - on P2 before B1: remove related data
  - on timeout before B1: remove related data
  - on B2 before B1: remove related data
  - on B1: check IP (spoofing not possible with modern TCP), returns temporary pri key
  - on P2: keep state(uuid/IP) = encrypted secrets, remove state(uuid/IP) = temporary pri key
  - on timeout before B2: remove related data
  - on B2: check IP, returns encrypted secrets
  - on short timeout after B2: remove related data
