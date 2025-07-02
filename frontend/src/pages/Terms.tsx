import { Link } from 'react-router-dom';

function Terms() {
  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto bg-white p-8 rounded-lg shadow">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Terms and Conditions</h1>
          <p className="mt-2 text-sm text-gray-600">Last updated: July 2, 2025</p>
        </div>

        <div className="prose prose-indigo max-w-none">
          <h2>1. Introduction</h2>
          <p>
            Welcome to FreeFileConverterZ. By using our service, you agree to these terms and conditions.
            Please read them carefully before using our service.
          </p>

          <h2>2. Use of Service</h2>
          <p>
            Our service allows you to convert files between different formats. You agree to use our service
            only for lawful purposes and in accordance with these Terms.
          </p>

          <h2>3. User Responsibilities</h2>
          <p>You agree not to use the service to:</p>
          <ul className="list-disc pl-5 space-y-2 mt-2">
            <li>Upload or share any content that is illegal, harmful, or infringes on intellectual property rights</li>
            <li>Attempt to gain unauthorized access to our systems or interfere with the service</li>
            <li>Use the service for any commercial purposes without our express written consent</li>
            <li>Upload files containing viruses or malicious code</li>
          </ul>

          <h2>4. Intellectual Property</h2>
          <p>
            The service and its original content, features, and functionality are owned by FreeFileConverterZ
            and are protected by international copyright, trademark, and other intellectual property laws.
          </p>

          <h2>5. Limitation of Liability</h2>
          <p>
            In no event shall FreeFileConverterZ be liable for any indirect, incidental, special,
            consequential or punitive damages, including without limitation, loss of profits, data, use,
            goodwill, or other intangible losses.
          </p>

          <h2>6. Changes to Terms</h2>
          <p>
            We reserve the right to modify these terms at any time. We will provide notice of any changes
            by posting the updated terms on our website.
          </p>

          <h2>7. Contact Us</h2>
          <p>
            If you have any questions about these Terms, please contact us at
            <a href="mailto:legal@freefileconverterz.com" className="text-indigo-600 hover:text-indigo-800">
              {' '}
              legal@freefileconverterz.com
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

export default Terms;
