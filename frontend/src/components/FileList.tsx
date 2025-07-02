import { ArrowPathIcon, CheckCircleIcon, XCircleIcon, XMarkIcon } from '@heroicons/react/24/outline';
import type { FileWithStatus } from '../types';

interface FileListProps {
  files: FileWithStatus[];
  onDownload: (file: FileWithStatus) => void;
  onRemove: (fileId: string) => void;
}

export default function FileList({ files, onDownload, onRemove }: FileListProps) {
  if (files.length === 0) {
    return null;
  }

  const getFileIcon = (fileName: string) => {
    const ext = fileName.split('.').pop()?.toLowerCase() || '';
    
    if (['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp', 'svg'].includes(ext)) {
      return '🖼️';
    } else if (['pdf', 'docx', 'doc', 'txt', 'rtf', 'odt'].includes(ext)) {
      return '📄';
    } else if (['mp4', 'avi', 'mov', 'wmv', 'mkv', 'webm'].includes(ext)) {
      return '🎥';
    } else if (['mp3', 'wav', 'aac', 'ogg', 'flac'].includes(ext)) {
      return '🎵';
    } else if (['zip', 'rar', '7z', 'tar', 'gz'].includes(ext)) {
      return '🗄️';
    } else {
      return '📁';
    }
  };

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  return (
    <div className="space-y-4">
      <h2 className="text-lg font-medium text-gray-900 dark:text-white">
        Files to convert ({files.length})
      </h2>
      
      <div className="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-md">
        <ul className="divide-y divide-gray-200 dark:divide-gray-700">
          {files.map((fileWithStatus) => (
            <li key={fileWithStatus.id}>
              <div className="px-4 py-4 flex items-center justify-between">
                <div className="flex items-center">
                  <div className="flex-shrink-0 h-10 w-10 rounded-md bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
                    <span className="text-xl">
                      {getFileIcon(fileWithStatus.file.name)}
                    </span>
                  </div>
                  <div className="ml-4">
                    <div className="text-sm font-medium text-gray-900 dark:text-white truncate max-w-xs">
                      {fileWithStatus.file.name}
                    </div>
                    <div className="text-xs text-gray-500 dark:text-gray-400">
                      {formatFileSize(fileWithStatus.file.size)}
                    </div>
                  </div>
                </div>
                
                <div className="flex items-center space-x-2">
                  {fileWithStatus.status === 'uploading' && (
                    <div className="flex items-center text-sm text-gray-500 dark:text-gray-400">
                      <ArrowPathIcon className="h-4 w-4 mr-1 animate-spin" />
                      Uploading...
                    </div>
                  )}
                  
                  {fileWithStatus.status === 'converting' && (
                    <div className="flex items-center text-sm text-blue-600 dark:text-blue-400">
                      <ArrowPathIcon className="h-4 w-4 mr-1 animate-spin" />
                      Converting...
                    </div>
                  )}
                  
                  {fileWithStatus.status === 'completed' && (
                    <button
                      type="button"
                      onClick={() => onDownload(fileWithStatus)}
                      className="inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                    >
                      <CheckCircleIcon className="h-4 w-4 mr-1" />
                      Download
                    </button>
                  )}
                  
                  {fileWithStatus.status === 'error' && (
                    <div className="flex items-center text-sm text-red-600 dark:text-red-400">
                      <XCircleIcon className="h-4 w-4 mr-1" />
                      {fileWithStatus.error || 'Error'}
                    </div>
                  )}
                  
                  <button
                    type="button"
                    onClick={() => onRemove(fileWithStatus.id)}
                    className="text-gray-400 hover:text-gray-500 dark:text-gray-500 dark:hover:text-gray-400"
                  >
                    <XMarkIcon className="h-5 w-5" />
                  </button>
                </div>
              </div>
              
              {/* Progress bar */}
              {(fileWithStatus.status === 'uploading' || fileWithStatus.status === 'converting') && (
                <div className="px-4 pb-2">
                  <div className="w-full bg-gray-200 rounded-full h-1.5 dark:bg-gray-700">
                    <div 
                      className="bg-blue-600 h-1.5 rounded-full dark:bg-blue-500" 
                      style={{ width: `${fileWithStatus.progress}%` }}
                    />
                  </div>
                </div>
              )}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
