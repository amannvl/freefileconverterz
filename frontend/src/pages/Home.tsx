import { Link } from 'react-router-dom';
import { 
  ArrowRightIcon, 
  LockClosedIcon, 
  ArrowsRightLeftIcon, 
  BoltIcon, 
  ShieldCheckIcon,
  DocumentTextIcon,
  DocumentArrowDownIcon,
  TableCellsIcon,
  PhotoIcon,
  FilmIcon,
  MusicalNoteIcon,
  BookOpenIcon,
  ArchiveBoxIcon
} from '@heroicons/react/24/outline';
import { motion } from 'framer-motion';

interface ConversionType {
  name: string;
  description: string;
  href: string;
  icon: JSX.Element;
  value: string;
  gradient: string;
}

const conversionTypes: ConversionType[] = [
  {
    name: 'PDF Tools',
    description: 'Convert, compress, merge, split and edit PDF files',
    href: '/convert/pdf',
    value: 'pdf',
    icon: <DocumentTextIcon className="h-6 w-6" />,
    gradient: 'from-blue-500 to-blue-600'
  },
  {
    name: 'Word',
    description: 'Convert Word documents to and from other formats',
    href: '/convert/word',
    value: 'word',
    icon: <DocumentArrowDownIcon className="h-6 w-6" />,
    gradient: 'from-indigo-500 to-indigo-600'
  },
  {
    name: 'Excel',
    description: 'Convert Excel spreadsheets to and from other formats',
    href: '/convert/excel',
    value: 'excel',
    icon: <TableCellsIcon className="h-6 w-6" />,
    gradient: 'from-green-500 to-green-600'
  },
  {
    name: 'Images',
    description: 'Convert between JPG, PNG, GIF, WebP and more',
    href: '/convert/image',
    value: 'image',
    icon: <PhotoIcon className="h-6 w-6" />,
    gradient: 'from-purple-500 to-purple-600'
  },
  {
    name: 'Video',
    description: 'Convert video files between different formats',
    href: '/convert/video',
    value: 'video',
    icon: <FilmIcon className="h-6 w-6" />,
    gradient: 'from-pink-500 to-pink-600'
  },
  {
    name: 'Audio',
    description: 'Convert audio files between different formats',
    href: '/convert/audio',
    value: 'audio',
    icon: <MusicalNoteIcon className="h-6 w-6" />,
    gradient: 'from-yellow-500 to-yellow-600'
  },
  {
    name: 'eBook',
    description: 'Convert between eBook formats like EPUB, MOBI, and more',
    href: '/convert/ebook',
    value: 'ebook',
    icon: <BookOpenIcon className="h-6 w-6" />,
    gradient: 'from-indigo-500 to-indigo-600'
  },
  {
    name: 'Archive',
    description: 'Create and extract ZIP, RAR, and other archive files',
    href: '/convert/archive',
    value: 'archive',
    icon: <ArchiveBoxIcon className="h-6 w-6" />,
    gradient: 'from-amber-500 to-amber-600'
  },
];

const features = [
  {
    name: 'Lightning Fast',
    description: 'Convert files in seconds with our high-speed processing',
    icon: BoltIcon,
  },
  {
    name: 'Secure & Private',
    description: 'Your files are automatically deleted after conversion',
    icon: ShieldCheckIcon,
  },
  {
    name: 'No Watermarks',
    description: 'Get clean, professional results without any branding',
    icon: LockClosedIcon,
  },
  {
    name: 'Easy to Use',
    description: 'Simple interface that works on any device',
    icon: ArrowsRightLeftIcon,
  },
];

// Animation variants
const fadeInUp = {
  hidden: { opacity: 0, y: 20 },
  visible: { 
    opacity: 1, 
    y: 0,
    transition: {
      duration: 0.6,
      ease: [0.6, -0.05, 0.01, 0.99]
    }
  }
};

const staggerContainer = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1
    }
  }
};

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 via-white to-gray-100 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900">
      {/* Hero Section */}
      <div className="relative overflow-hidden">
        {/* Animated background */}
        <div className="absolute inset-0">
          <div className="absolute inset-0 bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-indigo-100 to-transparent opacity-40 dark:from-indigo-900/20 dark:to-transparent"></div>
          <div className="absolute inset-0 bg-grid-gray-200 dark:bg-grid-gray-700/[0.2] bg-[length:40px_40px]"></div>
        </div>
        
        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 md:py-32">
          <div className="text-center">
            <motion.div
              initial="hidden"
              animate="visible"
              variants={staggerContainer}
            >
              <motion.div variants={fadeInUp} className="inline-flex items-center px-4 py-1.5 rounded-full text-sm font-medium bg-indigo-100 text-indigo-800 dark:bg-indigo-900/30 dark:text-indigo-200 mb-6">
                <span className="flex h-2 w-2 relative mr-2">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-indigo-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-2 w-2 bg-indigo-500"></span>
                </span>
                No registration required
              </motion.div>
              
              <motion.h1 
                className="text-4xl tracking-tight font-extrabold text-gray-900 dark:text-white sm:text-5xl md:text-6xl"
                variants={fadeInUp}
              >
                <span className="block">Convert Any File</span>
                <span className="block bg-clip-text text-transparent bg-gradient-to-r from-indigo-600 to-purple-600">
                  In Seconds
                </span>
              </motion.h1>
              
              <motion.p 
                className="mt-6 max-w-2xl mx-auto text-xl text-gray-600 dark:text-gray-300"
                variants={fadeInUp}
              >
                Transform your documents, images, videos, and more with our lightning-fast, secure, and free online file converter.
              </motion.p>
              
              <motion.div 
                className="mt-10 flex flex-col sm:flex-row gap-4 justify-center items-center"
                variants={fadeInUp}
              >
                <a
                  href="#formats-section"
                  className="group relative inline-flex items-center px-8 py-4 overflow-hidden text-lg font-bold text-white rounded-full bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200 transform hover:-translate-y-0.5 hover:shadow-xl"
                  onClick={(e) => {
                    e.preventDefault();
                    document.getElementById('formats-section')?.scrollIntoView({ behavior: 'smooth' });
                  }}
                >
                  <span className="absolute right-0 -mr-1 transform group-hover:translate-x-1 transition-transform duration-200">
                    <ArrowRightIcon className="h-5 w-5" />
                  </span>
                  <span className="relative mr-4">Start Converting Now</span>
                </a>
                
                <a
                  href="#features"
                  className="inline-flex items-center px-8 py-4 text-lg font-medium text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-800 rounded-full border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200"
                >
                  Learn More
                </a>
              </motion.div>
              
              <motion.div 
                className="mt-8 flex flex-wrap justify-center gap-4"
                variants={fadeInUp}
              >
                {['PDF', 'Word', 'Excel', 'Image', 'Video', 'Audio'].map((type) => (
                  <span key={type} className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-indigo-50 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-200">
                    {type}
                  </span>
                ))}
              </motion.div>
            </motion.div>
          </div>
        </div>
      </div>

      {/* Features Section */}
      <div id="features" className="py-16 bg-white dark:bg-gray-900">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="lg:text-center">
            <h2 className="text-base text-primary-600 dark:text-primary-400 font-semibold tracking-wide uppercase">Features</h2>
            <p className="mt-2 text-3xl leading-8 font-extrabold tracking-tight text-gray-900 dark:text-white sm:text-4xl">
              A better way to convert files
            </p>
            <p className="mt-4 max-w-2xl text-xl text-gray-500 dark:text-gray-300 lg:mx-auto">
              Everything you need to convert files quickly and securely
            </p>
          </div>

          <div className="mt-16">
            <div className="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-4">
              {features.map((feature) => (
                <motion.div 
                  key={feature.name}
                  className="pt-6"
                  whileHover={{ y: -5, transition: { duration: 0.2 } }}
                >
                  <div className="flow-root bg-gray-50 dark:bg-gray-800 rounded-lg px-6 pb-8 h-full">
                    <div className="-mt-6">
                      <div className="inline-flex items-center justify-center p-3 bg-primary-500 rounded-md shadow-lg">
                        <feature.icon className="h-6 w-6 text-white" aria-hidden="true" />
                      </div>
                      <h3 className="mt-4 text-lg font-medium text-gray-900 dark:text-white tracking-tight">
                        {feature.name}
                      </h3>
                      <p className="mt-2 text-base text-gray-500 dark:text-gray-400">
                        {feature.description}
                      </p>
                    </div>
                  </div>
                </motion.div>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Conversion Types */}
      <div id="formats-section" className="py-16 bg-gradient-to-b from-white to-gray-50 dark:from-gray-900 dark:to-gray-800">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <motion.div 
            className="text-center"
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
            variants={staggerContainer}
          >
            <motion.div variants={fadeInUp} className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-indigo-100 text-indigo-800 dark:bg-indigo-900/30 dark:text-indigo-200 mb-4">
              <span className="flex h-2 w-2 relative mr-2">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-indigo-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-2 w-2 bg-indigo-500"></span>
              </span>
              Supported Formats
            </motion.div>
            <motion.h2 
              className="text-3xl font-bold tracking-tight text-gray-900 dark:text-white sm:text-4xl"
              variants={fadeInUp}
            >
              Convert Between 100+ Formats
            </motion.h2>
            <motion.p 
              className="mt-6 max-w-2xl mx-auto text-xl text-gray-600 dark:text-gray-300"
              variants={fadeInUp}
            >
              High-quality conversions for all your document, image, video, and audio needs
            </motion.p>
          </motion.div>

          <motion.div 
            className="mt-16 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true, margin: "-100px" }}
            variants={staggerContainer}
          >
            {conversionTypes.map((type) => (
              <motion.div
                key={type.name}
                variants={fadeInUp}
                className="group relative flex flex-col rounded-2xl bg-white dark:bg-gray-800 shadow-lg overflow-hidden transition-all duration-300 hover:shadow-xl hover:-translate-y-1 border border-gray-100 dark:border-gray-700"
                whileHover={{ 
                  y: -5,
                  boxShadow: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)'
                }}
              >
                <div 
                  className={`h-2 ${
                    type.value === 'pdf' ? 'bg-red-500' :
                    type.value === 'word' ? 'bg-blue-500' :
                    type.value === 'excel' ? 'bg-green-500' :
                    type.value === 'image' ? 'bg-purple-500' :
                    type.value === 'video' ? 'bg-pink-500' :
                    type.value === 'audio' ? 'bg-yellow-500' :
                    type.value === 'ebook' ? 'bg-indigo-500' :
                    type.value === 'archive' ? 'bg-amber-500' : 'bg-gray-500'
                  }`}
                ></div>
                <div className="flex-1 p-6 flex flex-col">
                  <div className="flex items-center mb-4">
                    <div className={`p-3 rounded-lg ${
                      type.value === 'pdf' ? 'bg-red-100 text-red-600 dark:bg-red-900/30 dark:text-red-400' :
                      type.value === 'word' ? 'bg-blue-100 text-blue-600 dark:bg-blue-900/30 dark:text-blue-400' :
                      type.value === 'excel' ? 'bg-green-100 text-green-600 dark:bg-green-900/30 dark:text-green-400' :
                      type.value === 'image' ? 'bg-purple-100 text-purple-600 dark:bg-purple-900/30 dark:text-purple-400' :
                      type.value === 'video' ? 'bg-pink-100 text-pink-600 dark:bg-pink-900/30 dark:text-pink-400' :
                      type.value === 'audio' ? 'bg-yellow-100 text-yellow-600 dark:bg-yellow-900/30 dark:text-yellow-400' :
                      type.value === 'ebook' ? 'bg-indigo-100 text-indigo-600 dark:bg-indigo-900/30 dark:text-indigo-400' :
                      type.value === 'archive' ? 'bg-amber-100 text-amber-600 dark:bg-amber-900/30 dark:text-amber-400' :
                      'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300'
                    }`}>
                      <span className="text-2xl">{type.icon}</span>
                    </div>
                    <h3 className="ml-4 text-lg font-semibold text-gray-900 dark:text-white">
                      {type.name}
                    </h3>
                  </div>
                  <p className="text-gray-600 dark:text-gray-300 flex-1">
                    {type.description}
                  </p>
                  <div className="mt-6 pt-4 border-t border-gray-100 dark:border-gray-700">
                    <Link
                      to={`/convert/${type.value}`}
                      className="group inline-flex items-center text-base font-medium text-indigo-600 hover:text-indigo-500 dark:text-indigo-400 dark:hover:text-indigo-300 transition-colors duration-200"
                    >
                      Convert now
                      <svg 
                        className="ml-2 h-5 w-5 transition-transform duration-200 group-hover:translate-x-1" 
                        fill="none" 
                        viewBox="0 0 24 24" 
                        stroke="currentColor"
                      >
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14 5l7 7m0 0l-7 7m7-7H3" />
                      </svg>
                    </Link>
                  </div>
                </div>
              </motion.div>
            ))}
          </motion.div>
        </div>
      </div>

      {/* CTA Section */}
      <div className="bg-white dark:bg-gray-900">
        <div className="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:py-16 lg:px-8 lg:flex lg:items-center lg:justify-between">
          <h2 className="text-3xl font-extrabold tracking-tight text-gray-900 dark:text-white sm:text-4xl">
            <span className="block">Ready to get started?</span>
            <span className="block text-primary-600 dark:text-primary-400">Convert your files now.</span>
          </h2>
          <div className="mt-8 flex lg:mt-0 lg:flex-shrink-0">
            <div className="inline-flex rounded-md shadow">
              <Link
                to="/convert"
                className="inline-flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-full text-white bg-gradient-to-r from-primary-600 to-primary-500 hover:from-primary-700 hover:to-primary-600 md:py-4 md:text-lg md:px-10 transition-all duration-200 transform hover:-translate-y-1"
              >
                Get Started
                <ArrowRightIcon className="ml-2 -mr-1 h-5 w-5" aria-hidden="true" />
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
