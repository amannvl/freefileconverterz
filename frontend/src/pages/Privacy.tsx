import { Link } from 'react-router-dom';

function Privacy() {
  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto bg-white p-8 rounded-lg shadow">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Privacy Policy</h1>
          <p className="mt-2 text-sm text-gray-600">Last updated: July 2, 2025</p>
        </div>

        <div className="prose prose-indigo max-w-none">
          <h2>1. Information We Collect</h2>
          <p>
            When you use our service, we may collect the following information:
          </p>
          <ul className="list-disc pl-5 space-y-2 mt-2">
            <li>Files you upload for conversion</li>
            <li>IP address and browser information</li>
            <li>Usage data and analytics</li>
            <li>Email address (if provided for notifications)</li>
          </ul>

          <h2>2. How We Use Your Information</h2>
          <p>We use the information we collect to:</p>
          <ul className="list-disc pl-5 space-y-2 mt-2">
            <li>Provide and maintain our service</li>
            <li>Process file conversions</li>
            <li>Improve our service and user experience</li>
            <li>Monitor usage and analyze trends</li>
            <li>Communicate with you about your account or our services</li>
          </ul>

          <h2>3. Data Retention</h2>
          <p>
            We retain your uploaded files only for the duration necessary to complete the conversion process.
            All files are automatically deleted from our servers after 24 hours. We may retain metadata and
            usage statistics for analytical purposes, but this data is anonymized and cannot be used to
            identify individual users.
          </p>

          <h2>4. Data Security</h2>
          <p>
            We implement appropriate technical and organizational measures to protect your personal data
            against unauthorized or unlawful processing, accidental loss, destruction, or damage.
          </p>

          <h2>5. Third-Party Services</h2>
          <p>
            We may use third-party services to help operate our service, such as hosting providers and
            analytics services. These third parties have access to your information only to perform specific
            tasks on our behalf and are obligated not to disclose or use it for any other purpose.
          </p>

          <h2>6. Your Rights</h2>
          <p>You have the right to:</p>
          <ul className="list-disc pl-5 space-y-2 mt-2">
            <li>Access the personal data we hold about you</li>
            <li>Request correction or deletion of your personal data</li>
            <li>Object to processing of your personal data</li>
            <li>Request restriction of processing your personal data</li>
            <li>Request transfer of your personal data</li>
          </ul>

          <h2>7. Changes to This Policy</h2>
          <p>
            We may update our Privacy Policy from time to time. We will notify you of any changes by posting
            the new Privacy Policy on this page and updating the "Last updated" date.
          </p>

          <h2>8. Contact Us</h2>
          <p>
            If you have any questions about this Privacy Policy, please contact us at
            <a href="mailto:privacy@freefileconverterz.com" className="text-indigo-600 hover:text-indigo-800">
              {' '}
              privacy@freefileconverterz.com
            </a>.
          </p>
        </div>

        <div className="mt-8 pt-6 border-t border-gray-200">
          <Link
            to="/"
            className="inline-flex items-center text-indigo-600 hover:text-indigo-800"
          >
            <svg
              className="h-5 w-5 mr-2"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M10 19l-7-7m0 0l7-7m-7 7h18"
              />
            </svg>
            Back to Home
          </Link>
        </div>
      </div>
    </div>
  );
}

export default Privacy;
