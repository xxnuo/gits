# gits
Git over SSH.

# Usage Sample

1. Make sure `~/.ssh/config` contains your server config.

> Currently, only `IdentityFile` is supported.

```
Host server1
    HostName 192.168.1.100
    User root
    IdentityFile ~/.ssh/id_ed25519
```

1. Create `config.toml` in a empty directory as a git repo.

```bash
mkdir -p ~/work/llamacpp
cd ~/work/llamacpp
touch config.toml
```

3. Edit `config.toml` as below:

```toml
[ssh]
# host name in ssh config file
Host = "server1"

[repo]
# path to the repository on the remote server
Path = "~/work/llamacpp"
```

4. Run `gits` in the repo directory.

```bash
gits
```

5. Run `gits` with arguments.