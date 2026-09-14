# Contribution Guidelines & Legal Notices

**Organization:** FJ™ Cybertronic Systems®
**Repository:** FJ-cybetzilla/aud-dl
**Support & Security Desk:** [cyberzilla.systems@gmail.com](mailto:cyberzilla.systems@gmail.com)
**Application:** Aud-Dl (Enterprise Audio Downloader & Stream Engine)

---

## Part 1: Contribution Guidelines

Thank you for your interest in contributing to Aud-Dl! We maintain high standards for code quality, architectural consistency, and security across all internal Go packages.

### 1. Code Standards & Architecture
*   **Modularity:** All new features must be scoped within their appropriate `internal/` package. Do not bloat `cmd/app/main.go`.
*   **Architecture:** Maintain strict unidirectional layering: `cmd/app` -> `cli/` -> `task/` -> `vault/` / `processor/`.
*   **Go Style:** Adhere strictly to idiomatic Go 1.22+ patterns, explicit error management (custom typed errors), and zero global state.
*   **Context Propagation:** All long-running workers must accept and respect `context.Context` for immediate, graceful cancellation.
*   **Unicode Safety:** Ensure all filenames and metadata paths are normalized (NFC) via `golang.org/x/text/unicode/norm`.
*   **Testing:** All new modules, features, or bug fixes **must** include comprehensive unit tests (e.g., `_test.go` files corresponding to your package).

### 2. Pull Request Workflow
1.  **Fork & Branch:** Fork the repository and create your feature branch: `git checkout -b feature/my-feature`.
2.  **Validate:** Ensure code conforms to standards and passes all tests:
    ```bash
    go test ./... -v
    go build ./internal/...
    ```
3.  **Submit PR:** Open a pull request against the `main` branch with a detailed description of the changes.

---

## Part 2: Legal Notices & Proprietary Rights

### 1. Intellectual Property Notice
© 2026 FJ™ Cybertronic Systems®. All Rights Reserved.

Aud-Dl®, Cybertronic Systems®, and associated logos are proprietary trademarks of FJ™ Cybertronic Systems®. Unauthorized commercial distribution, reverse engineering, or sub-licensing of the core architecture without express written permission is strictly prohibited.

### 2. Disclaimer of Warranty
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS, COPYRIGHT HOLDERS, OR FJ™ CYBERTRONIC SYSTEMS® BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

### 3. Responsible Usage Policy
Aud-Dl is designed for enterprise media management, lawful content archiving, and stream transcoding. Users and contributors are solely responsible for ensuring that their utilization of download pipelines, search providers, and anti-bot defusers complies with applicable local laws, copyright regulations, and target platform terms of service.

For official inquiries or vulnerability disclosures, contact us securely at [cyberzilla.systems@gmail.com](mailto:cyberzilla.systems@gmail.com).
