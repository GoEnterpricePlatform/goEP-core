# Contributing

Thank you for your interest in contributing to this project! 🚀  
We welcome contributions of all kinds, including bug fixes, new features, documentation improvements, and suggestions.

Please make sure to follow the guidelines below to ensure a smooth collaboration process and maintain code quality across the project.

## Keep Your Fork Up to Date with the Original Repository
Be sure to regularly synchronize your fork with the original (upstream) repository to stay aligned with the latest changes.


### To have the changes in your branch of work, example:
- Verify whether your local branch is up to date with the remote repository.
- Stash any uncommitted changes to avoid conflicts.
- Switch to the `dev` branch.
- Pull the latest changes from the remote repository.
- Switch back to your working branch (`feat/something`).
- Merge the updated `dev` branch into your feature branch.
- Reapply the stashed changes.
- If any conflicts arise, resolve them manually and complete the merge.

## Getting Started
1. **Fork the Project**: Click the "Fork" button at the top right corner of the repository page to create a copy of the repository under your own GitHub account.

2. **Clone the Repository**: Clone your forked repository to your local machine.
      ```sh
      git clone https://github.com/Tmplx/go-cms-tmpl
      cd go-cms-tmpl
      ```

3. **Create a Branch**: Create a new branch for your feature or bugfix.
   ```sh
    git switch -c feat/your-feature-name
   ```
4. **Link the Related Issue**: If there is an existing issue related to your work, link it in your commit messages and pull requests. 
    ```sh 
    git commit -m "Fixes #[issue-number]: [description of the fix]"
    ````
   
5. **Commit and Push Your Changes**: Commit your changes and push the branch to your fork. 
    ```sh
    git add .
    git commit -m "Add new feature [feature-description]"
    git push origin feat/your-feature-name
    ```
6. **Create a Pull Request**: Go to the original repository and create a pull request from your forked branch to the develop branch of the original repository.
    Provide a clear description of the changes and link any related issues.

    If you haven't done it already, associate the issue with the pull request. You can do this by adding the issue number in the pull request description or comments or by going to the "Development" section of your pull request and choosing the issue from the "Linked Issues" dropdown.

7. **Wait for Review**: Your pull request will be reviewed by the maintainers. Be responsive to feedback and make necessary adjustments.

## How to Claim an Issue
1. **Check Existing Issues**: Before starting work, check the [Issues](https://github.com/jairogloz/go-l/issues) section to see if the issue you want to work on is already claimed or in progress.

2. **Comment on the Issue**: If the issue is not already assigned, comment on the issue expressing your interest in working on it. Example:

    ```
    @maintainers, I would like to work on this issue.
    ```

3. **Wait for Assignment**: Wait for a maintainer to assign the issue to you. This helps prevent multiple people from working on the same issue simultaneously.
4. **Start Working on the Issue**: Once the issue is assigned, feel free to start working on it. If you need help or have questions, don't hesitate to leave a comment on the issue.
