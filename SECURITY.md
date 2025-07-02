# Security Policy

## Supported Versions

We provide security updates for the following versions of FreeFileConverterZ:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0.0 | :x:                |

## Reporting a Vulnerability

We take all security vulnerabilities seriously. Thank you for improving the security of FreeFileConverterZ. We appreciate your efforts and responsible disclosure and will make every effort to acknowledge your contributions.

### How to Report a Security Vulnerability

If you discover a security vulnerability in FreeFileConverterZ, please report it by emailing [security@freefileconverterz.com](mailto:security@freefileconverterz.com) with the subject line "SECURITY: [brief description of the issue]".

Please include the following details in your report:
- A description of the vulnerability
- Steps to reproduce the issue
- The version of FreeFileConverterZ you are using
- Any potential impact of the vulnerability
- Your contact information (optional)

### What to Expect

- You will receive a response from our security team within 48 hours
- We will work with you to understand and validate the issue
- We will keep you informed of the progress towards resolving the issue
- Once the issue is resolved, we will publicly acknowledge your contribution (unless you prefer to remain anonymous)

### Security Best Practices

To help keep FreeFileConverterZ secure, we recommend the following:

1. **Keep your installation up to date**
   - Always run the latest stable version of FreeFileConverterZ
   - Subscribe to security announcements for the project

2. **Secure your server**
   - Keep your server's operating system and software up to date
   - Use strong, unique passwords for all accounts
   - Implement proper firewall rules
   - Use HTTPS with a valid SSL certificate
   - Regularly back up your data

3. **Secure your application**
   - Use strong, unique passwords for all user accounts
   - Enable two-factor authentication if available
   - Regularly review and update access controls
   - Monitor for suspicious activity

### Security Features

FreeFileConverterZ includes the following security features:

- Secure file upload handling
- Input validation and sanitization
- Protection against common web vulnerabilities (XSS, CSRF, SQL injection, etc.)
- Secure password hashing
- Rate limiting
- Secure HTTP headers
- Content Security Policy (CSP)
- Regular security audits

### Responsible Disclosure Policy

We follow a responsible disclosure policy for security vulnerabilities. This means:

1. Do not publicly disclose the vulnerability until we've had time to address it
2. Allow us a reasonable amount of time to fix the issue before making it public
3. Make a good faith effort to avoid privacy violations, destruction of data, and interruption or degradation of our service

### Security Updates

Security updates will be released as patch versions (e.g., 1.0.1, 1.0.2). We recommend always running the latest patch version of your installed major.minor version.

### Third-party Dependencies

We regularly update our dependencies to ensure known security vulnerabilities are addressed. You can check for known vulnerabilities in our dependencies using:

```bash
# For Go dependencies
go list -json -m all | nancy sleuth

# For Node.js dependencies
npm audit
```

### Additional Security Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Mozilla Web Security Guidelines](https://infosec.mozilla.org/guidelines/web_security)
- [CIS Benchmarks](https://www.cisecurity.org/cis-benchmarks/)

## Security Contact

For any security-related questions or concerns, please contact us at [security@freefileconverterz.com](mailto:security@freefileconverterz.com).

---

*Last updated: July 2, 2025*
