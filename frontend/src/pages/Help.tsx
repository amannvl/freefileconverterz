import { useState } from 'react';
import { Link } from 'react-router-dom';

const faqs = [
  {
    question: 'How do I convert a file?',
    answer: 'To convert a file, simply select the desired conversion type, upload your file, and click the convert button. Your file will be processed, and you can download the converted file once ready.'
  },
  {
    question: 'What file formats are supported?',
    answer: 'We support a wide range of formats including PDF, DOCX, XLSX, PPTX, JPG, PNG, MP4, MP3, and many more. Check our homepage for a complete list of supported formats.'
  },
  {
    question: 'Is there a file size limit?',
    answer: 'Yes, the maximum file size is 100MB for free users. For larger files, consider compressing them before uploading or contact us for premium options.'
  },
  {
    question: 'How long are my files stored?',
    answer: 'Your files are automatically deleted from our servers after 24 hours. We do not store your files longer than necessary.'
  },
  {
    question: 'Is my data secure?',
    answer: 'Yes, we take security seriously. All file transfers are encrypted with SSL, and we automatically delete files after processing.'
  }
];

function Help() {
  const [activeIndex, setActiveIndex] = useState<number | null>(null);

  const toggleAccordion = (index: number) => {
    setActiveIndex(activeIndex === index ? null : index);
  };

  return (
    <div className="max-w-4xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
      <div className="text-center mb-12">
        <h1 className="text-4xl font-extrabold text-gray-900 sm:text-5xl sm:tracking-tight lg:text-6xl">
          Help Center
        </h1>
        <p className="mt-3 max-w-2xl mx-auto text-xl text-gray-500 sm:mt-4">
          Find answers to common questions or contact our support team
        </p>
      </div>

      <div className="bg-white shadow overflow-hidden sm:rounded-lg mb-12">
        <div className="px-4 py-5 sm:px-6 bg-gray-50">
          <h2 className="text-lg leading-6 font-medium text-gray-900">
            Frequently Asked Questions
          </h2>
        </div>
        <div className="border-t border-gray-200 divide-y divide-gray-200">
          {faqs.map((faq, index) => (
            <div key={index} className="px-4 py-5 sm:px-6">
              <button
                onClick={() => toggleAccordion(index)}
                className="w-full flex justify-between items-start text-left focus:outline-none"
                aria-expanded={activeIndex === index}
              >
                <h3 className="text-lg font-medium text-gray-900">{faq.question}</h3>
                <span className="ml-6 flex items-center">
                  {activeIndex === index ? (
                    <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
                    </svg>
                  ) : (
                    <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                    </svg>
                  )}
                </span>
              </button>
              {activeIndex === index && (
                <div className="mt-4 text-gray-600">
                  {faq.answer}
                </div>
              )}
            </div>
          ))}
        </div>
      </div>

      <div className="bg-white shadow overflow-hidden sm:rounded-lg">
        <div className="px-4 py-5 sm:px-6">
          <h2 className="text-lg leading-6 font-medium text-gray-900">
            Still need help?
          </h2>
          <p className="mt-1 max-w-2xl text-sm text-gray-500">
            Can't find what you're looking for? Our support team is here to help.
          </p>
        </div>
        <div className="border-t border-gray-200 px-4 py-5 sm:px-6">
          <div className="grid grid-cols-1 gap-8 sm:grid-cols-2">
            <div>
              <h3 className="text-lg font-medium text-gray-900">Email Us</h3>
              <p className="mt-2 text-base text-gray-500">
                <a href="mailto:support@freefileconverterz.com" className="text-indigo-600 hover:text-indigo-500">
                  support@freefileconverterz.com
                </a>
              </p>
              <p className="mt-2 text-sm text-gray-500">
                We typically respond within 24 hours.
              </p>
            </div>
            <div>
              <h3 className="text-lg font-medium text-gray-900">Helpful Links</h3>
              <ul className="mt-2 space-y-2">
                <li>
                  <Link to="/about" className="text-indigo-600 hover:text-indigo-500">
                    About Us
                  </Link>
                </li>
                <li>
                  <Link to="/terms" className="text-indigo-600 hover:text-indigo-500">
                    Terms of Service
                  </Link>
                </li>
                <li>
                  <Link to="/privacy" className="text-indigo-600 hover:text-indigo-500">
                    Privacy Policy
                  </Link>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Help;
