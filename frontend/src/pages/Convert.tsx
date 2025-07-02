import { useState, useCallback, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { 
  XMarkIcon, 
  ArrowPathIcon,
  CheckCircleIcon,
  ExclamationCircleIcon
} from '@heroicons/react/24/outline';
import FileUploader from '../components/FileUploader';
import { motion } from 'framer-motion';
import type { FileWithStatus, SupportedFormats } from '../types';

type ConversionType = 'document' | 'image' | 'video' | 'audio';

const supportedFormats: SupportedFormats = {
  image: {
    input: ['jpg', 'jpeg', 'png', 'webp', 'gif', 'bmp', 'tiff', 'svg'],
    output: ['jpg', 'png', 'webp', 'gif', 'bmp']
  },
  document: {
    input: ['pdf', 'docx', 'doc', 'txt', 'rtf', 'odt'],
    output: ['pdf', 'docx', 'txt', 'rtf', 'odt']
  },
  video: {
    input: ['mp4', 'avi', 'mov', 'wmv', 'flv', 'mkv', 'webm'],
    output: ['mp4', 'avi', 'mov', 'webm']
  },
  audio: {
    input: ['mp3', 'wav', 'aac', 'ogg', 'wma', 'flac'],
    output: ['mp3', 'wav', 'aac', 'ogg', 'flac']
  }
};

export default function Convert() {
  const { conversionType } = useParams<{ conversionType?: string }>();
  const [files, setFiles] = useState<FileWithStatus[]>([]);
  const [conversionTypeState, setConversionTypeState] = useState<ConversionType>('document');
  const [outputFormat, setOutputFormat] = useState<string>('');

  // Get supported formats for the current conversion type
  const getSupportedFormats = useCallback((type: ConversionType): string[] => {
    return supportedFormats[type]?.output || [];
  }, []);

  // Get the current conversion type from the URL
  const getConversionTypeFromUrl = useCallback((pathname: string): ConversionType => {
    const parts = pathname.split('/');
    const type = parts[parts.length - 1];
    return Object.keys(supportedFormats).includes(type) 
      ? type as ConversionType 
      : 'document';
  }, []);

  // Set conversion type from URL parameter
  useEffect(() => {
    const type = getConversionTypeFromUrl(window.location.pathname);
    setConversionTypeState(type);
    
    // Set initial output format based on conversion type
    const formats = getSupportedFormats(type);
    if (formats.length > 0) {
      setOutputFormat(formats[0]);
    }
  }, [conversionType, getConversionTypeFromUrl, getSupportedFormats]);

  // Handle file upload
  const uploadFile = useCallback(async (fileWithStatus: FileWithStatus) => {
    try {
      // Simulate upload progress
      for (let i = 0; i <= 100; i += 10) {
        await new Promise(resolve => setTimeout(resolve, 100));
        setFiles(prevFiles => 
          prevFiles.map(f => 
            f.id === fileWithStatus.id 
              ? { ...f, progress: i, status: 'uploading' as const } 
              : f
          )
        );
      }

      // Simulate conversion
      setFiles(prevFiles => 
        prevFiles.map(f => 
          f.id === fileWithStatus.id 
            ? { ...f, status: 'converting' as const } 
            : f
        )
      );

      // Simulate completion after a delay
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      setFiles(prevFiles => 
        prevFiles.map(f => 
          f.id === fileWithStatus.id 
            ? { 
                ...f, 
                status: 'completed' as const, 
                progress: 100,
                downloadUrl: URL.createObjectURL(new Blob(['Simulated file content']))
              } 
            : f
        )
      );
    } catch (error) {
      setFiles(prevFiles => 
        prevFiles.map(f => 
          f.id === fileWithStatus.id 
            ? { 
                ...f, 
                status: 'error' as const, 
                error: error instanceof Error ? error.message : 'Upload failed'
              } 
            : f
        )
      );
    }
  }, []);

  // Handle file selection
  const handleFiles = useCallback((newFiles: File[]) => {
    const filesWithStatus: FileWithStatus[] = newFiles.map(file => ({
      id: Math.random().toString(36).substring(2, 9),
      file,
      status: 'uploading' as const,
      progress: 0
    }));

    setFiles(prevFiles => [...prevFiles, ...filesWithStatus]);
    
    // Start upload for each file
    filesWithStatus.forEach(uploadFile);
  }, [uploadFile]);

  // Handle file drop
  const handleDrop = useCallback((droppedFiles: File[]) => {
    handleFiles(droppedFiles);
  }, [handleFiles]);

  // Remove file from the list
  const removeFile = useCallback((fileId: string) => {
    setFiles(prevFiles => prevFiles.filter(file => file.id !== fileId));
  }, []);

  // Handle format selection
  const handleFormatSelect = useCallback((format: string) => {
    setOutputFormat(format);
  }, []);

  // Handle download
  const handleDownload = useCallback(async (file: FileWithStatus) => {
    if (!file.downloadUrl) return;

    try {
      const a = document.createElement('a');
      a.href = file.downloadUrl;
      
      // Get the file extension from the selected format or original file
      const fileExt = outputFormat || file.file.name.split('.').pop() || '';
      const fileName = file.file.name.split('.').slice(0, -1).join('.');
      
      a.download = fileExt ? `${fileName}.${fileExt}` : fileName;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
    } catch (error) {
      console.error('Error downloading file:', error);
    }
  }, [outputFormat]);

  // Get the current output formats based on conversion type
  const outputFormats = getSupportedFormats(conversionTypeState);

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-extrabold text-gray-900 sm:text-4xl">
            {conversionTypeState ? `Convert to ${conversionTypeState.charAt(0).toUpperCase() + conversionTypeState.slice(1)}` : 'File Converter'}
          </h1>
          <p className="mt-3 text-xl text-gray-500">
            Convert your files to {conversionTypeState} format quickly and easily
          </p>
        </div>

        <div className="bg-white shadow rounded-lg p-6 mb-8">
          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Output Format
            </label>
            <select
              className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
              value={outputFormat}
              onChange={(e) => handleFormatSelect(e.target.value)}
            >
              {outputFormats.map((format) => (
                <option key={format} value={format}>
                  {format.toUpperCase()}
                </option>
              ))}
            </select>
          </div>

          <FileUploader 
            onDrop={handleDrop}
            onFilesSelected={handleFiles}
            accept={supportedFormats[conversionTypeState]?.input.map(ext => `.${ext}`).join(',')}
            acceptedFormats={supportedFormats[conversionTypeState]?.input || []}
          />
        </div>

        {/* File list */}
        {files.length > 0 && (
          <div className="bg-white shadow overflow-hidden sm:rounded-md">
            <ul className="divide-y divide-gray-200">
              {files.map((file) => (
                <motion.li
                  key={file.id}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, x: -100 }}
                  transition={{ duration: 0.3 }}
                  className="px-4 py-4 sm:px-6"
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center min-w-0">
                      <div className="flex-shrink-0 h-10 w-10 flex items-center justify-center bg-indigo-100 rounded-full">
                        {file.status === 'completed' ? (
                          <CheckCircleIcon className="h-6 w-6 text-green-500" />
                        ) : file.status === 'error' ? (
                          <ExclamationCircleIcon className="h-6 w-6 text-red-500" />
                        ) : (
                          <ArrowPathIcon className="h-5 w-5 text-indigo-500 animate-spin" />
                        )}
                      </div>
                      <div className="ml-4 min-w-0 flex-1">
                        <p className="text-sm font-medium text-indigo-600 truncate">
                          {file.file.name}
                        </p>
                        <p className="text-sm text-gray-500">
                          {file.status === 'uploading' && 'Uploading...'}
                          {file.status === 'converting' && 'Converting...'}
                          {file.status === 'completed' && 'Conversion complete'}
                          {file.status === 'error' && file.error}
                        </p>
                        {file.status !== 'completed' && file.status !== 'error' && (
                          <div className="mt-1 w-full bg-gray-200 rounded-full h-2">
                            <div
                              className="bg-indigo-600 h-2 rounded-full"
                              style={{ width: `${file.progress}%` }}
                            />
                          </div>
                        )}
                      </div>
                    </div>
                    <div className="ml-4 flex-shrink-0 flex space-x-2">
                      {file.status === 'completed' && (
                        <button
                          type="button"
                          onClick={() => handleDownload(file)}
                          className="inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                        >
                          Download
                        </button>
                      )}
                      <button
                        type="button"
                        onClick={() => removeFile(file.id)}
                        className="inline-flex items-center p-1 border border-transparent rounded-full text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                      >
                        <XMarkIcon className="h-5 w-5" />
                      </button>
                    </div>
                  </div>
                </motion.li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}
