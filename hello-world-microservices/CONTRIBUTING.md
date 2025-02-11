# Contributing to Hello World Microservices

Thank you for considering contributing to the **Hello World Microservices** project! Your contributions help improve the quality, scalability, and usability of this project.

## How to Contribute

### 1. Fork the Repository
1. Navigate to the repository on GitHub.
2. Click on the **Fork** button.
3. Clone your forked repository:
   ```sh
   git clone https://github.com/your-username/hello-world-microservices.git
   cd hello-world-microservices
   ```

### 2. Create a Feature Branch
- Ensure you're on the `main` branch before creating a new branch.
- Use meaningful branch names such as `feature-logging-improvements`.
  ```sh
  git checkout -b feature-your-feature-name
  ```

### 3. Make Changes and Commit
- Follow **code style guidelines** for Golang, Python, and C#.
- Ensure your code is **well-documented** and **test-covered**.
- Format your code before committing.
  ```sh
  # For Golang
  go fmt ./...
  
  # For Python
  black .
  
  # For C#
  dotnet format
  ```
- Commit changes with a **clear message**:
  ```sh
  git commit -m "Add OpenTelemetry tracing to Golang service"
  ```

### 4. Push and Create a Pull Request (PR)
```sh
git push origin feature-your-feature-name
```
- Go to GitHub and open a **Pull Request (PR)**.
- Provide a **detailed description** of your changes.
- Link related **issues** if applicable.

## Code Style Guidelines
- Follow **industry best practices** for coding, logging, and observability.
- Use **structured logging** with OpenTelemetry.
- Ensure **metrics** conform to Prometheus naming conventions.

## Reporting Issues
If you encounter a bug, performance issue, or have an idea for improvement:
1. Open a **GitHub Issue**.
2. Provide **steps to reproduce**, expected behavior, and actual results.

## License
By contributing, you agree that your contributions will be licensed under the **MIT License**.

---

Thank you for contributing to **Hello World Microservices**!

