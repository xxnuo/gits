# gits

Git over SSH.

Anything over SSH.

## Usage

1. Make sure `~/.ssh/config` contains your server config.

> Currently, only `IdentityFile` is supported.

```
Host server1
    HostName 192.168.1.100
    User user1
    IdentityFile ~/.ssh/id_ed25519
```

2. Connect to the server at first to make sure the server is authorized.

```bash
ssh server1
```

3. Create `config.toml` in a empty directory as a git repo.

```bash
mkdir -p ~/work/llamacpp
cd ~/work/llamacpp
touch config.toml
```

4. Edit `config.toml` as below:

```toml
[ssh]
# host name in ssh config file
Host = "server1"

[repo]
# path to the repository on the remote server
Path = "~/work/llamacpp"
```

> If not provide `config.toml`, `gits` will call `git` directly.

5. Run `gits` in the repo directory.

```bash
gits
```

6. Done. It just like `git` command in local.

## Advanced

1. If you want to use other command in remote server, you can set `mode` to the command in `config.toml`.

```toml
[mode]
# default is "git"
exec = "uname"
```

Run `gits` got `uname` command in remote server.

```
Linux
```