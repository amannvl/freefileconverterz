import React from 'react';

interface FormatSelectorProps {
  formats: string[];
  selectedFormat: string;
  onSelectFormat: (format: string) => void;
}

const FormatSelector: React.FC<FormatSelectorProps> = ({ 
  formats, 
  selectedFormat, 
  onSelectFormat 
}) => {
  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
      <div className="mb-4">
        <label htmlFor="outputFormat" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Convert to:
        </label>
        <select
          id="outputFormat"
          value={selectedFormat}
          onChange={(e) => onSelectFormat(e.target.value)}
          className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md dark:bg-gray-700 dark:border-gray-600 dark:text-white"
        >
          {formats.map(format => (
            <option key={format} value={format}>
              {format.toUpperCase()}
            </option>
          ))}
        </select>
      </div>
    </div>
  );
};

export default FormatSelector;
