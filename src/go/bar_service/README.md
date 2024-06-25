# backup and restore service

This is a service intended to be run as a "sidecar" container inside the pod having data to be back up and restored.
This may be useful for pods using write-once PVs.

It provides a web service:

- POST /backup launch backup script
- POST /restore launch restore script
- GET /status provides status about operation in progress or last terminated
- GET /healthz health check

One operation maximum at a time is permitted.

Configuration done by environment:

- BAR_ADDRESS: host/port to listen to, defaults to :3000, format is [net.Dial](https://pkg.go.dev/net#Dial)
- BAR_BACKUP: backup command, defaults to `sh -c /etc/bar/backup.sh`
- BAR_RESTORE: restore command, defaults to `sh -c /etc/bar/restore.sh`
