function GDPR() {
  return (
    <div className="max-w-4xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
      <div className="text-center mb-12">
        <h1 className="text-4xl font-extrabold text-gray-900 sm:text-5xl sm:tracking-tight lg:text-6xl">
          GDPR Compliance
        </h1>
        <p className="mt-3 max-w-2xl mx-auto text-xl text-gray-500 sm:mt-4">
          General Data Protection Regulation
        </p>
      </div>

      <div className="prose prose-indigo max-w-none">
        <div className="bg-white shadow overflow-hidden sm:rounded-lg mb-12">
          <div className="px-4 py-5 sm:px-6">
            <h2 className="text-2xl leading-6 font-medium text-gray-900">
              Your Data Protection Rights
            </h2>
            <p className="mt-1 text-sm text-gray-500">
              Last updated: July 2, 2025
            </p>
          </div>
        </div>

        <div className="space-y-8">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">Your Rights Under GDPR</h2>
            <p className="mt-4">
              The General Data Protection Regulation (GDPR) provides you with specific rights regarding your personal data. 
              As a data subject, you have the following rights:
            </p>
            <ul className="list-disc pl-5 space-y-2 mt-4">
              <li><strong>Right to Access:</strong> You have the right to request copies of your personal data.</li>
              <li><strong>Right to Rectification:</strong> You have the right to request correction of any information you believe is inaccurate.</li>
              <li><strong>Right to Erasure:</strong> You have the right to request that we erase your personal data, under certain conditions.</li>
              <li><strong>Right to Restrict Processing:</strong> You have the right to request that we restrict the processing of your personal data.</li>
              <li><strong>Right to Object to Processing:</strong> You have the right to object to our processing of your personal data.</li>
              <li><strong>Right to Data Portability:</strong> You have the right to request that we transfer the data that we have collected to another organization, or directly to you.</li>
            </ul>
          </div>

          <div>
            <h2 className="text-2xl font-bold text-gray-900">How We Protect Your Data</h2>
            <p className="mt-4">
              We implement appropriate technical and organizational measures to ensure a level of security appropriate to the risk, including:
            </p>
            <ul className="list-disc pl-5 space-y-2 mt-2">
              <li>Encryption of personal data during transmission</li>
              <li>Regular testing and evaluation of security measures</li>
              <li>Restricted access to personal data on a need-to-know basis</li>
              <li>Regular staff training on data protection</li>
            </ul>
          </div>

          <div>
            <h2 className="text-2xl font-bold text-gray-900">Data Processing Agreements</h2>
            <p className="mt-4">
              We have Data Processing Agreements (DPAs) in place with all third-party service providers that process personal data on our behalf. 
              These agreements ensure that your data is handled in compliance with GDPR requirements.
            </p>
          </div>

          <div>
            <h2 className="text-2xl font-bold text-gray-900">Data Protection Officer</h2>
            <p className="mt-4">
              We have appointed a Data Protection Officer (DPO) who is responsible for overseeing questions about this privacy notice. 
              If you have any questions about this privacy notice or our privacy practices, please contact our DPO at:
            </p>
            <div className="mt-4 p-4 bg-gray-50 rounded-md">
              <p className="font-medium">Data Protection Officer</p>
              <p>Email: <a href="mailto:dpo@freefileconverterz.com" className="text-indigo-600 hover:text-indigo-500">dpo@freefileconverterz.com</a></p>
              <p>Address: 1234 Data Protection Street, San Francisco, CA 94107, USA</p>
            </div>
          </div>

          <div className="bg-indigo-50 p-6 rounded-lg">
            <h2 className="text-2xl font-bold text-gray-900">Contact Us</h2>
            <p className="mt-4">
              If you would like to exercise any of these rights or have any questions about our GDPR compliance, please contact us at:
            </p>
            <p className="mt-2">
              Email: <a href="mailto:privacy@freefileconverterz.com" className="text-indigo-600 hover:text-indigo-500">privacy@freefileconverterz.com</a>
            </p>
            <p className="mt-2">
              For more information about your data protection rights, please visit the 
              <a href="https://ec.europa.eu/info/law/law-topic/data-protection_en" target="_blank" rel="noopener noreferrer" className="text-indigo-600 hover:text-indigo-500">
                {' '}European Commission's Data Protection page
              </a>.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default GDPR;
