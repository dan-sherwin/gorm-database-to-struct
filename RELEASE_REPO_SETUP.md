# YUM/DNF Repository Publication

This repository publishes RPM packages to a Yum/DNF repo hosted via GitHub Pages under:

  https://dan-sherwin.github.io/gormdb2struct/rpm/$basearch/

Where `$basearch` is either `x86_64` or `aarch64`.

## How it works

- On every GitHub Release (published), the workflow `.github/workflows/publish_yum_repo.yml`:
  1. Downloads the release RPM artifacts for linux amd64/arm64.
  2. Organizes them into `public/rpm/x86_64` and `public/rpm/aarch64`.
  3. Runs `createrepo_c` to produce `repodata/` for each arch directory.
  4. Publishes the `public/` directory to the `gh-pages` branch.

- GitHub Pages must be enabled for this repository (Settings → Pages → Build from `gh-pages`).

## Client setup

Users can install via DNF by creating `/etc/yum.repos.d/gormdb2struct.repo`:

```
[gormdb2struct]
name=gormdb2struct
baseurl=https://dan-sherwin.github.io/gormdb2struct/rpm/$basearch/
enabled=1
gpgcheck=0
```

Then run:

```
sudo dnf clean all
sudo dnf makecache
sudo dnf install gormdb2struct
```

## Optional: GPG signing

If you want to sign RPMs and enable `gpgcheck=1`:
1. Generate a GPG key and store the private key in GitHub secrets.
2. Configure nfpm signing in `.goreleaser.yaml` (rpm.signature section) and pass secrets in the release workflow.
3. Publish the public key at `public.key` to the `gh-pages` branch, and update the repo file to include:

```
gpgcheck=1
gpgkey=https://dan-sherwin.github.io/gormdb2struct/public.key
```
