import React, { useCallback, useState } from 'react';
import { ArrowUpTrayIcon } from '@heroicons/react/24/outline';

interface FileUploaderProps {
  onFilesSelected: (files: File[]) => void;
  onDrop?: (files: File[]) => void;
  accept?: string; // MIME types or file extensions (e.g., '.jpg,.png' or 'image/*')
  acceptedFormats: string[];
  maxFiles?: number;
  maxSize?: number; // in bytes
  className?: string;
}

const FileUploader: React.FC<FileUploaderProps> = ({
  onFilesSelected, 
  onDrop,
  accept,
  acceptedFormats, 
  maxFiles = 10,
  maxSize = 100 * 1024 * 1024, // 100MB default
  className = ''
}) => {
  const [isDragging, setIsDragging] = useState(false);

  const handleDragOver = useCallback((e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  }, []);

  const validateFiles = useCallback((files: File[]): File[] => {
    const validFiles: File[] = [];
    const invalidFiles: string[] = [];
    
    // Check max files
    if (files.length > maxFiles) {
      alert(`You can only upload up to ${maxFiles} files at once.`);
      return [];
    }
    
    // Check file types and sizes
    files.forEach(file => {
      const fileExt = file.name.split('.').pop()?.toLowerCase() || '';
      
      if (!acceptedFormats.includes(fileExt)) {
        invalidFiles.push(`• ${file.name}: Unsupported file type`);
      } else if (file.size > maxSize) {
        invalidFiles.push(`• ${file.name}: File too large (max ${(maxSize / 1024 / 1024).toFixed(0)}MB)`);
      } else {
        validFiles.push(file);
      }
    });
    
    if (invalidFiles.length > 0) {
      alert(`Some files were not accepted:\n\n${invalidFiles.join('\n')}`);
    }
    
    return validFiles;
  }, [acceptedFormats, maxFiles, maxSize]);
  
  const handleDrop = useCallback((e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    
    const files = Array.from(e.dataTransfer.files);
    const validFiles = validateFiles(files);
    
    if (validFiles.length > 0) {
      if (onDrop) {
        onDrop(validFiles);
      } else {
        onFilesSelected(validFiles);
      }
    }
  }, [onFilesSelected, onDrop, validateFiles]);

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const validFiles = validateFiles(Array.from(e.target.files));
      if (validFiles.length > 0) {
        onFilesSelected(validFiles);
      }
      // Reset the input value to allow selecting the same file again
      e.target.value = '';
    }
  };

  return (
    <div 
      className={`border-2 border-dashed rounded-xl p-8 md:p-12 text-center transition-colors ${
        isDragging 
          ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20' 
          : 'border-gray-300 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-500'
      } ${className}`}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
    >
      <div className="flex flex-col items-center justify-center space-y-4">
        <div className="p-3 rounded-full bg-primary-100 dark:bg-primary-900/30">
          <ArrowUpTrayIcon className="h-8 w-8 text-primary-600 dark:text-primary-400" />
        </div>
        <div className="space-y-1">
          <p className="text-sm text-gray-600 dark:text-gray-300">
            <label 
              htmlFor="file-upload" 
              className="relative cursor-pointer font-medium text-primary-600 hover:text-primary-500 dark:text-primary-400 dark:hover:text-primary-300 focus-within:outline-none focus-within:ring-2 focus-within:ring-offset-2 focus-within:ring-primary-500 rounded-md"
            >
              <span>Upload a file</span>
            </label>{' '}
            or drag and drop
          </p>
          <p className="text-xs text-gray-500 dark:text-gray-400">
            {acceptedFormats.map(f => f.toUpperCase()).join(', ')} files up to {(maxSize / 1024 / 1024).toFixed(0)}MB
            {maxFiles > 1 ? ` (max ${maxFiles} files)` : ''}
          </p>
          <p className="mt-2 text-xs text-gray-500 dark:text-gray-400">
            Your files are automatically deleted after conversion. We do not store your files permanently.
          </p>
          <input
            type="file"
            multiple={maxFiles > 1}
            onChange={handleFileSelect}
            className="hidden"
            id="file-upload"
            accept={accept}
          />
        </div>
      </div>
    </div>
  );
};

export default FileUploader;
