function CookiePolicy() {
  return (
    <div className="max-w-4xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
      <div className="text-center mb-12">
        <h1 className="text-4xl font-extrabold text-gray-900 sm:text-5xl sm:tracking-tight lg:text-6xl">
          Cookie Policy
        </h1>
        <p className="mt-3 max-w-2xl mx-auto text-xl text-gray-500 sm:mt-4">
          Last updated: July 2, 2025
        </p>
      </div>

      <div className="prose prose-indigo max-w-none">
        <div className="bg-white shadow overflow-hidden sm:rounded-lg mb-12">
          <div className="px-4 py-5 sm:px-6">
            <h2 className="text-2xl leading-6 font-medium text-gray-900">
              Our Use of Cookies
            </h2>
            <p className="mt-1 text-sm text-gray-500">
              This policy explains how we use cookies and similar technologies on our website.
            </p>
          </div>
        </div>

        <div className="space-y-8">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">What are Cookies?</h2>
            <p className="mt-4">
              Cookies are small text files that are placed on your computer or mobile device when you visit a website. 
              They are widely used to make websites work more efficiently and to provide information to the site owners.
            </p>
          </div>

          <div>
            <h2 className="text-2xl font-bold text-gray-900">How We Use Cookies</h2>
            <p className="mt-4">
              We use cookies for several purposes:
            </p>
            <ul className="list-disc pl-5 space-y-2 mt-2">
              <li>
                <strong>Essential Cookies:</strong> These are necessary for the website to function and cannot be switched off.
              </li>
              <li>
                <strong>Performance Cookies:</strong> These allow us to count visits and traffic sources so we can measure and improve the performance of our site.
              </li>
              <li>
                <strong>Functional Cookies:</strong> These enable the website to provide enhanced functionality and personalization.
              </li>
              <li>
                <strong>Analytics Cookies:</strong> These help us understand how visitors interact with our website.
              </li>
            </ul>
          </div>

          <div>
            <h2 className="text-2xl font-bold text-gray-900">Managing Cookies</h2>
            <p className="mt-4">
              You can control and/or delete cookies as you wish. You can delete all cookies that are already on your computer 
              and you can set most browsers to prevent them from being placed. If you do this, however, you may have to 
              manually adjust some preferences every time you visit a site and some services and functionalities may not work.
            </p>
          </div>

          <div className="bg-indigo-50 p-6 rounded-lg">
            <h2 className="text-2xl font-bold text-gray-900">Need More Information?</h2>
            <p className="mt-4">
              If you have any questions about our use of cookies, please contact us at{' '}
              <a href="mailto:privacy@freefileconverterz.com" className="text-indigo-600 hover:text-indigo-500">
                privacy@freefileconverterz.com
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default CookiePolicy;
