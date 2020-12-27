# hermes
Hermes is a tool for deploying to and managing websites on S3+CloudFront

**HERMES IS NOT RELEASED. COME BACK LATER**

## Commands

### `hermes init`

Starts the wizard for a hermes setup. The wizard will create a new resource or use an existing one, if you want. The wizard works in the following order:
- Ask for domain(s), and ensure certificates are available, or request the certificates (then exit)
- Create S3 bucket and set correct permissions
- Create CloudFront distribution
- Deploy Lambda@Edge function for CloudFront
- Generate manifest and deploy default `index.html` and `404.html`


### `hermes list`

Lists all hermes setups within AWS account.


### `hermes verify`

Verifies a hermes setup by checking all resources are healthy


## `hermes deploy`

Deploys a folder to a hermes setup