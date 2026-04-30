# SSH Hardening

SSH changes can lock you out of a remote VPS. Keep an existing SSH session open and verify key-based login in a second session before applying hardening.

## Review current SSH posture

```bash
sudo proxctl ssh harden --plan
```

The plan checks the effective `sshd -T` config and recommends safer values.

Target posture:

```text
PasswordAuthentication no
PermitRootLogin prohibit-password
KbdInteractiveAuthentication no
X11Forwarding no
MaxAuthTries 3
```

## Apply hardening

```bash
sudo proxctl ssh harden --apply
```

You must type:

```text
SSH
```

Proxctl writes a drop-in config and runs `sshd -t` before reloading SSH.
