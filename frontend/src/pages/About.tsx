import { Link } from 'react-router-dom';

function About() {
  return (
    <div className="max-w-4xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
      <div className="text-center mb-12">
        <h1 className="text-4xl font-extrabold text-gray-900 sm:text-5xl sm:tracking-tight lg:text-6xl">
          About FreeFileConverterZ
        </h1>
        <p className="mt-3 max-w-2xl mx-auto text-xl text-gray-500 sm:mt-4">
          Your trusted solution for all file conversion needs
        </p>
      </div>

      <div className="prose prose-indigo prose-lg text-gray-500 mx-auto">
        <div className="bg-white shadow overflow-hidden sm:rounded-lg mb-12">
          <div className="px-4 py-5 sm:px-6">
            <h2 className="text-2xl leading-6 font-medium text-gray-900">Our Mission</h2>
            <p className="mt-1 max-w-2xl text-sm text-gray-500">
              Making file conversion simple, fast, and secure for everyone.
            </p>
          </div>
          <div className="border-t border-gray-200 px-4 py-5 sm:p-0">
            <dl className="sm:divide-y sm:divide-gray-200">
              <div className="py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                <dt className="text-sm font-medium text-gray-500">Founded</dt>
                <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">2023</dd>
              </div>

              <div className="py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                <dt className="text-sm font-medium text-gray-500">Supported Formats</dt>
                <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                  <div className="grid grid-cols-2 gap-2">
                    <div>
                      <h4 className="font-medium">Images:</h4>
                      <p>JPG, PNG, WebP, GIF, BMP, TIFF, SVG, HEIC, HEIF</p>
                    </div>
                    <div>
                      <h4 className="font-medium">Documents:</h4>
                      <p>PDF, DOCX, DOC, TXT, RTF, ODT</p>
                    </div>
                    <div>
                      <h4 className="font-medium">Videos:</h4>
                      <p>MP4, AVI, MOV, WMV, FLV, MKV, WebM, 3GP</p>
                    </div>
                    <div>
                      <h4 className="font-medium">Audio:</h4>
                      <p>MP3, WAV, AAC, OGG, WMA, FLAC</p>
                    </div>
                  </div>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <div className="space-y-8">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">Our Story</h2>
            <p className="mt-4">
              FreeFileConverterZ was born out of a simple idea: file conversion should be easy, fast, and accessible to everyone. 
              Frustrated with existing solutions that were either too complex or loaded with ads, our team set out to create 
              a better experience for users who need to convert files quickly and reliably.
            </p>
          </div>

          <div>
            <h2 className="text-2xl font-bold text-gray-900">Our Values</h2>
            <div className="mt-6 grid gap-8 lg:grid-cols-3">
              {[
                {
                  title: 'Simplicity',
                  description: 'We believe in creating intuitive experiences that make complex tasks simple.'
                },
                {
                  title: 'Privacy',
                  description: 'Your files are your business. We automatically delete them after processing and never store them longer than necessary.'
                },
                {
                  title: 'Quality',
                  'description': 'We strive to provide the highest quality conversions with support for a wide range of file formats.'
                }
              ].map((value, index) => (
                <div key={index} className="bg-white p-6 rounded-lg shadow">
                  <h3 className="text-lg font-medium text-gray-900">{value.title}</h3>
                  <p className="mt-2 text-gray-600">{value.description}</p>
                </div>
              ))}
            </div>
          </div>

          <div className="bg-indigo-50 p-6 rounded-lg">
            <h2 className="text-2xl font-bold text-gray-900">Get In Touch</h2>
            <p className="mt-4">
              Have questions or feedback? We'd love to hear from you! Reach out to us at{' '}
              <a href="mailto:support@freefileconverterz.com" className="text-indigo-600 hover:text-indigo-500">
                support@freefileconverterz.com
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default About;
