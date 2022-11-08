# github-template-default

This is a basic template for when creating a Github repository.

## Initialization instructions:

1. Setup [CODEOWNERS](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners)

    In a new branch Add you own team as the default code owner. Replace 
    `@coopnorge/engineering` with `@coopnorge/your-team-here` in
    `./CODEOWNERS`. Leave the rest of the file as is.

    ```CODEOWNERS
    * @coopnorge/<your-team-here>
    
    Dockerfile @coopnorge/github-review-bots @coopnorge/<your-team-here>
    go.mod @coopnorge/github-review-bots @coopnorge/<your-team-here>
    go.sum @coopnorge/github-review-bots @coopnorge/<your-team-here>
    vendor @coopnorge/github-review-bots @coopnorge/<your-team-here>
    pyproject.toml @coopnorge/github-review-bots @coopnorge/<your-team-here>
    poetry.lock @coopnorge/github-review-bots @coopnorge/<your-team-here>
    
    /CODEOWNERS @coopnorge/cloud-security
    .github/workflows/security-* @coopnorge/cloud-security
    ```

    Create a pull request, get it approved and merge it.

2. Setup and fix default [workflows](https://docs.github.com/en/actions/using-workflows)

    `.github/workflows/build.yaml` declares the default required GitHub Actions
    job `build`. The job will fail on all builds in all repositories, except
    <https://github.com/coopnorge/github-template-default>. Update the workflow
    to do something that actually validates the content of your repository.

3. Setup
    [dependabot](https://playbook.internal.coop/platforms/cloud_platform/dev_build_deploy/github/guide_github_dependabot.html)
    to update all dependencies from all ecosystems in the repository.

4. Create a new branch and start initializing your repository with the code you
   need.
