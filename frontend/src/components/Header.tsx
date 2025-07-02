import { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';

const conversionTypes = [
  { name: 'PDF Tools', value: 'pdf' },
  { name: 'Word', value: 'word' },
  { name: 'Images', value: 'image' },
  { name: 'Video', value: 'video' },
  { name: 'Audio', value: 'audio' },
  { name: 'Archive', value: 'archive' },
];

export default function Header() {
  const location = useLocation();
  const currentPath = location.pathname;
  const isConvertPage = currentPath.startsWith('/convert/');
  const currentType = isConvertPage ? currentPath.split('/')[2] : '';
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  return (
    <header className="bg-white shadow-sm sticky top-0 z-50">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          <div className="flex items-center">
            <Link to="/" className="flex items-center">
              <span className="text-xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
                FreeFileConverterZ
              </span>
            </Link>
            
            {/* Desktop Navigation */}
            <div className="hidden md:ml-10 md:flex md:space-x-1">
              <Link
                to="/"
                className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                  currentPath === '/'
                    ? 'bg-indigo-50 text-indigo-700'
                    : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                }`}
              >
                Home
              </Link>
              
              {conversionTypes.map((type) => (
                <Link
                  key={type.value}
                  to={`/convert/${type.value}`}
                  className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                    currentType === type.value
                      ? 'bg-indigo-50 text-indigo-700'
                      : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                  }`}
                >
                  {type.name}
                </Link>
              ))}
              
              <Link
                to="/about"
                className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                  currentPath === '/about'
                    ? 'bg-indigo-50 text-indigo-700'
                    : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                }`}
              >
                About
              </Link>
              <Link
                to="/help"
                className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                  currentPath === '/help'
                    ? 'bg-indigo-50 text-indigo-700'
                    : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                }`}
              >
                Help
              </Link>
              <Link
                to="/contact"
                className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                  currentPath === '/contact'
                    ? 'bg-indigo-50 text-indigo-700'
                    : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                }`}
              >
                Contact
              </Link>
            </div>
          </div>

          {/* Mobile menu button */}
          <div className="-mr-2 flex items-center md:hidden">
            <button
              type="button"
              className="inline-flex items-center justify-center rounded-md p-2 text-gray-400 hover:bg-gray-100 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500"
              onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            >
              <span className="sr-only">Open main menu</span>
              <svg
                className="block h-6 w-6"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d={isMobileMenuOpen ? "M6 18L18 6M6 6l12 12" : "M4 6h16M4 12h16M4 18h16"}
                />
              </svg>
            </button>
          </div>
        </div>
      </div>

      {/* Mobile menu */}
      <div className={`md:hidden transition-all duration-300 ease-in-out ${isMobileMenuOpen ? 'max-h-screen' : 'max-h-0 overflow-hidden'}`}>
        <div className="space-y-1 pb-3 pt-2 bg-white border-t border-gray-200">
          <Link
            to="/"
            className="block px-4 py-3 text-base font-medium text-gray-700 hover:bg-gray-50 hover:text-gray-900"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            Home
          </Link>
          
          <div className="px-4 py-2 font-medium text-gray-500 text-sm uppercase tracking-wider">
            Convert To:
          </div>
          
          <div className="grid grid-cols-2 gap-1 px-2">
            {conversionTypes.map((type) => (
              <Link
                key={type.value}
                to={`/convert/${type.value}`}
                className={`block px-4 py-3 rounded-md text-base font-medium ${
                  currentType === type.value
                    ? 'bg-indigo-50 text-indigo-700'
                    : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
                }`}
                onClick={() => setIsMobileMenuOpen(false)}
              >
                {type.name}
              </Link>
            ))}
          </div>
          
          <div className="border-t border-gray-200 my-2"></div>
          
          <Link
            to="/about"
            className="block px-4 py-3 text-base font-medium text-gray-700 hover:bg-gray-50 hover:text-gray-900"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            About
          </Link>
          <Link
            to="/help"
            className="block px-4 py-3 text-base font-medium text-gray-700 hover:bg-gray-50 hover:text-gray-900"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            Help
          </Link>
          <Link
            to="/contact"
            className="block px-4 py-3 text-base font-medium text-gray-700 hover:bg-gray-50 hover:text-gray-900"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            Contact
          </Link>
          
          <div className="border-t border-gray-200 my-2"></div>
          
          <div className="px-4 py-2 font-medium text-gray-500 text-sm uppercase tracking-wider">
            Legal
          </div>
          
          <Link
            to="/terms"
            className="block px-4 py-2 pl-8 text-base font-medium text-gray-700 hover:bg-gray-50 hover:text-gray-900"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            Terms of Service
          </Link>
          <Link
            to="/privacy"
            className="block px-4 py-2 pl-8 text-base font-medium text-gray-700 hover:bg-gray-50 hover:text-gray-900"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            Privacy Policy
          </Link>
          <Link
            to="/cookie-policy"
            className="block px-4 py-2 pl-8 text-base font-medium text-gray-700 hover:bg-gray-50 hover:text-gray-900"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            Cookie Policy
          </Link>
          <Link
            to="/gdpr"
            className="block px-4 py-2 pl-8 text-base font-medium text-gray-700 hover:bg-gray-50 hover:text-gray-900"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            GDPR
          </Link>
        </div>
      </div>
    </header>
  );
}
