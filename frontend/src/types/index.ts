export type FileStatus = 'uploading' | 'converting' | 'completed' | 'error';

export interface FileWithStatus {
  id: string;
  file: File;
  status: FileStatus;
  progress: number;
  error?: string;
  downloadUrl?: string;
}

export interface SupportedFormats {
  [key: string]: {
    input: string[];
    output: string[];
  };
}
