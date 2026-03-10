# Gitflow Workflow

This repository uses a lightweight Gitflow model focused on safe PR delivery.

## Branches
- `main`: production-ready branch. Merge by PR only.
- `develop`: integration branch for upcoming release work.
- `feature/<short-name>`: feature or improvement branches created from `develop`.
- `fix/<short-name>`: bugfix branches created from `develop`.
- `hotfix/<short-name>`: urgent fixes created from `main`.
- `release/<version>`: release hardening branches created from `develop`.

## Daily Flow
1. Sync local branches:
```bash
git checkout main && git pull origin main
git checkout develop && git pull origin develop
```
2. Create a work branch from `develop`:
```bash
git checkout -b feature/drawio-export-improvements develop
```
3. Commit in small, testable increments:
```bash
git add -A
git commit -m "feat: improve draw.io export icon mapping"
```
4. Push your branch:
```bash
git push -u origin feature/drawio-export-improvements
```
5. Open a PR to `develop`.

## PR Checklist
- All tests pass locally (`go test ./...`).
- Documentation and examples are updated when behavior changes.
- Commit messages are explicit (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`).
- At least 1 reviewer approves before merge.
- Use squash merge when the branch has many micro-commits.

## Release Flow
1. Create release branch from `develop`:
```bash
git checkout -b release/vX.Y.Z develop
```
2. Update `package.json` version to match the release tag:
```bash
npm version X.Y.Z --no-git-tag-version
```
3. Run final checks and bugfixes only.
4. Open PR `release/vX.Y.Z -> main`.
5. After merge to `main`, tag the release:
```bash
git checkout main
git pull origin main
git tag vX.Y.Z
git push origin vX.Y.Z
```
6. GitHub Actions automatically:
   - Builds cross-platform binaries and creates the GitHub Release
   - Publishes the npm package to https://www.npmjs.com/package/awsdac
7. Merge `main` back into `develop` to keep both branches synchronized.

## Hotfix Flow
1. Create branch from `main`:
```bash
git checkout -b hotfix/critical-icon-fix main
```
2. Commit and open PR to `main`.
3. After merge to `main`, merge/cherry-pick into `develop`.

## Protected Branch Rules (Recommended)
- Protect `main` and `develop`.
- Require PR before merge.
- Require at least 1 approval.
- Require status checks (`go test ./...`) to pass.
- Disallow force pushes.
